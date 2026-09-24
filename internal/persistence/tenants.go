package persistence

import (
	"context"
	"database/sql"
	"strings"
	"time"

	"github.com/azezaz/HostMK-Security-Framework/internal/models"
)

func (s *Store) CreateTenant(ctx context.Context, tenant *models.Tenant) (*models.Tenant, error) {
	tenant.Normalize()
	if tenant.ID == "" {
		tenant.ID = newID()
	}
	if tenant.CreatedAt.IsZero() {
		tenant.CreatedAt = time.Now().UTC()
	}
	tenant.UpdatedAt = time.Now().UTC()
	if err := tenant.Validate(); err != nil {
		return nil, err
	}

	_, err := s.db.ExecContext(ctx, `
		INSERT INTO tenants (id, name, slug, domain, status, owner_id, plan, timezone, country, quotas, branding, metadata, created_at, updated_at, suspended_at, expires_at)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15,$16)
	`, tenant.ID, tenant.Name, tenant.Slug, nullString(tenant.Domain), tenant.Status, nullString(tenant.OwnerID), tenant.Plan, nullString(tenant.Timezone), nullString(tenant.Country),
		marshalJSON(tenant.Quotas), marshalJSON(tenant.Branding), marshalJSON(tenant.Metadata), tenant.CreatedAt, tenant.UpdatedAt, tenant.SuspendedAt, tenant.ExpiresAt)
	if err != nil {
		return nil, err
	}
	return s.GetTenant(ctx, tenant.ID)
}

func (s *Store) GetTenant(ctx context.Context, id string) (*models.Tenant, error) {
	return s.scanTenant(s.db.QueryRowContext(ctx, tenantSelect+` WHERE id = $1`, id))
}

func (s *Store) GetTenantBySlug(ctx context.Context, slug string) (*models.Tenant, error) {
	return s.scanTenant(s.db.QueryRowContext(ctx, tenantSelect+` WHERE slug = $1`, strings.ToLower(strings.TrimSpace(slug))))
}

func (s *Store) ListTenants(ctx context.Context) ([]*models.Tenant, error) {
	rows, err := s.db.QueryContext(ctx, tenantSelect+` ORDER BY created_at DESC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []*models.Tenant{}
	for rows.Next() {
		tenant, err := scanTenantRows(rows)
		if err != nil {
			return nil, err
		}
		items = append(items, tenant)
	}
	return items, rows.Err()
}

func (s *Store) UpdateTenant(ctx context.Context, tenant *models.Tenant) (*models.Tenant, error) {
	tenant.Normalize()
	tenant.UpdatedAt = time.Now().UTC()
	if err := tenant.Validate(); err != nil {
		return nil, err
	}
	res, err := s.db.ExecContext(ctx, `
		UPDATE tenants SET name=$2, slug=$3, domain=$4, status=$5, owner_id=$6, plan=$7, timezone=$8, country=$9,
			quotas=$10, branding=$11, metadata=$12, updated_at=$13, suspended_at=$14, expires_at=$15
		WHERE id=$1
	`, tenant.ID, tenant.Name, tenant.Slug, nullString(tenant.Domain), tenant.Status, nullString(tenant.OwnerID), tenant.Plan, nullString(tenant.Timezone), nullString(tenant.Country),
		marshalJSON(tenant.Quotas), marshalJSON(tenant.Branding), marshalJSON(tenant.Metadata), tenant.UpdatedAt, tenant.SuspendedAt, tenant.ExpiresAt)
	if err != nil {
		return nil, err
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return nil, models.ErrNotFound
	}
	return s.GetTenant(ctx, tenant.ID)
}

const tenantSelect = `
	SELECT id, name, slug, domain, status, owner_id, plan, timezone, country, quotas, branding, metadata, created_at, updated_at, suspended_at, expires_at
	FROM tenants
`

func (s *Store) scanTenant(row *sql.Row) (*models.Tenant, error) {
	tenant, err := scanTenantRow(row)
	if err != nil {
		return nil, wrapNotFound(err)
	}
	return tenant, nil
}

type tenantScanner interface {
	Scan(dest ...any) error
}

func scanTenantRow(row tenantScanner) (*models.Tenant, error) {
	var tenant models.Tenant
	var domain, ownerID, timezone, country sql.NullString
	var quotas, branding, metadata []byte
	var suspendedAt, expiresAt sql.NullTime
	if err := row.Scan(&tenant.ID, &tenant.Name, &tenant.Slug, &domain, &tenant.Status, &ownerID, &tenant.Plan, &timezone, &country, &quotas, &branding, &metadata, &tenant.CreatedAt, &tenant.UpdatedAt, &suspendedAt, &expiresAt); err != nil {
		return nil, err
	}
	tenant.Domain = stringValue(domain)
	tenant.OwnerID = stringValue(ownerID)
	tenant.Timezone = stringValue(timezone)
	tenant.Country = stringValue(country)
	tenant.Quotas = unmarshalMap(quotas)
	tenant.Branding = unmarshalMap(branding)
	tenant.Metadata = unmarshalMap(metadata)
	tenant.SuspendedAt = timePtr(suspendedAt)
	tenant.ExpiresAt = timePtr(expiresAt)
	return &tenant, nil
}

func scanTenantRows(rows *sql.Rows) (*models.Tenant, error) {
	return scanTenantRow(rows)
}
