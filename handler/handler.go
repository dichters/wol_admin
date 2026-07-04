// Package handler implements the HTTP API endpoints.
package handler

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"wol_admin/antishake"
	"wol_admin/nas"
	"wol_admin/version"
)

// response is the unified API response structure.
type response struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// versionResponse is the response structure for the version endpoint.
type versionResponse struct {
	Code      int    `json:"code"`
	Version   string `json:"version"`
	Arch      string `json:"arch"`
	BuildTime string `json:"build_time"`
}

// APIHandler holds dependencies for API handlers.
type APIHandler struct {
	locker *antishake.Locker
}

// NewAPIHandler creates an APIHandler with the given anti-shake locker.
func NewAPIHandler(l *antishake.Locker) *APIHandler {
	return &APIHandler{locker: l}
}

// Version handles GET /wol/api/version — returns build version info.
func (h *APIHandler) Version(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, versionResponse{
		Code:      0,
		Version:   version.Version,
		Arch:      version.Arch,
		BuildTime: version.BuildTime,
	})
}

// WOL handles POST /wol/api/wol — sends a Wake-on-LAN packet to the NAS.
func (h *APIHandler) WOL(w http.ResponseWriter, r *http.Request) {
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

// Shutdown handles POST /wol/api/shutdown — sends SSH poweroff to the NAS.
func (h *APIHandler) Shutdown(w http.ResponseWriter, r *http.Request) {
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

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(v)
}
