package web

import (
	"context"
	"encoding/json"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/sine-io/cosbench-go/internal/controlplane"
	"github.com/sine-io/cosbench-go/internal/domain"
	"github.com/sine-io/cosbench-go/internal/snapshot"
)

func TestDashboardRenders(t *testing.T) {
	h := newTestHandler(t)
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	h.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d", rec.Code)
	}
	if !strings.Contains(rec.Body.String(), "Dashboard") {
		t.Fatalf("unexpected body: %s", rec.Body.String())
	}
}

func TestCreateWorkloadRedirectsToJob(t *testing.T) {
	h := newTestHandler(t)
	body := &strings.Builder{}
	writer := multipart.NewWriter(body)
	field, err := writer.CreateFormFile("workload", "sample.xml")
	if err != nil {
		t.Fatal(err)
	}
	_, _ = field.Write([]byte(`<?xml version="1.0" encoding="UTF-8"?><workload name="web-upload"><storage type="mock" /><workflow><workstage name="init"><work type="init" workers="1" config="cprefix=t;containers=c(1)" /></workstage></workflow></workload>`))
	_ = writer.Close()
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/workloads", strings.NewReader(body.String()))
	req.Header.Set("Content-Type", writer.FormDataContentType())
	h.ServeHTTP(rec, req)
	if rec.Code != http.StatusSeeOther {
		t.Fatalf("status = %d body=%s", rec.Code, rec.Body.String())
	}
	if !strings.HasPrefix(rec.Header().Get("Location"), "/jobs/job-") {
		t.Fatalf("location = %q", rec.Header().Get("Location"))
	}
}

func TestJobDetailShowsXMLAndErrors(t *testing.T) {
	h := newTestHandler(t)
	mgr := h.manager
	job, err := mgr.CreateJobFromXML([]byte(`<?xml version="1.0" encoding="UTF-8"?><workload name="detail"><storage type="mock" /><workflow><workstage name="main"><work name="missing" workers="1" totalOps="1"><operation type="read" ratio="100" config="cprefix=t;containers=c(1);objects=c(1)" /></work></workstage></workflow></workload>`), "")
	if err != nil {
		t.Fatal(err)
	}
	if err := mgr.StartJob(context.Background(), job.ID); err != nil {
		t.Fatal(err)
	}
	waitForCompletedJob(t, mgr, job.ID)
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/jobs/"+job.ID, nil)
	h.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d", rec.Code)
	}
	body := rec.Body.String()
	if !strings.Contains(body, "Uploaded XML") || !strings.Contains(body, "Execution note:") || !strings.Contains(body, "Work Summary") {
		t.Fatalf("unexpected job detail body: %s", body)
	}
}

func TestExportEndpoints(t *testing.T) {
	h := newTestHandler(t)
	mgr := h.manager
	job, err := mgr.CreateJobFromXML([]byte(`<?xml version="1.0" encoding="UTF-8"?><workload name="export"><storage type="mock" /><workflow><workstage name="prepare"><work type="prepare" workers="1" config="cprefix=t;containers=c(1);objects=s(1,2);sizes=c(1)KB" /></workstage></workflow></workload>`), "")
	if err != nil {
		t.Fatal(err)
	}
	if err := mgr.StartJob(context.Background(), job.ID); err != nil {
		t.Fatal(err)
	}
	waitForCompletedJob(t, mgr, job.ID)

	jsonRec := httptest.NewRecorder()
	jsonReq := httptest.NewRequest(http.MethodGet, "/exports/jobs/"+job.ID+"/result.json", nil)
	h.ServeHTTP(jsonRec, jsonReq)
	if jsonRec.Code != http.StatusOK {
		t.Fatalf("json status = %d body=%s", jsonRec.Code, jsonRec.Body.String())
	}
	var payload map[string]any
	if err := json.Unmarshal(jsonRec.Body.Bytes(), &payload); err != nil {
		t.Fatalf("json unmarshal: %v", err)
	}
	if payload["job"] == nil || payload["result"] == nil {
		t.Fatalf("unexpected export payload: %#v", payload)
	}
	resultPayload, ok := payload["result"].(map[string]any)
	if !ok {
		t.Fatalf("result payload = %#v", payload["result"])
	}
	stageTotals, ok := resultPayload["stage_totals"].([]any)
	if !ok || len(stageTotals) == 0 {
		t.Fatalf("stage totals = %#v", resultPayload["stage_totals"])
	}
	firstStage, ok := stageTotals[0].(map[string]any)
	if !ok {
		t.Fatalf("first stage = %#v", stageTotals[0])
	}
	workResults, ok := firstStage["work_results"].([]any)
	if !ok || len(workResults) == 0 {
		t.Fatalf("work results = %#v", firstStage["work_results"])
	}

	csvRec := httptest.NewRecorder()
	csvReq := httptest.NewRequest(http.MethodGet, "/exports/jobs/"+job.ID+"/result.csv", nil)
	h.ServeHTTP(csvRec, csvReq)
	if csvRec.Code != http.StatusOK {
		t.Fatalf("csv status = %d body=%s", csvRec.Code, csvRec.Body.String())
	}
	if !strings.Contains(csvRec.Body.String(), "scope,name,operation_count") || !strings.Contains(csvRec.Body.String(), ",work,") {
		t.Fatalf("unexpected csv export: %s", csvRec.Body.String())
	}
}

func TestRunningJobShowsCancelActionAndCanBeCancelled(t *testing.T) {
	h := newTestHandler(t)
	mgr := h.manager
	job, err := mgr.CreateJobFromXML([]byte(`<?xml version="1.0" encoding="UTF-8"?>
<workload name="cancel-web">
  <storage type="mock" />
  <workflow>
    <workstage name="main">
      <work name="main" workers="2" runtime="30">
        <operation type="write" ratio="100" config="cprefix=t;containers=c(1);objects=u(1,100);sizes=c(1)KB" />
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

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/jobs/"+job.ID, nil)
	h.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d", rec.Code)
	}
	if !strings.Contains(rec.Body.String(), "Cancel Job") {
		t.Fatalf("expected cancel action: %s", rec.Body.String())
	}

	cancelRec := httptest.NewRecorder()
	cancelReq := httptest.NewRequest(http.MethodPost, "/jobs/"+job.ID+"/cancel", nil)
	h.ServeHTTP(cancelRec, cancelReq)
	if cancelRec.Code != http.StatusSeeOther {
		t.Fatalf("cancel status = %d body=%s", cancelRec.Code, cancelRec.Body.String())
	}
	jobAfterCancel, _ := mgr.GetJob(job.ID)
	if jobAfterCancel.Status != domain.JobStatusCancelling && jobAfterCancel.Status != domain.JobStatusCancelled {
		t.Fatalf("job status after cancel request = %s", jobAfterCancel.Status)
	}
	waitForJobStatus(t, mgr, job.ID, domain.JobStatusCancelled)
}

func TestCancellingJobDetailHidesCancelAction(t *testing.T) {
	now := time.Now().UTC()
	store, err := snapshot.New(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	job := domain.Job{
		ID:        "job-cancelling",
		Name:      "cancelling",
		Status:    domain.JobStatusCancelling,
		Workload:  domain.Workload{Name: "cancelling"},
		CreatedAt: now,
		StartedAt: &now,
		Stages:    []domain.StageState{{Name: "main", Status: domain.JobStatusCancelling, StartedAt: &now}},
	}
	if err := store.SaveJob(job); err != nil {
		t.Fatal(err)
	}
	if err := store.SaveEvents(job.ID, []domain.JobEvent{{JobID: job.ID, OccurredAt: now, Level: domain.EventLevelInfo, Message: "job cancellation requested"}}); err != nil {
		t.Fatal(err)
	}
	mgr, err := controlplane.New(store)
	if err != nil {
		t.Fatal(err)
	}
	root, err := filepath.Abs(filepath.Join("..", ".."))
	if err != nil {
		t.Fatal(err)
	}
	h, err := NewHandler(mgr, filepath.Join(root, "web", "templates"))
	if err != nil {
		t.Fatal(err)
	}

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/jobs/"+job.ID, nil)
	h.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d", rec.Code)
	}
	body := rec.Body.String()
	if !strings.Contains(body, "cancelling") {
		t.Fatalf("expected cancelling status: %s", body)
	}
	if strings.Contains(body, "Cancel Job") {
		t.Fatalf("unexpected cancel action: %s", body)
	}
}

func TestControllerPagesRender(t *testing.T) {
	h := newTestHandler(t)
	job := createCompletedControllerAPIJob(t, h.manager)

	cases := []struct {
		path string
		want string
	}{
		{path: "/controller/matrix", want: "Controller Matrix"},
		{path: "/controller/jobs/" + job.ID + "/config", want: "Job Config"},
		{path: "/controller/jobs/" + job.ID + "/config/advanced", want: "Advanced Config"},
		{path: "/controller/jobs/" + job.ID + "/stages/main", want: "Stage Detail"},
		{path: "/controller/jobs/" + job.ID + "/timeline", want: "Job Timeline"},
	}

	for _, tc := range cases {
		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, tc.path, nil)
		h.ServeHTTP(rec, req)
		if rec.Code != http.StatusOK {
			t.Fatalf("%s status = %d body=%s", tc.path, rec.Code, rec.Body.String())
		}
		if !strings.Contains(rec.Body.String(), tc.want) {
			t.Fatalf("%s body = %s", tc.path, rec.Body.String())
		}
	}
}

func TestDriverPagesRender(t *testing.T) {
	h := newTestHandler(t)
	driver, mission := createDriverPageFixture(t, h.manager)

	cases := []struct {
		path string
		want string
	}{
		{path: "/driver", want: "Driver Dashboard"},
		{path: "/driver/missions", want: "Driver Missions"},
		{path: "/driver/missions/" + mission.ID, want: "Driver Mission Detail"},
		{path: "/driver/workers", want: "Driver Workers"},
	}

	for _, tc := range cases {
		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, tc.path, nil)
		h.ServeHTTP(rec, req)
		if rec.Code != http.StatusOK {
			t.Fatalf("%s status = %d body=%s", tc.path, rec.Code, rec.Body.String())
		}
		body := rec.Body.String()
		if !strings.Contains(body, tc.want) {
			t.Fatalf("%s body = %s", tc.path, body)
		}
		if !strings.Contains(body, driver.Name) {
			t.Fatalf("%s missing driver context: %s", tc.path, body)
		}
	}
}

func TestDriverLogsPageAndSharedNavigation(t *testing.T) {
	h := newTestHandler(t)
	driver, _ := createDriverPageFixture(t, h.manager)

	logsRec := httptest.NewRecorder()
	logsReq := httptest.NewRequest(http.MethodGet, "/driver/logs", nil)
	h.ServeHTTP(logsRec, logsReq)
	if logsRec.Code != http.StatusOK {
		t.Fatalf("logs status = %d body=%s", logsRec.Code, logsRec.Body.String())
	}
	logsBody := logsRec.Body.String()
	if !strings.Contains(logsBody, "Driver Logs") || !strings.Contains(logsBody, driver.Name) {
		t.Fatalf("logs body = %s", logsBody)
	}

	for _, path := range []string{"/", "/endpoints"} {
		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, path, nil)
		h.ServeHTTP(rec, req)
		if rec.Code != http.StatusOK {
			t.Fatalf("%s status = %d body=%s", path, rec.Code, rec.Body.String())
		}
		body := rec.Body.String()
		if !strings.Contains(body, "/driver") || !strings.Contains(body, "/controller/matrix") {
			t.Fatalf("%s body missing shared nav: %s", path, body)
		}
	}
}

func newTestHandler(t *testing.T) *Handler {
	t.Helper()
	store, err := snapshot.New(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	mgr, err := controlplane.New(store)
	if err != nil {
		t.Fatal(err)
	}
	_, err = mgr.CreateEndpoint(domain.EndpointConfig{Name: "mock", Type: domain.EndpointTypeMock})
	if err != nil {
		t.Fatal(err)
	}
	root, err := filepath.Abs(filepath.Join("..", ".."))
	if err != nil {
		t.Fatal(err)
	}
	h, err := NewHandler(mgr, filepath.Join(root, "web", "templates"))
	if err != nil {
		t.Fatal(err)
	}
	return h
}

func waitForCompletedJob(t *testing.T, mgr *controlplane.Manager, jobID string) {
	t.Helper()
	deadline := time.Now().Add(5 * time.Second)
	for time.Now().Before(deadline) {
		job, ok := mgr.GetJob(jobID)
		if !ok {
			t.Fatal("job disappeared")
		}
		if job.Status == domain.JobStatusSucceeded {
			return
		}
		time.Sleep(20 * time.Millisecond)
	}
	job, _ := mgr.GetJob(jobID)
	t.Fatalf("job did not complete: %#v", job)
}

func waitForJobStatus(t *testing.T, mgr *controlplane.Manager, jobID string, want domain.JobStatus) {
	t.Helper()
	deadline := time.Now().Add(5 * time.Second)
	for time.Now().Before(deadline) {
		job, ok := mgr.GetJob(jobID)
		if !ok {
			t.Fatal("job disappeared")
		}
		if job.Status == want {
			return
		}
		time.Sleep(20 * time.Millisecond)
	}
	job, _ := mgr.GetJob(jobID)
	t.Fatalf("job did not reach %s: %#v", want, job)
}

func createDriverPageFixture(t *testing.T, mgr *controlplane.Manager) (domain.DriverNode, domain.Mission) {
	t.Helper()
	driver, err := mgr.RegisterDriverNode(domain.DriverNode{Name: "driver-page", Mode: domain.DriverModeDriver})
	if err != nil {
		t.Fatal(err)
	}
	job, err := mgr.CreateJobFromXML([]byte(`<?xml version="1.0" encoding="UTF-8"?>
<workload name="driver-page-job">
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
	return driver, mission
}
