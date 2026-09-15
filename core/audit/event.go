package audit

import (
	"time"

	"github.com/google/uuid"
)

type Severity string

const (
	SeverityInfo     Severity = "info"
	SeverityLow      Severity = "low"
	SeverityMedium   Severity = "medium"
	SeverityHigh     Severity = "high"
	SeverityCritical Severity = "critical"
)

type Category string

const (
	CategoryAuthentication Category = "authentication"
	CategoryAuthorization  Category = "authorization"
	CategoryTenancy        Category = "tenancy"
	CategoryFraud          Category = "fraud"
	CategorySecrets        Category = "secrets"
	CategorySystem         Category = "system"
	CategoryAPI            Category = "api"
	CategoryWebhook        Category = "webhook"
	CategoryDatabase       Category = "database"
	CategorySecurity       Category = "security"
)

type Status string

const (
	StatusSuccess Status = "success"
	StatusFailed  Status = "failed"
	StatusBlocked Status = "blocked"
	StatusDenied  Status = "denied"
)

type Event struct {
	ID string `json:"id"`

	Timestamp time.Time `json:"timestamp"`

	Category Category `json:"category"`

	Action string `json:"action"`

	Status Status `json:"status"`

	Severity Severity `json:"severity"`

	TenantID string `json:"tenant_id,omitempty"`

	ActorID string `json:"actor_id,omitempty"`

	ActorType string `json:"actor_type,omitempty"`

	Resource string `json:"resource,omitempty"`

	ResourceID string `json:"resource_id,omitempty"`

	RequestID string `json:"request_id,omitempty"`

	SessionID string `json:"session_id,omitempty"`

	IPAddress string `json:"ip_address,omitempty"`

	UserAgent string `json:"user_agent,omitempty"`

	Message string `json:"message,omitempty"`

	Metadata map[string]any `json:"metadata,omitempty"`

	Tags []string `json:"tags,omitempty"`
}

func NewEvent(
	category Category,
	action string,
	severity Severity,
	status Status,
) *Event {
	return &Event{
		ID:        uuid.NewString(),
		Timestamp: time.Now().UTC(),
		Category:  category,
		Action:    action,
		Severity:  severity,
		Status:    status,
		Metadata:  make(map[string]any),
		Tags:      make([]string, 0),
	}
}

func (e *Event) AddTag(tag string) {
	if tag == "" {
		return
	}

	e.Tags = append(e.Tags, tag)
}

func (e *Event) SetMetadata(key string, value any) {
	if key == "" {
		return
	}

	if e.Metadata == nil {
		e.Metadata = make(map[string]any)
	}

	e.Metadata[key] = value
}

func (e *Event) IsCritical() bool {
	return e.Severity == SeverityCritical
}

func (e *Event) IsHighRisk() bool {
	return e.Severity == SeverityHigh ||
		e.Severity == SeverityCritical
}

func (e *Event) IsFailure() bool {
	return e.Status == StatusFailed ||
		e.Status == StatusBlocked ||
		e.Status == StatusDenied
}

func (e *Event) IsSuccess() bool {
	return e.Status == StatusSuccess
}

func (e *Event) Clone() *Event {
	metadata := make(map[string]any)

	for key, value := range e.Metadata {
		metadata[key] = value
	}

	tags := make([]string, len(e.Tags))
	copy(tags, e.Tags)

	return &Event{
		ID:         e.ID,
		Timestamp:  e.Timestamp,
		Category:   e.Category,
		Action:     e.Action,
		Status:     e.Status,
		Severity:   e.Severity,
		TenantID:   e.TenantID,
		ActorID:    e.ActorID,
		ActorType:  e.ActorType,
		Resource:   e.Resource,
		ResourceID: e.ResourceID,
		RequestID:  e.RequestID,
		SessionID:  e.SessionID,
		IPAddress:  e.IPAddress,
		UserAgent:  e.UserAgent,
		Message:    e.Message,
		Metadata:   metadata,
		Tags:       tags,
	}
}

func (e *Event) Validate() error {
	if e.ID == "" {
		return ErrInvalidEventID
	}

	if e.Action == "" {
		return ErrInvalidEventAction
	}

	if e.Category == "" {
		return ErrInvalidEventCategory
	}

	if e.Severity == "" {
		return ErrInvalidEventSeverity
	}

	if e.Status == "" {
		return ErrInvalidEventStatus
	}

	return nil
}
