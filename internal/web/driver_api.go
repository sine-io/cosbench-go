package web

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/sine-io/cosbench-go/internal/domain"
)

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
