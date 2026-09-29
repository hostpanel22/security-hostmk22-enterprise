package handlers

import (
	"net/http"
	"strings"

	"github.com/azezaz/HostMK-Security-Framework/internal/api/response"
)

func (a *API) Sessions(w http.ResponseWriter, r *http.Request) {
	id, ok := a.tenantScope(r)
	if !ok {
		response.Unauthorized(w, "tenant required")
		return
	}
	switch r.Method {
	case http.MethodGet:
		sessions, err := a.Store.ListSessions(r.Context(), id.TenantID)
		if err != nil {
			writeStoreError(w, err)
			return
		}
		response.OK(w, "sessions retrieved", map[string]any{"count": len(sessions), "sessions": sessions})
	case http.MethodDelete, http.MethodPost:
		sessionID := strings.TrimSpace(r.URL.Query().Get("id"))
		if sessionID == "" {
			response.BadRequest(w, "id is required")
			return
		}
		if r.Method == http.MethodPost && r.URL.Query().Get("action") != "revoke" {
			response.BadRequest(w, "unsupported action")
			return
		}
		if err := a.Store.RevokeSession(r.Context(), id.TenantID, sessionID); err != nil {
			writeStoreError(w, err)
			return
		}
		a.audit(r, id, "sessions.revoke", "sessions", sessionID)
		response.OK(w, "session revoked", map[string]string{"id": sessionID})
	default:
		response.Error(w, http.StatusMethodNotAllowed, "method not allowed")
	}
}
