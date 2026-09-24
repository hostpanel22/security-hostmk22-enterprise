package persistence

import (
	"context"
	"strings"
	"time"

	"github.com/azezaz/HostMK-Security-Framework/internal/models"
)

func (s *Store) CreateRole(ctx context.Context, role *models.Role) (*models.Role, error) {
	role.Normalize()
	if role.ID == "" {
		role.ID = newID()
	}
	role.UpdatedAt = time.Now().UTC()
	if err := role.Validate(); err != nil {
		return nil, err
	}
	_, err := s.db.ExecContext(ctx, `
		INSERT INTO roles (id, tenant_id, name, description, created_at, updated_at)
		VALUES ($1,$2,$3,$4,$5,$6)
	`, role.ID, role.TenantID, role.Name, role.Description, role.CreatedAt, role.UpdatedAt)
	if err != nil {
		return nil, err
	}
	if err := s.SetRolePermissions(ctx, role.ID, role.Permissions); err != nil {
		return nil, err
	}
	return s.GetRole(ctx, role.ID)
}

func (s *Store) GetRole(ctx context.Context, id string) (*models.Role, error) {
	role, err := scanRole(s.db.QueryRowContext(ctx, `SELECT id, tenant_id, name, description, created_at, updated_at FROM roles WHERE id = $1`, id))
	if err != nil {
		return nil, wrapNotFound(err)
	}
	names, err := s.ListRolePermissionNames(ctx, role.ID)
	if err != nil {
		return nil, err
	}
	role.Permissions = names
	return role, nil
}

func (s *Store) ListRoles(ctx context.Context, tenantID string) ([]*models.Role, error) {
	rows, err := s.db.QueryContext(ctx, `SELECT id, tenant_id, name, description, created_at, updated_at FROM roles WHERE tenant_id = $1 ORDER BY name`, tenantID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []*models.Role{}
	for rows.Next() {
		role, err := scanRole(rows)
		if err != nil {
			return nil, err
		}
		names, err := s.ListRolePermissionNames(ctx, role.ID)
		if err != nil {
			return nil, err
		}
		role.Permissions = names
		items = append(items, role)
	}
	return items, rows.Err()
}

func (s *Store) UpdateRole(ctx context.Context, role *models.Role) (*models.Role, error) {
	role.Normalize()
	role.UpdatedAt = time.Now().UTC()
	if err := role.Validate(); err != nil {
		return nil, err
	}
	res, err := s.db.ExecContext(ctx, `UPDATE roles SET name=$2, description=$3, updated_at=$4 WHERE id=$1`, role.ID, role.Name, role.Description, role.UpdatedAt)
	if err != nil {
		return nil, err
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return nil, models.ErrNotFound
	}
	if role.Permissions != nil {
		if err := s.SetRolePermissions(ctx, role.ID, role.Permissions); err != nil {
			return nil, err
		}
	}
	return s.GetRole(ctx, role.ID)
}

func (s *Store) DeleteRole(ctx context.Context, id string) error {
	res, err := s.db.ExecContext(ctx, `DELETE FROM roles WHERE id = $1`, id)
	if err != nil {
		return err
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return models.ErrNotFound
	}
	return nil
}

func (s *Store) CreatePermission(ctx context.Context, permission *models.Permission) (*models.Permission, error) {
	permission.Normalize()
	if permission.ID == "" {
		permission.ID = newID()
	}
	permission.UpdatedAt = time.Now().UTC()
	if err := permission.Validate(); err != nil {
		return nil, err
	}
	_, err := s.db.ExecContext(ctx, `
		INSERT INTO permissions (id, tenant_id, name, resource, action, description, created_at, updated_at)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8)
	`, permission.ID, permission.TenantID, permission.Name, permission.Resource, permission.Action, permission.Description, permission.CreatedAt, permission.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return s.GetPermission(ctx, permission.ID)
}

func (s *Store) GetPermission(ctx context.Context, id string) (*models.Permission, error) {
	permission, err := scanPermission(s.db.QueryRowContext(ctx, permissionSelect+` WHERE id = $1`, id))
	if err != nil {
		return nil, wrapNotFound(err)
	}
	return permission, nil
}

func (s *Store) ListPermissions(ctx context.Context, tenantID string) ([]*models.Permission, error) {
	rows, err := s.db.QueryContext(ctx, permissionSelect+` WHERE tenant_id = $1 ORDER BY name`, tenantID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []*models.Permission{}
	for rows.Next() {
		permission, err := scanPermission(rows)
		if err != nil {
			return nil, err
		}
		items = append(items, permission)
	}
	return items, rows.Err()
}

func (s *Store) UpdatePermission(ctx context.Context, permission *models.Permission) (*models.Permission, error) {
	permission.Normalize()
	permission.UpdatedAt = time.Now().UTC()
	if err := permission.Validate(); err != nil {
		return nil, err
	}
	res, err := s.db.ExecContext(ctx, `
		UPDATE permissions SET name=$2, resource=$3, action=$4, description=$5, updated_at=$6 WHERE id=$1
	`, permission.ID, permission.Name, permission.Resource, permission.Action, permission.Description, permission.UpdatedAt)
	if err != nil {
		return nil, err
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return nil, models.ErrNotFound
	}
	return s.GetPermission(ctx, permission.ID)
}

func (s *Store) DeletePermission(ctx context.Context, id string) error {
	res, err := s.db.ExecContext(ctx, `DELETE FROM permissions WHERE id = $1`, id)
	if err != nil {
		return err
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return models.ErrNotFound
	}
	return nil
}

func (s *Store) SetRolePermissions(ctx context.Context, roleID string, permissionNames []string) error {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer func() { _ = tx.Rollback() }()
	if _, err := tx.ExecContext(ctx, `DELETE FROM role_permissions WHERE role_id = $1`, roleID); err != nil {
		return err
	}
	for _, name := range permissionNames {
		name = strings.TrimSpace(strings.ToLower(name))
		if name == "" {
			continue
		}
		var permissionID string
		err := tx.QueryRowContext(ctx, `
			SELECT p.id FROM permissions p
			JOIN roles r ON r.tenant_id = p.tenant_id
			WHERE r.id = $1 AND lower(p.name) = $2
		`, roleID, name).Scan(&permissionID)
		if err != nil {
			return err
		}
		if _, err := tx.ExecContext(ctx, `INSERT INTO role_permissions (role_id, permission_id) VALUES ($1,$2) ON CONFLICT DO NOTHING`, roleID, permissionID); err != nil {
			return err
		}
	}
	return tx.Commit()
}

func (s *Store) ListRolePermissionNames(ctx context.Context, roleID string) ([]string, error) {
	rows, err := s.db.QueryContext(ctx, `
		SELECT p.name FROM permissions p
		JOIN role_permissions rp ON rp.permission_id = p.id
		WHERE rp.role_id = $1
		ORDER BY p.name
	`, roleID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	names := []string{}
	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err != nil {
			return nil, err
		}
		names = append(names, name)
	}
	return names, rows.Err()
}

const permissionSelect = `SELECT id, tenant_id, name, resource, action, description, created_at, updated_at FROM permissions`

type scanner interface {
	Scan(dest ...any) error
}

func scanRole(row scanner) (*models.Role, error) {
	var role models.Role
	if err := row.Scan(&role.ID, &role.TenantID, &role.Name, &role.Description, &role.CreatedAt, &role.UpdatedAt); err != nil {
		return nil, err
	}
	return &role, nil
}

func scanPermission(row scanner) (*models.Permission, error) {
	var permission models.Permission
	if err := row.Scan(&permission.ID, &permission.TenantID, &permission.Name, &permission.Resource, &permission.Action, &permission.Description, &permission.CreatedAt, &permission.UpdatedAt); err != nil {
		return nil, err
	}
	return &permission, nil
}
