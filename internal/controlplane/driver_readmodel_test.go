package controlplane

import (
	"testing"
	"time"

	"github.com/sine-io/cosbench-go/internal/domain"
	"github.com/sine-io/cosbench-go/internal/snapshot"
)

func TestDriverReadModelsExposeOverviewMissionsWorkersAndLogs(t *testing.T) {
	store, err := snapshot.New(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	mgr, err := New(store)
	if err != nil {
		t.Fatal(err)
	}

	driver, err := mgr.RegisterDriverNode(domain.DriverNode{Name: "driver-a", Mode: domain.DriverModeDriver})
	if err != nil {
		t.Fatal(err)
	}
	job, err := mgr.CreateJobFromXML([]byte(`<?xml version="1.0" encoding="UTF-8"?>
<workload name="driver-readmodel">
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
	mission, ok, err := mgr.ClaimMission(driver.ID, 30*time.Second)
	if err != nil || !ok {
		t.Fatalf("ClaimMission(): mission=%#v ok=%v err=%v", mission, ok, err)
	}
	if err := mgr.AppendMissionEvents(mission.ID, []domain.JobEvent{{OccurredAt: time.Now().UTC(), Level: domain.EventLevelInfo, Message: "mission started"}}); err != nil {
		t.Fatal(err)
	}

	overview, ok := mgr.GetDriverOverview(driver.ID)
	if !ok {
		t.Fatal("expected driver overview")
	}
	if overview.Driver.ID != driver.ID || overview.ActiveMissionCount != 1 || overview.LogCount == 0 {
		t.Fatalf("overview = %#v", overview)
	}

	missions := mgr.ListDriverMissions(driver.ID)
	if len(missions) != 1 || missions[0].ID != mission.ID {
		t.Fatalf("missions = %#v", missions)
	}

	workerState, ok := mgr.GetDriverWorkerState(driver.ID)
	if !ok || workerState.ActiveMissionCount != 1 {
		t.Fatalf("workerState = %#v", workerState)
	}

	logs := mgr.GetDriverLogs(driver.ID)
	if len(logs) == 0 {
		t.Fatal("expected driver logs")
	}
}
