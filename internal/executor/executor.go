package executor

import (
	"context"

	"github.com/sine-io/cosbench-go/internal/application/ports"
	"github.com/sine-io/cosbench-go/internal/domain"
	legacyexec "github.com/sine-io/cosbench-go/internal/domain/execution"
	"github.com/sine-io/cosbench-go/internal/reporting"
)

type StageExecutor struct {
	Storage ports.StorageAdapter
}

type StageResult struct {
	Summary domain.MetricsSummary
	Works   []WorkResult
	Err     error
}

type WorkResult struct {
	WorkName string
	Summary domain.MetricsSummary
	Err     error
}

func (e StageExecutor) RunWork(ctx context.Context, work domain.Work) WorkResult {
	engine := &legacyexec.Engine{Work: work.ToLegacy(), Storage: e.Storage}
	result := engine.Run(ctx)
	summary := reporting.Summarize(result.Samples)
	if result.Err != nil {
		return WorkResult{WorkName: work.Name, Summary: summary, Err: result.Err}
	}
	return WorkResult{WorkName: work.Name, Summary: summary}
}

func (e StageExecutor) RunStage(ctx context.Context, stage domain.Stage) StageResult {
	results := make([]WorkResult, 0, len(stage.Works))
	parts := make([]domain.MetricsSummary, 0, len(stage.Works))
	for _, work := range stage.Works {
		result := e.RunWork(ctx, work)
		results = append(results, result)
		if result.Err != nil {
			return StageResult{Works: results, Summary: reporting.Merge(parts...), Err: result.Err}
		}
		parts = append(parts, result.Summary)
	}
	return StageResult{Works: results, Summary: reporting.Merge(parts...)}
}

func ValidateWork(work domain.Work) error {
	legacy := work.ToLegacy()
	storageRaw := ""
	if legacy.Storage != nil {
		storageRaw = legacy.Storage.Config
	}
	for _, op := range legacy.Operations {
		if err := legacyexec.ValidateOperation(op, storageRaw); err != nil {
			return err
		}
	}
	return nil
}
