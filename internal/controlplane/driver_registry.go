package controlplane

import (
	"errors"
	"strings"
	"time"

	"github.com/sine-io/cosbench-go/internal/domain"
)

func (m *Manager) RegisterDriverNode(input domain.DriverNode) (domain.DriverNode, error) {
	name := strings.TrimSpace(input.Name)
	if name == "" {
		return domain.DriverNode{}, errors.New("driver name is required")
	}
	now := time.Now().UTC()
	if input.ID == "" {
		input.ID = newID("drv")
	}
	if input.Mode == "" {
		input.Mode = domain.DriverModeDriver
	}
	input.Name = name
	input.Status = domain.DriverStatusHealthy
	input.RegisteredAt = now
	input.LastHeartbeatAt = &now
	if err := m.PutDriverNode(input); err != nil {
		return domain.DriverNode{}, err
	}
	return input, nil
}

func (m *Manager) GetDriverNode(id string) (domain.DriverNode, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	driver, ok := m.drivers[id]
	return driver, ok
}

func (m *Manager) RecordDriverHeartbeat(id string, at time.Time) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	driver, ok := m.drivers[id]
	if !ok {
		return errors.New("driver not found")
	}
	driver.Status = domain.DriverStatusHealthy
	driver.LastHeartbeatAt = &at
	m.drivers[id] = driver
	return m.store.SaveDriverNode(driver)
}
