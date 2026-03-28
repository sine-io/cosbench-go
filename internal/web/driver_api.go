package web

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

	legacyexec "github.com/sine-io/cosbench-go/internal/domain/execution"
	"github.com/sine-io/cosbench-go/internal/domain"
)

func (h *Handler) driverSelf(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	driverID := strings.TrimSpace(r.URL.Query().Get("driver_id"))
	overview, ok := h.manager.GetDriverOverview(driverID)
	if !ok {
		http.NotFound(w, r)
		return
	}
	writeJSON(w, http.StatusOK, overview)
}

func (h *Handler) driverMissions(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	driverID := strings.TrimSpace(r.URL.Query().Get("driver_id"))
	writeJSON(w, http.StatusOK, h.manager.ListDriverMissions(driverID))
}

func (h *Handler) driverWorkers(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	driverID := strings.TrimSpace(r.URL.Query().Get("driver_id"))
	state, ok := h.manager.GetDriverWorkerState(driverID)
	if !ok {
		http.NotFound(w, r)
		return
	}
	writeJSON(w, http.StatusOK, state)
}

func (h *Handler) driverLogs(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	driverID := strings.TrimSpace(r.URL.Query().Get("driver_id"))
	writeJSON(w, http.StatusOK, h.manager.GetDriverLogs(driverID))
}

func (h *Handler) driverRegister(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	var input struct {
		Name string          `json:"name"`
		Mode domain.DriverMode `json:"mode"`
	}
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	driver, err := h.manager.RegisterDriverNode(domain.DriverNode{Name: input.Name, Mode: input.Mode})
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	writeJSON(w, http.StatusOK, driver)
}

func (h *Handler) driverHeartbeat(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	var input struct {
		DriverID    string `json:"driver_id"`
		HeartbeatAt string `json:"heartbeat_at"`
	}
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	at, err := time.Parse(time.RFC3339Nano, input.HeartbeatAt)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if err := h.manager.RecordDriverHeartbeat(input.DriverID, at); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	driver, _ := h.manager.GetDriverNode(input.DriverID)
	writeJSON(w, http.StatusOK, driver)
}

func (h *Handler) driverClaimMission(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	var input struct {
		DriverID        string `json:"driver_id"`
		LeaseDurationMs int    `json:"lease_duration_ms"`
	}
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	mission, ok, err := h.manager.ClaimMission(input.DriverID, time.Duration(input.LeaseDurationMs)*time.Millisecond)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if !ok {
		w.WriteHeader(http.StatusNoContent)
		return
	}
	writeJSON(w, http.StatusOK, mission)
}

func (h *Handler) driverMissionRoute(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimPrefix(r.URL.Path, "/api/driver/missions/")
	path = strings.Trim(path, "/")
	parts := strings.Split(path, "/")
	if len(parts) == 1 && r.Method == http.MethodGet {
		mission, ok := h.manager.GetMission(parts[0])
		if !ok {
			http.NotFound(w, r)
			return
		}
		writeJSON(w, http.StatusOK, mission)
		return
	}
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	if len(parts) != 2 {
		http.NotFound(w, r)
		return
	}
	missionID, action := parts[0], parts[1]
	switch action {
	case "events":
		var payload struct {
			BatchID string            `json:"batch_id"`
			Events  []domain.JobEvent `json:"events"`
		}
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		if err := h.manager.AppendMissionEventsBatch(missionID, payload.BatchID, payload.Events); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
	case "samples":
		var payload struct {
			BatchID string              `json:"batch_id"`
			Samples []legacyexec.Sample `json:"samples"`
		}
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		if err := h.manager.AppendMissionSamplesBatch(missionID, payload.BatchID, payload.Samples); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
	case "complete":
		var input struct {
			Status       domain.MissionStatus `json:"status"`
			ErrorMessage string               `json:"error_message"`
		}
		if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		if err := h.manager.CompleteMission(missionID, input.Status, input.ErrorMessage); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
	default:
		http.NotFound(w, r)
	}
}
