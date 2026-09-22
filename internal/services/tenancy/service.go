package tenancy

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/azezaz/HostMK-Security-Framework/internal/models"
	repo "github.com/azezaz/HostMK-Security-Framework/internal/repositories/tenancy"
)

type Service interface {
	Create(ctx context.Context, tenant *models.Tenant) (*models.Tenant, error)
	GetByID(ctx context.Context, id string) (*models.Tenant, error)
	GetBySlug(ctx context.Context, slug string) (*models.Tenant, error)
	GetByDomain(ctx context.Context, domain string) (*models.Tenant, error)
	List(ctx context.Context) ([]*models.Tenant, error)
	Update(ctx context.Context, tenant *models.Tenant) (*models.Tenant, error)
	Delete(ctx context.Context, id string) error
}

type service struct {
	repository repo.Repository
}

func NewService(repository repo.Repository) *service {
	return &service{repository: repository}
}

func (s *service) Create(ctx context.Context, tenant *models.Tenant) (*models.Tenant, error) {
	if s == nil {
		return nil, errors.New("tenant service is nil")
	}
	if s.repository == nil {
		return nil, errors.New("tenant repository is nil")
	}
	if tenant == nil {
		return nil, errors.New("tenant is required")
	}
	tenant.Normalize()
	if strings.TrimSpace(tenant.ID) == "" {
		tenant.ID = "tenant_" + time.Now().UTC().Format("20060102150405.000000000")
	}
	return s.repository.Create(ctx, tenant)
}

func (s *service) GetByID(ctx context.Context, id string) (*models.Tenant, error) {
	if s == nil {
		return nil, errors.New("tenant service is nil")
	}
	if s.repository == nil {
		return nil, errors.New("tenant repository is nil")
	}
	if strings.TrimSpace(id) == "" {
		return nil, errors.New("tenant id is required")
	}
	return s.repository.GetByID(ctx, id)
}

func (s *service) GetBySlug(ctx context.Context, slug string) (*models.Tenant, error) {
	if s == nil {
		return nil, errors.New("tenant service is nil")
	}
	if s.repository == nil {
		return nil, errors.New("tenant repository is nil")
	}
	if strings.TrimSpace(slug) == "" {
		return nil, errors.New("tenant slug is required")
	}
	return s.repository.GetBySlug(ctx, slug)
}

func (s *service) GetByDomain(ctx context.Context, domain string) (*models.Tenant, error) {
	if s == nil {
		return nil, errors.New("tenant service is nil")
	}
	if s.repository == nil {
		return nil, errors.New("tenant repository is nil")
	}
	if strings.TrimSpace(domain) == "" {
		return nil, errors.New("tenant domain is required")
	}
	return s.repository.GetByDomain(ctx, domain)
}

func (s *service) List(ctx context.Context) ([]*models.Tenant, error) {
	if s == nil {
		return nil, errors.New("tenant service is nil")
	}
	if s.repository == nil {
		return nil, errors.New("tenant repository is nil")
	}
	return s.repository.List(ctx)
}

func (s *service) Update(ctx context.Context, tenant *models.Tenant) (*models.Tenant, error) {
	if s == nil {
		return nil, errors.New("tenant service is nil")
	}
	if s.repository == nil {
		return nil, errors.New("tenant repository is nil")
	}
	if tenant == nil {
		return nil, errors.New("tenant is required")
	}
	tenant.Normalize()
	return s.repository.Update(ctx, tenant)
}

func (s *service) Delete(ctx context.Context, id string) error {
	if s == nil {
		return errors.New("tenant service is nil")
	}
	if s.repository == nil {
		return errors.New("tenant repository is nil")
	}
	if strings.TrimSpace(id) == "" {
		return errors.New("tenant id is required")
	}
	return s.repository.Delete(ctx, id)
}
