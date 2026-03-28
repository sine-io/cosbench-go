package web

import (
	"net/http"
	"strings"
)

func (h *Handler) requireDriverWriteAuth(w http.ResponseWriter, r *http.Request) bool {
	if strings.TrimSpace(h.driverSharedToken) == "" {
		http.Error(w, "driver shared token is not configured", http.StatusServiceUnavailable)
		return false
	}
	authz := strings.TrimSpace(r.Header.Get("Authorization"))
	if authz == "" {
		http.Error(w, "missing authorization header", http.StatusUnauthorized)
		return false
	}
	if !strings.HasPrefix(authz, "Bearer ") {
		http.Error(w, "malformed authorization header", http.StatusUnauthorized)
		return false
	}
	token := strings.TrimSpace(strings.TrimPrefix(authz, "Bearer "))
	if token == "" {
		http.Error(w, "malformed authorization header", http.StatusUnauthorized)
		return false
	}
	if token != h.driverSharedToken {
		http.Error(w, "invalid driver shared token", http.StatusForbidden)
		return false
	}
	return true
}
