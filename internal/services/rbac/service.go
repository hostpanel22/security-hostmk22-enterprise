package rbac

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/azezaz/HostMK-Security-Framework/internal/models"
	repo "github.com/azezaz/HostMK-Security-Framework/internal/repositories/rbac"
)

type Service interface {
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

type service struct {
	repository repo.Repository
}

func NewService(repository repo.Repository) *service {
	return &service{repository: repository}
}

func (s *service) CreateRole(ctx context.Context, role *models.Role) (*models.Role, error) {
	if s == nil {
		return nil, errors.New("rbac service is nil")
	}
	if s.repository == nil {
		return nil, errors.New("rbac repository is nil")
	}
	if role == nil {
		return nil, errors.New("role is required")
	}
	role.Normalize()
	if strings.TrimSpace(role.ID) == "" {
		role.ID = "role_" + time.Now().UTC().Format("20060102150405.000000000")
	}
	return s.repository.CreateRole(ctx, role)
}

func (s *service) GetRoleByID(ctx context.Context, id string) (*models.Role, error) {
	if s == nil {
		return nil, errors.New("rbac service is nil")
	}
	if s.repository == nil {
		return nil, errors.New("rbac repository is nil")
	}
	if strings.TrimSpace(id) == "" {
		return nil, errors.New("role id is required")
	}
	return s.repository.GetRoleByID(ctx, id)
}

func (s *service) GetRoleByName(ctx context.Context, tenantID, name string) (*models.Role, error) {
	if s == nil {
		return nil, errors.New("rbac service is nil")
	}
	if s.repository == nil {
		return nil, errors.New("rbac repository is nil")
	}
	if strings.TrimSpace(tenantID) == "" || strings.TrimSpace(name) == "" {
		return nil, errors.New("tenant id and role name are required")
	}
	return s.repository.GetRoleByName(ctx, tenantID, name)
}

func (s *service) ListRoles(ctx context.Context, tenantID string) ([]*models.Role, error) {
	if s == nil {
		return nil, errors.New("rbac service is nil")
	}
	if s.repository == nil {
		return nil, errors.New("rbac repository is nil")
	}
	return s.repository.ListRoles(ctx, tenantID)
}

func (s *service) UpdateRole(ctx context.Context, role *models.Role) (*models.Role, error) {
	if s == nil {
		return nil, errors.New("rbac service is nil")
	}
	if s.repository == nil {
		return nil, errors.New("rbac repository is nil")
	}
	if role == nil {
		return nil, errors.New("role is required")
	}
	role.Normalize()
	return s.repository.UpdateRole(ctx, role)
}

func (s *service) DeleteRole(ctx context.Context, id string) error {
	if s == nil {
		return errors.New("rbac service is nil")
	}
	if s.repository == nil {
		return errors.New("rbac repository is nil")
	}
	if strings.TrimSpace(id) == "" {
		return errors.New("role id is required")
	}
	return s.repository.DeleteRole(ctx, id)
}

func (s *service) CreatePermission(ctx context.Context, permission *models.Permission) (*models.Permission, error) {
	if s == nil {
		return nil, errors.New("rbac service is nil")
	}
	if s.repository == nil {
		return nil, errors.New("rbac repository is nil")
	}
	if permission == nil {
		return nil, errors.New("permission is required")
	}
	permission.Normalize()
	if strings.TrimSpace(permission.ID) == "" {
		permission.ID = "perm_" + time.Now().UTC().Format("20060102150405.000000000")
	}
	return s.repository.CreatePermission(ctx, permission)
}

func (s *service) GetPermissionByID(ctx context.Context, id string) (*models.Permission, error) {
	if s == nil {
		return nil, errors.New("rbac service is nil")
	}
	if s.repository == nil {
		return nil, errors.New("rbac repository is nil")
	}
	if strings.TrimSpace(id) == "" {
		return nil, errors.New("permission id is required")
	}
	return s.repository.GetPermissionByID(ctx, id)
}

func (s *service) GetPermissionByName(ctx context.Context, tenantID, name string) (*models.Permission, error) {
	if s == nil {
		return nil, errors.New("rbac service is nil")
	}
	if s.repository == nil {
		return nil, errors.New("rbac repository is nil")
	}
	if strings.TrimSpace(tenantID) == "" || strings.TrimSpace(name) == "" {
		return nil, errors.New("tenant id and permission name are required")
	}
	return s.repository.GetPermissionByName(ctx, tenantID, name)
}

func (s *service) ListPermissions(ctx context.Context, tenantID string) ([]*models.Permission, error) {
	if s == nil {
		return nil, errors.New("rbac service is nil")
	}
	if s.repository == nil {
		return nil, errors.New("rbac repository is nil")
	}
	return s.repository.ListPermissions(ctx, tenantID)
}

func (s *service) UpdatePermission(ctx context.Context, permission *models.Permission) (*models.Permission, error) {
	if s == nil {
		return nil, errors.New("rbac service is nil")
	}
	if s.repository == nil {
		return nil, errors.New("rbac repository is nil")
	}
	if permission == nil {
		return nil, errors.New("permission is required")
	}
	permission.Normalize()
	return s.repository.UpdatePermission(ctx, permission)
}

func (s *service) DeletePermission(ctx context.Context, id string) error {
	if s == nil {
		return errors.New("rbac service is nil")
	}
	if s.repository == nil {
		return errors.New("rbac repository is nil")
	}
	if strings.TrimSpace(id) == "" {
		return errors.New("permission id is required")
	}
	return s.repository.DeletePermission(ctx, id)
}
