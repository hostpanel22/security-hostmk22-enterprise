package persistence

import (
	"context"
	"database/sql"
	"strings"
	"time"

	"github.com/azezaz/HostMK-Security-Framework/internal/models"
)

func (s *Store) CountUsers(ctx context.Context) (int, error) {
	var n int
	if err := s.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM users`).Scan(&n); err != nil {
		return 0, err
	}
	return n, nil
}

func (s *Store) CreateUser(ctx context.Context, user *models.User) (*models.User, error) {
	user.Normalize()
	if user.ID == "" {
		user.ID = newID()
	}
	user.UpdatedAt = time.Now().UTC()
	if err := user.Validate(); err != nil {
		return nil, err
	}
	_, err := s.db.ExecContext(ctx, `
		INSERT INTO users (id, tenant_id, email, password_hash, name, status, email_verified_at, mfa_enabled, mfa_secret_enc, last_login, created_at, updated_at)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12)
	`, user.ID, user.TenantID, user.Email, user.PasswordHash, user.Name, user.Status, user.EmailVerifiedAt, user.MFAEnabled, user.MFASecretEnc, user.LastLogin, user.CreatedAt, user.UpdatedAt)
	if err != nil {
		return nil, err
	}
	if len(user.RoleIDs) > 0 {
		if err := s.SetUserRoles(ctx, user.ID, user.RoleIDs); err != nil {
			return nil, err
		}
	}
	return s.GetUser(ctx, user.ID)
}

func (s *Store) GetUser(ctx context.Context, id string) (*models.User, error) {
	user, err := scanUser(s.db.QueryRowContext(ctx, userSelect+` WHERE id = $1`, id))
	if err != nil {
		return nil, wrapNotFound(err)
	}
	roles, err := s.ListUserRoleIDs(ctx, user.ID)
	if err != nil {
		return nil, err
	}
	user.RoleIDs = roles
	return user, nil
}

func (s *Store) GetUserByEmail(ctx context.Context, tenantID, email string) (*models.User, error) {
	user, err := scanUser(s.db.QueryRowContext(ctx, userSelect+` WHERE tenant_id = $1 AND email = $2`, tenantID, strings.ToLower(strings.TrimSpace(email))))
	if err != nil {
		return nil, wrapNotFound(err)
	}
	roles, err := s.ListUserRoleIDs(ctx, user.ID)
	if err != nil {
		return nil, err
	}
	user.RoleIDs = roles
	return user, nil
}

func (s *Store) ListUsers(ctx context.Context, tenantID string) ([]*models.User, error) {
	rows, err := s.db.QueryContext(ctx, userSelect+` WHERE tenant_id = $1 ORDER BY created_at DESC`, tenantID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []*models.User{}
	for rows.Next() {
		user, err := scanUser(rows)
		if err != nil {
			return nil, err
		}
		roles, err := s.ListUserRoleIDs(ctx, user.ID)
		if err != nil {
			return nil, err
		}
		user.RoleIDs = roles
		items = append(items, user)
	}
	return items, rows.Err()
}

func (s *Store) UpdateUser(ctx context.Context, user *models.User) (*models.User, error) {
	user.Normalize()
	user.UpdatedAt = time.Now().UTC()
	if err := user.Validate(); err != nil {
		return nil, err
	}
	res, err := s.db.ExecContext(ctx, `
		UPDATE users SET email=$2, password_hash=$3, name=$4, status=$5, email_verified_at=$6, mfa_enabled=$7, mfa_secret_enc=$8, last_login=$9, updated_at=$10
		WHERE id=$1
	`, user.ID, user.Email, user.PasswordHash, user.Name, user.Status, user.EmailVerifiedAt, user.MFAEnabled, user.MFASecretEnc, user.LastLogin, user.UpdatedAt)
	if err != nil {
		return nil, err
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return nil, models.ErrNotFound
	}
	if user.RoleIDs != nil {
		if err := s.SetUserRoles(ctx, user.ID, user.RoleIDs); err != nil {
			return nil, err
		}
	}
	return s.GetUser(ctx, user.ID)
}

func (s *Store) DeleteUser(ctx context.Context, id string) error {
	res, err := s.db.ExecContext(ctx, `DELETE FROM users WHERE id = $1`, id)
	if err != nil {
		return err
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return models.ErrNotFound
	}
	return nil
}

func (s *Store) SetUserRoles(ctx context.Context, userID string, roleIDs []string) error {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer func() { _ = tx.Rollback() }()
	if _, err := tx.ExecContext(ctx, `DELETE FROM user_roles WHERE user_id = $1`, userID); err != nil {
		return err
	}
	for _, roleID := range roleIDs {
		roleID = strings.TrimSpace(roleID)
		if roleID == "" {
			continue
		}
		if _, err := tx.ExecContext(ctx, `INSERT INTO user_roles (user_id, role_id) VALUES ($1,$2) ON CONFLICT DO NOTHING`, userID, roleID); err != nil {
			return err
		}
	}
	return tx.Commit()
}

func (s *Store) ListUserRoleIDs(ctx context.Context, userID string) ([]string, error) {
	rows, err := s.db.QueryContext(ctx, `SELECT role_id FROM user_roles WHERE user_id = $1`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	ids := []string{}
	for rows.Next() {
		var id string
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}
		ids = append(ids, id)
	}
	return ids, rows.Err()
}

func (s *Store) UserHasPermission(ctx context.Context, userID, permissionName string) (bool, error) {
	var exists bool
	err := s.db.QueryRowContext(ctx, `
		SELECT EXISTS (
			SELECT 1
			FROM user_roles ur
			JOIN role_permissions rp ON rp.role_id = ur.role_id
			JOIN permissions p ON p.id = rp.permission_id
			WHERE ur.user_id = $1 AND lower(p.name) = lower($2)
		)
	`, userID, permissionName).Scan(&exists)
	return exists, err
}

func (s *Store) UserHasResourceAction(ctx context.Context, userID, resource, action string) (bool, error) {
	var exists bool
	err := s.db.QueryRowContext(ctx, `
		SELECT EXISTS (
			SELECT 1
			FROM user_roles ur
			JOIN role_permissions rp ON rp.role_id = ur.role_id
			JOIN permissions p ON p.id = rp.permission_id
			WHERE ur.user_id = $1 AND lower(p.resource) = lower($2) AND lower(p.action) = lower($3)
		)
	`, userID, resource, action).Scan(&exists)
	return exists, err
}

const userSelect = `
	SELECT id, tenant_id, email, password_hash, name, status, email_verified_at, mfa_enabled, mfa_secret_enc, last_login, created_at, updated_at
	FROM users
`

type userScanner interface {
	Scan(dest ...any) error
}

func scanUser(row userScanner) (*models.User, error) {
	var user models.User
	var verified, lastLogin sql.NullTime
	if err := row.Scan(&user.ID, &user.TenantID, &user.Email, &user.PasswordHash, &user.Name, &user.Status, &verified, &user.MFAEnabled, &user.MFASecretEnc, &lastLogin, &user.CreatedAt, &user.UpdatedAt); err != nil {
		return nil, err
	}
	user.EmailVerifiedAt = timePtr(verified)
	user.LastLogin = timePtr(lastLogin)
	return &user, nil
}
