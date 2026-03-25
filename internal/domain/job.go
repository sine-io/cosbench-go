package domain

import "time"

type JobStatus string

const (
	JobStatusCreated     JobStatus = "created"
	JobStatusRunning     JobStatus = "running"
	JobStatusCancelling  JobStatus = "cancelling"
	JobStatusCancelled   JobStatus = "cancelled"
	JobStatusSucceeded   JobStatus = "succeeded"
	JobStatusFailed      JobStatus = "failed"
	JobStatusInterrupted JobStatus = "interrupted"
)

type EventLevel string

const (
	EventLevelInfo  EventLevel = "info"
	EventLevelError EventLevel = "error"
)

type Job struct {
	ID           string        `json:"id"`
	Name         string        `json:"name"`
	Status       JobStatus     `json:"status"`
	Workload     Workload      `json:"workload"`
	RawXML       string        `json:"raw_xml,omitempty"`
	EndpointID   string        `json:"endpoint_id,omitempty"`
	EndpointName string        `json:"endpoint_name,omitempty"`
	CreatedAt    time.Time     `json:"created_at"`
	StartedAt    *time.Time    `json:"started_at,omitempty"`
	FinishedAt   *time.Time    `json:"finished_at,omitempty"`
	ErrorMessage string        `json:"error_message,omitempty"`
	Stages       []StageState  `json:"stages"`
	Metrics      MetricsSummary `json:"metrics"`
}

type StageState struct {
	Name         string         `json:"name"`
	Status       JobStatus      `json:"status"`
	StartedAt    *time.Time     `json:"started_at,omitempty"`
	FinishedAt   *time.Time     `json:"finished_at,omitempty"`
	ErrorMessage string         `json:"error_message,omitempty"`
	Metrics      MetricsSummary `json:"metrics"`
	Works        int            `json:"works"`
	WorkResults  []WorkSummary  `json:"work_results,omitempty"`
}

type WorkSummary struct {
	Name         string         `json:"name"`
	ErrorMessage string         `json:"error_message,omitempty"`
	Metrics      MetricsSummary `json:"metrics"`
}

type JobEvent struct {
	JobID      string     `json:"job_id"`
	OccurredAt time.Time  `json:"occurred_at"`
	Level      EventLevel `json:"level"`
	Message    string     `json:"message"`
}

type JobResult struct {
	JobID       string         `json:"job_id"`
	UpdatedAt   time.Time      `json:"updated_at"`
	Metrics     MetricsSummary `json:"metrics"`
	StageTotals []StageState   `json:"stage_totals"`
}

type MetricsSummary struct {
	OperationCount int64              `json:"operation_count"`
	ByteCount      int64              `json:"byte_count"`
	ErrorCount     int64              `json:"error_count"`
	TotalLatencyMs int64              `json:"total_latency_ms"`
	DurationMs     int64              `json:"duration_ms"`
	AvgLatencyMs   float64            `json:"avg_latency_ms"`
	P50LatencyMs   float64            `json:"p50_latency_ms"`
	P95LatencyMs   float64            `json:"p95_latency_ms"`
	P99LatencyMs   float64            `json:"p99_latency_ms"`
	OpsPerSecond   float64            `json:"ops_per_second"`
	ByOperation    []OperationMetrics `json:"by_operation,omitempty"`
}

type OperationMetrics struct {
	Operation      string  `json:"operation"`
	OperationCount int64   `json:"operation_count"`
	ByteCount      int64   `json:"byte_count"`
	ErrorCount     int64   `json:"error_count"`
	AvgLatencyMs   float64 `json:"avg_latency_ms"`
	P50LatencyMs   float64 `json:"p50_latency_ms"`
	P95LatencyMs   float64 `json:"p95_latency_ms"`
	P99LatencyMs   float64 `json:"p99_latency_ms"`
}

func NewStageStates(w Workload) []StageState {
	states := make([]StageState, 0, len(w.Workflow.Stages))
	for _, stage := range w.Workflow.Stages {
		states = append(states, StageState{
			Name:   stage.Name,
			Status: JobStatusCreated,
			Works:  len(stage.Works),
		})
	}
	return states
}
