package controlplane

import (
	"sort"

	"github.com/sine-io/cosbench-go/internal/domain"
)

func (m *Manager) PutWorkUnit(unit domain.WorkUnit) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.workUnits[unit.ID] = unit
	return m.store.SaveWorkUnit(unit)
}

func (m *Manager) ListWorkUnits(jobID, stageName, workName string) []domain.WorkUnit {
	m.mu.RLock()
	defer m.mu.RUnlock()
	items := make([]domain.WorkUnit, 0)
	for _, unit := range m.workUnits {
		if unit.JobID != jobID || unit.StageName != stageName || unit.WorkName != workName {
			continue
		}
		items = append(items, unit)
	}
	sort.Slice(items, func(i, j int) bool {
		if items[i].UnitIndex == items[j].UnitIndex {
			return items[i].ID < items[j].ID
		}
		return items[i].UnitIndex < items[j].UnitIndex
	})
	return items
}

func (m *Manager) ListMissionAttempts() []domain.MissionAttempt {
	m.mu.RLock()
	defer m.mu.RUnlock()
	items := make([]domain.MissionAttempt, 0, len(m.missions))
	for _, attempt := range m.missions {
		items = append(items, attempt)
	}
	sort.Slice(items, func(i, j int) bool {
		if items[i].CreatedAt.Equal(items[j].CreatedAt) {
			return items[i].ID < items[j].ID
		}
		return items[i].CreatedAt.Before(items[j].CreatedAt)
	})
	return items
}

func listWorkUnitsLocked(units map[string]domain.WorkUnit, jobID, stageName, workName string) []domain.WorkUnit {
	items := make([]domain.WorkUnit, 0)
	for _, unit := range units {
		if unit.JobID != jobID || unit.StageName != stageName || unit.WorkName != workName {
			continue
		}
		items = append(items, unit)
	}
	sort.Slice(items, func(i, j int) bool {
		if items[i].UnitIndex == items[j].UnitIndex {
			return items[i].ID < items[j].ID
		}
		return items[i].UnitIndex < items[j].UnitIndex
	})
	return items
}

func latestAttemptForUnitLocked(attempts map[string]domain.Mission, workUnitID string) (domain.MissionAttempt, bool) {
	var latest domain.MissionAttempt
	found := false
	for _, attempt := range attempts {
		if attempt.WorkUnitID != workUnitID {
			continue
		}
		if !found || attempt.Attempt > latest.Attempt || (attempt.Attempt == latest.Attempt && attempt.UpdatedAt.After(latest.UpdatedAt)) {
			latest = attempt
			found = true
		}
	}
	return latest, found
}
