package app

import (
	"context"
	"net/http/httptest"
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

func TestDriverOnlyModeProcessesControllerMissions(t *testing.T) {
	viewDir, err := filepath.Abs(filepath.Join("..", "..", "web", "templates"))
	if err != nil {
		t.Fatal(err)
	}

	controllerApp, err := New(Config{DataDir: t.TempDir(), ViewDir: viewDir, Mode: ModeControllerOnly, DriverSharedToken: "shared-token"})
	if err != nil {
		t.Fatalf("New(controller): %v", err)
	}
	server := httptest.NewServer(controllerApp.Handler)
	defer server.Close()

	driverApp, err := New(Config{
		DataDir:           t.TempDir(),
		ViewDir:           viewDir,
		Mode:              ModeDriverOnly,
		DriverSharedToken: "shared-token",
		ControllerURL:     server.URL,
		DriverName:        "driver-only-a",
	})
	if err != nil {
		t.Fatalf("New(driver): %v", err)
	}
	bgCtx, cancel := context.WithCancel(context.Background())
	defer cancel()
	if err := driverApp.StartBackground(bgCtx); err != nil {
		t.Fatalf("StartBackground(): %v", err)
	}

	job, err := controllerApp.Manager.CreateJobFromXML([]byte(`<?xml version="1.0" encoding="UTF-8"?>
<workload name="driver-only-remote">
  <storage type="mock" />
  <workflow>
    <workstage name="main">
      <work name="main" workers="1" totalOps="1">
        <operation type="write" ratio="100" config="cprefix=t;containers=c(1);objects=s(1,1);sizes=c(1)KB" />
      </work>
    </workstage>
  </workflow>
</workload>`), "")
	if err != nil {
		t.Fatal(err)
	}
	if err := controllerApp.Manager.StartJob(context.Background(), job.ID); err != nil {
		t.Fatalf("StartJob(): %v", err)
	}

	deadline := time.Now().Add(3 * time.Second)
	for time.Now().Before(deadline) {
		result, ok := controllerApp.Manager.GetJobResult(job.ID)
		if ok && result.Metrics.OperationCount > 0 {
			if len(controllerApp.Manager.ListMissionAttempts()) == 0 {
				t.Fatalf("expected remote mission attempts, got local-only result %#v", result)
			}
			if len(controllerApp.Manager.ListWorkUnits(job.ID, "main", "main")) == 0 {
				t.Fatalf("expected work units for remote execution, got %#v", result)
			}
			return
		}
		time.Sleep(25 * time.Millisecond)
	}
	result, _ := controllerApp.Manager.GetJobResult(job.ID)
	t.Fatalf("expected driver-only mode to process remote mission, got %#v", result)
}

func TestControllerBackgroundSweepMarksStaleDriverUnhealthy(t *testing.T) {
	viewDir, err := filepath.Abs(filepath.Join("..", "..", "web", "templates"))
	if err != nil {
		t.Fatal(err)
	}
	controllerApp, err := New(Config{
		DataDir:            t.TempDir(),
		ViewDir:            viewDir,
		Mode:               ModeControllerOnly,
		DriverSharedToken:  "shared-token",
		LeaseSweepInterval: 10 * time.Millisecond,
	})
	if err != nil {
		t.Fatalf("New(controller): %v", err)
	}

	driver, err := controllerApp.Manager.RegisterDriverNode(domain.DriverNode{Name: "driver-stale", Mode: domain.DriverModeDriver})
	if err != nil {
		t.Fatalf("RegisterDriverNode(): %v", err)
	}
	staleAt := time.Now().UTC().Add(-2 * time.Minute)
	driver.LastHeartbeatAt = &staleAt
	driver.Status = domain.DriverStatusHealthy
	if err := controllerApp.Manager.PutDriverNode(driver); err != nil {
		t.Fatalf("PutDriverNode(): %v", err)
	}

	bgCtx, cancel := context.WithCancel(context.Background())
	defer cancel()
	if err := controllerApp.StartBackground(bgCtx); err != nil {
		t.Fatalf("StartBackground(): %v", err)
	}

	deadline := time.Now().Add(500 * time.Millisecond)
	for time.Now().Before(deadline) {
		loaded, ok := controllerApp.Manager.GetDriverNode(driver.ID)
		if ok && loaded.Status == domain.DriverStatusUnhealthy {
			return
		}
		time.Sleep(10 * time.Millisecond)
	}
	loaded, _ := controllerApp.Manager.GetDriverNode(driver.ID)
	t.Fatalf("expected stale driver to become unhealthy, got %#v", loaded)
}

func TestControllerBackgroundSweepEmitsDriverUnhealthyJobEvent(t *testing.T) {
	viewDir, err := filepath.Abs(filepath.Join("..", "..", "web", "templates"))
	if err != nil {
		t.Fatal(err)
	}
	controllerApp, err := New(Config{
		DataDir:            t.TempDir(),
		ViewDir:            viewDir,
		Mode:               ModeControllerOnly,
		DriverSharedToken:  "shared-token",
		LeaseSweepInterval: 10 * time.Millisecond,
	})
	if err != nil {
		t.Fatalf("New(controller): %v", err)
	}

	driver, err := controllerApp.Manager.RegisterDriverNode(domain.DriverNode{Name: "driver-stale", Mode: domain.DriverModeDriver})
	if err != nil {
		t.Fatalf("RegisterDriverNode(): %v", err)
	}
	job, err := controllerApp.Manager.CreateJobFromXML([]byte(`<?xml version="1.0" encoding="UTF-8"?>
<workload name="driver-unhealthy-event-app">
  <storage type="mock" />
  <workflow>
    <workstage name="main">
      <work name="main" workers="1" totalOps="1">
        <operation type="write" ratio="100" config="cprefix=t;containers=c(1);objects=s(1,1);sizes=c(1)KB" />
      </work>
    </workstage>
  </workflow>
</workload>`), "")
	if err != nil {
		t.Fatal(err)
	}
	if _, err := controllerApp.Manager.ScheduleJobStage(job.ID); err != nil {
		t.Fatalf("ScheduleJobStage(): %v", err)
	}
	if _, ok, err := controllerApp.Manager.ClaimMission(driver.ID, 30*time.Second); err != nil || !ok {
		t.Fatalf("ClaimMission(): ok=%v err=%v", ok, err)
	}

	staleAt := time.Now().UTC().Add(-2 * time.Minute)
	driver.LastHeartbeatAt = &staleAt
	driver.Status = domain.DriverStatusHealthy
	if err := controllerApp.Manager.PutDriverNode(driver); err != nil {
		t.Fatalf("PutDriverNode(): %v", err)
	}

	bgCtx, cancel := context.WithCancel(context.Background())
	defer cancel()
	if err := controllerApp.StartBackground(bgCtx); err != nil {
		t.Fatalf("StartBackground(): %v", err)
	}

	deadline := time.Now().Add(500 * time.Millisecond)
	for time.Now().Before(deadline) {
		events := controllerApp.Manager.GetJobEvents(job.ID)
		for _, event := range events {
			if event.Message == "driver driver-stale marked unhealthy by heartbeat timeout" {
				return
			}
		}
		time.Sleep(10 * time.Millisecond)
	}
	events := controllerApp.Manager.GetJobEvents(job.ID)
	t.Fatalf("expected driver unhealthy event, events=%#v", events)
}

func TestControllerBackgroundSweepUsesConfiguredDriverHeartbeatTimeoutForJobEvent(t *testing.T) {
	viewDir, err := filepath.Abs(filepath.Join("..", "..", "web", "templates"))
	if err != nil {
		t.Fatal(err)
	}
	controllerApp, err := New(Config{
		DataDir:                t.TempDir(),
		ViewDir:                viewDir,
		Mode:                   ModeControllerOnly,
		DriverSharedToken:      "shared-token",
		LeaseSweepInterval:     10 * time.Millisecond,
		DriverHeartbeatTimeout: 20 * time.Millisecond,
	})
	if err != nil {
		t.Fatalf("New(controller): %v", err)
	}

	driver, err := controllerApp.Manager.RegisterDriverNode(domain.DriverNode{Name: "driver-stale", Mode: domain.DriverModeDriver})
	if err != nil {
		t.Fatalf("RegisterDriverNode(): %v", err)
	}
	job, err := controllerApp.Manager.CreateJobFromXML([]byte(`<?xml version="1.0" encoding="UTF-8"?>
<workload name="driver-unhealthy-event-app-custom-timeout">
  <storage type="mock" />
  <workflow>
    <workstage name="main">
      <work name="main" workers="1" totalOps="1">
        <operation type="write" ratio="100" config="cprefix=t;containers=c(1);objects=s(1,1);sizes=c(1)KB" />
      </work>
    </workstage>
  </workflow>
</workload>`), "")
	if err != nil {
		t.Fatal(err)
	}
	if _, err := controllerApp.Manager.ScheduleJobStage(job.ID); err != nil {
		t.Fatalf("ScheduleJobStage(): %v", err)
	}
	if _, ok, err := controllerApp.Manager.ClaimMission(driver.ID, 30*time.Second); err != nil || !ok {
		t.Fatalf("ClaimMission(): ok=%v err=%v", ok, err)
	}

	staleAt := time.Now().UTC().Add(-40 * time.Millisecond)
	driver.LastHeartbeatAt = &staleAt
	driver.Status = domain.DriverStatusHealthy
	if err := controllerApp.Manager.PutDriverNode(driver); err != nil {
		t.Fatalf("PutDriverNode(): %v", err)
	}

	bgCtx, cancel := context.WithCancel(context.Background())
	defer cancel()
	if err := controllerApp.StartBackground(bgCtx); err != nil {
		t.Fatalf("StartBackground(): %v", err)
	}

	deadline := time.Now().Add(500 * time.Millisecond)
	for time.Now().Before(deadline) {
		events := controllerApp.Manager.GetJobEvents(job.ID)
		for _, event := range events {
			if event.Message == "driver driver-stale marked unhealthy by heartbeat timeout" {
				return
			}
		}
		time.Sleep(10 * time.Millisecond)
	}
	events := controllerApp.Manager.GetJobEvents(job.ID)
	t.Fatalf("expected driver unhealthy event with custom timeout, events=%#v", events)
}

func TestCombinedModeProgressesAcrossMultipleStages(t *testing.T) {
	viewDir, err := filepath.Abs(filepath.Join("..", "..", "web", "templates"))
	if err != nil {
		t.Fatal(err)
	}
	application, err := New(Config{DataDir: t.TempDir(), ViewDir: viewDir, Mode: ModeCombined, DriverSharedToken: "shared-token"})
	if err != nil {
		t.Fatalf("New(): %v", err)
	}

	job, err := application.Manager.CreateJobFromXML([]byte(`<?xml version="1.0" encoding="UTF-8"?>
<workload name="combined-multistage">
  <storage type="mock" />
  <workflow>
    <workstage name="stage-a">
      <work name="main-a" workers="1" totalOps="1">
        <operation type="write" ratio="100" config="cprefix=t;containers=c(1);objects=s(1,1);sizes=c(1)KB" />
      </work>
    </workstage>
    <workstage name="stage-b">
      <work name="main-b" workers="1" totalOps="1">
        <operation type="write" ratio="100" config="cprefix=t;containers=c(1);objects=s(2,2);sizes=c(1)KB" />
      </work>
    </workstage>
  </workflow>
</workload>`), "")
	if err != nil {
		t.Fatal(err)
	}
	if err := application.Manager.StartJob(context.Background(), job.ID); err != nil {
		t.Fatalf("StartJob(): %v", err)
	}

	deadline := time.Now().Add(5 * time.Second)
	for time.Now().Before(deadline) {
		if _, err := application.ProcessCombinedMission(context.Background()); err != nil {
			t.Fatalf("ProcessCombinedMission(): %v", err)
		}
		result, ok := application.Manager.GetJobResult(job.ID)
		if ok && result.Metrics.OperationCount == 2 {
			loaded, _ := application.Manager.GetJob(job.ID)
			if loaded.Status == domain.JobStatusSucceeded && len(loaded.Stages) == 2 && loaded.Stages[0].Status == domain.JobStatusSucceeded && loaded.Stages[1].Status == domain.JobStatusSucceeded {
				return
			}
		}
		time.Sleep(25 * time.Millisecond)
	}
	result, _ := application.Manager.GetJobResult(job.ID)
	t.Fatalf("expected combined-mode multistage result, got %#v", result)
}

func TestDriverOnlyModeProgressesAcrossMultipleStages(t *testing.T) {
	viewDir, err := filepath.Abs(filepath.Join("..", "..", "web", "templates"))
	if err != nil {
		t.Fatal(err)
	}
	controllerApp, err := New(Config{DataDir: t.TempDir(), ViewDir: viewDir, Mode: ModeControllerOnly, DriverSharedToken: "shared-token"})
	if err != nil {
		t.Fatalf("New(controller): %v", err)
	}
	server := httptest.NewServer(controllerApp.Handler)
	defer server.Close()

	driverApp, err := New(Config{
		DataDir:           t.TempDir(),
		ViewDir:           viewDir,
		Mode:              ModeDriverOnly,
		DriverSharedToken: "shared-token",
		ControllerURL:     server.URL,
		DriverName:        "driver-multistage",
	})
	if err != nil {
		t.Fatalf("New(driver): %v", err)
	}
	bgCtx, cancel := context.WithCancel(context.Background())
	defer cancel()
	if err := driverApp.StartBackground(bgCtx); err != nil {
		t.Fatalf("StartBackground(): %v", err)
	}

	job, err := controllerApp.Manager.CreateJobFromXML([]byte(`<?xml version="1.0" encoding="UTF-8"?>
<workload name="driver-only-multistage">
  <storage type="mock" />
  <workflow>
    <workstage name="stage-a">
      <work name="main-a" workers="1" totalOps="1">
        <operation type="write" ratio="100" config="cprefix=t;containers=c(1);objects=s(1,1);sizes=c(1)KB" />
      </work>
    </workstage>
    <workstage name="stage-b">
      <work name="main-b" workers="1" totalOps="1">
        <operation type="write" ratio="100" config="cprefix=t;containers=c(1);objects=s(2,2);sizes=c(1)KB" />
      </work>
    </workstage>
  </workflow>
</workload>`), "")
	if err != nil {
		t.Fatal(err)
	}
	if err := controllerApp.Manager.StartJob(context.Background(), job.ID); err != nil {
		t.Fatalf("StartJob(): %v", err)
	}

	deadline := time.Now().Add(5 * time.Second)
	for time.Now().Before(deadline) {
		result, ok := controllerApp.Manager.GetJobResult(job.ID)
		if ok && result.Metrics.OperationCount == 2 {
			loaded, _ := controllerApp.Manager.GetJob(job.ID)
			if loaded.Status == domain.JobStatusSucceeded && len(loaded.Stages) == 2 && loaded.Stages[0].Status == domain.JobStatusSucceeded && loaded.Stages[1].Status == domain.JobStatusSucceeded {
				return
			}
		}
		time.Sleep(25 * time.Millisecond)
	}
	result, _ := controllerApp.Manager.GetJobResult(job.ID)
	t.Fatalf("expected driver-only multistage result, got %#v", result)
}

func TestControllerBackgroundSweepRequeuesExpiredMissionForDriverOnlyMode(t *testing.T) {
	viewDir, err := filepath.Abs(filepath.Join("..", "..", "web", "templates"))
	if err != nil {
		t.Fatal(err)
	}
	controllerApp, err := New(Config{
		DataDir:            t.TempDir(),
		ViewDir:            viewDir,
		Mode:               ModeControllerOnly,
		DriverSharedToken:  "shared-token",
		LeaseSweepInterval: 10 * time.Millisecond,
	})
	if err != nil {
		t.Fatalf("New(controller): %v", err)
	}
	server := httptest.NewServer(controllerApp.Handler)
	defer server.Close()

	stalledDriver, err := controllerApp.Manager.RegisterDriverNode(domain.DriverNode{Name: "driver-stalled", Mode: domain.DriverModeDriver})
	if err != nil {
		t.Fatalf("RegisterDriverNode(stalled): %v", err)
	}

	job, err := controllerApp.Manager.CreateJobFromXML([]byte(`<?xml version="1.0" encoding="UTF-8"?>
<workload name="driver-only-requeue">
  <storage type="mock" />
  <workflow>
    <workstage name="main">
      <work name="main" workers="1" totalOps="1">
        <operation type="write" ratio="100" config="cprefix=t;containers=c(1);objects=s(1,1);sizes=c(1)KB" />
      </work>
    </workstage>
  </workflow>
</workload>`), "")
	if err != nil {
		t.Fatal(err)
	}
	if err := controllerApp.Manager.StartJob(context.Background(), job.ID); err != nil {
		t.Fatalf("StartJob(): %v", err)
	}

	claimed, ok, err := controllerApp.Manager.ClaimMission(stalledDriver.ID, 50*time.Millisecond)
	if err != nil || !ok {
		t.Fatalf("ClaimMission(stalled): mission=%#v ok=%v err=%v", claimed, ok, err)
	}

	driverApp, err := New(Config{
		DataDir:            t.TempDir(),
		ViewDir:            viewDir,
		Mode:               ModeDriverOnly,
		DriverSharedToken:  "shared-token",
		ControllerURL:      server.URL,
		DriverName:         "driver-live",
		DriverPollInterval: 20 * time.Millisecond,
	})
	if err != nil {
		t.Fatalf("New(driver): %v", err)
	}
	bgCtx, cancel := context.WithCancel(context.Background())
	defer cancel()
	if err := controllerApp.StartBackground(bgCtx); err != nil {
		t.Fatalf("controller StartBackground(): %v", err)
	}
	if err := driverApp.StartBackground(bgCtx); err != nil {
		t.Fatalf("driver StartBackground(): %v", err)
	}

	deadline := time.Now().Add(5 * time.Second)
	for time.Now().Before(deadline) {
		result, ok := controllerApp.Manager.GetJobResult(job.ID)
		if ok && result.Metrics.OperationCount == 1 {
			attempts := controllerApp.Manager.ListMissionAttempts()
			if len(attempts) >= 2 {
				return
			}
		}
		time.Sleep(25 * time.Millisecond)
	}
	result, _ := controllerApp.Manager.GetJobResult(job.ID)
	attempts := controllerApp.Manager.ListMissionAttempts()
	t.Fatalf("expected background sweep to requeue and finish job, result=%#v attempts=%#v", result, attempts)
}
