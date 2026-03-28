package controlplane

import (
	"sort"
	"time"

	"github.com/sine-io/cosbench-go/internal/domain"
)

func (m *Manager) GetDriverOverview(driverID string) (domain.DriverOverview, bool) {
	m.mu.Lock()
	m.refreshDriverHealthLocked(time.Now().UTC())
	driver, ok := m.drivers[driverID]
	m.mu.Unlock()
	if !ok {
		return domain.DriverOverview{}, false
	}
	missions := m.ListDriverMissions(driverID)
	logs := m.GetDriverLogs(driverID)
	active := 0
	for _, mission := range missions {
		if mission.Status == domain.MissionStatusClaimed || mission.Status == domain.MissionStatusRunning {
			active++
		}
	}
	return domain.DriverOverview{
		Driver:             driver,
		ActiveMissionCount: active,
		MissionCount:       len(missions),
		LogCount:           len(logs),
	}, true
}

func (m *Manager) ListDriverMissions(driverID string) []domain.Mission {
	m.mu.RLock()
	defer m.mu.RUnlock()
	items := make([]domain.Mission, 0)
	for _, mission := range m.missions {
		if mission.Lease != nil && mission.Lease.DriverID == driverID {
			items = append(items, mission)
		}
	}
	sort.Slice(items, func(i, j int) bool {
		if items[i].UpdatedAt.Equal(items[j].UpdatedAt) {
			return items[i].ID > items[j].ID
		}
		return items[i].UpdatedAt.After(items[j].UpdatedAt)
	})
	return items
}

func (m *Manager) GetDriverWorkerState(driverID string) (domain.DriverWorkerState, bool) {
	if _, ok := m.GetDriverNode(driverID); !ok {
		return domain.DriverWorkerState{}, false
	}
	active := 0
	for _, mission := range m.ListDriverMissions(driverID) {
		if mission.Status == domain.MissionStatusClaimed || mission.Status == domain.MissionStatusRunning {
			active++
		}
	}
	return domain.DriverWorkerState{DriverID: driverID, ActiveMissionCount: active}, true
}

func (m *Manager) GetDriverLogs(driverID string) []domain.JobEvent {
	missions := m.ListDriverMissions(driverID)
	jobSeen := map[string]struct{}{}
	logs := make([]domain.JobEvent, 0)
	for _, mission := range missions {
		if _, ok := jobSeen[mission.JobID]; ok {
			continue
		}
		jobSeen[mission.JobID] = struct{}{}
		logs = append(logs, m.GetJobEvents(mission.JobID)...)
	}
	sort.Slice(logs, func(i, j int) bool {
		if logs[i].OccurredAt.Equal(logs[j].OccurredAt) {
			return logs[i].Message > logs[j].Message
		}
		return logs[i].OccurredAt.After(logs[j].OccurredAt)
	})
	return logs
}

func latestMissionUpdate(missions []domain.Mission) *time.Time {
	if len(missions) == 0 {
		return nil
	}
	latest := missions[0].UpdatedAt
	for _, mission := range missions[1:] {
		if mission.UpdatedAt.After(latest) {
			latest = mission.UpdatedAt
		}
	}
	return &latest
}
