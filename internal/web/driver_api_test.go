package web

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	legacyexec "github.com/sine-io/cosbench-go/internal/domain/execution"
	"github.com/sine-io/cosbench-go/internal/domain"
)

func TestDriverAPIRegisterHeartbeatAndClaim(t *testing.T) {
	h := newTestHandler(t)

	registerBody, _ := json.Marshal(map[string]any{
		"name": "driver-a",
		"mode": string(domain.DriverModeDriver),
	})
	registerRec := httptest.NewRecorder()
	registerReq := httptest.NewRequest(http.MethodPost, "/api/driver/register", bytes.NewReader(registerBody))
	registerReq.Header.Set("Content-Type", "application/json")
	h.ServeHTTP(registerRec, registerReq)
	if registerRec.Code != http.StatusOK {
		t.Fatalf("register status = %d body=%s", registerRec.Code, registerRec.Body.String())
	}
	var driver domain.DriverNode
	if err := json.Unmarshal(registerRec.Body.Bytes(), &driver); err != nil {
		t.Fatalf("register unmarshal: %v", err)
	}
	if driver.ID == "" || driver.Status != domain.DriverStatusHealthy {
		t.Fatalf("register payload = %#v", driver)
	}

	heartbeatAt := time.Now().UTC()
	heartbeatBody, _ := json.Marshal(map[string]any{
		"driver_id":     driver.ID,
		"heartbeat_at": heartbeatAt.Format(time.RFC3339Nano),
	})
	heartbeatRec := httptest.NewRecorder()
	heartbeatReq := httptest.NewRequest(http.MethodPost, "/api/driver/heartbeat", bytes.NewReader(heartbeatBody))
	heartbeatReq.Header.Set("Content-Type", "application/json")
	h.ServeHTTP(heartbeatRec, heartbeatReq)
	if heartbeatRec.Code != http.StatusOK {
		t.Fatalf("heartbeat status = %d body=%s", heartbeatRec.Code, heartbeatRec.Body.String())
	}

	job, err := h.manager.CreateJobFromXML([]byte(`<?xml version="1.0" encoding="UTF-8"?>
<workload name="driver-claim">
  <storage type="mock" />
  <workflow>
    <workstage name="main">
      <work name="work-a" workers="1" totalOps="1">
        <operation type="write" ratio="100" config="cprefix=t;containers=c(1);objects=c(1);sizes=c(1)KB" />
      </work>
    </workstage>
  </workflow>
</workload>`), "")
	if err != nil {
		t.Fatal(err)
	}
	if _, err := h.manager.ScheduleJobStage(job.ID); err != nil {
		t.Fatalf("ScheduleJobStage(): %v", err)
	}

	claimBody, _ := json.Marshal(map[string]any{
		"driver_id":         driver.ID,
		"lease_duration_ms": 30000,
	})
	claimRec := httptest.NewRecorder()
	claimReq := httptest.NewRequest(http.MethodPost, "/api/driver/missions/claim", bytes.NewReader(claimBody))
	claimReq.Header.Set("Content-Type", "application/json")
	h.ServeHTTP(claimRec, claimReq)
	if claimRec.Code != http.StatusOK {
		t.Fatalf("claim status = %d body=%s", claimRec.Code, claimRec.Body.String())
	}
	var mission domain.Mission
	if err := json.Unmarshal(claimRec.Body.Bytes(), &mission); err != nil {
		t.Fatalf("claim unmarshal: %v", err)
	}
	if mission.ID == "" || mission.Lease == nil || mission.Lease.DriverID != driver.ID {
		t.Fatalf("claim payload = %#v", mission)
	}
}

func TestDriverAPIReadEndpoints(t *testing.T) {
	h := newTestHandler(t)

	registerBody, _ := json.Marshal(map[string]any{
		"name": "driver-read",
		"mode": string(domain.DriverModeDriver),
	})
	registerRec := httptest.NewRecorder()
	registerReq := httptest.NewRequest(http.MethodPost, "/api/driver/register", bytes.NewReader(registerBody))
	registerReq.Header.Set("Content-Type", "application/json")
	h.ServeHTTP(registerRec, registerReq)
	if registerRec.Code != http.StatusOK {
		t.Fatalf("register status = %d body=%s", registerRec.Code, registerRec.Body.String())
	}
	var driver domain.DriverNode
	if err := json.Unmarshal(registerRec.Body.Bytes(), &driver); err != nil {
		t.Fatalf("register unmarshal: %v", err)
	}

	job, err := h.manager.CreateJobFromXML([]byte(`<?xml version="1.0" encoding="UTF-8"?>
<workload name="driver-api-read">
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
	if _, err := h.manager.ScheduleJobStage(job.ID); err != nil {
		t.Fatal(err)
	}
	mission, ok, err := h.manager.ClaimMission(driver.ID, 30*time.Second)
	if err != nil || !ok {
		t.Fatalf("ClaimMission(): mission=%#v ok=%v err=%v", mission, ok, err)
	}
	if err := h.manager.AppendMissionEvents(mission.ID, []domain.JobEvent{{OccurredAt: time.Now().UTC(), Level: domain.EventLevelInfo, Message: "mission started"}}); err != nil {
		t.Fatal(err)
	}

	selfRec := httptest.NewRecorder()
	selfReq := httptest.NewRequest(http.MethodGet, "/api/driver/self?driver_id="+driver.ID, nil)
	h.ServeHTTP(selfRec, selfReq)
	if selfRec.Code != http.StatusOK {
		t.Fatalf("self status = %d body=%s", selfRec.Code, selfRec.Body.String())
	}
	var overview domain.DriverOverview
	if err := json.Unmarshal(selfRec.Body.Bytes(), &overview); err != nil {
		t.Fatalf("self unmarshal: %v", err)
	}
	if overview.Driver.ID != driver.ID || overview.ActiveMissionCount != 1 {
		t.Fatalf("overview = %#v", overview)
	}

	missionsRec := httptest.NewRecorder()
	missionsReq := httptest.NewRequest(http.MethodGet, "/api/driver/missions?driver_id="+driver.ID, nil)
	h.ServeHTTP(missionsRec, missionsReq)
	if missionsRec.Code != http.StatusOK {
		t.Fatalf("missions status = %d body=%s", missionsRec.Code, missionsRec.Body.String())
	}
	var missions []domain.Mission
	if err := json.Unmarshal(missionsRec.Body.Bytes(), &missions); err != nil {
		t.Fatalf("missions unmarshal: %v", err)
	}
	if len(missions) != 1 || missions[0].ID != mission.ID {
		t.Fatalf("missions = %#v", missions)
	}

	detailRec := httptest.NewRecorder()
	detailReq := httptest.NewRequest(http.MethodGet, "/api/driver/missions/"+mission.ID, nil)
	h.ServeHTTP(detailRec, detailReq)
	if detailRec.Code != http.StatusOK {
		t.Fatalf("mission detail status = %d body=%s", detailRec.Code, detailRec.Body.String())
	}

	workersRec := httptest.NewRecorder()
	workersReq := httptest.NewRequest(http.MethodGet, "/api/driver/workers?driver_id="+driver.ID, nil)
	h.ServeHTTP(workersRec, workersReq)
	if workersRec.Code != http.StatusOK {
		t.Fatalf("workers status = %d body=%s", workersRec.Code, workersRec.Body.String())
	}
	var workerState domain.DriverWorkerState
	if err := json.Unmarshal(workersRec.Body.Bytes(), &workerState); err != nil {
		t.Fatalf("workers unmarshal: %v", err)
	}
	if workerState.DriverID != driver.ID || workerState.ActiveMissionCount != 1 {
		t.Fatalf("workerState = %#v", workerState)
	}

	logsRec := httptest.NewRecorder()
	logsReq := httptest.NewRequest(http.MethodGet, "/api/driver/logs?driver_id="+driver.ID, nil)
	h.ServeHTTP(logsRec, logsReq)
	if logsRec.Code != http.StatusOK {
		t.Fatalf("logs status = %d body=%s", logsRec.Code, logsRec.Body.String())
	}
	var logs []domain.JobEvent
	if err := json.Unmarshal(logsRec.Body.Bytes(), &logs); err != nil {
		t.Fatalf("logs unmarshal: %v", err)
	}
	if len(logs) == 0 {
		t.Fatal("expected driver logs")
	}
}

func TestDriverReportingEndpointsAreIdempotent(t *testing.T) {
	h := newTestHandler(t)

	registerBody, _ := json.Marshal(map[string]any{
		"name": "driver-idempotent",
		"mode": string(domain.DriverModeDriver),
	})
	registerRec := httptest.NewRecorder()
	registerReq := httptest.NewRequest(http.MethodPost, "/api/driver/register", bytes.NewReader(registerBody))
	registerReq.Header.Set("Content-Type", "application/json")
	h.ServeHTTP(registerRec, registerReq)
	if registerRec.Code != http.StatusOK {
		t.Fatalf("register status = %d body=%s", registerRec.Code, registerRec.Body.String())
	}
	var driver domain.DriverNode
	if err := json.Unmarshal(registerRec.Body.Bytes(), &driver); err != nil {
		t.Fatalf("register unmarshal: %v", err)
	}

	job, err := h.manager.CreateJobFromXML([]byte(`<?xml version="1.0" encoding="UTF-8"?>
<workload name="driver-idempotent-job">
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
	if _, err := h.manager.ScheduleJobStage(job.ID); err != nil {
		t.Fatal(err)
	}
	mission, ok, err := h.manager.ClaimMission(driver.ID, 30*time.Second)
	if err != nil || !ok {
		t.Fatalf("ClaimMission(): mission=%#v ok=%v err=%v", mission, ok, err)
	}

	eventPayload, _ := json.Marshal(map[string]any{
		"batch_id": "events-1",
		"events": []domain.JobEvent{{
			OccurredAt: time.Now().UTC(),
			Level:      domain.EventLevelInfo,
			Message:    "driver-http-event",
		}},
	})
	for range 2 {
		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPost, "/api/driver/missions/"+mission.ID+"/events", bytes.NewReader(eventPayload))
		req.Header.Set("Content-Type", "application/json")
		h.ServeHTTP(rec, req)
		if rec.Code != http.StatusOK {
			t.Fatalf("event status = %d body=%s", rec.Code, rec.Body.String())
		}
	}

	samplePayload, _ := json.Marshal(map[string]any{
		"batch_id": "samples-1",
		"samples": []legacyexec.Sample{{
			Timestamp:   time.Now().UTC(),
			OpType:      "write",
			OpCount:     1,
			ByteCount:   1000,
			TotalTimeMs: 10,
		}},
	})
	for range 2 {
		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPost, "/api/driver/missions/"+mission.ID+"/samples", bytes.NewReader(samplePayload))
		req.Header.Set("Content-Type", "application/json")
		h.ServeHTTP(rec, req)
		if rec.Code != http.StatusOK {
			t.Fatalf("sample status = %d body=%s", rec.Code, rec.Body.String())
		}
	}

	completePayload, _ := json.Marshal(map[string]any{
		"status":        domain.MissionStatusSucceeded,
		"error_message": "",
	})
	for range 2 {
		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPost, "/api/driver/missions/"+mission.ID+"/complete", bytes.NewReader(completePayload))
		req.Header.Set("Content-Type", "application/json")
		h.ServeHTTP(rec, req)
		if rec.Code != http.StatusOK {
			t.Fatalf("complete status = %d body=%s", rec.Code, rec.Body.String())
		}
	}

	logs := h.manager.GetJobEvents(job.ID)
	count := 0
	for _, event := range logs {
		if event.Message == "driver-http-event" {
			count++
		}
	}
	if count != 1 {
		t.Fatalf("driver-http-event count = %d logs=%#v", count, logs)
	}
	result, ok := h.manager.GetJobResult(job.ID)
	if !ok {
		t.Fatal("expected result")
	}
	if result.Metrics.OperationCount != 1 || result.Metrics.ByteCount != 1000 {
		t.Fatalf("result = %#v", result)
	}
}
