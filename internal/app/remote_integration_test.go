package app

import (
	"context"
	"path/filepath"
	"testing"
	"time"

	"github.com/sine-io/cosbench-go/internal/domain"
)

func TestNewAppSupportsExplicitModes(t *testing.T) {
	viewDir, err := filepath.Abs(filepath.Join("..", "..", "web", "templates"))
	if err != nil {
		t.Fatal(err)
	}
	for _, mode := range []Mode{ModeControllerOnly, ModeDriverOnly, ModeCombined} {
		application, err := New(Config{DataDir: t.TempDir(), ViewDir: viewDir, Mode: mode, DriverSharedToken: "shared-token"})
		if err != nil {
			t.Fatalf("New(%s): %v", mode, err)
		}
		if application.Mode != mode {
			t.Fatalf("application mode = %s want %s", application.Mode, mode)
		}
		if application.Handler == nil || application.Manager == nil {
			t.Fatalf("unexpected app for %s: %#v", mode, application)
		}
	}
}

func TestCombinedModeProcessesMissionViaLoopback(t *testing.T) {
	viewDir, err := filepath.Abs(filepath.Join("..", "..", "web", "templates"))
	if err != nil {
		t.Fatal(err)
	}
	application, err := New(Config{DataDir: t.TempDir(), ViewDir: viewDir, Mode: ModeCombined, DriverSharedToken: "shared-token"})
	if err != nil {
		t.Fatalf("New(): %v", err)
	}

	job, err := application.Manager.CreateJobFromXML([]byte(`<?xml version="1.0" encoding="UTF-8"?>
<workload name="combined-loopback">
  <storage type="mock" />
  <workflow>
    <workstage name="main">
      <work name="main" workers="1" totalOps="2">
        <operation type="write" ratio="100" config="cprefix=t;containers=c(1);objects=s(1,2);sizes=c(1)KB" />
      </work>
    </workstage>
  </workflow>
</workload>`), "")
	if err != nil {
		t.Fatal(err)
	}
	if _, err := application.Manager.ScheduleJobStage(job.ID); err != nil {
		t.Fatalf("ScheduleJobStage(): %v", err)
	}

	processed, err := application.ProcessCombinedMission(context.Background())
	if err != nil {
		t.Fatalf("ProcessCombinedMission(): %v", err)
	}
	if !processed {
		t.Fatal("expected combined loopback to process a mission")
	}

	missions := application.Manager.ListMissions()
	if len(missions) != 1 || missions[0].Status != domain.MissionStatusSucceeded {
		t.Fatalf("missions = %#v", missions)
	}

	deadline := time.Now().Add(2 * time.Second)
	for time.Now().Before(deadline) {
		result, ok := application.Manager.GetJobResult(job.ID)
		if ok && result.Metrics.OperationCount > 0 {
			return
		}
		time.Sleep(20 * time.Millisecond)
	}
	result, _ := application.Manager.GetJobResult(job.ID)
	t.Fatalf("expected combined loopback result, got %#v", result)
}
