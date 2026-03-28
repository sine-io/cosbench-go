package controlplane

import (
	"testing"
	"time"

	"github.com/sine-io/cosbench-go/internal/domain"
	legacyexec "github.com/sine-io/cosbench-go/internal/domain/execution"
	"github.com/sine-io/cosbench-go/internal/snapshot"
)

func TestManagerRegistersDriversSchedulesMissionsAndClaimsWork(t *testing.T) {
	store, err := snapshot.New(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	mgr, err := New(store)
	if err != nil {
		t.Fatal(err)
	}

	driverA, err := mgr.RegisterDriverNode(domain.DriverNode{Name: "driver-a", Mode: domain.DriverModeDriver})
	if err != nil {
		t.Fatalf("RegisterDriverNode(driver-a): %v", err)
	}
	driverB, err := mgr.RegisterDriverNode(domain.DriverNode{Name: "driver-b", Mode: domain.DriverModeDriver})
	if err != nil {
		t.Fatalf("RegisterDriverNode(driver-b): %v", err)
	}
	if driverA.Status != domain.DriverStatusHealthy || driverB.Status != domain.DriverStatusHealthy {
		t.Fatalf("drivers not healthy: %#v %#v", driverA, driverB)
	}

	now := time.Now().UTC()
	if err := mgr.RecordDriverHeartbeat(driverA.ID, now); err != nil {
		t.Fatalf("RecordDriverHeartbeat(): %v", err)
	}
	loadedA, ok := mgr.GetDriverNode(driverA.ID)
	if !ok || loadedA.LastHeartbeatAt == nil || !loadedA.LastHeartbeatAt.Equal(now) {
		t.Fatalf("driver heartbeat = %#v", loadedA)
	}

	job, err := mgr.CreateJobFromXML([]byte(`<?xml version="1.0" encoding="UTF-8"?>
<workload name="remote-schedule">
  <storage type="mock" />
  <workflow>
    <workstage name="main">
      <work name="work-a" workers="1" totalOps="1">
        <operation type="write" ratio="100" config="cprefix=t;containers=c(1);objects=c(1);sizes=c(1)KB" />
      </work>
      <work name="work-b" workers="1" totalOps="1">
        <operation type="write" ratio="100" config="cprefix=t;containers=c(1);objects=c(2);sizes=c(1)KB" />
      </work>
    </workstage>
  </workflow>
</workload>`), "")
	if err != nil {
		t.Fatal(err)
	}

	missions, err := mgr.ScheduleJobStage(job.ID)
	if err != nil {
		t.Fatalf("ScheduleJobStage(): %v", err)
	}
	if len(missions) != 2 {
		t.Fatalf("missions = %#v", missions)
	}

	claimA, ok, err := mgr.ClaimMission(driverA.ID, 30*time.Second)
	if err != nil {
		t.Fatalf("ClaimMission(driver-a): %v", err)
	}
	if !ok || claimA.Lease == nil || claimA.Lease.DriverID != driverA.ID || claimA.Status != domain.MissionStatusClaimed {
		t.Fatalf("claimA = %#v", claimA)
	}

	claimB, ok, err := mgr.ClaimMission(driverB.ID, 30*time.Second)
	if err != nil {
		t.Fatalf("ClaimMission(driver-b): %v", err)
	}
	if !ok || claimB.Lease == nil || claimB.Lease.DriverID != driverB.ID || claimB.ID == claimA.ID {
		t.Fatalf("claimB = %#v", claimB)
	}
}

func TestClaimMissionReclaimsExpiredLeaseButNotActiveLease(t *testing.T) {
	store, err := snapshot.New(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	mgr, err := New(store)
	if err != nil {
		t.Fatal(err)
	}

	driverA, err := mgr.RegisterDriverNode(domain.DriverNode{Name: "driver-a", Mode: domain.DriverModeDriver})
	if err != nil {
		t.Fatal(err)
	}
	driverB, err := mgr.RegisterDriverNode(domain.DriverNode{Name: "driver-b", Mode: domain.DriverModeDriver})
	if err != nil {
		t.Fatal(err)
	}
	job, err := mgr.CreateJobFromXML([]byte(`<?xml version="1.0" encoding="UTF-8"?>
<workload name="remote-expiry">
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
	if len(missions) != 1 {
		t.Fatalf("missions = %#v", missions)
	}

	claimed, ok, err := mgr.ClaimMission(driverA.ID, 30*time.Second)
	if err != nil || !ok {
		t.Fatalf("ClaimMission(driver-a): mission=%#v ok=%v err=%v", claimed, ok, err)
	}
	unexpired, ok, err := mgr.ClaimMission(driverB.ID, 30*time.Second)
	if err != nil {
		t.Fatalf("ClaimMission(driver-b/unexpired): %v", err)
	}
	if ok {
		t.Fatalf("expected no claim while lease active, got %#v", unexpired)
	}

	claimed.Status = domain.MissionStatusRunning
	claimed.Lease.ExpiresAt = time.Now().UTC().Add(-time.Second)
	claimed.UpdatedAt = time.Now().UTC().Add(-time.Second)
	if err := mgr.PutMission(claimed); err != nil {
		t.Fatalf("PutMission(): %v", err)
	}

	reclaimed, ok, err := mgr.ClaimMission(driverB.ID, 30*time.Second)
	if err != nil {
		t.Fatalf("ClaimMission(driver-b/reclaimed): %v", err)
	}
	if !ok {
		t.Fatal("expected reclaimed mission to be available")
	}
	if reclaimed.ID != claimed.ID || reclaimed.Lease == nil || reclaimed.Lease.DriverID != driverB.ID {
		t.Fatalf("reclaimed = %#v", reclaimed)
	}
}

func TestMissionReportingIsIdempotentPerBatch(t *testing.T) {
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
<workload name="remote-idempotent">
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

	events := []domain.JobEvent{{OccurredAt: time.Now().UTC(), Level: domain.EventLevelInfo, Message: "dedupe-event"}}
	if err := mgr.AppendMissionEventsBatch(mission.ID, "events-1", events); err != nil {
		t.Fatal(err)
	}
	if err := mgr.AppendMissionEventsBatch(mission.ID, "events-1", events); err != nil {
		t.Fatal(err)
	}

	samples := []legacyexec.Sample{{Timestamp: time.Now().UTC(), OpType: "write", OpCount: 1, ByteCount: 1000, TotalTimeMs: 10}}
	if err := mgr.AppendMissionSamplesBatch(mission.ID, "samples-1", samples); err != nil {
		t.Fatal(err)
	}
	if err := mgr.AppendMissionSamplesBatch(mission.ID, "samples-1", samples); err != nil {
		t.Fatal(err)
	}

	if err := mgr.CompleteMission(mission.ID, domain.MissionStatusSucceeded, ""); err != nil {
		t.Fatal(err)
	}
	if err := mgr.CompleteMission(mission.ID, domain.MissionStatusSucceeded, ""); err != nil {
		t.Fatal(err)
	}

	logs := mgr.GetJobEvents(job.ID)
	count := 0
	for _, event := range logs {
		if event.Message == "dedupe-event" {
			count++
		}
	}
	if count != 1 {
		t.Fatalf("dedupe-event count = %d logs=%#v", count, logs)
	}
	result, ok := mgr.GetJobResult(job.ID)
	if !ok {
		t.Fatal("expected job result")
	}
	if result.Metrics.OperationCount != 1 || result.Metrics.ByteCount != 1000 {
		t.Fatalf("result = %#v", result)
	}
}
