package controlplane

import (
	"sort"
	"time"

	"github.com/sine-io/cosbench-go/internal/domain"
)

func (m *Manager) PutDriverNode(driver domain.DriverNode) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.drivers[driver.ID] = driver
	return m.store.SaveDriverNode(driver)
}

func (m *Manager) SetRemoteScheduling(enabled bool) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.remoteScheduling = enabled
}

func (m *Manager) SetDriverHeartbeatTimeout(timeout time.Duration) {
	m.mu.Lock()
	defer m.mu.Unlock()
	if timeout <= 0 {
		m.driverHeartbeatTimeout = defaultDriverHeartbeatTimeout
		return
	}
	m.driverHeartbeatTimeout = timeout
}

func (m *Manager) ListDriverNodes() []domain.DriverNode {
	m.mu.Lock()
	m.refreshDriverHealthLocked(time.Now().UTC())
	items := make([]domain.DriverNode, 0, len(m.drivers))
	for _, driver := range m.drivers {
		items = append(items, driver)
	}
	m.mu.Unlock()
	sort.Slice(items, func(i, j int) bool {
		if items[i].RegisteredAt.Equal(items[j].RegisteredAt) {
			return items[i].ID > items[j].ID
		}
		return items[i].RegisteredAt.After(items[j].RegisteredAt)
	})
	return items
}

func (m *Manager) PutMission(mission domain.Mission) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.missions[mission.ID] = mission
	return m.store.SaveMission(mission)
}

func (m *Manager) GetMission(id string) (domain.Mission, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	mission, ok := m.missions[id]
	return mission, ok
}

func (m *Manager) ListMissions() []domain.Mission {
	m.mu.RLock()
	defer m.mu.RUnlock()
	items := make([]domain.Mission, 0, len(m.missions))
	for _, mission := range m.missions {
		items = append(items, mission)
	}
	sort.Slice(items, func(i, j int) bool {
		if items[i].UpdatedAt.Equal(items[j].UpdatedAt) {
			return items[i].ID > items[j].ID
		}
		return items[i].UpdatedAt.After(items[j].UpdatedAt)
	})
	return items
}
