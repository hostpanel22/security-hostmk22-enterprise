package persistence

import "context"

// UserHasPermissionForTenant checks permission membership while binding every
// authorization join to the authenticated tenant.
func (s *Store) UserHasPermissionForTenant(ctx context.Context, userID, tenantID, permissionName string) (bool, error) {
	var exists bool
	err := s.db.QueryRowContext(ctx, `
		SELECT EXISTS (
			SELECT 1
			FROM users u
			JOIN user_roles ur ON ur.user_id = u.id
			JOIN roles r ON r.id = ur.role_id AND r.tenant_id = u.tenant_id
			JOIN role_permissions rp ON rp.role_id = r.id
			JOIN permissions p ON p.id = rp.permission_id AND p.tenant_id = u.tenant_id
			WHERE u.id = $1 AND u.tenant_id = $2 AND lower(p.name) = lower($3)
		)
	`, userID, tenantID, permissionName).Scan(&exists)
	return exists, err
}

// UserHasResourceActionForTenant checks a resource/action grant within one tenant.
func (s *Store) UserHasResourceActionForTenant(ctx context.Context, userID, tenantID, resource, action string) (bool, error) {
	var exists bool
	err := s.db.QueryRowContext(ctx, `
		SELECT EXISTS (
			SELECT 1
			FROM users u
			JOIN user_roles ur ON ur.user_id = u.id
			JOIN roles r ON r.id = ur.role_id AND r.tenant_id = u.tenant_id
			JOIN role_permissions rp ON rp.role_id = r.id
			JOIN permissions p ON p.id = rp.permission_id AND p.tenant_id = u.tenant_id
			WHERE u.id = $1 AND u.tenant_id = $2
			  AND lower(p.resource) = lower($3) AND lower(p.action) = lower($4)
		)
	`, userID, tenantID, resource, action).Scan(&exists)
	return exists, err
}
