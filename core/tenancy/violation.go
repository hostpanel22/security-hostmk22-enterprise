package tenancy

import (
	"time"

	"github.com/google/uuid"
)

/*
|--------------------------------------------------------------------------
| Severity Levels
|--------------------------------------------------------------------------
*/

const (
	SeverityInfo     = "info"
	SeverityLow      = "low"
	SeverityMedium   = "medium"
	SeverityHigh     = "high"
	SeverityCritical = "critical"
)

/*
|--------------------------------------------------------------------------
| Violation Types
|--------------------------------------------------------------------------
*/

const (
	ViolationCrossTenantAccess = "cross_tenant_access"

	ViolationCrossTenantQuery = "cross_tenant_query"

	ViolationCrossTenantResource = "cross_tenant_resource"

	ViolationInvalidTenantContext = "invalid_tenant_context"

	ViolationUnauthorizedTenantAccess = "unauthorized_tenant_access"

	ViolationTenantMismatch = "tenant_mismatch"

	ViolationSuspiciousActivity = "suspicious_activity"

	ViolationSecurityPolicy = "security_policy_violation"
)

/*
|--------------------------------------------------------------------------
| Security Violation
|--------------------------------------------------------------------------
|
| Represents a detected security incident.
| Every violation can be persisted, exported,
| audited and forwarded to external monitoring
| systems such as SIEM platforms.
|
*/

type Violation struct {
	ID string `json:"id"`

	Type string `json:"type"`

	Severity string `json:"severity"`

	TenantID string `json:"tenant_id"`

	ActorID string `json:"actor_id"`

	RequestID string `json:"request_id"`

	Resource string `json:"resource"`

	Action string `json:"action"`

	IPAddress string `json:"ip_address"`

	UserAgent string `json:"user_agent"`

	Message string `json:"message"`

	Metadata map[string]any `json:"metadata,omitempty"`

	DetectedAt time.Time `json:"detected_at"`
}

/*
|--------------------------------------------------------------------------
| Constructors
|--------------------------------------------------------------------------
*/

/*
NewViolation creates a new security violation.
*/
func NewViolation(
	violationType string,
	severity string,
	tenantID string,
	message string,
) *Violation {

	return &Violation{
		ID:         uuid.NewString(),
		Type:       violationType,
		Severity:   severity,
		TenantID:   tenantID,
		Message:    message,
		DetectedAt: time.Now().UTC(),
		Metadata:   make(map[string]any),
	}
}

/*
|--------------------------------------------------------------------------
| Severity Helpers
|--------------------------------------------------------------------------
*/

func (v *Violation) IsCritical() bool {
	return v.Severity == SeverityCritical
}

func (v *Violation) IsHigh() bool {
	return v.Severity == SeverityHigh
}

func (v *Violation) IsMedium() bool {
	return v.Severity == SeverityMedium
}

func (v *Violation) IsLow() bool {
	return v.Severity == SeverityLow
}

/*
|--------------------------------------------------------------------------
| Classification Helpers
|--------------------------------------------------------------------------
*/

func (v *Violation) IsCrossTenantViolation() bool {
	switch v.Type {

	case ViolationCrossTenantAccess:
		return true

	case ViolationCrossTenantQuery:
		return true

	case ViolationCrossTenantResource:
		return true

	default:
		return false
	}
}

func (v *Violation) RequiresImmediateAction() bool {
	return v.IsCritical() ||
		v.IsHigh()
}

/*
|--------------------------------------------------------------------------
| Metadata
|--------------------------------------------------------------------------
*/

func (v *Violation) AddMetadata(
	key string,
	value any,
) {
	if v.Metadata == nil {
		v.Metadata = make(map[string]any)
	}

	v.Metadata[key] = value
}

/*
|--------------------------------------------------------------------------
| Serialization
|--------------------------------------------------------------------------
*/

func (v *Violation) ToMap() map[string]any {
	return map[string]any{
		"id":          v.ID,
		"type":        v.Type,
		"severity":    v.Severity,
		"tenant_id":   v.TenantID,
		"actor_id":    v.ActorID,
		"request_id":  v.RequestID,
		"resource":    v.Resource,
		"action":      v.Action,
		"ip_address":  v.IPAddress,
		"user_agent":  v.UserAgent,
		"message":     v.Message,
		"metadata":    v.Metadata,
		"detected_at": v.DetectedAt,
	}
}

/*
|--------------------------------------------------------------------------
| Audit Event Conversion
|--------------------------------------------------------------------------
|
| This structure can later be forwarded
| to audit engines, SIEM systems,
| security dashboards or alerting services.
|
*/

func (v *Violation) ToAuditEvent() map[string]any {
	return map[string]any{
		"event":      "security_violation",
		"violation":  v.Type,
		"severity":   v.Severity,
		"tenant_id":  v.TenantID,
		"actor_id":   v.ActorID,
		"resource":   v.Resource,
		"action":     v.Action,
		"timestamp":  v.DetectedAt,
		"request_id": v.RequestID,
	}
}
