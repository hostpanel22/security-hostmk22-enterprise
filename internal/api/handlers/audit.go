package handlers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/google/uuid"

	"github.com/azezaz/HostMK-Security-Framework/internal/api/response"
)

type AuditEvent struct {
	ID         string    `json:"id"`
	TenantID   string    `json:"tenant_id"`
	UserID     string    `json:"user_id"`
	Action     string    `json:"action"`
	Resource   string    `json:"resource"`
	Severity   string    `json:"severity"`
	IPAddress  string    `json:"ip_address"`
	UserAgent  string    `json:"user_agent"`
	OccurredAt time.Time `json:"occurred_at"`
}

func AuditEvents(
	w http.ResponseWriter,
	r *http.Request,
) {
	switch r.Method {

	case http.MethodGet:
		getAuditEvents(w, r)

	case http.MethodPost:
		createAuditEvent(w, r)

	default:
		response.Error(
			w,
			http.StatusMethodNotAllowed,
			"method not allowed",
		)
	}
}

func getAuditEvents(
	w http.ResponseWriter,
	r *http.Request,
) {
	limit := 50

	if value := r.URL.Query().Get("limit"); value != "" {

		parsed, err := strconv.Atoi(value)

		if err == nil && parsed > 0 && parsed <= 500 {
			limit = parsed
		}
	}

	events := []AuditEvent{
		{
			ID:         uuid.NewString(),
			TenantID:   "tenant-001",
			UserID:     "user-001",
			Action:     "LOGIN_SUCCESS",
			Resource:   "authentication",
			Severity:   "info",
			IPAddress:  "127.0.0.1",
			UserAgent:  "HostMK Security Framework",
			OccurredAt: time.Now().UTC(),
		},
	}

	response.OK(
		w,
		"audit events retrieved",
		map[string]any{
			"count":  len(events),
			"limit":  limit,
			"events": events,
		},
	)
}

func createAuditEvent(
	w http.ResponseWriter,
	r *http.Request,
) {
	event := AuditEvent{
		ID:         uuid.NewString(),
		TenantID:   "default",
		UserID:     "system",
		Action:     "MANUAL_EVENT",
		Resource:   "audit",
		Severity:   "info",
		IPAddress:  r.RemoteAddr,
		UserAgent:  r.UserAgent(),
		OccurredAt: time.Now().UTC(),
	}

	response.Created(
		w,
		"audit event created",
		event,
	)
}
