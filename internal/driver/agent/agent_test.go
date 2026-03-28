package agent

import (
	"context"
	"net/http/httptest"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/sine-io/cosbench-go/internal/controlplane"
	"github.com/sine-io/cosbench-go/internal/domain"
	"github.com/sine-io/cosbench-go/internal/snapshot"
	"github.com/sine-io/cosbench-go/internal/web"
)

func TestAgentProcessesClaimedMissionAndReportsBack(t *testing.T) {
	store, err := snapshot.New(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	mgr, err := controlplane.New(store)
	if err != nil {
		t.Fatal(err)
	}

	job, err := mgr.CreateJobFromXML([]byte(`<?xml version="1.0" encoding="UTF-8"?>
<workload name="agent-run">
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
	if _, err := mgr.ScheduleJobStage(job.ID); err != nil {
		t.Fatalf("ScheduleJobStage(): %v", err)
	}

	root, err := filepath.Abs(filepath.Join("..", "..", ".."))
	if err != nil {
		t.Fatal(err)
	}
	handler, err := web.NewHandler(mgr, filepath.Join(root, "web", "templates"), "shared-token")
	if err != nil {
		t.Fatal(err)
	}
	server := httptest.NewServer(handler)
	defer server.Close()

	agent := &Agent{
		Client: &HTTPClient{BaseURL: server.URL, SharedToken: "shared-token"},
		Name:   "driver-a",
		Mode:   domain.DriverModeDriver,
	}
	processed, err := agent.ProcessOne(context.Background())
	if err != nil {
		t.Fatalf("ProcessOne(): %v", err)
	}
	if !processed {
		t.Fatal("expected a claimed mission to be processed")
	}
	if agent.DriverID == "" {
		t.Fatal("expected registered driver id")
	}

	missions := mgr.ListMissions()
	if len(missions) != 1 {
		t.Fatalf("missions = %#v", missions)
	}
	mission := missions[0]
	if mission.Status != domain.MissionStatusSucceeded {
		t.Fatalf("mission = %#v", mission)
	}

	result, ok := mgr.GetJobResult(job.ID)
	if !ok {
		t.Fatal("expected job result")
	}
	if result.Metrics.OperationCount == 0 || len(result.StageTotals) == 0 {
		t.Fatalf("result = %#v", result)
	}

	timeline, ok := mgr.GetJobTimeline(job.ID)
	if !ok || len(timeline.Job) == 0 {
		t.Fatalf("timeline = %#v", timeline)
	}

	events := mgr.GetJobEvents(job.ID)
	foundStart := false
	foundFinish := false
	for _, event := range events {
		msg := strings.ToLower(event.Message)
		if strings.Contains(msg, "mission") && strings.Contains(msg, "started") {
			foundStart = true
		}
		if strings.Contains(msg, "mission") && strings.Contains(msg, "finished") {
			foundFinish = true
		}
	}
	if !foundStart || !foundFinish {
		t.Fatalf("events = %#v", events)
	}
}

func TestAgentProcessOneReturnsFalseWhenNoMissionExists(t *testing.T) {
	store, err := snapshot.New(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	mgr, err := controlplane.New(store)
	if err != nil {
		t.Fatal(err)
	}
	root, err := filepath.Abs(filepath.Join("..", "..", ".."))
	if err != nil {
		t.Fatal(err)
	}
	handler, err := web.NewHandler(mgr, filepath.Join(root, "web", "templates"), "shared-token")
	if err != nil {
		t.Fatal(err)
	}
	server := httptest.NewServer(handler)
	defer server.Close()

	agent := &Agent{
		Client: &HTTPClient{BaseURL: server.URL, SharedToken: "shared-token"},
		Name:   "driver-idle",
		Mode:   domain.DriverModeDriver,
	}
	processed, err := agent.ProcessOne(context.Background())
	if err != nil {
		t.Fatalf("ProcessOne(): %v", err)
	}
	if processed {
		t.Fatal("expected no mission to be processed")
	}
}

func TestAgentProcessOneCanReclaimExpiredMissionLease(t *testing.T) {
	store, err := snapshot.New(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	mgr, err := controlplane.New(store)
	if err != nil {
		t.Fatal(err)
	}
	job, err := mgr.CreateJobFromXML([]byte(`<?xml version="1.0" encoding="UTF-8"?>
<workload name="agent-expiry">
  <storage type="mock" />
  <workflow>
    <workstage name="main">
      <work name="main" workers="1" totalOps="1">
        <operation type="write" ratio="100" config="cprefix=t;containers=c(1);objects=c(1);sizes=c(1)KB" />
      </work>
    </workstage>
  </workflow>
</workload>`), "")
	if err != nil {
		t.Fatal(err)
	}
	missions, err := mgr.ScheduleJobStage(job.ID)
	if err != nil {
		t.Fatal(err)
	}
	driver, err := mgr.RegisterDriverNode(domain.DriverNode{Name: "expired-owner", Mode: domain.DriverModeDriver})
	if err != nil {
		t.Fatal(err)
	}
	claimed, ok, err := mgr.ClaimMission(driver.ID, 30*time.Second)
	if err != nil || !ok {
		t.Fatalf("ClaimMission(): mission=%#v ok=%v err=%v", claimed, ok, err)
	}
	claimed.Status = domain.MissionStatusRunning
	claimed.Lease.ExpiresAt = time.Now().UTC().Add(-time.Second)
	claimed.UpdatedAt = time.Now().UTC().Add(-time.Second)
	if err := mgr.PutMission(claimed); err != nil {
		t.Fatal(err)
	}
	if len(missions) != 1 {
		t.Fatalf("missions = %#v", missions)
	}

	root, err := filepath.Abs(filepath.Join("..", "..", ".."))
	if err != nil {
		t.Fatal(err)
	}
	handler, err := web.NewHandler(mgr, filepath.Join(root, "web", "templates"), "shared-token")
	if err != nil {
		t.Fatal(err)
	}
	server := httptest.NewServer(handler)
	defer server.Close()

	agent := &Agent{
		Client: &HTTPClient{BaseURL: server.URL, SharedToken: "shared-token"},
		Name:   "driver-reclaimer",
		Mode:   domain.DriverModeDriver,
	}
	processed, err := agent.ProcessOne(context.Background())
	if err != nil {
		t.Fatalf("ProcessOne(): %v", err)
	}
	if !processed {
		t.Fatal("expected expired mission to be reclaimed and processed")
	}
	reloaded, ok := mgr.GetMission(claimed.ID)
	if !ok || reloaded.Status != domain.MissionStatusSucceeded {
		t.Fatalf("reloaded mission = %#v", reloaded)
	}
}

func TestAgentFailsProtectedWritesWithoutSharedToken(t *testing.T) {
	store, err := snapshot.New(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	mgr, err := controlplane.New(store)
	if err != nil {
		t.Fatal(err)
	}
	job, err := mgr.CreateJobFromXML([]byte(`<?xml version="1.0" encoding="UTF-8"?>
<workload name="agent-auth">
  <storage type="mock" />
  <workflow>
    <workstage name="main">
      <work name="main" workers="1" totalOps="1">
        <operation type="write" ratio="100" config="cprefix=t;containers=c(1);objects=c(1);sizes=c(1)KB" />
      </work>
    </workstage>
  </workflow>
</workload>`), "")
	if err != nil {
		t.Fatal(err)
	}
	if _, err := mgr.ScheduleJobStage(job.ID); err != nil {
		t.Fatal(err)
	}
	root, err := filepath.Abs(filepath.Join("..", "..", ".."))
	if err != nil {
		t.Fatal(err)
	}
	handler, err := web.NewHandler(mgr, filepath.Join(root, "web", "templates"), "shared-token")
	if err != nil {
		t.Fatal(err)
	}
	server := httptest.NewServer(handler)
	defer server.Close()

	agent := &Agent{
		Client: &HTTPClient{BaseURL: server.URL},
		Name:   "driver-unauthorized",
		Mode:   domain.DriverModeDriver,
	}
	if _, err := agent.ProcessOne(context.Background()); err == nil {
		t.Fatal("expected protected write failure without shared token")
	}
}

func TestAgentProcessesOnlyClaimedWorkerSlice(t *testing.T) {
	store, err := snapshot.New(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	mgr, err := controlplane.New(store)
	if err != nil {
		t.Fatal(err)
	}

	job, err := mgr.CreateJobFromXML([]byte(`<?xml version="1.0" encoding="UTF-8"?>
<workload name="agent-unit-slice">
  <storage type="mock" />
  <workflow>
    <workstage name="main">
      <work name="main" workers="3" totalOps="3">
        <operation type="write" ratio="100" config="cprefix=t;containers=c(1);objects=s(1,3);sizes=c(1)KB" />
      </work>
    </workstage>
  </workflow>
</workload>`), "")
	if err != nil {
		t.Fatal(err)
	}
	if _, err := mgr.ScheduleJobStage(job.ID); err != nil {
		t.Fatalf("ScheduleJobStage(): %v", err)
	}

	root, err := filepath.Abs(filepath.Join("..", "..", ".."))
	if err != nil {
		t.Fatal(err)
	}
	handler, err := web.NewHandler(mgr, filepath.Join(root, "web", "templates"), "shared-token")
	if err != nil {
		t.Fatal(err)
	}
	server := httptest.NewServer(handler)
	defer server.Close()

	agent := &Agent{
		Client: &HTTPClient{BaseURL: server.URL, SharedToken: "shared-token"},
		Name:   "driver-slice",
		Mode:   domain.DriverModeDriver,
	}
	processed, err := agent.ProcessOne(context.Background())
	if err != nil {
		t.Fatalf("ProcessOne(): %v", err)
	}
	if !processed {
		t.Fatal("expected a unit to be processed")
	}

	result, ok := mgr.GetJobResult(job.ID)
	if !ok {
		t.Fatal("expected job result")
	}
	if result.Metrics.OperationCount != 1 {
		t.Fatalf("expected one-op unit execution, got %#v", result)
	}
}
