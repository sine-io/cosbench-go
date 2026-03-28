package domain

import "time"

type DriverMode string

const (
	DriverModeController DriverMode = "controller"
	DriverModeDriver     DriverMode = "driver"
	DriverModeCombined   DriverMode = "combined"
)

type DriverStatus string

const (
	DriverStatusUnknown  DriverStatus = "unknown"
	DriverStatusHealthy  DriverStatus = "healthy"
	DriverStatusUnhealthy DriverStatus = "unhealthy"
)

type DriverNode struct {
	ID              string       `json:"id"`
	Name            string       `json:"name"`
	Mode            DriverMode   `json:"mode"`
	Status          DriverStatus `json:"status"`
	RegisteredAt    time.Time    `json:"registered_at"`
	LastHeartbeatAt *time.Time   `json:"last_heartbeat_at,omitempty"`
}

type DriverOverview struct {
	Driver             DriverNode `json:"driver"`
	ActiveMissionCount int        `json:"active_mission_count"`
	MissionCount       int        `json:"mission_count"`
	LogCount           int        `json:"log_count"`
}

type DriverWorkerState struct {
	DriverID           string `json:"driver_id"`
	ActiveMissionCount int    `json:"active_mission_count"`
}

type WorkUnitStatus string

const (
	WorkUnitStatusPending   WorkUnitStatus = "pending"
	WorkUnitStatusClaimed   WorkUnitStatus = "claimed"
	WorkUnitStatusRunning   WorkUnitStatus = "running"
	WorkUnitStatusSucceeded WorkUnitStatus = "succeeded"
	WorkUnitStatusFailed    WorkUnitStatus = "failed"
)

type WorkUnitSlice struct {
	WorkerIndex int `json:"worker_index"`
	WorkerCount int `json:"worker_count"`
}

type WorkUnit struct {
	ID        string         `json:"id"`
	JobID     string         `json:"job_id"`
	StageName string         `json:"stage_name"`
	WorkName  string         `json:"work_name"`
	UnitIndex int            `json:"unit_index"`
	UnitCount int            `json:"unit_count"`
	Work      Work           `json:"work"`
	Slice     WorkUnitSlice  `json:"slice"`
	Status    WorkUnitStatus `json:"status"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
}

type MissionStatus string

const (
	MissionAttemptStatusPending   MissionStatus = "pending"
	MissionStatusCreated   MissionStatus = "created"
	MissionStatusClaimed   MissionStatus = "claimed"
	MissionStatusRunning   MissionStatus = "running"
	MissionStatusSucceeded MissionStatus = "succeeded"
	MissionStatusFailed    MissionStatus = "failed"
	MissionStatusExpired   MissionStatus = "expired"
)

type MissionLease struct {
	DriverID  string     `json:"driver_id"`
	ClaimedAt *time.Time `json:"claimed_at,omitempty"`
	ExpiresAt time.Time  `json:"expires_at"`
}

type MissionAttempt struct {
	ID        string         `json:"id"`
	WorkUnitID string        `json:"work_unit_id"`
	WorkUnit  WorkUnit       `json:"work_unit"`
	JobID     string         `json:"job_id"`
	StageName string         `json:"stage_name"`
	WorkName  string         `json:"work_name"`
	Work      Work           `json:"work"`
	Attempt   int            `json:"attempt"`
	Status    MissionStatus  `json:"status"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	ErrorMessage string      `json:"error_message,omitempty"`
	ReceivedEventBatches  []string `json:"received_event_batches,omitempty"`
	ReceivedSampleBatches []string `json:"received_sample_batches,omitempty"`
	Lease     *MissionLease  `json:"lease,omitempty"`
}

type Mission = MissionAttempt
