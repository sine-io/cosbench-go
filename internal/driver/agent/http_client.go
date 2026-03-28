package agent

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	legacyexec "github.com/sine-io/cosbench-go/internal/domain/execution"
	"github.com/sine-io/cosbench-go/internal/domain"
)

type HTTPClient struct {
	BaseURL    string
	HTTPClient *http.Client
}

func (c *HTTPClient) client() *http.Client {
	if c.HTTPClient != nil {
		return c.HTTPClient
	}
	return http.DefaultClient
}

func (c *HTTPClient) RegisterDriver(name string, mode domain.DriverMode) (domain.DriverNode, error) {
	var driver domain.DriverNode
	err := c.postJSON("/api/driver/register", map[string]any{
		"name": name,
		"mode": mode,
	}, &driver)
	return driver, err
}

func (c *HTTPClient) Heartbeat(driverID string, at time.Time) (domain.DriverNode, error) {
	var driver domain.DriverNode
	err := c.postJSON("/api/driver/heartbeat", map[string]any{
		"driver_id":    driverID,
		"heartbeat_at": at.Format(time.RFC3339Nano),
	}, &driver)
	return driver, err
}

func (c *HTTPClient) ClaimMission(driverID string, leaseTTL time.Duration) (domain.Mission, bool, error) {
	var mission domain.Mission
	resp, err := c.doPost("/api/driver/missions/claim", map[string]any{
		"driver_id":         driverID,
		"lease_duration_ms": int(leaseTTL / time.Millisecond),
	})
	if err != nil {
		return domain.Mission{}, false, err
	}
	defer resp.Body.Close()
	if resp.StatusCode == http.StatusNoContent {
		return domain.Mission{}, false, nil
	}
	if resp.StatusCode != http.StatusOK {
		return domain.Mission{}, false, fmt.Errorf("claim mission: unexpected status %d", resp.StatusCode)
	}
	if err := json.NewDecoder(resp.Body).Decode(&mission); err != nil {
		return domain.Mission{}, false, err
	}
	return mission, true, nil
}

func (c *HTTPClient) UploadEvents(missionID string, events []domain.JobEvent) error {
	return c.UploadEventsBatch(missionID, fmt.Sprintf("events-%d", time.Now().UnixNano()), events)
}

func (c *HTTPClient) UploadSamples(missionID string, samples []legacyexec.Sample) error {
	return c.UploadSamplesBatch(missionID, fmt.Sprintf("samples-%d", time.Now().UnixNano()), samples)
}

func (c *HTTPClient) CompleteMission(missionID string, status domain.MissionStatus, errorMessage string) error {
	return c.postJSON("/api/driver/missions/"+missionID+"/complete", map[string]any{
		"status":        status,
		"error_message": errorMessage,
	}, nil)
}

func (c *HTTPClient) UploadEventsBatch(missionID string, batchID string, events []domain.JobEvent) error {
	return c.postJSON("/api/driver/missions/"+missionID+"/events", map[string]any{
		"batch_id": batchID,
		"events":   events,
	}, nil)
}

func (c *HTTPClient) UploadSamplesBatch(missionID string, batchID string, samples []legacyexec.Sample) error {
	return c.postJSON("/api/driver/missions/"+missionID+"/samples", map[string]any{
		"batch_id": batchID,
		"samples":  samples,
	}, nil)
}

func (c *HTTPClient) postJSON(path string, payload any, out any) error {
	resp, err := c.doPost(path, payload)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("post %s: unexpected status %d", path, resp.StatusCode)
	}
	if out == nil {
		return nil
	}
	return json.NewDecoder(resp.Body).Decode(out)
}

func (c *HTTPClient) doPost(path string, payload any) (*http.Response, error) {
	body, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest(http.MethodPost, strings.TrimRight(c.BaseURL, "/")+path, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	return c.client().Do(req)
}
