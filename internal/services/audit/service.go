package audit

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/azezaz/HostMK-Security-Framework/internal/models"
	repo "github.com/azezaz/HostMK-Security-Framework/internal/repositories/audit"
)

type Service interface {
	Create(ctx context.Context, event *models.AuditEvent) (*models.AuditEvent, error)
	GetByID(ctx context.Context, id string) (*models.AuditEvent, error)
	List(ctx context.Context, tenantID string, limit, offset int) ([]*models.AuditEvent, error)
	Update(ctx context.Context, event *models.AuditEvent) (*models.AuditEvent, error)
	Delete(ctx context.Context, id string) error
}

type service struct {
	repository repo.Repository
}

func NewService(repository repo.Repository) *service {
	return &service{repository: repository}
}

func (s *service) Create(ctx context.Context, event *models.AuditEvent) (*models.AuditEvent, error) {
	if s == nil {
		return nil, errors.New("audit service is nil")
	}
	if s.repository == nil {
		return nil, errors.New("audit repository is nil")
	}
	if event == nil {
		return nil, errors.New("audit event is required")
	}
	event.Normalize()
	if strings.TrimSpace(event.ID) == "" {
		event.ID = "evt_" + time.Now().UTC().Format("20060102150405.000000000")
	}
	return s.repository.Create(ctx, event)
}

func (s *service) GetByID(ctx context.Context, id string) (*models.AuditEvent, error) {
	if s == nil {
		return nil, errors.New("audit service is nil")
	}
	if s.repository == nil {
		return nil, errors.New("audit repository is nil")
	}
	if strings.TrimSpace(id) == "" {
		return nil, errors.New("audit event id is required")
	}
	return s.repository.GetByID(ctx, id)
}

func (s *service) List(ctx context.Context, tenantID string, limit, offset int) ([]*models.AuditEvent, error) {
	if s == nil {
		return nil, errors.New("audit service is nil")
	}
	if s.repository == nil {
		return nil, errors.New("audit repository is nil")
	}
	return s.repository.List(ctx, tenantID, limit, offset)
}

func (s *service) Update(ctx context.Context, event *models.AuditEvent) (*models.AuditEvent, error) {
	if s == nil {
		return nil, errors.New("audit service is nil")
	}
	if s.repository == nil {
		return nil, errors.New("audit repository is nil")
	}
	if event == nil {
		return nil, errors.New("audit event is required")
	}
	event.Normalize()
	return s.repository.Update(ctx, event)
}

func (s *service) Delete(ctx context.Context, id string) error {
	if s == nil {
		return errors.New("audit service is nil")
	}
	if s.repository == nil {
		return errors.New("audit repository is nil")
	}
	if strings.TrimSpace(id) == "" {
		return errors.New("audit event id is required")
	}
	return s.repository.Delete(ctx, id)
}
