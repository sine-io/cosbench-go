package controlplane

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/sine-io/cosbench-go/internal/domain"
	"github.com/sine-io/cosbench-go/internal/executor"
	"github.com/sine-io/cosbench-go/internal/infrastructure/config"
	storagefactory "github.com/sine-io/cosbench-go/internal/infrastructure/storage"
	"github.com/sine-io/cosbench-go/internal/reporting"
	"github.com/sine-io/cosbench-go/internal/snapshot"
	"github.com/sine-io/cosbench-go/internal/workloadxml"
)

type Manager struct {
	mu        sync.RWMutex
	store     *snapshot.Store
	jobs      map[string]domain.Job
	results   map[string]domain.JobResult
	events    map[string][]domain.JobEvent
	endpoints map[string]domain.EndpointConfig
	running   map[string]context.CancelFunc
}

func New(store *snapshot.Store) (*Manager, error) {
	m := &Manager{
		store:     store,
		jobs:      map[string]domain.Job{},
		results:   map[string]domain.JobResult{},
		events:    map[string][]domain.JobEvent{},
		endpoints: map[string]domain.EndpointConfig{},
		running:   map[string]context.CancelFunc{},
	}
	if err := m.loadSnapshots(); err != nil {
		return nil, err
	}
	return m, nil
}

func (m *Manager) loadSnapshots() error {
	endpoints, err := m.store.LoadEndpoints()
	if err != nil {
		return err
	}
	for _, endpoint := range endpoints {
		m.endpoints[endpoint.ID] = endpoint
	}
	jobs, err := m.store.LoadJobs()
	if err != nil {
		return err
	}
	for _, job := range jobs {
		if job.Status == domain.JobStatusRunning {
			job.Status = domain.JobStatusInterrupted
			job.ErrorMessage = "job was running before restart"
			for i := range job.Stages {
				if job.Stages[i].Status == domain.JobStatusRunning {
					job.Stages[i].Status = domain.JobStatusInterrupted
					job.Stages[i].ErrorMessage = job.ErrorMessage
				}
			}
		} else if job.Status == domain.JobStatusCancelling {
			job.Status = domain.JobStatusCancelled
			job.ErrorMessage = "job cancellation was in progress before restart"
			for i := range job.Stages {
				if job.Stages[i].Status == domain.JobStatusCancelling {
					job.Stages[i].Status = domain.JobStatusCancelled
					job.Stages[i].ErrorMessage = job.ErrorMessage
				}
			}
		}
		m.jobs[job.ID] = job
		if result, err := m.store.LoadResult(job.ID); err == nil {
			m.results[job.ID] = result
		}
		if events, err := m.store.LoadEvents(job.ID); err == nil {
			m.events[job.ID] = events
		}
	}
	return nil
}

func (m *Manager) CreateEndpoint(input domain.EndpointConfig) (domain.EndpointConfig, error) {
	if err := input.Validate(); err != nil {
		return domain.EndpointConfig{}, err
	}
	now := time.Now().UTC()
	input.ID = newID("ep")
	input.CreatedAt = now
	input.UpdatedAt = now
	m.mu.Lock()
	m.endpoints[input.ID] = input
	m.mu.Unlock()
	return input, m.store.SaveEndpoint(input)
}

func (m *Manager) ListEndpoints() []domain.EndpointConfig {
	m.mu.RLock()
	defer m.mu.RUnlock()
	items := make([]domain.EndpointConfig, 0, len(m.endpoints))
	for _, endpoint := range m.endpoints {
		items = append(items, endpoint)
	}
	sortEndpoints(items)
	return items
}

func (m *Manager) CreateJobFromXML(raw []byte, endpointID string) (domain.Job, error) {
	workload, err := workloadxml.Parse(raw)
	if err != nil {
		return domain.Job{}, err
	}
	job := domain.Job{
		ID:         newID("job"),
		Name:       workload.Name,
		Status:     domain.JobStatusCreated,
		Workload:   workload,
		RawXML:     string(raw),
		CreatedAt:  time.Now().UTC(),
		Stages:     domain.NewStageStates(workload),
		EndpointID: endpointID,
	}
	if endpointID != "" {
		endpoint, ok := m.GetEndpoint(endpointID)
		if !ok {
			return domain.Job{}, errors.New("endpoint not found")
		}
		job.EndpointName = endpoint.Name
	}
	m.mu.Lock()
	m.jobs[job.ID] = job
	m.events[job.ID] = []domain.JobEvent{newEvent(job.ID, domain.EventLevelInfo, "job created")}
	m.mu.Unlock()
	if err := m.store.SaveJob(job); err != nil {
		return domain.Job{}, err
	}
	if err := m.store.SaveEvents(job.ID, m.events[job.ID]); err != nil {
		return domain.Job{}, err
	}
	return job, nil
}

func (m *Manager) StartJob(ctx context.Context, jobID string) error {
	m.mu.RLock()
	job, ok := m.jobs[jobID]
	endpoint, _ := m.endpoints[job.EndpointID]
	m.mu.RUnlock()
	if !ok {
		return errors.New("job not found")
	}
	if job.Status == domain.JobStatusRunning {
		return errors.New("job already running")
	}
	if err := preflightJob(job, endpoint); err != nil {
		return err
	}

	m.mu.Lock()
	job, ok = m.jobs[jobID]
	if !ok {
		m.mu.Unlock()
		return errors.New("job not found")
	}
	if job.Status == domain.JobStatusRunning {
		m.mu.Unlock()
		return errors.New("job already running")
	}
	now := time.Now().UTC()
	runCtx, cancel := context.WithCancel(ctx)
	job.Status = domain.JobStatusRunning
	job.StartedAt = &now
	m.jobs[jobID] = job
	m.running[jobID] = cancel
	m.appendEventLocked(jobID, domain.EventLevelInfo, "job started")
	m.persistLocked(jobID)
	m.mu.Unlock()
	go m.runJob(runCtx, jobID)
	return nil
}

func (m *Manager) CancelJob(jobID string) error {
	m.mu.Lock()
	job, ok := m.jobs[jobID]
	if !ok {
		m.mu.Unlock()
		return errors.New("job not found")
	}
	if job.Status != domain.JobStatusRunning {
		m.mu.Unlock()
		return errors.New("job is not running")
	}
	cancel := m.running[jobID]
	if cancel == nil {
		m.mu.Unlock()
		return errors.New("job has no active cancel function")
	}
	job.Status = domain.JobStatusCancelling
	for i := range job.Stages {
		if job.Stages[i].Status == domain.JobStatusRunning {
			job.Stages[i].Status = domain.JobStatusCancelling
			job.Stages[i].ErrorMessage = "job cancellation requested"
			break
		}
	}
	job.ErrorMessage = "job cancellation requested"
	m.jobs[jobID] = job
	m.appendEventLocked(jobID, domain.EventLevelInfo, "job cancellation requested")
	m.persistLocked(jobID)
	m.mu.Unlock()
	cancel()
	return nil
}

func preflightJob(job domain.Job, endpoint domain.EndpointConfig) error {
	for _, stage := range job.Workload.Workflow.Stages {
		for _, work := range stage.Works {
			storageType, rawConfig := resolveStorage(work.Storage, endpoint)
			if storageType == "" {
				return fmt.Errorf("stage %s work %s: no storage configured", stage.Name, work.Name)
			}
			adapter, err := storagefactory.NewAdapter(storageType, rawConfig)
			if err != nil {
				return fmt.Errorf("stage %s work %s: %w", stage.Name, work.Name, err)
			}
			if err := adapter.Init(config.ParseKVConfig(rawConfig)); err != nil {
				return fmt.Errorf("stage %s work %s: %w", stage.Name, work.Name, err)
			}
			if err := executor.ValidateWork(work); err != nil {
				_ = adapter.Dispose()
				return fmt.Errorf("stage %s work %s: %w", stage.Name, work.Name, err)
			}
			if err := adapter.Dispose(); err != nil {
				return fmt.Errorf("stage %s work %s: %w", stage.Name, work.Name, err)
			}
		}
	}
	return nil
}

func (m *Manager) runJob(ctx context.Context, jobID string) {
	defer m.clearRunning(jobID)
	m.mu.RLock()
	job := m.jobs[jobID]
	endpoint, _ := m.endpoints[job.EndpointID]
	m.mu.RUnlock()

	stageStates := make([]domain.StageState, len(job.Stages))
	copy(stageStates, job.Stages)
	stageTotals := make([]domain.StageState, 0, len(job.Workload.Workflow.Stages))
	var overall []domain.MetricsSummary
	adapters := storagefactory.NewRunAdapters()
	defer adapters.Close()

	for stageIndex, stage := range job.Workload.Workflow.Stages {
		stageStart := time.Now().UTC()
		stageState := &stageStates[stageIndex]
		stageState.Status = domain.JobStatusRunning
		stageState.StartedAt = &stageStart
		m.updateJobStages(jobID, stageStates, domain.EventLevelInfo, fmt.Sprintf("stage %s started", stage.Name))

		stageSummaries := make([]domain.MetricsSummary, 0, len(stage.Works))
		workSummaries := make([]domain.WorkSummary, 0, len(stage.Works))
		for _, work := range stage.Works {
			storageType, rawConfig := resolveStorage(work.Storage, endpoint)
			if storageType == "" {
				m.failJob(jobID, stageStates, stageIndex, fmt.Sprintf("stage %s: no storage configured", stage.Name))
				return
			}
			adapter, shared, err := adapters.Acquire(storageType, rawConfig)
			if err != nil {
				m.failJob(jobID, stageStates, stageIndex, err.Error())
				return
			}
			stageExecutor := executor.StageExecutor{Storage: adapter}
			workResult := stageExecutor.RunWork(ctx, work)
			if !shared {
				_ = adapter.Dispose()
			}
			if workResult.Err != nil {
				if errors.Is(workResult.Err, context.Canceled) {
					stageSummaries = append(stageSummaries, workResult.Summary)
					workSummaries = append(workSummaries, domain.WorkSummary{Name: work.Name, Metrics: workResult.Summary, ErrorMessage: "job cancelled by user"})
					stageState.WorkResults = workSummaries
					stageState.Metrics = reporting.Merge(stageSummaries...)
					m.cancelJob(jobID, stageStates, stageIndex, "job cancelled by user")
					return
				}
				stageState.WorkResults = append(workSummaries, domain.WorkSummary{Name: work.Name, ErrorMessage: workResult.Err.Error()})
				m.failJob(jobID, stageStates, stageIndex, workResult.Err.Error())
				return
			}
			stageSummaries = append(stageSummaries, workResult.Summary)
			workSummaries = append(workSummaries, domain.WorkSummary{Name: work.Name, Metrics: workResult.Summary})
			m.appendEvent(jobID, domain.EventLevelInfo, fmt.Sprintf("work %s finished", work.Name))
		}
		stageEnd := time.Now().UTC()
		stageState.Status = domain.JobStatusSucceeded
		stageState.FinishedAt = &stageEnd
		stageState.Metrics = reporting.Merge(stageSummaries...)
		stageState.WorkResults = workSummaries
		stageTotals = append(stageTotals, *stageState)
		overall = append(overall, stageState.Metrics)
		m.updateJobStages(jobID, stageStates, domain.EventLevelInfo, fmt.Sprintf("stage %s finished", stage.Name))
	}

	finished := time.Now().UTC()
	m.mu.Lock()
	job = m.jobs[jobID]
	job.Status = domain.JobStatusSucceeded
	job.FinishedAt = &finished
	job.Stages = stageStates
	job.Metrics = reporting.Merge(overall...)
	if job.Metrics.ErrorCount > 0 {
		job.ErrorMessage = "job completed with operation errors"
		for i := range job.Stages {
			if job.Stages[i].Metrics.ErrorCount > 0 {
				if job.Stages[i].ErrorMessage == "" {
					job.Stages[i].ErrorMessage = "stage completed with operation errors"
				}
			}
		}
		m.appendEventLocked(jobID, domain.EventLevelError, job.ErrorMessage)
	}
	m.jobs[jobID] = job
	m.results[jobID] = domain.JobResult{JobID: jobID, UpdatedAt: finished, Metrics: job.Metrics, StageTotals: stageTotals}
	m.appendEventLocked(jobID, domain.EventLevelInfo, "job finished")
	m.persistLocked(jobID)
	m.mu.Unlock()
}

func (m *Manager) failJob(jobID string, stages []domain.StageState, stageIndex int, message string) {
	finished := time.Now().UTC()
	stages[stageIndex].Status = domain.JobStatusFailed
	stages[stageIndex].FinishedAt = &finished
	stages[stageIndex].ErrorMessage = message
	m.mu.Lock()
	job := m.jobs[jobID]
	job.Status = domain.JobStatusFailed
	job.ErrorMessage = message
	job.FinishedAt = &finished
	job.Stages = stages
	job.Metrics = reporting.Merge(extractStageMetrics(stages)...)
	m.jobs[jobID] = job
	m.results[jobID] = domain.JobResult{JobID: jobID, UpdatedAt: finished, Metrics: job.Metrics, StageTotals: stages}
	m.appendEventLocked(jobID, domain.EventLevelError, message)
	m.persistLocked(jobID)
	m.mu.Unlock()
}

func (m *Manager) cancelJob(jobID string, stages []domain.StageState, stageIndex int, message string) {
	finished := time.Now().UTC()
	stages[stageIndex].Status = domain.JobStatusCancelled
	stages[stageIndex].FinishedAt = &finished
	stages[stageIndex].ErrorMessage = message
	m.mu.Lock()
	job := m.jobs[jobID]
	job.Status = domain.JobStatusCancelled
	job.ErrorMessage = message
	job.FinishedAt = &finished
	job.Stages = stages
	job.Metrics = reporting.Merge(extractStageMetrics(stages)...)
	m.jobs[jobID] = job
	m.results[jobID] = domain.JobResult{JobID: jobID, UpdatedAt: finished, Metrics: job.Metrics, StageTotals: stages}
	m.appendEventLocked(jobID, domain.EventLevelInfo, message)
	m.persistLocked(jobID)
	m.mu.Unlock()
}

func (m *Manager) updateJobStages(jobID string, stages []domain.StageState, level domain.EventLevel, message string) {
	m.mu.Lock()
	job := m.jobs[jobID]
	job.Stages = append([]domain.StageState(nil), stages...)
	job.Metrics = reporting.Merge(extractStageMetrics(stages)...)
	m.jobs[jobID] = job
	m.appendEventLocked(jobID, level, message)
	m.persistLocked(jobID)
	m.mu.Unlock()
}

func (m *Manager) appendEvent(jobID string, level domain.EventLevel, message string) {
	m.mu.Lock()
	m.appendEventLocked(jobID, level, message)
	m.persistLocked(jobID)
	m.mu.Unlock()
}

func (m *Manager) appendEventLocked(jobID string, level domain.EventLevel, message string) {
	m.events[jobID] = append(m.events[jobID], newEvent(jobID, level, message))
}

func (m *Manager) clearRunning(jobID string) {
	m.mu.Lock()
	delete(m.running, jobID)
	m.mu.Unlock()
}

func (m *Manager) persistLocked(jobID string) {
	_ = m.store.SaveJob(m.jobs[jobID])
	if result, ok := m.results[jobID]; ok {
		_ = m.store.SaveResult(result)
	}
	_ = m.store.SaveEvents(jobID, m.events[jobID])
}

func (m *Manager) ListJobs() []domain.Job {
	m.mu.RLock()
	defer m.mu.RUnlock()
	items := make([]domain.Job, 0, len(m.jobs))
	for _, job := range m.jobs {
		items = append(items, job)
	}
	sortJobs(items)
	return items
}

func (m *Manager) GetJob(jobID string) (domain.Job, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	job, ok := m.jobs[jobID]
	return job, ok
}

func (m *Manager) GetJobResult(jobID string) (domain.JobResult, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	result, ok := m.results[jobID]
	return result, ok
}

func (m *Manager) GetJobEvents(jobID string) []domain.JobEvent {
	m.mu.RLock()
	defer m.mu.RUnlock()
	items := append([]domain.JobEvent(nil), m.events[jobID]...)
	return items
}

func (m *Manager) GetEndpoint(id string) (domain.EndpointConfig, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	endpoint, ok := m.endpoints[id]
	return endpoint, ok
}

func resolveStorage(storage *domain.StorageSpec, endpoint domain.EndpointConfig) (string, string) {
	if endpoint.ID != "" {
		storageType := string(endpoint.Type)
		raw := endpoint.RawConfig()
		if storage != nil && storage.Config != "" {
			if raw != "" {
				raw += ";"
			}
			raw += storage.Config
		}
		return storageType, raw
	}
	if storage == nil {
		return "", ""
	}
	return storage.Type, storage.Config
}

func extractStageMetrics(stages []domain.StageState) []domain.MetricsSummary {
	parts := make([]domain.MetricsSummary, 0, len(stages))
	for _, stage := range stages {
		parts = append(parts, stage.Metrics)
	}
	return parts
}

func newEvent(jobID string, level domain.EventLevel, message string) domain.JobEvent {
	return domain.JobEvent{JobID: jobID, OccurredAt: time.Now().UTC(), Level: level, Message: message}
}
