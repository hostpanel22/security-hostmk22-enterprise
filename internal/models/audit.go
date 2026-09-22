package models

import (
	"fmt"
	"strings"
	"time"
)

type AuditSeverity string

const (
	AuditSeverityInfo     AuditSeverity = "info"
	AuditSeverityLow      AuditSeverity = "low"
	AuditSeverityMedium   AuditSeverity = "medium"
	AuditSeverityHigh     AuditSeverity = "high"
	AuditSeverityCritical AuditSeverity = "critical"
)

type AuditStatus string

const (
	AuditStatusSuccess AuditStatus = "success"
	AuditStatusFailed  AuditStatus = "failed"
	AuditStatusBlocked AuditStatus = "blocked"
	AuditStatusDenied  AuditStatus = "denied"
)

type AuditEvent struct {
	ID        string         `json:"id" db:"id"`
	TenantID  string         `json:"tenant_id" db:"tenant_id"`
	UserID    string         `json:"user_id" db:"user_id"`
	Action    string         `json:"action" db:"action"`
	Resource  string         `json:"resource" db:"resource"`
	Severity  AuditSeverity  `json:"severity" db:"severity"`
	Status    AuditStatus    `json:"status" db:"status"`
	IPAddress string         `json:"ip_address,omitempty" db:"ip_address"`
	UserAgent string         `json:"user_agent,omitempty" db:"user_agent"`
	Message   string         `json:"message,omitempty" db:"message"`
	Metadata  map[string]any `json:"metadata,omitempty" db:"metadata"`
	CreatedAt time.Time      `json:"created_at" db:"created_at"`
	UpdatedAt time.Time      `json:"updated_at" db:"updated_at"`
}

func (e *AuditEvent) Normalize() {
	if e == nil {
		return
	}

	e.TenantID = strings.TrimSpace(e.TenantID)
	e.UserID = strings.TrimSpace(e.UserID)
	e.Action = strings.TrimSpace(e.Action)
	e.Resource = strings.TrimSpace(e.Resource)
	e.Message = strings.TrimSpace(e.Message)
	e.IPAddress = strings.TrimSpace(e.IPAddress)
	e.UserAgent = strings.TrimSpace(e.UserAgent)

	if e.Severity == "" {
		e.Severity = AuditSeverityInfo
	}
	if e.Status == "" {
		e.Status = AuditStatusSuccess
	}
	if e.CreatedAt.IsZero() {
		e.CreatedAt = time.Now().UTC()
	}
	if e.UpdatedAt.IsZero() {
		e.UpdatedAt = e.CreatedAt
	}
}

func (e *AuditEvent) Validate() error {
	if e == nil {
		return fmt.Errorf("%w: audit event is required", ErrValidation)
	}
	if strings.TrimSpace(e.ID) == "" {
		return fmt.Errorf("%w: audit event id is required", ErrValidation)
	}
	if strings.TrimSpace(e.Action) == "" {
		return fmt.Errorf("%w: audit event action is required", ErrValidation)
	}
	if strings.TrimSpace(e.Resource) == "" {
		return fmt.Errorf("%w: audit event resource is required", ErrValidation)
	}
	if e.Severity == "" {
		return fmt.Errorf("%w: audit event severity is required", ErrValidation)
	}
	if e.Status == "" {
		return fmt.Errorf("%w: audit event status is required", ErrValidation)
	}

	switch e.Severity {
	case AuditSeverityInfo, AuditSeverityLow, AuditSeverityMedium, AuditSeverityHigh, AuditSeverityCritical:
	default:
		return fmt.Errorf("%w: unsupported audit severity %q", ErrValidation, e.Severity)
	}

	switch e.Status {
	case AuditStatusSuccess, AuditStatusFailed, AuditStatusBlocked, AuditStatusDenied:
	default:
		return fmt.Errorf("%w: unsupported audit status %q", ErrValidation, e.Status)
	}

	return nil
}

func (e *AuditEvent) IsCritical() bool {
	return e != nil && e.Severity == AuditSeverityCritical
}

func (e *AuditEvent) IsHighRisk() bool {
	return e != nil && (e.Severity == AuditSeverityHigh || e.Severity == AuditSeverityCritical)
}

func (e *AuditEvent) IsFailure() bool {
	return e != nil && (e.Status == AuditStatusFailed || e.Status == AuditStatusBlocked || e.Status == AuditStatusDenied)
}

func (e *AuditEvent) Clone() AuditEvent {
	if e == nil {
		return AuditEvent{}
	}

	clone := *e
	if e.Metadata != nil {
		clone.Metadata = make(map[string]any, len(e.Metadata))
		for k, v := range e.Metadata {
			clone.Metadata[k] = v
		}
	}
	return clone
}

func (e *AuditEvent) SetMetadata(key string, value any) {
	if e == nil || strings.TrimSpace(key) == "" {
		return
	}
	if e.Metadata == nil {
		e.Metadata = make(map[string]any)
	}
	e.Metadata[key] = value
	e.UpdatedAt = time.Now().UTC()
}

func (e *AuditEvent) AddTag(tag string) {
	if e == nil || strings.TrimSpace(tag) == "" {
		return
	}
	if e.Metadata == nil {
		e.Metadata = make(map[string]any)
	}
	if _, ok := e.Metadata["tags"]; !ok {
		e.Metadata["tags"] = []string{}
	}
	current, _ := e.Metadata["tags"].([]string)
	for _, existing := range current {
		if existing == tag {
			return
		}
	}
	e.Metadata["tags"] = append(current, tag)
	e.UpdatedAt = time.Now().UTC()
}
