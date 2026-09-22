package permissions

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/azezaz/HostMK-Security-Framework/internal/models"
	repo "github.com/azezaz/HostMK-Security-Framework/internal/repositories/permissions"
)

type Service interface {
	Create(ctx context.Context, permission *models.Permission) (*models.Permission, error)
	GetByID(ctx context.Context, id string) (*models.Permission, error)
	GetByKey(ctx context.Context, tenantID, resource, action string) (*models.Permission, error)
	ListByTenant(ctx context.Context, tenantID string) ([]*models.Permission, error)
	Update(ctx context.Context, permission *models.Permission) (*models.Permission, error)
	Delete(ctx context.Context, id string) error
}

type service struct {
	repository repo.Repository
}

func NewService(repository repo.Repository) *service {
	return &service{repository: repository}
}

func (s *service) Create(ctx context.Context, permission *models.Permission) (*models.Permission, error) {
	if s == nil {
		return nil, errors.New("permission service is nil")
	}
	if s.repository == nil {
		return nil, errors.New("permission repository is nil")
	}
	if permission == nil {
		return nil, errors.New("permission is required")
	}
	permission.Normalize()
	if strings.TrimSpace(permission.ID) == "" {
		permission.ID = "perm_" + time.Now().UTC().Format("20060102150405.000000000")
	}
	return s.repository.Create(ctx, permission)
}

func (s *service) GetByID(ctx context.Context, id string) (*models.Permission, error) {
	if s == nil {
		return nil, errors.New("permission service is nil")
	}
	if s.repository == nil {
		return nil, errors.New("permission repository is nil")
	}
	if strings.TrimSpace(id) == "" {
		return nil, errors.New("permission id is required")
	}
	return s.repository.GetByID(ctx, id)
}

func (s *service) GetByKey(ctx context.Context, tenantID, resource, action string) (*models.Permission, error) {
	if s == nil {
		return nil, errors.New("permission service is nil")
	}
	if s.repository == nil {
		return nil, errors.New("permission repository is nil")
	}
	if strings.TrimSpace(tenantID) == "" || strings.TrimSpace(resource) == "" || strings.TrimSpace(action) == "" {
		return nil, errors.New("tenant id, resource and action are required")
	}
	return s.repository.GetByKey(ctx, tenantID, resource, action)
}

func (s *service) ListByTenant(ctx context.Context, tenantID string) ([]*models.Permission, error) {
	if s == nil {
		return nil, errors.New("permission service is nil")
	}
	if s.repository == nil {
		return nil, errors.New("permission repository is nil")
	}
	return s.repository.ListByTenant(ctx, tenantID)
}

func (s *service) Update(ctx context.Context, permission *models.Permission) (*models.Permission, error) {
	if s == nil {
		return nil, errors.New("permission service is nil")
	}
	if s.repository == nil {
		return nil, errors.New("permission repository is nil")
	}
	if permission == nil {
		return nil, errors.New("permission is required")
	}
	permission.Normalize()
	return s.repository.Update(ctx, permission)
}

func (s *service) Delete(ctx context.Context, id string) error {
	if s == nil {
		return errors.New("permission service is nil")
	}
	if s.repository == nil {
		return errors.New("permission repository is nil")
	}
	if strings.TrimSpace(id) == "" {
		return errors.New("permission id is required")
	}
	return s.repository.Delete(ctx, id)
}
