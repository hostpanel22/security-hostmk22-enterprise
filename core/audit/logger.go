package audit

import (
	"context"
	"log"
	"sync"
)

type Logger interface {
	Log(ctx context.Context, event *Event) error
}

type AuditLogger struct {
	store Store

	mu sync.RWMutex

	enabled bool
}

func NewAuditLogger(store Store) *AuditLogger {
	return &AuditLogger{
		store:   store,
		enabled: true,
	}
}

func (l *AuditLogger) Enable() {
	l.mu.Lock()
	defer l.mu.Unlock()

	l.enabled = true
}

func (l *AuditLogger) Disable() {
	l.mu.Lock()
	defer l.mu.Unlock()

	l.enabled = false
}

func (l *AuditLogger) IsEnabled() bool {
	l.mu.RLock()
	defer l.mu.RUnlock()

	return l.enabled
}

func (l *AuditLogger) Log(
	ctx context.Context,
	event *Event,
) error {

	if event == nil {
		return ErrNilEvent
	}

	if err := event.Validate(); err != nil {
		return err
	}

	if !l.IsEnabled() {
		return nil
	}

	if l.store == nil {
		log.Printf(
			"[AUDIT] category=%s action=%s severity=%s status=%s tenant=%s actor=%s message=%s",
			event.Category,
			event.Action,
			event.Severity,
			event.Status,
			event.TenantID,
			event.ActorID,
			event.Message,
		)

		return nil
	}

	return l.store.Save(
		ctx,
		event,
	)
}

func (l *AuditLogger) LogSuccess(
	ctx context.Context,
	category Category,
	action string,
	message string,
) error {

	event := NewEvent(
		category,
		action,
		SeverityInfo,
		StatusSuccess,
	)

	event.Message = message

	return l.Log(
		ctx,
		event,
	)
}

func (l *AuditLogger) LogFailure(
	ctx context.Context,
	category Category,
	action string,
	message string,
) error {

	event := NewEvent(
		category,
		action,
		SeverityMedium,
		StatusFailed,
	)

	event.Message = message

	return l.Log(
		ctx,
		event,
	)
}

func (l *AuditLogger) LogSecurityViolation(
	ctx context.Context,
	category Category,
	action string,
	message string,
) error {

	event := NewEvent(
		category,
		action,
		SeverityCritical,
		StatusBlocked,
	)

	event.Message = message

	event.AddTag("security")
	event.AddTag("violation")

	return l.Log(
		ctx,
		event,
	)
}

func (l *AuditLogger) LogAuthentication(
	ctx context.Context,
	tenantID string,
	actorID string,
	success bool,
	message string,
) error {

	status := StatusSuccess

	if !success {
		status = StatusFailed
	}

	event := NewEvent(
		CategoryAuthentication,
		"user_login",
		SeverityLow,
		status,
	)

	event.TenantID = tenantID
	event.ActorID = actorID
	event.Message = message

	return l.Log(
		ctx,
		event,
	)
}

func (l *AuditLogger) LogAuthorization(
	ctx context.Context,
	tenantID string,
	actorID string,
	resource string,
	action string,
	allowed bool,
) error {

	status := StatusSuccess

	if !allowed {
		status = StatusDenied
	}

	event := NewEvent(
		CategoryAuthorization,
		action,
		SeverityMedium,
		status,
	)

	event.TenantID = tenantID
	event.ActorID = actorID
	event.Resource = resource

	return l.Log(
		ctx,
		event,
	)
}

func (l *AuditLogger) LogTenantViolation(
	ctx context.Context,
	tenantID string,
	actorID string,
	resource string,
	message string,
) error {

	event := NewEvent(
		CategoryTenancy,
		"cross_tenant_access",
		SeverityCritical,
		StatusBlocked,
	)

	event.TenantID = tenantID
	event.ActorID = actorID
	event.Resource = resource
	event.Message = message

	event.AddTag("tenancy")
	event.AddTag("critical")

	return l.Log(
		ctx,
		event,
	)
}
