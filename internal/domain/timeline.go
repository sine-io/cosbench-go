package domain

import "time"

type TimelinePoint struct {
	Timestamp      time.Time `json:"timestamp"`
	OperationCount int64     `json:"operation_count"`
	ByteCount      int64     `json:"byte_count"`
	ErrorCount     int64     `json:"error_count"`
	AvgLatencyMs   float64   `json:"avg_latency_ms"`
}

type JobTimeline struct {
	JobID     string                     `json:"job_id"`
	UpdatedAt time.Time                  `json:"updated_at"`
	Job       []TimelinePoint            `json:"job"`
	Stages    map[string][]TimelinePoint `json:"stages,omitempty"`
}

type JobMatrixRow struct {
	JobID          string    `json:"job_id"`
	Name           string    `json:"name"`
	Status         JobStatus `json:"status"`
	StageCount     int       `json:"stage_count"`
	OperationCount int64     `json:"operation_count"`
	ByteCount      int64     `json:"byte_count"`
	ErrorCount     int64     `json:"error_count"`
	UpdatedAt      time.Time `json:"updated_at"`
}
