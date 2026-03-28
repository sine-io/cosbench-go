package controlplane

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/sine-io/cosbench-go/internal/domain"
	"github.com/sine-io/cosbench-go/internal/snapshot"
)

func TestManagerLifecycle(t *testing.T) {
	store, err := snapshot.New(t.TempDir())
	if err != nil {
		t.Fatalf("snapshot.New(): %v", err)
	}
	mgr, err := New(store)
	if err != nil {
		t.Fatalf("New(): %v", err)
	}
	endpoint, err := mgr.CreateEndpoint(domain.EndpointConfig{Name: "mock", Type: domain.EndpointTypeMock})
	if err != nil {
		t.Fatalf("CreateEndpoint(): %v", err)
	}
	job, err := mgr.CreateJobFromXML([]byte(testWorkloadXML), endpoint.ID)
	if err != nil {
		t.Fatalf("CreateJobFromXML(): %v", err)
	}
	if err := mgr.StartJob(context.Background(), job.ID); err != nil {
		t.Fatalf("StartJob(): %v", err)
	}
	deadline := time.Now().Add(5 * time.Second)
	for time.Now().Before(deadline) {
		updated, ok := mgr.GetJob(job.ID)
		if !ok {
			t.Fatal("job disappeared")
		}
		if updated.Status == domain.JobStatusSucceeded {
			if updated.Metrics.OperationCount == 0 {
				t.Fatalf("expected metrics, got %#v", updated.Metrics)
			}
			result, ok := mgr.GetJobResult(job.ID)
			if !ok || result.JobID != job.ID {
				t.Fatalf("result = %#v", result)
			}
			if len(result.StageTotals) != 3 {
				t.Fatalf("stage totals = %#v", result.StageTotals)
			}
			if len(result.StageTotals[0].WorkResults) != 1 || result.StageTotals[0].WorkResults[0].Name != "init" {
				t.Fatalf("missing init work result: %#v", result.StageTotals[0].WorkResults)
			}
			if len(result.StageTotals[2].WorkResults) != 1 || result.StageTotals[2].WorkResults[0].Name != "main" {
				t.Fatalf("missing main work result: %#v", result.StageTotals[2].WorkResults)
			}
			events := mgr.GetJobEvents(job.ID)
			if len(events) < 3 {
				t.Fatalf("events = %#v", events)
			}
			return
		}
		time.Sleep(20 * time.Millisecond)
	}
	updated, _ := mgr.GetJob(job.ID)
	t.Fatalf("job did not finish: %#v", updated)
}

func TestManagerRejectsInvalidXML(t *testing.T) {
	store, _ := snapshot.New(t.TempDir())
	mgr, _ := New(store)
	_, err := mgr.CreateJobFromXML([]byte("<workload"), "")
	if err == nil || !strings.Contains(err.Error(), "parse workload xml") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestManagerRestartRecoveryMarksRunningJobInterrupted(t *testing.T) {
	store, err := snapshot.New(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	job := domain.Job{
		ID: "job-running",
		Name: "running",
		Status: domain.JobStatusRunning,
		Workload: domain.Workload{Name: "running"},
		Stages: []domain.StageState{{Name: "main", Status: domain.JobStatusRunning}},
		CreatedAt: time.Now().UTC(),
	}
	if err := store.SaveJob(job); err != nil {
		t.Fatal(err)
	}
	mgr, err := New(store)
	if err != nil {
		t.Fatal(err)
	}
	loaded, ok := mgr.GetJob(job.ID)
	if !ok {
		t.Fatal("job not recovered")
	}
	if loaded.Status != domain.JobStatusInterrupted || loaded.Stages[0].Status != domain.JobStatusInterrupted {
		t.Fatalf("unexpected recovered state: %#v", loaded)
	}
}

func TestManagerFixtureEndToEndAndFailurePath(t *testing.T) {
	store, err := snapshot.New(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	mgr, err := New(store)
	if err != nil {
		t.Fatal(err)
	}
	endpoint, err := mgr.CreateEndpoint(domain.EndpointConfig{Name: "mock", Type: domain.EndpointTypeMock})
	if err != nil {
		t.Fatal(err)
	}

	successXML, err := os.ReadFile(filepath.Join("..", "..", "testdata", "workloads", "s3-active-subset.xml"))
	if err != nil {
		t.Fatal(err)
	}
	successJob, err := mgr.CreateJobFromXML(successXML, endpoint.ID)
	if err != nil {
		t.Fatal(err)
	}
	if err := mgr.StartJob(context.Background(), successJob.ID); err != nil {
		t.Fatal(err)
	}
	waitForJobStatus(t, mgr, successJob.ID, domain.JobStatusSucceeded)
	loadedSuccess, _ := mgr.GetJob(successJob.ID)
	if loadedSuccess.Metrics.OperationCount == 0 || loadedSuccess.Metrics.ByOperation == nil {
		t.Fatalf("unexpected success summary: %#v", loadedSuccess.Metrics)
	}

	failureXML, err := os.ReadFile(filepath.Join("..", "..", "testdata", "workloads", "mock-failure.xml"))
	if err != nil {
		t.Fatal(err)
	}
	failureJob, err := mgr.CreateJobFromXML(failureXML, endpoint.ID)
	if err != nil {
		t.Fatal(err)
	}
	if err := mgr.StartJob(context.Background(), failureJob.ID); err != nil {
		t.Fatal(err)
	}
	waitForJobStatus(t, mgr, failureJob.ID, domain.JobStatusSucceeded)
	loadedFailure, _ := mgr.GetJob(failureJob.ID)
	if loadedFailure.ErrorMessage == "" || loadedFailure.Status != domain.JobStatusSucceeded || loadedFailure.Metrics.ErrorCount == 0 {
		t.Fatalf("expected succeeded job with errors: %#v", loadedFailure)
	}
	events := mgr.GetJobEvents(failureJob.ID)
	foundError := false
	for _, event := range events {
		if event.Level == domain.EventLevelError {
			foundError = true
			break
		}
	}
	if !foundError {
		t.Fatalf("expected error event: %#v", events)
	}
}

func TestStartJobRejectsInvalidOperationConfigDuringPreflight(t *testing.T) {
	store, err := snapshot.New(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	mgr, err := New(store)
	if err != nil {
		t.Fatal(err)
	}
	job, err := mgr.CreateJobFromXML([]byte(`<?xml version="1.0" encoding="UTF-8"?>
<workload name="bad-config">
  <storage type="mock" />
  <workflow>
    <workstage name="main">
      <work name="main" workers="1" totalOps="1">
        <operation type="write" ratio="100" config="containers=bad;objects=c(1);sizes=c(1)KB" />
      </work>
    </workstage>
  </workflow>
</workload>`), "")
	if err != nil {
		t.Fatal(err)
	}

	err = mgr.StartJob(context.Background(), job.ID)
	if err == nil || !strings.Contains(err.Error(), "containers") {
		t.Fatalf("unexpected preflight error: %v", err)
	}
	loaded, _ := mgr.GetJob(job.ID)
	if loaded.Status != domain.JobStatusCreated {
		t.Fatalf("job status = %s", loaded.Status)
	}
}

func TestStartJobRejectsUnreadableMFileWriteInputDuringPreflight(t *testing.T) {
	store, err := snapshot.New(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	mgr, err := New(store)
	if err != nil {
		t.Fatal(err)
	}
	job, err := mgr.CreateJobFromXML([]byte(`<?xml version="1.0" encoding="UTF-8"?>
<workload name="bad-file">
  <storage type="sio" config="accesskey=test;secretkey=test;endpoint=http://127.0.0.1:9000;path_style_access=true" />
  <workflow>
    <workstage name="main">
      <work name="main" workers="1" totalOps="1">
        <operation type="mfilewrite" ratio="100" config="containers=c(1);objects=c(1);files=/definitely/missing/payload.bin" />
      </work>
    </workstage>
  </workflow>
</workload>`), "")
	if err != nil {
		t.Fatal(err)
	}

	err = mgr.StartJob(context.Background(), job.ID)
	if err == nil || !strings.Contains(err.Error(), "payload.bin") {
		t.Fatalf("unexpected preflight error: %v", err)
	}
	loaded, _ := mgr.GetJob(job.ID)
	if loaded.Status != domain.JobStatusCreated {
		t.Fatalf("job status = %s", loaded.Status)
	}
}

func TestStartJobRejectsUnreadableFileWriteInputDuringPreflight(t *testing.T) {
	store, err := snapshot.New(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	mgr, err := New(store)
	if err != nil {
		t.Fatal(err)
	}
	job, err := mgr.CreateJobFromXML([]byte(`<?xml version="1.0" encoding="UTF-8"?>
<workload name="bad-filewrite">
  <storage type="sio" config="accesskey=test;secretkey=test;endpoint=http://127.0.0.1:9000;path_style_access=true" />
  <workflow>
    <workstage name="main">
      <work name="main" workers="1" totalOps="1">
        <auth type="basic" config="username=work;password=secret" />
        <operation type="filewrite" ratio="100" config="containers=c(1);objects=c(1);files=/definitely/missing/filewrite.bin" />
      </work>
    </workstage>
  </workflow>
</workload>`), "")
	if err != nil {
		t.Fatal(err)
	}

	err = mgr.StartJob(context.Background(), job.ID)
	if err == nil || !strings.Contains(err.Error(), "filewrite.bin") {
		t.Fatalf("unexpected preflight error: %v", err)
	}
	loaded, _ := mgr.GetJob(job.ID)
	if loaded.Status != domain.JobStatusCreated {
		t.Fatalf("job status = %s", loaded.Status)
	}
}

func TestManagerCanCancelRunningJob(t *testing.T) {
	store, err := snapshot.New(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	mgr, err := New(store)
	if err != nil {
		t.Fatal(err)
	}
	endpoint, err := mgr.CreateEndpoint(domain.EndpointConfig{Name: "mock", Type: domain.EndpointTypeMock})
	if err != nil {
		t.Fatal(err)
	}
	job, err := mgr.CreateJobFromXML([]byte(`<?xml version="1.0" encoding="UTF-8"?>
<workload name="cancel-me">
  <storage type="mock" />
  <workflow>
    <workstage name="main">
      <work name="main" workers="2" runtime="30">
        <operation type="write" ratio="100" config="cprefix=t;containers=c(1);objects=u(1,100);sizes=c(1)KB" />
      </work>
    </workstage>
  </workflow>
</workload>`), endpoint.ID)
	if err != nil {
		t.Fatal(err)
	}
	if err := mgr.StartJob(context.Background(), job.ID); err != nil {
		t.Fatal(err)
	}

	if err := mgr.CancelJob(job.ID); err != nil {
		t.Fatalf("CancelJob(): %v", err)
	}
	mid, _ := mgr.GetJob(job.ID)
	if mid.Status != domain.JobStatusCancelling {
		t.Fatalf("job status before completion = %s", mid.Status)
	}
	waitForJobStatus(t, mgr, job.ID, domain.JobStatusCancelled)

	loaded, _ := mgr.GetJob(job.ID)
	if loaded.Status != domain.JobStatusCancelled {
		t.Fatalf("job status = %s", loaded.Status)
	}
	result, ok := mgr.GetJobResult(job.ID)
	if !ok {
		t.Fatal("expected partial result")
	}
	if result.JobID != job.ID {
		t.Fatalf("unexpected result: %#v", result)
	}
	events := mgr.GetJobEvents(job.ID)
	found := false
	for _, event := range events {
		if strings.Contains(strings.ToLower(event.Message), "cancel") {
			found = true
			break
		}
	}
	if !found {
		t.Fatalf("expected cancellation event: %#v", events)
	}
}

func TestManagerRestartRecoveryTurnsCancellingIntoCancelled(t *testing.T) {
	store, err := snapshot.New(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	job := domain.Job{
		ID:           "job-cancelling",
		Name:         "cancelling",
		Status:       domain.JobStatusCancelling,
		Workload:     domain.Workload{Name: "cancelling"},
		ErrorMessage: "job cancellation requested",
		Stages:       []domain.StageState{{Name: "main", Status: domain.JobStatusCancelling, ErrorMessage: "job cancellation requested"}},
		CreatedAt:    time.Now().UTC(),
	}
	if err := store.SaveJob(job); err != nil {
		t.Fatal(err)
	}
	mgr, err := New(store)
	if err != nil {
		t.Fatal(err)
	}
	loaded, ok := mgr.GetJob(job.ID)
	if !ok {
		t.Fatal("job not recovered")
	}
	if loaded.Status != domain.JobStatusCancelled || loaded.Stages[0].Status != domain.JobStatusCancelled {
		t.Fatalf("unexpected recovered state: %#v", loaded)
	}
	if !strings.Contains(strings.ToLower(loaded.ErrorMessage), "restart") {
		t.Fatalf("expected restart note: %#v", loaded)
	}
}

func TestManagerMockStageAwareFixtureSucceeds(t *testing.T) {
	store, err := snapshot.New(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	mgr, err := New(store)
	if err != nil {
		t.Fatal(err)
	}
	endpoint, err := mgr.CreateEndpoint(domain.EndpointConfig{Name: "mock", Type: domain.EndpointTypeMock})
	if err != nil {
		t.Fatal(err)
	}
	raw, err := os.ReadFile(filepath.Join("..", "..", "testdata", "workloads", "mock-stage-aware.xml"))
	if err != nil {
		t.Fatal(err)
	}
	job, err := mgr.CreateJobFromXML(raw, endpoint.ID)
	if err != nil {
		t.Fatal(err)
	}
	if err := mgr.StartJob(context.Background(), job.ID); err != nil {
		t.Fatal(err)
	}
	waitForJobStatus(t, mgr, job.ID, domain.JobStatusSucceeded)
	loaded, _ := mgr.GetJob(job.ID)
	if loaded.Metrics.ErrorCount != 0 {
		t.Fatalf("unexpected error count: %#v", loaded.Metrics)
	}
	if len(loaded.Stages) != 6 || loaded.Stages[2].Metrics.OperationCount == 0 || loaded.Stages[3].Metrics.OperationCount == 0 {
		t.Fatalf("unexpected stage metrics: %#v", loaded.Stages)
	}
}

func TestManagerReuseDataFixtureSucceeds(t *testing.T) {
	store, err := snapshot.New(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	mgr, err := New(store)
	if err != nil {
		t.Fatal(err)
	}
	endpoint, err := mgr.CreateEndpoint(domain.EndpointConfig{Name: "mock", Type: domain.EndpointTypeMock})
	if err != nil {
		t.Fatal(err)
	}
	raw, err := os.ReadFile(filepath.Join("..", "..", "testdata", "workloads", "mock-reusedata-subset.xml"))
	if err != nil {
		t.Fatal(err)
	}
	job, err := mgr.CreateJobFromXML(raw, endpoint.ID)
	if err != nil {
		t.Fatal(err)
	}
	if err := mgr.StartJob(context.Background(), job.ID); err != nil {
		t.Fatal(err)
	}
	waitForJobStatus(t, mgr, job.ID, domain.JobStatusSucceeded)
	loaded, _ := mgr.GetJob(job.ID)
	if loaded.Metrics.ErrorCount != 0 {
		t.Fatalf("unexpected error count: %#v", loaded.Metrics)
	}
	if len(loaded.Stages) != 6 {
		t.Fatalf("unexpected stage count: %#v", loaded.Stages)
	}
	if loaded.Stages[2].Metrics.OperationCount == 0 {
		t.Fatalf("main-read metrics = %#v", loaded.Stages[2].Metrics)
	}
	if loaded.Stages[3].Metrics.ByteCount == 0 {
		t.Fatalf("main-list metrics = %#v", loaded.Stages[3].Metrics)
	}
}

func waitForJobStatus(t *testing.T, mgr *Manager, jobID string, want domain.JobStatus) {
	t.Helper()
	deadline := time.Now().Add(5 * time.Second)
	for time.Now().Before(deadline) {
		job, ok := mgr.GetJob(jobID)
		if !ok {
			t.Fatalf("job disappeared: %s", jobID)
		}
		if job.Status == want {
			return
		}
		time.Sleep(20 * time.Millisecond)
	}
	job, _ := mgr.GetJob(jobID)
	t.Fatalf("job did not reach %s: %#v", want, job)
}

const testWorkloadXML = `<?xml version="1.0" encoding="UTF-8"?>
<workload name="test" description="minimal">
  <storage type="mock" />
  <workflow>
    <workstage name="init">
      <work type="init" workers="1" config="cprefix=t;containers=c(1)" />
    </workstage>
    <workstage name="prepare">
      <work type="prepare" workers="1" config="cprefix=t;containers=c(1);objects=s(1,2);sizes=c(1)KB" />
    </workstage>
    <workstage name="main">
      <work name="main" workers="1" totalOps="3">
        <operation type="write" ratio="100" config="cprefix=t;containers=c(1);objects=s(3,5);sizes=c(1)KB" />
      </work>
    </workstage>
  </workflow>
</workload>`
