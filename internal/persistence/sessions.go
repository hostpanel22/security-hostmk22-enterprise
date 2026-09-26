package persistence

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"time"

	"github.com/azezaz/HostMK-Security-Framework/internal/models"
)

func (s *Store) CreateSession(ctx context.Context, session *models.Session) error {
	if session.ID == "" {
		session.ID = newID()
	}
	if session.CreatedAt.IsZero() {
		session.CreatedAt = time.Now().UTC()
	}
	_, err := s.db.ExecContext(ctx, `
		INSERT INTO sessions (id, tenant_id, user_id, user_agent, ip_address, expires_at, created_at)
		VALUES ($1,$2,$3,$4,$5,$6,$7)
	`, session.ID, session.TenantID, session.UserID, nullString(session.UserAgent), nullString(session.IPAddress), session.ExpiresAt, session.CreatedAt)
	return err
}

func (s *Store) CreateRefreshToken(ctx context.Context, token *models.RefreshToken) error {
	if token.ID == "" {
		token.ID = newID()
	}
	if token.CreatedAt.IsZero() {
		token.CreatedAt = time.Now().UTC()
	}
	_, err := s.db.ExecContext(ctx, `
		INSERT INTO refresh_tokens (id, tenant_id, user_id, session_id, token_hash, expires_at, created_at)
		VALUES ($1,$2,$3,$4,$5,$6,$7)
	`, token.ID, token.TenantID, token.UserID, nullString(token.SessionID), token.TokenHash, token.ExpiresAt, token.CreatedAt)
	return err
}

func (s *Store) GetRefreshToken(ctx context.Context, tokenHash string) (*models.RefreshToken, error) {
	var token models.RefreshToken
	err := s.db.QueryRowContext(ctx, `
		SELECT id, tenant_id, user_id, session_id, token_hash, expires_at, revoked_at, created_at
		FROM refresh_tokens WHERE token_hash = $1
	`, tokenHash).Scan(&token.ID, &token.TenantID, &token.UserID, &token.SessionID, &token.TokenHash, &token.ExpiresAt, &token.RevokedAt, &token.CreatedAt)
	if err != nil {
		return nil, wrapNotFound(err)
	}
	return &token, nil
}

func (s *Store) ListSessions(ctx context.Context, tenantID string) ([]*models.Session, error) {
	rows, err := s.db.QueryContext(ctx, `
		SELECT id, tenant_id, user_id, user_agent, ip_address, expires_at, revoked_at, created_at
		FROM sessions WHERE tenant_id = $1 ORDER BY created_at DESC
	`, tenantID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items := []*models.Session{}
	for rows.Next() {
		var session models.Session
		var userAgent, ipAddress sql.NullString
		var revokedAt sql.NullTime
		if err := rows.Scan(&session.ID, &session.TenantID, &session.UserID, &userAgent, &ipAddress, &session.ExpiresAt, &revokedAt, &session.CreatedAt); err != nil {
			return nil, err
		}
		session.UserAgent = stringValue(userAgent)
		session.IPAddress = stringValue(ipAddress)
		session.RevokedAt = timePtr(revokedAt)
		items = append(items, &session)
	}
	return items, rows.Err()
}

func (s *Store) RevokeSession(ctx context.Context, tenantID, sessionID string) error {
	now := time.Now().UTC()
	res, err := s.db.ExecContext(ctx, `
		UPDATE sessions
		SET revoked_at = $3
		WHERE tenant_id = $1 AND id = $2 AND revoked_at IS NULL
	`, tenantID, sessionID, now)
	if err != nil {
		return err
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return models.ErrNotFound
	}
	_, err = s.db.ExecContext(ctx, `
		UPDATE refresh_tokens
		SET revoked_at = $3
		WHERE tenant_id = $1 AND session_id = $2 AND revoked_at IS NULL
	`, tenantID, sessionID, now)
	return err
}

func (s *Store) RevokeRefreshToken(ctx context.Context, id string) error {
	_, err := s.db.ExecContext(ctx, `UPDATE refresh_tokens SET revoked_at = NOW() WHERE id = $1`, id)
	return err
}

func (s *Store) RevokeUserSessions(ctx context.Context, userID string) error {
	if _, err := s.db.ExecContext(ctx, `UPDATE sessions SET revoked_at = NOW() WHERE user_id = $1 AND revoked_at IS NULL`, userID); err != nil {
		return err
	}
	_, err := s.db.ExecContext(ctx, `UPDATE refresh_tokens SET revoked_at = NOW() WHERE user_id = $1 AND revoked_at IS NULL`, userID)
	return err
}

func RandomToken() (plain, hash string, err error) {
	buf := make([]byte, 32)
	if _, err = rand.Read(buf); err != nil {
		return "", "", err
	}
	plain = hex.EncodeToString(buf)
	sum := sha256.Sum256([]byte(plain))
	return plain, hex.EncodeToString(sum[:]), nil
}

func HashToken(plain string) string {
	sum := sha256.Sum256([]byte(plain))
	return hex.EncodeToString(sum[:])
}
