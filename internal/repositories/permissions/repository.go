package permissions

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
	Create(ctx context.Context, permission *models.Permission) (*models.Permission, error)
	GetByID(ctx context.Context, id string) (*models.Permission, error)
	GetByKey(ctx context.Context, tenantID, resource, action string) (*models.Permission, error)
	ListByTenant(ctx context.Context, tenantID string) ([]*models.Permission, error)
	Update(ctx context.Context, permission *models.Permission) (*models.Permission, error)
	Delete(ctx context.Context, id string) error
}

type InMemoryRepository struct {
	mu          sync.RWMutex
	permissions map[string]*models.Permission
}

func NewRepository() *InMemoryRepository {
	return &InMemoryRepository{permissions: make(map[string]*models.Permission)}
}

func (r *InMemoryRepository) Create(ctx context.Context, permission *models.Permission) (*models.Permission, error) {
	if r == nil {
		return nil, errors.New("permission repository is nil")
	}
	if permission == nil {
		return nil, errors.New("permission is required")
	}
	if err := permission.Validate(); err != nil {
		return nil, err
	}
	permission.Normalize()
	if strings.TrimSpace(permission.ID) == "" {
		permission.ID = "perm_" + time.Now().UTC().Format("20060102150405.000000000")
	}

	r.mu.Lock()
	defer r.mu.Unlock()
	clone := permission.Clone()
	r.permissions[clone.ID] = &clone
	stored := r.permissions[clone.ID]
	result := stored.Clone()
	return &result, nil
}

func (r *InMemoryRepository) GetByID(ctx context.Context, id string) (*models.Permission, error) {
	if r == nil {
		return nil, errors.New("permission repository is nil")
	}
	if strings.TrimSpace(id) == "" {
		return nil, errors.New("permission id is required")
	}

	r.mu.RLock()
	defer r.mu.RUnlock()
	permission, ok := r.permissions[id]
	if !ok || permission == nil {
		return nil, models.ErrNotFound
	}
	clone := permission.Clone()
	return &clone, nil
}

func (r *InMemoryRepository) GetByKey(ctx context.Context, tenantID, resource, action string) (*models.Permission, error) {
	if r == nil {
		return nil, errors.New("permission repository is nil")
	}
	if strings.TrimSpace(tenantID) == "" {
		return nil, errors.New("tenant id is required")
	}
	if strings.TrimSpace(resource) == "" || strings.TrimSpace(action) == "" {
		return nil, errors.New("resource and action are required")
	}

	r.mu.RLock()
	defer r.mu.RUnlock()
	for _, permission := range r.permissions {
		if permission == nil {
			continue
		}
		if permission.TenantID == tenantID && permission.Matches(resource, action) {
			clone := permission.Clone()
			return &clone, nil
		}
	}
	return nil, models.ErrNotFound
}

func (r *InMemoryRepository) ListByTenant(ctx context.Context, tenantID string) ([]*models.Permission, error) {
	if r == nil {
		return nil, errors.New("permission repository is nil")
	}

	r.mu.RLock()
	defer r.mu.RUnlock()
	items := make([]*models.Permission, 0, len(r.permissions))
	for _, permission := range r.permissions {
		if permission == nil {
			continue
		}
		if strings.TrimSpace(tenantID) != "" && permission.TenantID != tenantID {
			continue
		}
		clone := permission.Clone()
		items = append(items, &clone)
	}
	sort.Slice(items, func(i, j int) bool {
		if items[i].Resource == items[j].Resource {
			return items[i].Action < items[j].Action
		}
		return items[i].Resource < items[j].Resource
	})
	return items, nil
}

func (r *InMemoryRepository) Update(ctx context.Context, permission *models.Permission) (*models.Permission, error) {
	if r == nil {
		return nil, errors.New("permission repository is nil")
	}
	if permission == nil {
		return nil, errors.New("permission is required")
	}
	if err := permission.Validate(); err != nil {
		return nil, err
	}
	permission.Normalize()

	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.permissions[permission.ID]; !ok {
		return nil, models.ErrNotFound
	}
	clone := permission.Clone()
	r.permissions[permission.ID] = &clone
	stored := r.permissions[permission.ID]
	result := stored.Clone()
	return &result, nil
}

func (r *InMemoryRepository) Delete(ctx context.Context, id string) error {
	if r == nil {
		return errors.New("permission repository is nil")
	}
	if strings.TrimSpace(id) == "" {
		return errors.New("permission id is required")
	}

	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.permissions[id]; !ok {
		return models.ErrNotFound
	}
	delete(r.permissions, id)
	return nil
}
