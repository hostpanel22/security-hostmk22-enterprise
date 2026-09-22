package handlers

import (
	"net/http"
	"sync/atomic"
	"time"

	"github.com/azezaz/HostMK-Security-Framework/internal/api/response"
)

var applicationStartTime = time.Now().UTC()

var healthRequestCounter uint64

type ComponentHealth struct {
	Name      string `json:"name"`
	Status    string `json:"status"`
	Message   string `json:"message"`
	CheckedAt string `json:"checked_at"`
}

type HealthData struct {
	Status          string            `json:"status"`
	Service         string            `json:"service"`
	Version         string            `json:"version"`
	Environment     string            `json:"environment"`
	Timestamp       string            `json:"timestamp"`
	UptimeSeconds   int64             `json:"uptime_seconds"`
	RequestsHandled uint64            `json:"requests_handled"`
	Components      []ComponentHealth `json:"components"`
}

func Health(
	w http.ResponseWriter,
	r *http.Request,
) {
	if r.Method != http.MethodGet {
		response.Error(
			w,
			http.StatusMethodNotAllowed,
			"method not allowed",
		)
		return
	}

	atomic.AddUint64(
		&healthRequestCounter,
		1,
	)

	response.OK(
		w,
		"service healthy",
		buildHealthData(),
	)
}

func buildHealthData() HealthData {
	now := time.Now().UTC()

	return HealthData{
		Status:        "healthy",
		Service:       "HostMK Security Framework",
		Version:       "0.1.0",
		Environment:   "production",
		Timestamp:     now.Format(time.RFC3339),
		UptimeSeconds: uptimeSeconds(),
		RequestsHandled: atomic.LoadUint64(
			&healthRequestCounter,
		),
		Components: []ComponentHealth{
			buildAPIHealth(),
			buildAuditHealth(),
			buildFraudHealth(),
			buildRBACHealth(),
			buildPermissionsHealth(),
			buildSecretsHealth(),
			buildTenancyHealth(),
		},
	}
}

func uptimeSeconds() int64 {
	return int64(
		time.Since(
			applicationStartTime,
		).Seconds(),
	)
}

func buildAPIHealth() ComponentHealth {
	return ComponentHealth{
		Name:      "api",
		Status:    "healthy",
		Message:   "http server operational",
		CheckedAt: time.Now().UTC().Format(time.RFC3339),
	}
}

func buildAuditHealth() ComponentHealth {
	return ComponentHealth{
		Name:      "audit",
		Status:    "healthy",
		Message:   "audit module loaded",
		CheckedAt: time.Now().UTC().Format(time.RFC3339),
	}
}

func buildFraudHealth() ComponentHealth {
	return ComponentHealth{
		Name:      "fraud",
		Status:    "healthy",
		Message:   "fraud engine loaded",
		CheckedAt: time.Now().UTC().Format(time.RFC3339),
	}
}

func buildRBACHealth() ComponentHealth {
	return ComponentHealth{
		Name:      "rbac",
		Status:    "healthy",
		Message:   "rbac module loaded",
		CheckedAt: time.Now().UTC().Format(time.RFC3339),
	}
}

func buildPermissionsHealth() ComponentHealth {
	return ComponentHealth{
		Name:      "permissions",
		Status:    "healthy",
		Message:   "permission module loaded",
		CheckedAt: time.Now().UTC().Format(time.RFC3339),
	}
}

func buildSecretsHealth() ComponentHealth {
	return ComponentHealth{
		Name:      "secrets",
		Status:    "healthy",
		Message:   "secret manager loaded",
		CheckedAt: time.Now().UTC().Format(time.RFC3339),
	}
}

func buildTenancyHealth() ComponentHealth {
	return ComponentHealth{
		Name:      "tenancy",
		Status:    "healthy",
		Message:   "multi-tenant module loaded",
		CheckedAt: time.Now().UTC().Format(time.RFC3339),
	}
}
