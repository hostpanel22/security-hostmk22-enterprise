package rbac

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
	CreateRole(ctx context.Context, role *models.Role) (*models.Role, error)
	GetRoleByID(ctx context.Context, id string) (*models.Role, error)
	GetRoleByName(ctx context.Context, tenantID, name string) (*models.Role, error)
	ListRoles(ctx context.Context, tenantID string) ([]*models.Role, error)
	UpdateRole(ctx context.Context, role *models.Role) (*models.Role, error)
	DeleteRole(ctx context.Context, id string) error

	CreatePermission(ctx context.Context, permission *models.Permission) (*models.Permission, error)
	GetPermissionByID(ctx context.Context, id string) (*models.Permission, error)
	GetPermissionByName(ctx context.Context, tenantID, name string) (*models.Permission, error)
	ListPermissions(ctx context.Context, tenantID string) ([]*models.Permission, error)
	UpdatePermission(ctx context.Context, permission *models.Permission) (*models.Permission, error)
	DeletePermission(ctx context.Context, id string) error
}

type InMemoryRepository struct {
	mu          sync.RWMutex
	roles       map[string]*models.Role
	permissions map[string]*models.Permission
}

func NewRepository() *InMemoryRepository {
	return &InMemoryRepository{
		roles:       make(map[string]*models.Role),
		permissions: make(map[string]*models.Permission),
	}
}

func (r *InMemoryRepository) CreateRole(ctx context.Context, role *models.Role) (*models.Role, error) {
	if r == nil {
		return nil, errors.New("rbac repository is nil")
	}
	if role == nil {
		return nil, errors.New("role is required")
	}
	if err := role.Validate(); err != nil {
		return nil, err
	}
	role.Normalize()
	if strings.TrimSpace(role.ID) == "" {
		role.ID = "role_" + time.Now().UTC().Format("20060102150405.000000000")
	}

	r.mu.Lock()
	defer r.mu.Unlock()
	clone := role.Clone()
	r.roles[clone.ID] = &clone
	stored := r.roles[clone.ID]
	result := stored.Clone()
	return &result, nil
}

func (r *InMemoryRepository) GetRoleByID(ctx context.Context, id string) (*models.Role, error) {
	if r == nil {
		return nil, errors.New("rbac repository is nil")
	}
	if strings.TrimSpace(id) == "" {
		return nil, errors.New("role id is required")
	}

	r.mu.RLock()
	defer r.mu.RUnlock()
	role, ok := r.roles[id]
	if !ok || role == nil {
		return nil, models.ErrNotFound
	}
	clone := role.Clone()
	return &clone, nil
}

func (r *InMemoryRepository) GetRoleByName(ctx context.Context, tenantID, name string) (*models.Role, error) {
	if r == nil {
		return nil, errors.New("rbac repository is nil")
	}
	if strings.TrimSpace(tenantID) == "" || strings.TrimSpace(name) == "" {
		return nil, errors.New("tenant id and role name are required")
	}

	r.mu.RLock()
	defer r.mu.RUnlock()
	for _, role := range r.roles {
		if role == nil {
			continue
		}
		if role.TenantID == tenantID && strings.EqualFold(role.Name, name) {
			clone := role.Clone()
			return &clone, nil
		}
	}
	return nil, models.ErrNotFound
}

func (r *InMemoryRepository) ListRoles(ctx context.Context, tenantID string) ([]*models.Role, error) {
	if r == nil {
		return nil, errors.New("rbac repository is nil")
	}

	r.mu.RLock()
	defer r.mu.RUnlock()
	items := make([]*models.Role, 0, len(r.roles))
	for _, role := range r.roles {
		if role == nil {
			continue
		}
		if strings.TrimSpace(tenantID) != "" && role.TenantID != tenantID {
			continue
		}
		clone := role.Clone()
		items = append(items, &clone)
	}
	sort.Slice(items, func(i, j int) bool { return items[i].Name < items[j].Name })
	return items, nil
}

func (r *InMemoryRepository) UpdateRole(ctx context.Context, role *models.Role) (*models.Role, error) {
	if r == nil {
		return nil, errors.New("rbac repository is nil")
	}
	if role == nil {
		return nil, errors.New("role is required")
	}
	if err := role.Validate(); err != nil {
		return nil, err
	}
	role.Normalize()

	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.roles[role.ID]; !ok {
		return nil, models.ErrNotFound
	}
	clone := role.Clone()
	r.roles[role.ID] = &clone
	stored := r.roles[role.ID]
	result := stored.Clone()
	return &result, nil
}

func (r *InMemoryRepository) DeleteRole(ctx context.Context, id string) error {
	if r == nil {
		return errors.New("rbac repository is nil")
	}
	if strings.TrimSpace(id) == "" {
		return errors.New("role id is required")
	}

	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.roles[id]; !ok {
		return models.ErrNotFound
	}
	delete(r.roles, id)
	return nil
}

func (r *InMemoryRepository) CreatePermission(ctx context.Context, permission *models.Permission) (*models.Permission, error) {
	if r == nil {
		return nil, errors.New("rbac repository is nil")
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

func (r *InMemoryRepository) GetPermissionByID(ctx context.Context, id string) (*models.Permission, error) {
	if r == nil {
		return nil, errors.New("rbac repository is nil")
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

func (r *InMemoryRepository) GetPermissionByName(ctx context.Context, tenantID, name string) (*models.Permission, error) {
	if r == nil {
		return nil, errors.New("rbac repository is nil")
	}
	if strings.TrimSpace(tenantID) == "" || strings.TrimSpace(name) == "" {
		return nil, errors.New("tenant id and permission name are required")
	}

	r.mu.RLock()
	defer r.mu.RUnlock()
	for _, permission := range r.permissions {
		if permission == nil {
			continue
		}
		if permission.TenantID == tenantID && strings.EqualFold(permission.Name, name) {
			clone := permission.Clone()
			return &clone, nil
		}
	}
	return nil, models.ErrNotFound
}

func (r *InMemoryRepository) ListPermissions(ctx context.Context, tenantID string) ([]*models.Permission, error) {
	if r == nil {
		return nil, errors.New("rbac repository is nil")
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
	sort.Slice(items, func(i, j int) bool { return items[i].Name < items[j].Name })
	return items, nil
}

func (r *InMemoryRepository) UpdatePermission(ctx context.Context, permission *models.Permission) (*models.Permission, error) {
	if r == nil {
		return nil, errors.New("rbac repository is nil")
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

func (r *InMemoryRepository) DeletePermission(ctx context.Context, id string) error {
	if r == nil {
		return errors.New("rbac repository is nil")
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
