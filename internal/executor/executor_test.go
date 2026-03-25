package executor

import (
	"context"
	"testing"

	"github.com/sine-io/cosbench-go/internal/domain"
	"github.com/sine-io/cosbench-go/internal/infrastructure/storage/mock"
)

func TestStageExecutorRunWork(t *testing.T) {
	storage := mock.New()
	if err := storage.Init(nil); err != nil {
		t.Fatalf("Init(): %v", err)
	}
	e := StageExecutor{Storage: storage}
	result := e.RunWork(context.Background(), domain.Work{
		Name:     "write-main",
		Type:     "normal",
		Workers:  1,
		TotalOps: 3,
		Storage:  &domain.StorageSpec{Type: "mock"},
		Operations: []domain.Operation{{
			Type:   "write",
			Ratio:  100,
			Config: "cprefix=t;containers=c(1);objects=s(1,3);sizes=c(1)KB",
		}},
	})
	if result.Err != nil {
		t.Fatalf("RunWork(): %v", result.Err)
	}
	if result.Summary.OperationCount != 3 {
		t.Fatalf("operation count = %d", result.Summary.OperationCount)
	}
	if result.Summary.P95LatencyMs < 0 {
		t.Fatalf("unexpected percentile summary: %#v", result.Summary)
	}
}

func TestStageExecutorRunStage(t *testing.T) {
	storage := mock.New()
	if err := storage.Init(nil); err != nil {
		t.Fatalf("Init(): %v", err)
	}
	e := StageExecutor{Storage: storage}
	result := e.RunStage(context.Background(), domain.Stage{
		Name: "main",
		Works: []domain.Work{
			{
				Name:     "init",
				Type:     "init",
				Workers:  1,
				TotalOps: 1,
				Storage:  &domain.StorageSpec{Type: "mock"},
				Operations: []domain.Operation{{Type: "init", Ratio: 100, Config: "cprefix=t;containers=c(1)"}},
			},
			{
				Name:     "write",
				Type:     "normal",
				Workers:  1,
				TotalOps: 2,
				Storage:  &domain.StorageSpec{Type: "mock"},
				Operations: []domain.Operation{{Type: "write", Ratio: 100, Config: "cprefix=t;containers=c(1);objects=s(1,2);sizes=c(1)KB"}},
			},
		},
	})
	if result.Err != nil {
		t.Fatalf("RunStage(): %v", result.Err)
	}
	if len(result.Works) != 2 || result.Summary.OperationCount != 3 {
		t.Fatalf("unexpected stage result: %#v", result)
	}
}
