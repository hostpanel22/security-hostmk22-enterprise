package tenancy

import (
	"context"
	"errors"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/azezaz/HostMK-Security-Framework/internal/models"
)

type Repository interface {
	Create(ctx context.Context, tenant *models.Tenant) (*models.Tenant, error)
	GetByID(ctx context.Context, id string) (*models.Tenant, error)
	GetBySlug(ctx context.Context, slug string) (*models.Tenant, error)
	GetByDomain(ctx context.Context, domain string) (*models.Tenant, error)
	List(ctx context.Context) ([]*models.Tenant, error)
	Update(ctx context.Context, tenant *models.Tenant) (*models.Tenant, error)
	Delete(ctx context.Context, id string) error
}

type InMemoryRepository struct {
	mu       sync.RWMutex
	tenants  map[string]*models.Tenant
	bySlug   map[string]string
	byDomain map[string]string
}

func NewRepository() *InMemoryRepository {
	return &InMemoryRepository{
		tenants:  make(map[string]*models.Tenant),
		bySlug:   make(map[string]string),
		byDomain: make(map[string]string),
	}
}

func (r *InMemoryRepository) Create(ctx context.Context, tenant *models.Tenant) (*models.Tenant, error) {
	if r == nil {
		return nil, errors.New("tenant repository is nil")
	}
	if tenant == nil {
		return nil, errors.New("tenant is required")
	}
	if err := tenant.Validate(); err != nil {
		return nil, err
	}
	tenant.Normalize()
	if strings.TrimSpace(tenant.ID) == "" {
		tenant.ID = "tenant_" + time.Now().UTC().Format("20060102150405.000000000")
	}

	r.mu.Lock()
	defer r.mu.Unlock()
	clone := tenant.Clone()
	r.tenants[clone.ID] = &clone
	if strings.TrimSpace(clone.Slug) != "" {
		r.bySlug[strings.ToLower(clone.Slug)] = clone.ID
	}
	if strings.TrimSpace(clone.Domain) != "" {
		r.byDomain[strings.ToLower(clone.Domain)] = clone.ID
	}
	stored := r.tenants[clone.ID]
	result := stored.Clone()
	return &result, nil
}

func (r *InMemoryRepository) GetByID(ctx context.Context, id string) (*models.Tenant, error) {
	if r == nil {
		return nil, errors.New("tenant repository is nil")
	}
	if strings.TrimSpace(id) == "" {
		return nil, errors.New("tenant id is required")
	}

	r.mu.RLock()
	defer r.mu.RUnlock()
	tenant, ok := r.tenants[id]
	if !ok || tenant == nil {
		return nil, models.ErrNotFound
	}
	clone := tenant.Clone()
	return &clone, nil
}

func (r *InMemoryRepository) GetBySlug(ctx context.Context, slug string) (*models.Tenant, error) {
	if r == nil {
		return nil, errors.New("tenant repository is nil")
	}
	if strings.TrimSpace(slug) == "" {
		return nil, errors.New("tenant slug is required")
	}
	normalized := strings.ToLower(strings.TrimSpace(slug))

	r.mu.RLock()
	defer r.mu.RUnlock()
	id, ok := r.bySlug[normalized]
	if !ok {
		return nil, models.ErrNotFound
	}
	tenant, exists := r.tenants[id]
	if !exists || tenant == nil {
		return nil, models.ErrNotFound
	}
	clone := tenant.Clone()
	return &clone, nil
}

func (r *InMemoryRepository) GetByDomain(ctx context.Context, domain string) (*models.Tenant, error) {
	if r == nil {
		return nil, errors.New("tenant repository is nil")
	}
	if strings.TrimSpace(domain) == "" {
		return nil, errors.New("tenant domain is required")
	}
	normalized := strings.ToLower(strings.TrimSpace(domain))

	r.mu.RLock()
	defer r.mu.RUnlock()
	id, ok := r.byDomain[normalized]
	if !ok {
		return nil, models.ErrNotFound
	}
	tenant, exists := r.tenants[id]
	if !exists || tenant == nil {
		return nil, models.ErrNotFound
	}
	clone := tenant.Clone()
	return &clone, nil
}

func (r *InMemoryRepository) List(ctx context.Context) ([]*models.Tenant, error) {
	if r == nil {
		return nil, errors.New("tenant repository is nil")
	}

	r.mu.RLock()
	defer r.mu.RUnlock()
	items := make([]*models.Tenant, 0, len(r.tenants))
	for _, tenant := range r.tenants {
		if tenant == nil {
			continue
		}
		clone := tenant.Clone()
		items = append(items, &clone)
	}
	sort.Slice(items, func(i, j int) bool { return items[i].Name < items[j].Name })
	return items, nil
}

func (r *InMemoryRepository) Update(ctx context.Context, tenant *models.Tenant) (*models.Tenant, error) {
	if r == nil {
		return nil, errors.New("tenant repository is nil")
	}
	if tenant == nil {
		return nil, errors.New("tenant is required")
	}
	if err := tenant.Validate(); err != nil {
		return nil, err
	}
	tenant.Normalize()

	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.tenants[tenant.ID]; !ok {
		return nil, models.ErrNotFound
	}
	clone := tenant.Clone()
	r.tenants[tenant.ID] = &clone
	if strings.TrimSpace(clone.Slug) != "" {
		r.bySlug[strings.ToLower(clone.Slug)] = clone.ID
	}
	if strings.TrimSpace(clone.Domain) != "" {
		r.byDomain[strings.ToLower(clone.Domain)] = clone.ID
	}
	stored := r.tenants[tenant.ID]
	result := stored.Clone()
	return &result, nil
}

func (r *InMemoryRepository) Delete(ctx context.Context, id string) error {
	if r == nil {
		return errors.New("tenant repository is nil")
	}
	if strings.TrimSpace(id) == "" {
		return errors.New("tenant id is required")
	}

	r.mu.Lock()
	defer r.mu.Unlock()
	tenant, ok := r.tenants[id]
	if !ok || tenant == nil {
		return models.ErrNotFound
	}
	delete(r.tenants, id)
	if strings.TrimSpace(tenant.Slug) != "" {
		delete(r.bySlug, strings.ToLower(tenant.Slug))
	}
	if strings.TrimSpace(tenant.Domain) != "" {
		delete(r.byDomain, strings.ToLower(tenant.Domain))
	}
	return nil
}
