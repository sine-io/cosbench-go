package controlplane

import (
	"errors"
	"sort"
	"time"

	"github.com/sine-io/cosbench-go/internal/domain"
	legacyexec "github.com/sine-io/cosbench-go/internal/domain/execution"
	"github.com/sine-io/cosbench-go/internal/reporting"
)

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
		mission := domain.Mission{
			ID:        newID("mission"),
			JobID:     jobID,
			StageName: stage.Name,
			WorkName:  work.Name,
			Work:      resolvedWork,
			Status:    domain.MissionStatusCreated,
			CreatedAt: now,
			UpdatedAt: now,
		}
		if err := m.PutMission(mission); err != nil {
			return nil, err
		}
		missions = append(missions, mission)
	}
	return missions, nil
}

func (m *Manager) ClaimMission(driverID string, leaseTTL time.Duration) (domain.Mission, bool, error) {
	if leaseTTL <= 0 {
		leaseTTL = 30 * time.Second
	}
	if _, ok := m.GetDriverNode(driverID); !ok {
		return domain.Mission{}, false, errors.New("driver not found")
	}

	m.mu.Lock()
	defer m.mu.Unlock()

	candidates := make([]domain.Mission, 0, len(m.missions))
	for _, mission := range m.missions {
		if mission.Status == domain.MissionStatusCreated || mission.Status == domain.MissionStatusExpired {
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
	now := time.Now().UTC()
	mission := candidates[0]
	mission.Status = domain.MissionStatusClaimed
	mission.UpdatedAt = now
	mission.Lease = &domain.MissionLease{
		DriverID:  driverID,
		ClaimedAt: &now,
		ExpiresAt: now.Add(leaseTTL),
	}
	m.missions[mission.ID] = mission
	if err := m.store.SaveMission(mission); err != nil {
		return domain.Mission{}, false, err
	}
	return mission, true, nil
}

func (m *Manager) AppendMissionEvents(missionID string, events []domain.JobEvent) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	mission, ok := m.missions[missionID]
	if !ok {
		return errors.New("mission not found")
	}
	if mission.Status == domain.MissionStatusClaimed {
		mission.Status = domain.MissionStatusRunning
	}
	mission.UpdatedAt = time.Now().UTC()
	m.missions[missionID] = mission
	for _, event := range events {
		event.JobID = mission.JobID
		m.events[mission.JobID] = append(m.events[mission.JobID], event)
	}
	return m.store.SaveMission(mission)
}

func (m *Manager) AppendMissionSamples(missionID string, samples []legacyexec.Sample) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	mission, ok := m.missions[missionID]
	if !ok {
		return errors.New("mission not found")
	}
	if mission.Status == domain.MissionStatusClaimed {
		mission.Status = domain.MissionStatusRunning
	}
	mission.UpdatedAt = time.Now().UTC()
	m.missions[missionID] = mission
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
	mission.Status = status
	mission.ErrorMessage = errorMessage
	mission.UpdatedAt = time.Now().UTC()
	m.missions[missionID] = mission
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
			mission, ok := findMissionForWork(m.missions, jobID, stage.Name, work.Name)
			if !ok {
				stageAllSucceeded = false
				continue
			}
			stageAnyStarted = true
			samples := append([]legacyexec.Sample(nil), m.missionSamples[mission.ID]...)
			summary := reporting.Summarize(samples)
			stageParts = append(stageParts, summary)
			workSummaries = append(workSummaries, domain.WorkSummary{
				Name:         work.Name,
				ErrorMessage: mission.ErrorMessage,
				Metrics:      summary,
			})
			stageSamples = append(stageSamples, samples...)
			jobSamples = append(jobSamples, samples...)

			switch mission.Status {
			case domain.MissionStatusSucceeded:
			case domain.MissionStatusFailed:
				stageAnyFailed = true
				anyFailed = true
				stageAllSucceeded = false
			case domain.MissionStatusClaimed, domain.MissionStatusRunning:
				stageAnyRunning = true
				anyRunning = true
				stageAllSucceeded = false
			default:
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
