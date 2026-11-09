package persistence

import (
	"context"
	"database/sql"
	"strings"
	"time"

	"github.com/azezaz/HostMK-Security-Framework/internal/crypto"
	"github.com/azezaz/HostMK-Security-Framework/internal/models"
)

func (s *Store) GetAPIKeyByID(ctx context.Context, id string) (*models.APIKey, error) {
	var key models.APIKey
	var scopes []byte
	err := s.db.QueryRowContext(ctx, `
		SELECT id, tenant_id, name, prefix, key_hash, scopes, rate_limit, expires_at, revoked_at, last_used_at, created_at, updated_at
		FROM api_keys WHERE id = $1
	`, id).Scan(&key.ID, &key.TenantID, &key.Name, &key.Prefix, &key.KeyHash, &scopes, &key.RateLimit, &key.ExpiresAt, &key.RevokedAt, &key.LastUsedAt, &key.CreatedAt, &key.UpdatedAt)
	if err != nil {
		return nil, wrapNotFound(err)
	}
	key.Scopes = unmarshalStringSlice(scopes)
	return &key, nil
}

func (s *Store) CreateSecret(ctx context.Context, aead *crypto.AEAD, secret *models.Secret, actorID string) (*models.Secret, error) {
	secret.Normalize()
	if secret.ID == "" {
		secret.ID = newID()
	}
	secret.UpdatedAt = time.Now().UTC()
	if err := secret.Validate(); err != nil {
		return nil, err
	}
	cipher, err := aead.Encrypt([]byte(secret.Value))
	if err != nil {
		return nil, err
	}
	_, err = s.db.ExecContext(ctx, `
		INSERT INTO secrets (id, tenant_id, name, type, ciphertext, masked_value, description, created_at, updated_at)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9)
	`, secret.ID, secret.TenantID, secret.Name, secret.Type, cipher, secret.MaskedValue, secret.Description, secret.CreatedAt, secret.UpdatedAt)
	if err != nil {
		return nil, err
	}
	if _, err := s.db.ExecContext(ctx, `
		INSERT INTO secret_versions (id, secret_id, ciphertext, masked_value, created_at, created_by)
		VALUES ($1,$2,$3,$4,$5,$6)
	`, newID(), secret.ID, cipher, secret.MaskedValue, time.Now().UTC(), nullString(actorID)); err != nil {
		return nil, err
	}
	secret.Value = ""
	return secret, nil
}
