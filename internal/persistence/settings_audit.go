package persistence

import (
	"context"
	"database/sql"
	"encoding/json"
	"strconv"
	"strings"
	"time"

	"github.com/azezaz/HostMK-Security-Framework/internal/models"
)

func (s *Store) UpsertSetting(ctx context.Context, setting *models.Setting) (*models.Setting, error) {
	setting.Normalize()
	if setting.ID == "" {
		setting.ID = newID()
	}
	setting.UpdatedAt = time.Now().UTC()
	if err := setting.Validate(); err != nil {
		return nil, err
	}

	var tenant any
	if setting.TenantID == "" {
		tenant = nil
	} else {
		tenant = setting.TenantID
	}

	res, err := s.db.ExecContext(ctx, `
		UPDATE settings SET value = $3, updated_at = $4
		WHERE key = $1 AND (
			($2::uuid IS NULL AND tenant_id IS NULL) OR tenant_id = $2
		)
	`, setting.Key, tenant, setting.Value, setting.UpdatedAt)
	if err != nil {
		return nil, err
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		if _, err := s.db.ExecContext(ctx, `
			INSERT INTO settings (id, tenant_id, key, value, updated_at)
			VALUES ($1,$2,$3,$4,$5)
		`, setting.ID, tenant, setting.Key, setting.Value, setting.UpdatedAt); err != nil {
			return nil, err
		}
	}
	return s.GetSetting(ctx, setting.TenantID, setting.Key)
}

func (s *Store) GetSetting(ctx context.Context, tenantID, key string) (*models.Setting, error) {
	var row *sql.Row
	if strings.TrimSpace(tenantID) == "" {
		row = s.db.QueryRowContext(ctx, `SELECT id, tenant_id, key, value, updated_at FROM settings WHERE tenant_id IS NULL AND key = $1`, key)
	} else {
		row = s.db.QueryRowContext(ctx, `SELECT id, tenant_id, key, value, updated_at FROM settings WHERE tenant_id = $1 AND key = $2`, tenantID, key)
	}
	setting, err := scanSetting(row)
	if err != nil {
		return nil, wrapNotFound(err)
	}
	return setting, nil
}

func (s *Store) ListSettings(ctx context.Context, tenantID string) ([]*models.Setting, error) {
	var rows *sql.Rows
	var err error
	if strings.TrimSpace(tenantID) == "" {
		rows, err = s.db.QueryContext(ctx, `SELECT id, tenant_id, key, value, updated_at FROM settings WHERE tenant_id IS NULL ORDER BY key`)
	} else {
		rows, err = s.db.QueryContext(ctx, `SELECT id, tenant_id, key, value, updated_at FROM settings WHERE tenant_id = $1 OR tenant_id IS NULL ORDER BY key`, tenantID)
	}
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []*models.Setting{}
	for rows.Next() {
		setting, err := scanSetting(rows)
		if err != nil {
			return nil, err
		}
		items = append(items, setting)
	}
	return items, rows.Err()
}

func (s *Store) IntSetting(ctx context.Context, tenantID, key string, fallback int) int {
	setting, err := s.GetSetting(ctx, tenantID, key)
	if err != nil && tenantID != "" {
		setting, err = s.GetSetting(ctx, "", key)
	}
	if err != nil || setting == nil {
		return fallback
	}
	var n int
	if err := json.Unmarshal(setting.Value, &n); err == nil {
		return n
	}
	var asString string
	if err := json.Unmarshal(setting.Value, &asString); err == nil {
		if parsed, err := strconv.Atoi(asString); err == nil {
			return parsed
		}
	}
	return fallback
}

func scanSetting(row scanner) (*models.Setting, error) {
	var setting models.Setting
	var tenantID sql.NullString
	if err := row.Scan(&setting.ID, &tenantID, &setting.Key, &setting.Value, &setting.UpdatedAt); err != nil {
		return nil, err
	}
	setting.TenantID = stringValue(tenantID)
	return &setting, nil
}

func (s *Store) CreateAuditEvent(ctx context.Context, event *models.AuditEvent) (*models.AuditEvent, error) {
	event.Normalize()
	if event.ID == "" {
		event.ID = newID()
	}
	if err := event.Validate(); err != nil {
		return nil, err
	}
	_, err := s.db.ExecContext(ctx, `
		INSERT INTO audit_events (id, tenant_id, user_id, action, resource, severity, status, ip_address, user_agent, message, metadata, created_at, updated_at)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13)
	`, event.ID, nullString(event.TenantID), nullString(event.UserID), event.Action, event.Resource, event.Severity, event.Status, nullString(event.IPAddress), nullString(event.UserAgent), nullString(event.Message), marshalJSON(event.Metadata), event.CreatedAt, event.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return event, nil
}

func (s *Store) ListAuditEvents(ctx context.Context, tenantID string, limit int) ([]*models.AuditEvent, error) {
	if limit <= 0 || limit > 500 {
		limit = 50
	}
	rows, err := s.db.QueryContext(ctx, `
		SELECT id, tenant_id, user_id, action, resource, severity, status, ip_address, user_agent, message, metadata, created_at, updated_at
		FROM audit_events
		WHERE tenant_id = $1
		ORDER BY created_at DESC
		LIMIT $2
	`, tenantID, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []*models.AuditEvent{}
	for rows.Next() {
		event, err := scanAudit(rows)
		if err != nil {
			return nil, err
		}
		items = append(items, event)
	}
	return items, rows.Err()
}

func (s *Store) PurgeAuditOlderThan(ctx context.Context, tenantID string, olderThan time.Duration) (int64, error) {
	cutoff := time.Now().UTC().Add(-olderThan)
	res, err := s.db.ExecContext(ctx, `DELETE FROM audit_events WHERE tenant_id = $1 AND created_at < $2`, tenantID, cutoff)
	if err != nil {
		return 0, err
	}
	return res.RowsAffected()
}

func scanAudit(row scanner) (*models.AuditEvent, error) {
	var event models.AuditEvent
	var tenantID, userID, ip, ua, message sql.NullString
	var metadata []byte
	if err := row.Scan(&event.ID, &tenantID, &userID, &event.Action, &event.Resource, &event.Severity, &event.Status, &ip, &ua, &message, &metadata, &event.CreatedAt, &event.UpdatedAt); err != nil {
		return nil, err
	}
	event.TenantID = stringValue(tenantID)
	event.UserID = stringValue(userID)
	event.IPAddress = stringValue(ip)
	event.UserAgent = stringValue(ua)
	event.Message = stringValue(message)
	event.Metadata = unmarshalMap(metadata)
	return &event, nil
}
