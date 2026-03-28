package web

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

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
