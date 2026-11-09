package persistence

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"strings"
	"time"

	"github.com/azezaz/HostMK-Security-Framework/internal/crypto"
	"github.com/azezaz/HostMK-Security-Framework/internal/models"
)

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

func (s *Store) ListSecrets(ctx context.Context, tenantID string) ([]*models.Secret, error) {
	rows, err := s.db.QueryContext(ctx, `
		SELECT id, tenant_id, name, type, masked_value, description, revoked_at, created_at, updated_at
		FROM secrets WHERE tenant_id = $1 ORDER BY name
	`, tenantID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []*models.Secret{}
	for rows.Next() {
		var secret models.Secret
		if err := rows.Scan(&secret.ID, &secret.TenantID, &secret.Name, &secret.Type, &secret.MaskedValue, &secret.Description, &secret.RevokedAt, &secret.CreatedAt, &secret.UpdatedAt); err != nil {
			return nil, err
		}
		items = append(items, &secret)
	}
	return items, rows.Err()
}

func (s *Store) RotateSecret(ctx context.Context, aead *crypto.AEAD, id, value, actorID string) (*models.Secret, error) {
	secret, err := s.getSecretRow(ctx, id)
	if err != nil {
		return nil, wrapNotFound(err)
	}
	cipher, err := aead.Encrypt([]byte(value))
	if err != nil {
		return nil, err
	}
	secret.Value = value
	secret.Normalize()
	secret.UpdatedAt = time.Now().UTC()
	if _, err := s.db.ExecContext(ctx, `UPDATE secrets SET ciphertext=$2, masked_value=$3, updated_at=$4 WHERE id=$1`, id, cipher, secret.MaskedValue, secret.UpdatedAt); err != nil {
		return nil, err
	}
	if _, err := s.db.ExecContext(ctx, `
		INSERT INTO secret_versions (id, secret_id, ciphertext, masked_value, created_at, created_by)
		VALUES ($1,$2,$3,$4,$5,$6)
	`, newID(), id, cipher, secret.MaskedValue, time.Now().UTC(), nullString(actorID)); err != nil {
		return nil, err
	}
	secret.Value = ""
	return secret, nil
}

func (s *Store) RevokeSecret(ctx context.Context, id string) error {
	now := time.Now().UTC()
	res, err := s.db.ExecContext(ctx, `UPDATE secrets SET revoked_at=$2, updated_at=$2 WHERE id=$1`, id, now)
	if err != nil {
		return err
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return models.ErrNotFound
	}
	return nil
}

func (s *Store) ListSecretVersions(ctx context.Context, secretID string) ([]*models.SecretVersion, error) {
	rows, err := s.db.QueryContext(ctx, `
		SELECT id, secret_id, masked_value, created_at, created_by
		FROM secret_versions WHERE secret_id = $1 ORDER BY created_at DESC
	`, secretID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []*models.SecretVersion{}
	for rows.Next() {
		var item models.SecretVersion
		var createdBy sql.NullString
		if err := rows.Scan(&item.ID, &item.SecretID, &item.MaskedValue, &item.CreatedAt, &createdBy); err != nil {
			return nil, err
		}
		item.CreatedBy = stringValue(createdBy)
		items = append(items, &item)
	}
	return items, rows.Err()
}

func (s *Store) getSecretRow(ctx context.Context, id string) (*models.Secret, error) {
	var secret models.Secret
	err := s.db.QueryRowContext(ctx, `
		SELECT id, tenant_id, name, type, masked_value, description, revoked_at, created_at, updated_at
		FROM secrets WHERE id = $1
	`, id).Scan(&secret.ID, &secret.TenantID, &secret.Name, &secret.Type, &secret.MaskedValue, &secret.Description, &secret.RevokedAt, &secret.CreatedAt, &secret.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &secret, nil
}

func (s *Store) CreateAPIKey(ctx context.Context, key *models.APIKey) (*models.APIKey, error) {
	key.Normalize()
	if key.ID == "" {
		key.ID = newID()
	}
	plain, prefix, hash, err := generateAPIKey()
	if err != nil {
		return nil, err
	}
	key.Prefix = prefix
	key.KeyHash = hash
	key.Plaintext = plain
	key.UpdatedAt = time.Now().UTC()
	if err := key.Validate(); err != nil {
		return nil, err
	}
	_, err = s.db.ExecContext(ctx, `
		INSERT INTO api_keys (id, tenant_id, name, prefix, key_hash, scopes, rate_limit, expires_at, created_at, updated_at)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10)
	`, key.ID, key.TenantID, key.Name, key.Prefix, key.KeyHash, marshalStringSlice(key.Scopes), key.RateLimit, key.ExpiresAt, key.CreatedAt, key.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return key, nil
}

func (s *Store) ListAPIKeys(ctx context.Context, tenantID string) ([]*models.APIKey, error) {
	rows, err := s.db.QueryContext(ctx, `
		SELECT id, tenant_id, name, prefix, scopes, rate_limit, expires_at, revoked_at, last_used_at, created_at, updated_at
		FROM api_keys WHERE tenant_id = $1 ORDER BY created_at DESC
	`, tenantID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []*models.APIKey{}
	for rows.Next() {
		var key models.APIKey
		var scopes []byte
		if err := rows.Scan(&key.ID, &key.TenantID, &key.Name, &key.Prefix, &scopes, &key.RateLimit, &key.ExpiresAt, &key.RevokedAt, &key.LastUsedAt, &key.CreatedAt, &key.UpdatedAt); err != nil {
			return nil, err
		}
		key.Scopes = unmarshalStringSlice(scopes)
		items = append(items, &key)
	}
	return items, rows.Err()
}

func (s *Store) GetAPIKeyByPlain(ctx context.Context, plain string) (*models.APIKey, error) {
	hash := hashAPIKey(plain)
	var key models.APIKey
	var scopes []byte
	err := s.db.QueryRowContext(ctx, `
		SELECT id, tenant_id, name, prefix, key_hash, scopes, rate_limit, expires_at, revoked_at, last_used_at, created_at, updated_at
		FROM api_keys WHERE key_hash = $1
	`, hash).Scan(&key.ID, &key.TenantID, &key.Name, &key.Prefix, &key.KeyHash, &scopes, &key.RateLimit, &key.ExpiresAt, &key.RevokedAt, &key.LastUsedAt, &key.CreatedAt, &key.UpdatedAt)
	if err != nil {
		return nil, wrapNotFound(err)
	}
	key.Scopes = unmarshalStringSlice(scopes)
	return &key, nil
}

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

func (s *Store) RevokeAPIKey(ctx context.Context, id string) error {
	now := time.Now().UTC()
	res, err := s.db.ExecContext(ctx, `UPDATE api_keys SET revoked_at=$2, updated_at=$2 WHERE id=$1`, id, now)
	if err != nil {
		return err
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return models.ErrNotFound
	}
	return nil
}

func (s *Store) TouchAPIKey(ctx context.Context, id string) error {
	_, err := s.db.ExecContext(ctx, `UPDATE api_keys SET last_used_at=$2 WHERE id=$1`, id, time.Now().UTC())
	return err
}

func (s *Store) RotateAPIKey(ctx context.Context, id string) (*models.APIKey, error) {
	plain, prefix, hash, err := generateAPIKey()
	if err != nil {
		return nil, err
	}
	now := time.Now().UTC()
	res, err := s.db.ExecContext(ctx, `UPDATE api_keys SET prefix=$2, key_hash=$3, updated_at=$4, revoked_at=NULL WHERE id=$1`, id, prefix, hash, now)
	if err != nil {
		return nil, err
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return nil, models.ErrNotFound
	}
	return &models.APIKey{ID: id, Prefix: prefix, Plaintext: plain, UpdatedAt: now}, nil
}

func generateAPIKey() (plain, prefix, hash string, err error) {
	buf := make([]byte, 24)
	if _, err = rand.Read(buf); err != nil {
		return "", "", "", err
	}
	plain = "hmk_" + hex.EncodeToString(buf)
	prefix = plain[:8]
	hash = hashAPIKey(plain)
	return plain, prefix, hash, nil
}

func hashAPIKey(plain string) string {
	sum := sha256.Sum256([]byte(strings.TrimSpace(plain)))
	return hex.EncodeToString(sum[:])
}
