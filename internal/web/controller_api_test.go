package web

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/sine-io/cosbench-go/internal/controlplane"
	"github.com/sine-io/cosbench-go/internal/domain"
)

func TestControllerAPIJobsListAndDetail(t *testing.T) {
	h := newTestHandler(t)
	job := createCompletedControllerAPIJob(t, h.manager)

	listRec := httptest.NewRecorder()
	listReq := httptest.NewRequest(http.MethodGet, "/api/controller/jobs", nil)
	h.ServeHTTP(listRec, listReq)
	if listRec.Code != http.StatusOK {
		t.Fatalf("jobs status = %d body=%s", listRec.Code, listRec.Body.String())
	}
	var rows []domain.JobMatrixRow
	if err := json.Unmarshal(listRec.Body.Bytes(), &rows); err != nil {
		t.Fatalf("jobs unmarshal: %v", err)
	}
	if len(rows) == 0 || rows[0].OperationCount == 0 {
		t.Fatalf("jobs rows = %#v", rows)
	}

	detailRec := httptest.NewRecorder()
	detailReq := httptest.NewRequest(http.MethodGet, "/api/controller/jobs/"+job.ID, nil)
	h.ServeHTTP(detailRec, detailReq)
	if detailRec.Code != http.StatusOK {
		t.Fatalf("detail status = %d body=%s", detailRec.Code, detailRec.Body.String())
	}
	var payload struct {
		Job    domain.Job       `json:"job"`
		Result domain.JobResult `json:"result"`
		Events []domain.JobEvent `json:"events"`
	}
	if err := json.Unmarshal(detailRec.Body.Bytes(), &payload); err != nil {
		t.Fatalf("detail unmarshal: %v", err)
	}
	if payload.Job.ID != job.ID || payload.Result.JobID != job.ID || len(payload.Events) == 0 {
		t.Fatalf("detail payload = %#v", payload)
	}
}

func TestControllerAPIConfigAndStageEndpoints(t *testing.T) {
	h := newTestHandler(t)
	job := createCompletedControllerAPIJob(t, h.manager)

	configRec := httptest.NewRecorder()
	configReq := httptest.NewRequest(http.MethodGet, "/api/controller/jobs/"+job.ID+"/config", nil)
	h.ServeHTTP(configRec, configReq)
	if configRec.Code != http.StatusOK {
		t.Fatalf("config status = %d body=%s", configRec.Code, configRec.Body.String())
	}
	var configPayload struct {
		Job   domain.Job `json:"job"`
		RawXML string    `json:"raw_xml"`
	}
	if err := json.Unmarshal(configRec.Body.Bytes(), &configPayload); err != nil {
		t.Fatalf("config unmarshal: %v", err)
	}
	if configPayload.Job.ID != job.ID || configPayload.RawXML == "" {
		t.Fatalf("config payload = %#v", configPayload)
	}

	advancedRec := httptest.NewRecorder()
	advancedReq := httptest.NewRequest(http.MethodGet, "/api/controller/jobs/"+job.ID+"/config/advanced", nil)
	h.ServeHTTP(advancedRec, advancedReq)
	if advancedRec.Code != http.StatusOK {
		t.Fatalf("advanced status = %d body=%s", advancedRec.Code, advancedRec.Body.String())
	}
	var advancedPayload struct {
		Job                domain.Job      `json:"job"`
		NormalizedWorkload domain.Workload `json:"normalized_workload"`
	}
	if err := json.Unmarshal(advancedRec.Body.Bytes(), &advancedPayload); err != nil {
		t.Fatalf("advanced unmarshal: %v", err)
	}
	if advancedPayload.Job.ID != job.ID || len(advancedPayload.NormalizedWorkload.Workflow.Stages) == 0 {
		t.Fatalf("advanced payload = %#v", advancedPayload)
	}

	stageRec := httptest.NewRecorder()
	stageReq := httptest.NewRequest(http.MethodGet, "/api/controller/jobs/"+job.ID+"/stages/main", nil)
	h.ServeHTTP(stageRec, stageReq)
	if stageRec.Code != http.StatusOK {
		t.Fatalf("stage status = %d body=%s", stageRec.Code, stageRec.Body.String())
	}
	var stagePayload struct {
		JobID string            `json:"job_id"`
		Stage domain.StageState `json:"stage"`
	}
	if err := json.Unmarshal(stageRec.Body.Bytes(), &stagePayload); err != nil {
		t.Fatalf("stage unmarshal: %v", err)
	}
	if stagePayload.JobID != job.ID || stagePayload.Stage.Name != "main" {
		t.Fatalf("stage payload = %#v", stagePayload)
	}
}

func TestControllerAPITimelineJSONAndCSV(t *testing.T) {
	h := newTestHandler(t)
	job := createCompletedControllerAPIJob(t, h.manager)

	jsonRec := httptest.NewRecorder()
	jsonReq := httptest.NewRequest(http.MethodGet, "/api/controller/jobs/"+job.ID+"/timeline", nil)
	h.ServeHTTP(jsonRec, jsonReq)
	if jsonRec.Code != http.StatusOK {
		t.Fatalf("timeline status = %d body=%s", jsonRec.Code, jsonRec.Body.String())
	}
	var timeline domain.JobTimeline
	if err := json.Unmarshal(jsonRec.Body.Bytes(), &timeline); err != nil {
		t.Fatalf("timeline unmarshal: %v", err)
	}
	if timeline.JobID != job.ID || len(timeline.Job) == 0 || len(timeline.Stages["main"]) == 0 {
		t.Fatalf("timeline payload = %#v", timeline)
	}

	csvRec := httptest.NewRecorder()
	csvReq := httptest.NewRequest(http.MethodGet, "/api/controller/jobs/"+job.ID+"/timeline.csv", nil)
	h.ServeHTTP(csvRec, csvReq)
	if csvRec.Code != http.StatusOK {
		t.Fatalf("timeline csv status = %d body=%s", csvRec.Code, csvRec.Body.String())
	}
	if body := csvRec.Body.String(); body == "" || !strings.HasPrefix(body, "scope,name,timestamp,operation_count") {
		t.Fatalf("unexpected timeline csv: %s", body)
	}
}

func TestControllerAPIPrometheusAndArtifacts(t *testing.T) {
	h := newTestHandler(t)
	job := createCompletedControllerAPIJob(t, h.manager)

	promRec := httptest.NewRecorder()
	promReq := httptest.NewRequest(http.MethodGet, "/api/controller/metrics/prometheus", nil)
	h.ServeHTTP(promRec, promReq)
	if promRec.Code != http.StatusOK {
		t.Fatalf("prometheus status = %d body=%s", promRec.Code, promRec.Body.String())
	}
	if body := promRec.Body.String(); !strings.Contains(body, "cosbench_job_operation_count") {
		t.Fatalf("unexpected prometheus body: %s", body)
	}

	configRec := httptest.NewRecorder()
	configReq := httptest.NewRequest(http.MethodGet, "/api/controller/jobs/"+job.ID+"/artifacts/config", nil)
	h.ServeHTTP(configRec, configReq)
	if configRec.Code != http.StatusOK {
		t.Fatalf("config artifact status = %d body=%s", configRec.Code, configRec.Body.String())
	}
	if body := configRec.Body.String(); !strings.Contains(body, "<workload") {
		t.Fatalf("unexpected config artifact: %s", body)
	}

	logRec := httptest.NewRecorder()
	logReq := httptest.NewRequest(http.MethodGet, "/api/controller/jobs/"+job.ID+"/artifacts/log", nil)
	h.ServeHTTP(logRec, logReq)
	if logRec.Code != http.StatusOK {
		t.Fatalf("log artifact status = %d body=%s", logRec.Code, logRec.Body.String())
	}
	if body := logRec.Body.String(); !strings.Contains(strings.ToLower(body), "job finished") {
		t.Fatalf("unexpected log artifact: %s", body)
	}
}

func createCompletedControllerAPIJob(t *testing.T, mgr *controlplane.Manager) domain.Job {
	t.Helper()
	job, err := mgr.CreateJobFromXML([]byte(`<?xml version="1.0" encoding="UTF-8"?>
<workload name="controller-api">
  <storage type="mock" />
  <workflow>
    <workstage name="prepare">
      <work type="prepare" workers="1" config="cprefix=t;containers=c(1);objects=s(1,2);sizes=c(1)KB" />
    </workstage>
    <workstage name="main">
      <work name="main" workers="1" totalOps="2">
        <operation type="read" ratio="100" config="cprefix=t;containers=c(1);objects=s(1,2)" />
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
	waitForCompletedJob(t, mgr, job.ID)
	loaded, ok := mgr.GetJob(job.ID)
	if !ok {
		t.Fatal("job disappeared")
	}
	return loaded
}
