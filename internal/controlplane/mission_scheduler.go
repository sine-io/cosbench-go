package controlplane

import (
	"errors"
	"sort"
	"time"

	"github.com/sine-io/cosbench-go/internal/domain"
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
	now := time.Now().UTC()
	missions := make([]domain.Mission, 0, len(stage.Works))
	for _, work := range stage.Works {
		mission := domain.Mission{
			ID:        newID("mission"),
			JobID:     jobID,
			StageName: stage.Name,
			WorkName:  work.Name,
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
