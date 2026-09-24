package persistence

import (
	"context"
	"database/sql"
	"time"

	"github.com/azezaz/HostMK-Security-Framework/internal/models"
)

func (s *Store) CreateFraudRule(ctx context.Context, rule *models.FraudRule) (*models.FraudRule, error) {
	rule.Normalize()
	if rule.ID == "" {
		rule.ID = newID()
	}
	rule.UpdatedAt = time.Now().UTC()
	if err := rule.Validate(); err != nil {
		return nil, err
	}
	_, err := s.db.ExecContext(ctx, `
		INSERT INTO fraud_rules (id, tenant_id, name, enabled, field, operator, threshold, window_seconds, action, score, metadata, created_at, updated_at)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13)
	`, rule.ID, rule.TenantID, rule.Name, rule.Enabled, rule.Field, rule.Operator, rule.Threshold, rule.WindowSeconds, rule.Action, rule.Score, marshalJSON(rule.Metadata), rule.CreatedAt, rule.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return s.GetFraudRule(ctx, rule.ID)
}

func (s *Store) GetFraudRule(ctx context.Context, id string) (*models.FraudRule, error) {
	rule, err := scanFraudRule(s.db.QueryRowContext(ctx, fraudRuleSelect+` WHERE id = $1`, id))
	if err != nil {
		return nil, wrapNotFound(err)
	}
	return rule, nil
}

func (s *Store) ListFraudRules(ctx context.Context, tenantID string) ([]*models.FraudRule, error) {
	rows, err := s.db.QueryContext(ctx, fraudRuleSelect+` WHERE tenant_id = $1 ORDER BY name`, tenantID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []*models.FraudRule{}
	for rows.Next() {
		rule, err := scanFraudRule(rows)
		if err != nil {
			return nil, err
		}
		items = append(items, rule)
	}
	return items, rows.Err()
}

func (s *Store) UpdateFraudRule(ctx context.Context, rule *models.FraudRule) (*models.FraudRule, error) {
	rule.Normalize()
	rule.UpdatedAt = time.Now().UTC()
	if err := rule.Validate(); err != nil {
		return nil, err
	}
	res, err := s.db.ExecContext(ctx, `
		UPDATE fraud_rules SET name=$2, enabled=$3, field=$4, operator=$5, threshold=$6, window_seconds=$7, action=$8, score=$9, metadata=$10, updated_at=$11
		WHERE id=$1
	`, rule.ID, rule.Name, rule.Enabled, rule.Field, rule.Operator, rule.Threshold, rule.WindowSeconds, rule.Action, rule.Score, marshalJSON(rule.Metadata), rule.UpdatedAt)
	if err != nil {
		return nil, err
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return nil, models.ErrNotFound
	}
	return s.GetFraudRule(ctx, rule.ID)
}

func (s *Store) DeleteFraudRule(ctx context.Context, id string) error {
	res, err := s.db.ExecContext(ctx, `DELETE FROM fraud_rules WHERE id = $1`, id)
	if err != nil {
		return err
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return models.ErrNotFound
	}
	return nil
}

func (s *Store) CreateFraudEvent(ctx context.Context, event *models.FraudViolation) (*models.FraudViolation, error) {
	event.Normalize()
	if event.ID == "" {
		event.ID = newID()
	}
	if err := event.Validate(); err != nil {
		return nil, err
	}
	_, err := s.db.ExecContext(ctx, `
		INSERT INTO fraud_events (id, tenant_id, user_id, rule_id, signal_id, score, level, reason, status, metadata, created_at, updated_at)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12)
	`, event.ID, event.TenantID, nullString(event.UserID), nullString(event.RuleID), nullString(event.SignalID), event.Score, event.Level, event.Reason, event.Status, marshalJSON(event.Metadata), event.CreatedAt, event.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return event, nil
}

func (s *Store) ListFraudEvents(ctx context.Context, tenantID string, limit int) ([]*models.FraudViolation, error) {
	if limit <= 0 {
		limit = 50
	}
	rows, err := s.db.QueryContext(ctx, `
		SELECT id, tenant_id, user_id, rule_id, signal_id, score, level, reason, status, metadata, created_at, updated_at
		FROM fraud_events WHERE tenant_id = $1 ORDER BY created_at DESC LIMIT $2
	`, tenantID, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []*models.FraudViolation{}
	for rows.Next() {
		event, err := scanFraudEvent(rows)
		if err != nil {
			return nil, err
		}
		items = append(items, event)
	}
	return items, rows.Err()
}

const fraudRuleSelect = `SELECT id, tenant_id, name, enabled, field, operator, threshold, window_seconds, action, score, metadata, created_at, updated_at FROM fraud_rules`

func scanFraudRule(row scanner) (*models.FraudRule, error) {
	var rule models.FraudRule
	var window sql.NullInt64
	var metadata []byte
	if err := row.Scan(&rule.ID, &rule.TenantID, &rule.Name, &rule.Enabled, &rule.Field, &rule.Operator, &rule.Threshold, &window, &rule.Action, &rule.Score, &metadata, &rule.CreatedAt, &rule.UpdatedAt); err != nil {
		return nil, err
	}
	rule.WindowSeconds = intPtr(window)
	rule.Metadata = unmarshalMap(metadata)
	return &rule, nil
}

func scanFraudEvent(row scanner) (*models.FraudViolation, error) {
	var event models.FraudViolation
	var userID, ruleID, signalID sql.NullString
	var metadata []byte
	if err := row.Scan(&event.ID, &event.TenantID, &userID, &ruleID, &signalID, &event.Score, &event.Level, &event.Reason, &event.Status, &metadata, &event.CreatedAt, &event.UpdatedAt); err != nil {
		return nil, err
	}
	event.UserID = stringValue(userID)
	event.RuleID = stringValue(ruleID)
	event.SignalID = stringValue(signalID)
	event.Metadata = unmarshalMap(metadata)
	return &event, nil
}
