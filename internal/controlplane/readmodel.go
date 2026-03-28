package controlplane

import (
	"sort"
	"time"

	"github.com/sine-io/cosbench-go/internal/domain"
)

func (m *Manager) ListJobMatrix() []domain.JobMatrixRow {
	m.mu.RLock()
	defer m.mu.RUnlock()

	rows := make([]domain.JobMatrixRow, 0, len(m.jobs))
	for _, job := range m.jobs {
		updated := job.CreatedAt
		if job.FinishedAt != nil {
			updated = *job.FinishedAt
		} else if job.StartedAt != nil {
			updated = *job.StartedAt
		}
		rows = append(rows, domain.JobMatrixRow{
			JobID:          job.ID,
			Name:           job.Name,
			Status:         job.Status,
			StageCount:     len(job.Stages),
			OperationCount: job.Metrics.OperationCount,
			ByteCount:      job.Metrics.ByteCount,
			ErrorCount:     job.Metrics.ErrorCount,
			UpdatedAt:      updated,
		})
	}
	sort.Slice(rows, func(i, j int) bool {
		if rows[i].UpdatedAt.Equal(rows[j].UpdatedAt) {
			return rows[i].JobID > rows[j].JobID
		}
		return rows[i].UpdatedAt.After(rows[j].UpdatedAt)
	})
	return rows
}

func (m *Manager) GetJobTimeline(jobID string) (domain.JobTimeline, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	timeline, ok := m.timelines[jobID]
	return timeline, ok
}

func currentTimelineUpdate(now time.Time, points []domain.TimelinePoint) time.Time {
	if len(points) == 0 {
		return now
	}
	last := points[len(points)-1].Timestamp
	if last.After(now) {
		return last
	}
	return now
}

func updateTimelineUpdatedAt(now time.Time, stages map[string][]domain.TimelinePoint, job []domain.TimelinePoint) time.Time {
	updated := currentTimelineUpdate(now, job)
	for _, points := range stages {
		stageUpdated := currentTimelineUpdate(now, points)
		if stageUpdated.After(updated) {
			updated = stageUpdated
		}
	}
	return updated
}

func cloneStageTimelines(src map[string][]domain.TimelinePoint) map[string][]domain.TimelinePoint {
	if len(src) == 0 {
		return nil
	}
	dst := make(map[string][]domain.TimelinePoint, len(src))
	for name, points := range src {
		dst[name] = append([]domain.TimelinePoint(nil), points...)
	}
	return dst
}

func nowUTC() time.Time {
	return time.Now().UTC()
}
