package controlplane

import (
	"context"
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
	if reclaimed.ID == claimed.ID || reclaimed.WorkUnitID != claimed.WorkUnitID || reclaimed.Attempt != claimed.Attempt+1 || reclaimed.Lease == nil || reclaimed.Lease.DriverID != driverB.ID {
		t.Fatalf("reclaimed = %#v", reclaimed)
	}
}

func TestSweepExpiredLeasesRequeuesWithoutNewClaim(t *testing.T) {
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
	job, err := mgr.CreateJobFromXML([]byte(`<?xml version="1.0" encoding="UTF-8"?>
<workload name="remote-sweep-requeue">
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

	claimed, ok, err := mgr.ClaimMission(driverA.ID, 30*time.Second)
	if err != nil || !ok {
		t.Fatalf("ClaimMission(driver-a): mission=%#v ok=%v err=%v", claimed, ok, err)
	}
	expiredAt := time.Now().UTC().Add(-time.Second)
	claimed.Status = domain.MissionStatusRunning
	claimed.Lease.ExpiresAt = expiredAt
	claimed.UpdatedAt = expiredAt
	if err := mgr.PutMission(claimed); err != nil {
		t.Fatalf("PutMission(): %v", err)
	}

	mgr.SweepExpiredLeases(time.Now().UTC())

	attempts := mgr.ListMissionAttempts()
	if len(attempts) != 2 {
		t.Fatalf("attempts = %#v", attempts)
	}
	foundExpired := false
	foundRetry := false
	for _, attempt := range attempts {
		switch {
		case attempt.ID == claimed.ID:
			foundExpired = attempt.Status == domain.MissionStatusExpired
		case attempt.WorkUnitID == claimed.WorkUnitID && attempt.Attempt == claimed.Attempt+1:
			foundRetry = attempt.Status == domain.MissionAttemptStatusPending && attempt.Lease == nil
		}
	}
	if !foundExpired || !foundRetry {
		t.Fatalf("attempts after sweep = %#v", attempts)
	}
	units := mgr.ListWorkUnits(job.ID, "main", "main")
	if len(units) != 1 || units[0].Status != domain.WorkUnitStatusPending {
		t.Fatalf("units after sweep = %#v", units)
	}
	events := mgr.GetJobEvents(job.ID)
	count := 0
	for _, event := range events {
		if event.Message == "mission lease expired" {
			count++
		}
	}
	if count != 1 {
		t.Fatalf("mission lease expired events = %d events=%#v", count, events)
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

func TestClaimMissionRejectsStaleUnhealthyDriver(t *testing.T) {
	store, err := snapshot.New(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	mgr, err := New(store)
	if err != nil {
		t.Fatal(err)
	}
	staleDriver, err := mgr.RegisterDriverNode(domain.DriverNode{Name: "driver-stale", Mode: domain.DriverModeDriver})
	if err != nil {
		t.Fatal(err)
	}
	freshDriver, err := mgr.RegisterDriverNode(domain.DriverNode{Name: "driver-fresh", Mode: domain.DriverModeDriver})
	if err != nil {
		t.Fatal(err)
	}
	staleAt := time.Now().UTC().Add(-2 * driverHeartbeatTimeout)
	staleDriver.LastHeartbeatAt = &staleAt
	staleDriver.Status = domain.DriverStatusHealthy
	if err := mgr.PutDriverNode(staleDriver); err != nil {
		t.Fatal(err)
	}

	job, err := mgr.CreateJobFromXML([]byte(`<?xml version="1.0" encoding="UTF-8"?>
<workload name="remote-stale-driver">
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

	if _, ok, err := mgr.ClaimMission(staleDriver.ID, 30*time.Second); err == nil {
		t.Fatalf("expected stale driver claim error, ok=%v", ok)
	}

	claimed, ok, err := mgr.ClaimMission(freshDriver.ID, 30*time.Second)
	if err != nil || !ok {
		t.Fatalf("fresh driver claim: mission=%#v ok=%v err=%v", claimed, ok, err)
	}
}

func TestClaimMissionDistributesDifferentUnitsAcrossDrivers(t *testing.T) {
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
<workload name="multi-driver-units">
  <storage type="mock" />
  <workflow>
    <workstage name="main">
      <work name="fanout" workers="2" totalOps="2">
        <operation type="write" ratio="100" config="cprefix=t;containers=c(1);objects=s(1,2);sizes=c(1)KB" />
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

	claimA, ok, err := mgr.ClaimMission(driverA.ID, 30*time.Second)
	if err != nil || !ok {
		t.Fatalf("ClaimMission(driver-a): mission=%#v ok=%v err=%v", claimA, ok, err)
	}
	claimB, ok, err := mgr.ClaimMission(driverB.ID, 30*time.Second)
	if err != nil || !ok {
		t.Fatalf("ClaimMission(driver-b): mission=%#v ok=%v err=%v", claimB, ok, err)
	}
	if claimA.WorkUnitID == "" || claimB.WorkUnitID == "" {
		t.Fatalf("claims missing work unit ids: %#v %#v", claimA, claimB)
	}
	if claimA.WorkUnitID == claimB.WorkUnitID {
		t.Fatalf("claims reused same work unit: %#v %#v", claimA, claimB)
	}
	if claimA.WorkUnit.Slice.WorkerIndex == claimB.WorkUnit.Slice.WorkerIndex {
		t.Fatalf("claims reused same worker slice: %#v %#v", claimA, claimB)
	}
}

func TestFailedWorkUnitRetriesUpToCeilingThenFailsJob(t *testing.T) {
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
<workload name="unit-retry-ceiling">
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

	var attempts []domain.MissionAttempt
	for i := 1; i <= 3; i++ {
		attempt, ok, err := mgr.ClaimMission(driver.ID, 30*time.Second)
		if err != nil || !ok {
			t.Fatalf("ClaimMission attempt %d: mission=%#v ok=%v err=%v", i, attempt, ok, err)
		}
		if attempt.Attempt != i {
			t.Fatalf("attempt number = %d want %d", attempt.Attempt, i)
		}
		attempts = append(attempts, attempt)
		if err := mgr.CompleteMission(attempt.ID, domain.MissionStatusFailed, "boom"); err != nil {
			t.Fatalf("CompleteMission attempt %d: %v", i, err)
		}
	}

	if _, ok, err := mgr.ClaimMission(driver.ID, 30*time.Second); err != nil || ok {
		t.Fatalf("expected no more retries, ok=%v err=%v", ok, err)
	}

	allAttempts := mgr.ListMissionAttempts()
	if len(allAttempts) != 3 {
		t.Fatalf("all attempts = %#v", allAttempts)
	}
	result, ok := mgr.GetJobResult(job.ID)
	if !ok || result.StageTotals[0].Status != domain.JobStatusFailed {
		t.Fatalf("result = %#v", result)
	}
	loaded, ok := mgr.GetJob(job.ID)
	if !ok || loaded.Status != domain.JobStatusFailed {
		t.Fatalf("job = %#v", loaded)
	}
}

func TestSuccessfulUnitsAggregateBackToWorkStageAndJob(t *testing.T) {
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
<workload name="unit-aggregate">
  <storage type="mock" />
  <workflow>
    <workstage name="main">
      <work name="fanout" workers="2" totalOps="2">
        <operation type="write" ratio="100" config="cprefix=t;containers=c(1);objects=s(1,2);sizes=c(1)KB" />
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

	claimA, ok, err := mgr.ClaimMission(driverA.ID, 30*time.Second)
	if err != nil || !ok {
		t.Fatalf("ClaimMission(driver-a): mission=%#v ok=%v err=%v", claimA, ok, err)
	}
	claimB, ok, err := mgr.ClaimMission(driverB.ID, 30*time.Second)
	if err != nil || !ok {
		t.Fatalf("ClaimMission(driver-b): mission=%#v ok=%v err=%v", claimB, ok, err)
	}

	sampleA := []legacyexec.Sample{{Timestamp: time.Now().UTC(), OpType: "write", OpCount: 1, ByteCount: 1000, TotalTimeMs: 10}}
	sampleB := []legacyexec.Sample{{Timestamp: time.Now().UTC(), OpType: "write", OpCount: 1, ByteCount: 1000, TotalTimeMs: 20}}
	if err := mgr.AppendMissionSamplesBatch(claimA.ID, "samples-a", sampleA); err != nil {
		t.Fatal(err)
	}
	if err := mgr.AppendMissionSamplesBatch(claimB.ID, "samples-b", sampleB); err != nil {
		t.Fatal(err)
	}
	if err := mgr.CompleteMission(claimA.ID, domain.MissionStatusSucceeded, ""); err != nil {
		t.Fatal(err)
	}
	if err := mgr.CompleteMission(claimB.ID, domain.MissionStatusSucceeded, ""); err != nil {
		t.Fatal(err)
	}

	result, ok := mgr.GetJobResult(job.ID)
	if !ok {
		t.Fatal("expected job result")
	}
	if result.Metrics.OperationCount != 2 || result.Metrics.ByteCount != 2000 {
		t.Fatalf("result = %#v", result)
	}
	if result.StageTotals[0].Status != domain.JobStatusSucceeded {
		t.Fatalf("stage totals = %#v", result.StageTotals)
	}
}

func TestRemoteJobProgressesToNextStageOnlyAfterCurrentStageSucceeds(t *testing.T) {
	store, err := snapshot.New(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	mgr, err := New(store)
	if err != nil {
		t.Fatal(err)
	}
	mgr.SetRemoteScheduling(true)
	driverA, err := mgr.RegisterDriverNode(domain.DriverNode{Name: "driver-a", Mode: domain.DriverModeDriver})
	if err != nil {
		t.Fatal(err)
	}
	driverB, err := mgr.RegisterDriverNode(domain.DriverNode{Name: "driver-b", Mode: domain.DriverModeDriver})
	if err != nil {
		t.Fatal(err)
	}
	job, err := mgr.CreateJobFromXML([]byte(`<?xml version="1.0" encoding="UTF-8"?>
<workload name="remote-multistage">
  <storage type="mock" />
  <workflow>
    <workstage name="stage-a">
      <work name="fanout-a" workers="2" totalOps="2">
        <operation type="write" ratio="100" config="cprefix=t;containers=c(1);objects=s(1,2);sizes=c(1)KB" />
      </work>
    </workstage>
    <workstage name="stage-b">
      <work name="fanout-b" workers="2" totalOps="2">
        <operation type="write" ratio="100" config="cprefix=t;containers=c(1);objects=s(3,4);sizes=c(1)KB" />
      </work>
    </workstage>
  </workflow>
</workload>`), "")
	if err != nil {
		t.Fatal(err)
	}
	if err := mgr.StartJob(context.Background(), job.ID); err != nil {
		t.Fatal(err)
	}

	firstA, ok, err := mgr.ClaimMission(driverA.ID, 30*time.Second)
	if err != nil || !ok {
		t.Fatalf("ClaimMission(driver-a/stage-a): mission=%#v ok=%v err=%v", firstA, ok, err)
	}
	firstB, ok, err := mgr.ClaimMission(driverB.ID, 30*time.Second)
	if err != nil || !ok {
		t.Fatalf("ClaimMission(driver-b/stage-a): mission=%#v ok=%v err=%v", firstB, ok, err)
	}
	if firstA.StageName != "stage-a" || firstB.StageName != "stage-a" {
		t.Fatalf("unexpected first-stage claims: %#v %#v", firstA, firstB)
	}
	if _, ok, err := mgr.ClaimMission(driverA.ID, 30*time.Second); err != nil || ok {
		t.Fatalf("expected no stage-b claim before stage-a completion, ok=%v err=%v", ok, err)
	}

	sampleA := []legacyexec.Sample{{Timestamp: time.Now().UTC(), OpType: "write", OpCount: 1, ByteCount: 1000, TotalTimeMs: 10}}
	sampleB := []legacyexec.Sample{{Timestamp: time.Now().UTC(), OpType: "write", OpCount: 1, ByteCount: 1000, TotalTimeMs: 10}}
	if err := mgr.AppendMissionSamplesBatch(firstA.ID, "a-1", sampleA); err != nil {
		t.Fatal(err)
	}
	if err := mgr.AppendMissionSamplesBatch(firstB.ID, "b-1", sampleB); err != nil {
		t.Fatal(err)
	}
	if err := mgr.CompleteMission(firstA.ID, domain.MissionStatusSucceeded, ""); err != nil {
		t.Fatal(err)
	}
	if _, ok, err := mgr.ClaimMission(driverA.ID, 30*time.Second); err != nil || ok {
		t.Fatalf("expected no stage-b claim while stage-a still has running units, ok=%v err=%v", ok, err)
	}
	if err := mgr.CompleteMission(firstB.ID, domain.MissionStatusSucceeded, ""); err != nil {
		t.Fatal(err)
	}

	secondA, ok, err := mgr.ClaimMission(driverA.ID, 30*time.Second)
	if err != nil || !ok {
		t.Fatalf("ClaimMission(driver-a/stage-b): mission=%#v ok=%v err=%v", secondA, ok, err)
	}
	secondB, ok, err := mgr.ClaimMission(driverB.ID, 30*time.Second)
	if err != nil || !ok {
		t.Fatalf("ClaimMission(driver-b/stage-b): mission=%#v ok=%v err=%v", secondB, ok, err)
	}
	if secondA.StageName != "stage-b" || secondB.StageName != "stage-b" {
		t.Fatalf("unexpected second-stage claims: %#v %#v", secondA, secondB)
	}
	if err := mgr.CompleteMission(secondA.ID, domain.MissionStatusSucceeded, ""); err != nil {
		t.Fatal(err)
	}
	if err := mgr.CompleteMission(secondB.ID, domain.MissionStatusSucceeded, ""); err != nil {
		t.Fatal(err)
	}

	loaded, ok := mgr.GetJob(job.ID)
	if !ok {
		t.Fatal("expected loaded job")
	}
	if len(loaded.Stages) != 2 {
		t.Fatalf("stages = %#v", loaded.Stages)
	}
	if loaded.Stages[0].StartedAt == nil || loaded.Stages[0].FinishedAt == nil {
		t.Fatalf("stage-a timestamps = %#v", loaded.Stages[0])
	}
	if loaded.Stages[1].StartedAt == nil || loaded.Stages[1].FinishedAt == nil {
		t.Fatalf("stage-b timestamps = %#v", loaded.Stages[1])
	}
	if loaded.Stages[0].FinishedAt.After(*loaded.Stages[1].StartedAt) {
		t.Fatalf("stage ordering violated: stage-a=%#v stage-b=%#v", loaded.Stages[0], loaded.Stages[1])
	}
}

func TestRemoteJobDoesNotProgressPastFailedStage(t *testing.T) {
	store, err := snapshot.New(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	mgr, err := New(store)
	if err != nil {
		t.Fatal(err)
	}
	mgr.SetRemoteScheduling(true)
	driver, err := mgr.RegisterDriverNode(domain.DriverNode{Name: "driver-a", Mode: domain.DriverModeDriver})
	if err != nil {
		t.Fatal(err)
	}
	job, err := mgr.CreateJobFromXML([]byte(`<?xml version="1.0" encoding="UTF-8"?>
<workload name="remote-stop-on-failure">
  <storage type="mock" />
  <workflow>
    <workstage name="stage-a">
      <work name="fanout-a" workers="1" totalOps="1">
        <operation type="write" ratio="100" config="cprefix=t;containers=c(1);objects=s(1,1);sizes=c(1)KB" />
      </work>
    </workstage>
    <workstage name="stage-b">
      <work name="fanout-b" workers="1" totalOps="1">
        <operation type="write" ratio="100" config="cprefix=t;containers=c(1);objects=s(2,2);sizes=c(1)KB" />
      </work>
    </workstage>
  </workflow>
</workload>`), "")
	if err != nil {
		t.Fatal(err)
	}
	if err := mgr.StartJob(context.Background(), job.ID); err != nil {
		t.Fatal(err)
	}

	for attemptNo := 1; attemptNo <= maxAttemptsPerWorkUnit; attemptNo++ {
		claim, ok, err := mgr.ClaimMission(driver.ID, 30*time.Second)
		if err != nil || !ok {
			t.Fatalf("ClaimMission attempt %d: mission=%#v ok=%v err=%v", attemptNo, claim, ok, err)
		}
		if claim.StageName != "stage-a" {
			t.Fatalf("unexpected stage claim: %#v", claim)
		}
		if err := mgr.CompleteMission(claim.ID, domain.MissionStatusFailed, "boom"); err != nil {
			t.Fatal(err)
		}
	}

	if _, ok, err := mgr.ClaimMission(driver.ID, 30*time.Second); err != nil || ok {
		t.Fatalf("expected no stage-b claim after failed stage-a, ok=%v err=%v", ok, err)
	}
	loaded, ok := mgr.GetJob(job.ID)
	if !ok || loaded.Status != domain.JobStatusFailed {
		t.Fatalf("job = %#v", loaded)
	}
	if len(loaded.Stages) < 2 || loaded.Stages[1].Status != domain.JobStatusCreated {
		t.Fatalf("job stages = %#v", loaded.Stages)
	}
}
