// Package handler implements the HTTP API endpoints.
package handler

import (
	"encoding/json"
	"log/slog"
	"net"
	"net/http"
	"wol_admin/antishake"
	"wol_admin/nas"
)

// response is the unified API response structure.
type response struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// APIHandler holds dependencies for API handlers.
type APIHandler struct {
	locker *antishake.Locker
}

// NewAPIHandler creates an APIHandler with the given anti-shake locker.
func NewAPIHandler(l *antishake.Locker) *APIHandler {
	return &APIHandler{locker: l}
}

// WOL handles POST /api/wol — sends a Wake-on-LAN packet to the NAS.
func (h *APIHandler) WOL(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeJSON(w, http.StatusMethodNotAllowed, response{Code: -1, Message: "method not allowed"})
		return
	}

	clientID := r.Header.Get("X-Client-ID")
	if clientID == "" {
		clientID = r.RemoteAddr
	}

	if !h.locker.TryLock(clientID, "wol") {
		slog.Warn("WOL request rejected by anti-shake", "client", clientID)
		writeJSON(w, http.StatusTooManyRequests, response{Code: -1, Message: "request too frequent"})
		return
	}

	if err := nas.WOL(); err != nil {
		writeJSON(w, http.StatusInternalServerError, response{Code: -1, Message: err.Error()})
		return
	}

	writeJSON(w, http.StatusOK, response{Code: 0, Message: "WOL packet sent"})
}

// Shutdown handles POST /api/shutdown — sends SSH poweroff to the NAS.
func (h *APIHandler) Shutdown(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeJSON(w, http.StatusMethodNotAllowed, response{Code: -1, Message: "method not allowed"})
		return
	}

	clientID := r.Header.Get("X-Client-ID")
	if clientID == "" {
		clientID = r.RemoteAddr
	}

	if !h.locker.TryLock(clientID, "shutdown") {
		slog.Warn("shutdown request rejected by anti-shake", "client", clientID)
		writeJSON(w, http.StatusTooManyRequests, response{Code: -1, Message: "request too frequent"})
		return
	}

	if err := nas.Shutdown(); err != nil {
		writeJSON(w, http.StatusInternalServerError, response{Code: -1, Message: err.Error()})
		return
	}

	writeJSON(w, http.StatusOK, response{Code: 0, Message: "shutdown command sent"})
}

// IsLocalRequest checks if the request comes from localhost (127.0.0.1 or ::1).
func IsLocalRequest(r *http.Request) bool {
	host, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return false
	}
	return host == "127.0.0.1" || host == "::1"
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(v)
}
