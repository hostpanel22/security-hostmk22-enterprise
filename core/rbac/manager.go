package rbac

import (
	"context"
)

type Manager struct {
	store Store
}

func NewManager(
	store Store,
) (*Manager, error) {

	if store == nil {
		return nil, ErrStoreUnavailable
	}

	return &Manager{
		store: store,
	}, nil
}

func (m *Manager) CreateRole(
	ctx context.Context,
	role *Role,
) error {

	if m == nil {
		return ErrRBACUnavailable
	}

	if role == nil {
		return ErrInvalidRole
	}

	if err := role.Validate(); err != nil {
		return err
	}

	_, err := m.store.GetRole(
		ctx,
		role.Name,
	)

	if err == nil {
		return ErrRoleAlreadyExists
	}

	return m.store.SaveRole(
		ctx,
		role,
	)
}

func (m *Manager) CreatePermission(
	ctx context.Context,
	permission *Permission,
) error {

	if m == nil {
		return ErrRBACUnavailable
	}

	if permission == nil {
		return ErrInvalidPermission
	}

	if err := permission.Validate(); err != nil {
		return err
	}

	_, err := m.store.GetPermission(
		ctx,
		permission.Name,
	)

	if err == nil {
		return ErrPermissionAlreadyExists
	}

	return m.store.SavePermission(
		ctx,
		permission,
	)
}

func (m *Manager) GetRole(
	ctx context.Context,
	name string,
) (*Role, error) {

	if m == nil {
		return nil, ErrRBACUnavailable
	}

	return m.store.GetRole(
		ctx,
		name,
	)
}

func (m *Manager) GetPermission(
	ctx context.Context,
	name string,
) (*Permission, error) {

	if m == nil {
		return nil, ErrRBACUnavailable
	}

	return m.store.GetPermission(
		ctx,
		name,
	)
}

func (m *Manager) DeleteRole(
	ctx context.Context,
	name string,
) error {

	if m == nil {
		return ErrRBACUnavailable
	}

	return m.store.DeleteRole(
		ctx,
		name,
	)
}

func (m *Manager) DeletePermission(
	ctx context.Context,
	name string,
) error {

	if m == nil {
		return ErrRBACUnavailable
	}

	return m.store.DeletePermission(
		ctx,
		name,
	)
}

func (m *Manager) AssignPermissionToRole(
	ctx context.Context,
	roleName string,
	permissionName string,
) error {

	role, err := m.store.GetRole(
		ctx,
		roleName,
	)

	if err != nil {
		return err
	}

	_, err = m.store.GetPermission(
		ctx,
		permissionName,
	)

	if err != nil {
		return err
	}

	if err := role.AddPermission(
		permissionName,
	); err != nil {

		return err
	}

	return m.store.SaveRole(
		ctx,
		role,
	)
}

func (m *Manager) RemovePermissionFromRole(
	ctx context.Context,
	roleName string,
	permissionName string,
) error {

	role, err := m.store.GetRole(
		ctx,
		roleName,
	)

	if err != nil {
		return err
	}

	if err := role.RemovePermission(
		permissionName,
	); err != nil {

		return err
	}

	return m.store.SaveRole(
		ctx,
		role,
	)
}

func (m *Manager) RoleHasPermission(
	ctx context.Context,
	roleName string,
	permissionName string,
) (bool, error) {

	role, err := m.store.GetRole(
		ctx,
		roleName,
	)

	if err != nil {
		return false, err
	}

	return role.HasPermission(
		permissionName,
	), nil
}

func (m *Manager) AssignParentRole(
	ctx context.Context,
	roleName string,
	parentRole string,
) error {

	role, err := m.store.GetRole(
		ctx,
		roleName,
	)

	if err != nil {
		return err
	}

	if _, err = m.store.GetRole(
		ctx,
		parentRole,
	); err != nil {

		return err
	}

	if err := role.AddParentRole(
		parentRole,
	); err != nil {

		return err
	}

	return m.store.SaveRole(
		ctx,
		role,
	)
}

func (m *Manager) RemoveParentRole(
	ctx context.Context,
	roleName string,
	parentRole string,
) error {

	role, err := m.store.GetRole(
		ctx,
		roleName,
	)

	if err != nil {
		return err
	}

	if err := role.RemoveParentRole(
		parentRole,
	); err != nil {

		return err
	}

	return m.store.SaveRole(
		ctx,
		role,
	)
}

func (m *Manager) ListRoles(
	ctx context.Context,
) ([]*Role, error) {

	if m == nil {
		return nil, ErrRBACUnavailable
	}

	return m.store.ListRoles(
		ctx,
	)
}

func (m *Manager) ListPermissions(
	ctx context.Context,
) ([]*Permission, error) {

	if m == nil {
		return nil, ErrRBACUnavailable
	}

	return m.store.ListPermissions(
		ctx,
	)
}

func (m *Manager) Clear(
	ctx context.Context,
) error {

	if m == nil {
		return ErrRBACUnavailable
	}

	return m.store.Clear(
		ctx,
	)
}

func (m *Manager) Store() Store {

	if m == nil {
		return nil
	}

	return m.store
}
