package controlplane

import (
	"errors"
	"sort"
	"time"

	"github.com/sine-io/cosbench-go/internal/domain"
	legacyexec "github.com/sine-io/cosbench-go/internal/domain/execution"
	"github.com/sine-io/cosbench-go/internal/reporting"
)

const maxAttemptsPerWorkUnit = 3

func (m *Manager) ScheduleJobStage(jobID string) ([]domain.Mission, error) {
	job, ok := m.GetJob(jobID)
	if !ok {
		return nil, errors.New("job not found")
	}
	if len(job.Workload.Workflow.Stages) == 0 {
		return nil, errors.New("job has no stages")
	}
	stage := job.Workload.Workflow.Stages[0]
	endpoint, _ := m.GetEndpoint(job.EndpointID)
	now := time.Now().UTC()
	missions := make([]domain.Mission, 0, len(stage.Works))
	for _, work := range stage.Works {
		storageType, rawConfig := resolveStorage(work.Storage, endpoint)
		resolvedWork := work
		resolvedWork.Storage = &domain.StorageSpec{Type: storageType, Config: rawConfig}
		unitCount := work.Workers
		if unitCount <= 0 {
			unitCount = 1
		}
		for unitIndex := 1; unitIndex <= unitCount; unitIndex++ {
			unitWork := sliceWork(resolvedWork, unitIndex, unitCount)
			unit := domain.WorkUnit{
				ID:        newID("unit"),
				JobID:     jobID,
				StageName: stage.Name,
				WorkName:  work.Name,
				UnitIndex: unitIndex,
				UnitCount: unitCount,
				Work:      unitWork,
				Slice: domain.WorkUnitSlice{
					WorkerIndex: unitIndex,
					WorkerCount: unitCount,
				},
				Status:    domain.WorkUnitStatusPending,
				CreatedAt: now,
				UpdatedAt: now,
			}
			if err := m.PutWorkUnit(unit); err != nil {
				return nil, err
			}
			mission := m.createMissionAttemptLocked(unit, now, 1)
			if err := m.PutMission(mission); err != nil {
				return nil, err
			}
			missions = append(missions, mission)
		}
	}
	return missions, nil
}

func (m *Manager) ClaimMission(driverID string, leaseTTL time.Duration) (domain.Mission, bool, error) {
	if leaseTTL <= 0 {
		leaseTTL = 30 * time.Second
	}

	m.mu.Lock()
	defer m.mu.Unlock()
	now := time.Now().UTC()
	m.expireLeasesLocked(now)
	m.refreshDriverHealthLocked(now)
	driver, ok := m.drivers[driverID]
	if !ok {
		return domain.Mission{}, false, errors.New("driver not found")
	}
	if driver.Status != domain.DriverStatusHealthy {
		return domain.Mission{}, false, errors.New("driver is not healthy")
	}

	candidates := make([]domain.Mission, 0, len(m.missions))
	for _, mission := range m.missions {
		if mission.Status == domain.MissionAttemptStatusPending {
			candidates = append(candidates, mission)
		}
	}
	sort.Slice(candidates, func(i, j int) bool {
		if candidates[i].CreatedAt.Equal(candidates[j].CreatedAt) {
			return candidates[i].ID < candidates[j].ID
		}
		return candidates[i].CreatedAt.Before(candidates[j].CreatedAt)
	})
	if len(candidates) == 0 {
		return domain.Mission{}, false, nil
	}
	mission := candidates[0]
	mission.Status = domain.MissionStatusClaimed
	mission.UpdatedAt = now
	mission.Lease = &domain.MissionLease{
		DriverID:  driverID,
		ClaimedAt: &now,
		ExpiresAt: now.Add(leaseTTL),
	}
	m.missions[mission.ID] = mission
	if unit, ok := m.workUnits[mission.WorkUnitID]; ok {
		unit.Status = domain.WorkUnitStatusClaimed
		unit.UpdatedAt = now
		m.workUnits[unit.ID] = unit
		_ = m.store.SaveWorkUnit(unit)
		mission.WorkUnit = unit
	}
	if err := m.store.SaveMission(mission); err != nil {
		return domain.Mission{}, false, err
	}
	return mission, true, nil
}

func (m *Manager) expireLeasesLocked(now time.Time) {
	affectedJobs := map[string]struct{}{}
	for missionID, mission := range m.missions {
		if mission.Lease == nil {
			continue
		}
		if mission.Status != domain.MissionStatusClaimed && mission.Status != domain.MissionStatusRunning {
			continue
		}
		if mission.Lease.ExpiresAt.After(now) {
			continue
		}
		mission.Status = domain.MissionStatusExpired
		mission.ErrorMessage = "mission lease expired"
		mission.Lease = nil
		mission.UpdatedAt = now
		m.missions[missionID] = mission
		if unit, ok := m.workUnits[mission.WorkUnitID]; ok {
			unit.Status = domain.WorkUnitStatusPending
			unit.UpdatedAt = now
			m.workUnits[unit.ID] = unit
			_ = m.store.SaveWorkUnit(unit)
			if mission.Attempt < maxAttemptsPerWorkUnit {
				retry := m.createMissionAttemptLocked(unit, now, mission.Attempt+1)
				m.missions[retry.ID] = retry
				_ = m.store.SaveMission(retry)
			}
		}
		delete(m.missionSamples, missionID)
		_ = m.store.SaveMission(mission)
		affectedJobs[mission.JobID] = struct{}{}
	}
	for jobID := range affectedJobs {
		m.refreshJobFromMissionsLocked(jobID)
	}
}

func (m *Manager) AppendMissionEvents(missionID string, events []domain.JobEvent) error {
	return m.AppendMissionEventsBatch(missionID, newID("event-batch"), events)
}

func (m *Manager) AppendMissionEventsBatch(missionID string, batchID string, events []domain.JobEvent) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	mission, ok := m.missions[missionID]
	if !ok {
		return errors.New("mission not found")
	}
	if batchID != "" && containsBatchID(mission.ReceivedEventBatches, batchID) {
		return nil
	}
	if mission.Status == domain.MissionStatusClaimed {
		mission.Status = domain.MissionStatusRunning
	}
	mission.UpdatedAt = time.Now().UTC()
	if batchID != "" {
		mission.ReceivedEventBatches = append(mission.ReceivedEventBatches, batchID)
	}
	m.missions[missionID] = mission
	if unit, ok := m.workUnits[mission.WorkUnitID]; ok {
		unit.Status = domain.WorkUnitStatusRunning
		unit.UpdatedAt = mission.UpdatedAt
		m.workUnits[unit.ID] = unit
		_ = m.store.SaveWorkUnit(unit)
	}
	for _, event := range events {
		event.JobID = mission.JobID
		m.events[mission.JobID] = append(m.events[mission.JobID], event)
	}
	return m.store.SaveMission(mission)
}

func (m *Manager) AppendMissionSamples(missionID string, samples []legacyexec.Sample) error {
	return m.AppendMissionSamplesBatch(missionID, newID("sample-batch"), samples)
}

func (m *Manager) AppendMissionSamplesBatch(missionID string, batchID string, samples []legacyexec.Sample) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	mission, ok := m.missions[missionID]
	if !ok {
		return errors.New("mission not found")
	}
	if batchID != "" && containsBatchID(mission.ReceivedSampleBatches, batchID) {
		return nil
	}
	if mission.Status == domain.MissionStatusClaimed {
		mission.Status = domain.MissionStatusRunning
	}
	mission.UpdatedAt = time.Now().UTC()
	if batchID != "" {
		mission.ReceivedSampleBatches = append(mission.ReceivedSampleBatches, batchID)
	}
	m.missions[missionID] = mission
	if unit, ok := m.workUnits[mission.WorkUnitID]; ok {
		unit.Status = domain.WorkUnitStatusRunning
		unit.UpdatedAt = mission.UpdatedAt
		m.workUnits[unit.ID] = unit
		_ = m.store.SaveWorkUnit(unit)
	}
	m.missionSamples[missionID] = append(m.missionSamples[missionID], samples...)
	return m.store.SaveMission(mission)
}

func (m *Manager) CompleteMission(missionID string, status domain.MissionStatus, errorMessage string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	mission, ok := m.missions[missionID]
	if !ok {
		return errors.New("mission not found")
	}
	if mission.Status == domain.MissionStatusSucceeded || mission.Status == domain.MissionStatusFailed {
		return nil
	}
	mission.Status = status
	mission.ErrorMessage = errorMessage
	mission.UpdatedAt = time.Now().UTC()
	m.missions[missionID] = mission
	if unit, ok := m.workUnits[mission.WorkUnitID]; ok {
		switch status {
		case domain.MissionStatusSucceeded:
			unit.Status = domain.WorkUnitStatusSucceeded
		case domain.MissionStatusFailed:
			if mission.Attempt >= maxAttemptsPerWorkUnit {
				unit.Status = domain.WorkUnitStatusFailed
			} else {
				unit.Status = domain.WorkUnitStatusPending
			}
		default:
			unit.Status = domain.WorkUnitStatusRunning
		}
		unit.UpdatedAt = mission.UpdatedAt
		m.workUnits[unit.ID] = unit
		_ = m.store.SaveWorkUnit(unit)
		if status == domain.MissionStatusFailed && mission.Attempt < maxAttemptsPerWorkUnit {
			retry := m.createMissionAttemptLocked(unit, mission.UpdatedAt, mission.Attempt+1)
			m.missions[retry.ID] = retry
			_ = m.store.SaveMission(retry)
		}
	}
	if err := m.store.SaveMission(mission); err != nil {
		return err
	}
	m.refreshJobFromMissionsLocked(mission.JobID)
	return nil
}

func (m *Manager) refreshJobFromMissionsLocked(jobID string) {
	job, ok := m.jobs[jobID]
	if !ok {
		return
	}
	stageStates := domain.NewStageStates(job.Workload)
	stageTotals := make([]domain.StageState, 0, len(job.Workload.Workflow.Stages))
	stageTimelines := make(map[string][]domain.TimelinePoint)
	jobSamples := make([]legacyexec.Sample, 0)
	anyRunning := false
	anyFailed := false
	allSucceeded := len(job.Workload.Workflow.Stages) > 0

	for stageIndex, stage := range job.Workload.Workflow.Stages {
		stageState := stageStates[stageIndex]
		stageSamples := make([]legacyexec.Sample, 0)
		workSummaries := make([]domain.WorkSummary, 0, len(stage.Works))
		stageParts := make([]domain.MetricsSummary, 0, len(stage.Works))
		stageAllSucceeded := len(stage.Works) > 0
		stageAnyStarted := false
		stageAnyFailed := false
		stageAnyRunning := false

		for _, work := range stage.Works {
			units := listWorkUnitsLocked(m.workUnits, jobID, stage.Name, work.Name)
			if len(units) == 0 {
				stageAllSucceeded = false
				continue
			}
			workSamples := make([]legacyexec.Sample, 0)
			workSummary := domain.WorkSummary{Name: work.Name}
			workAllSucceeded := len(units) > 0
			for _, unit := range units {
				attempt, ok := latestAttemptForUnitLocked(m.missions, unit.ID)
				if ok {
					stageAnyStarted = true
					samples := append([]legacyexec.Sample(nil), m.missionSamples[attempt.ID]...)
					workSamples = append(workSamples, samples...)
					stageSamples = append(stageSamples, samples...)
					jobSamples = append(jobSamples, samples...)
					if attempt.ErrorMessage != "" {
						workSummary.ErrorMessage = attempt.ErrorMessage
					}
				}
				switch unit.Status {
				case domain.WorkUnitStatusSucceeded:
				case domain.WorkUnitStatusFailed:
					stageAnyFailed = true
					anyFailed = true
					workAllSucceeded = false
				case domain.WorkUnitStatusClaimed, domain.WorkUnitStatusRunning:
					stageAnyRunning = true
					anyRunning = true
					workAllSucceeded = false
				default:
					workAllSucceeded = false
				}
			}
			workSummary.Metrics = reporting.Summarize(workSamples)
			stageParts = append(stageParts, workSummary.Metrics)
			workSummaries = append(workSummaries, workSummary)
			if !workAllSucceeded && workSummary.Metrics.OperationCount == 0 {
				stageAllSucceeded = false
			}
		}

		stageState.WorkResults = workSummaries
		stageState.Metrics = reporting.Merge(stageParts...)
		stageTimelines[stage.Name] = reporting.BuildTimeline(stageSamples, time.Second)
		if stageAnyFailed {
			stageState.Status = domain.JobStatusFailed
		} else if stageAllSucceeded {
			stageState.Status = domain.JobStatusSucceeded
		} else if stageAnyRunning || stageAnyStarted {
			stageState.Status = domain.JobStatusRunning
		} else {
			stageState.Status = domain.JobStatusCreated
			allSucceeded = false
		}
		if stageState.Status != domain.JobStatusSucceeded {
			allSucceeded = false
		}
		stageTotals = append(stageTotals, stageState)
		stageStates[stageIndex] = stageState
	}

	job.Stages = stageStates
	job.Metrics = reporting.Merge(extractStageMetrics(stageStates)...)
	switch {
	case anyFailed:
		job.Status = domain.JobStatusFailed
	case allSucceeded:
		job.Status = domain.JobStatusSucceeded
	case anyRunning:
		job.Status = domain.JobStatusRunning
	default:
		job.Status = domain.JobStatusCreated
	}
	m.jobs[jobID] = job
	m.results[jobID] = domain.JobResult{
		JobID:       jobID,
		UpdatedAt:   time.Now().UTC(),
		Metrics:     job.Metrics,
		StageTotals: stageTotals,
	}
	m.timelines[jobID] = domain.JobTimeline{
		JobID:     jobID,
		UpdatedAt: time.Now().UTC(),
		Job:       reporting.BuildTimeline(jobSamples, time.Second),
		Stages:    cloneStageTimelines(stageTimelines),
	}
	m.persistLocked(jobID)
}

func findMissionForWork(missions map[string]domain.Mission, jobID, stageName, workName string) (domain.Mission, bool) {
	for _, mission := range missions {
		if mission.JobID == jobID && mission.StageName == stageName && mission.WorkName == workName {
			return mission, true
		}
	}
	return domain.Mission{}, false
}

func sliceWork(work domain.Work, unitIndex, unitCount int) domain.Work {
	sliced := work
	sliced.Workers = 1
	if work.TotalOps > 0 {
		sliced.TotalOps = splitShare(work.TotalOps, unitIndex, unitCount)
	}
	if work.TotalBytes > 0 {
		sliced.TotalBytes = int64(splitShareInt64(work.TotalBytes, unitIndex, unitCount))
	}
	return sliced
}

func splitShare(total int, unitIndex, unitCount int) int {
	if unitCount <= 0 {
		return total
	}
	base := total / unitCount
	remainder := total % unitCount
	share := base
	if unitIndex <= remainder {
		share++
	}
	if share < 0 {
		return 0
	}
	return share
}

func splitShareInt64(total int64, unitIndex, unitCount int) int64 {
	if unitCount <= 0 {
		return total
	}
	base := total / int64(unitCount)
	remainder := total % int64(unitCount)
	share := base
	if int64(unitIndex) <= remainder {
		share++
	}
	if share < 0 {
		return 0
	}
	return share
}

func containsBatchID(items []string, batchID string) bool {
	for _, item := range items {
		if item == batchID {
			return true
		}
	}
	return false
}

func (m *Manager) createMissionAttemptLocked(unit domain.WorkUnit, now time.Time, attempt int) domain.Mission {
	return domain.Mission{
		ID:         newID("mission"),
		WorkUnitID: unit.ID,
		WorkUnit:   unit,
		JobID:      unit.JobID,
		StageName:  unit.StageName,
		WorkName:   unit.WorkName,
		Work:       unit.Work,
		Attempt:    attempt,
		Status:     domain.MissionAttemptStatusPending,
		CreatedAt:  now,
		UpdatedAt:  now,
	}
}
