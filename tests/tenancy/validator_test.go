package tenancy_test

import (
	"testing"
	"time"

	"github.com/azezaz/HostMK-Security-Framework/core/tenancy"
)

func TestValidatorValidateActiveTenant(t *testing.T) {
	validator := tenancy.NewValidator()

	tenant := &tenancy.Tenant{
		ID:     "tenant-1",
		Name:   "Tenant One",
		Slug:   "tenant-one",
		Status: tenancy.StatusActive,
	}

	result := validator.Validate(tenant)

	if !result.Valid {
		t.Fatalf("expected tenant to be valid, got: %v", result.Error)
	}
}

func TestValidatorRejectsNilTenant(t *testing.T) {
	validator := tenancy.NewValidator()

	result := validator.Validate(nil)

	if result.Valid {
		t.Fatal("expected validation failure")
	}

	if result.Error != tenancy.ErrTenantNotFound {
		t.Fatalf(
			"expected ErrTenantNotFound, got %v",
			result.Error,
		)
	}
}

func TestValidatorRejectsDeletedTenant(t *testing.T) {
	validator := tenancy.NewValidator()

	tenant := &tenancy.Tenant{
		ID:     "tenant-1",
		Status: tenancy.StatusDeleted,
	}

	result := validator.Validate(tenant)

	if result.Valid {
		t.Fatal("expected deleted tenant to fail validation")
	}

	if result.Error != tenancy.ErrTenantDeleted {
		t.Fatalf(
			"expected ErrTenantDeleted, got %v",
			result.Error,
		)
	}

	if result.Violation == nil {
		t.Fatal("expected security violation")
	}
}

func TestValidatorRejectsDisabledTenant(t *testing.T) {
	validator := tenancy.NewValidator()

	tenant := &tenancy.Tenant{
		ID:     "tenant-1",
		Status: tenancy.StatusDisabled,
	}

	result := validator.Validate(tenant)

	if result.Valid {
		t.Fatal("expected disabled tenant to fail validation")
	}

	if result.Error != tenancy.ErrTenantDisabled {
		t.Fatalf(
			"expected ErrTenantDisabled, got %v",
			result.Error,
		)
	}
}

func TestValidatorRejectsSuspendedTenant(t *testing.T) {
	validator := tenancy.NewValidator()

	tenant := &tenancy.Tenant{
		ID:     "tenant-1",
		Status: tenancy.StatusSuspended,
	}

	result := validator.Validate(tenant)

	if result.Valid {
		t.Fatal("expected suspended tenant to fail validation")
	}

	if result.Error != tenancy.ErrTenantSuspended {
		t.Fatalf(
			"expected ErrTenantSuspended, got %v",
			result.Error,
		)
	}
}

func TestValidatorRejectsExpiredTenant(t *testing.T) {
	validator := tenancy.NewValidator()

	expiredAt := time.Now().UTC().Add(-24 * time.Hour)

	tenant := &tenancy.Tenant{
		ID:        "tenant-1",
		Status:    tenancy.StatusActive,
		ExpiresAt: &expiredAt,
	}

	result := validator.Validate(tenant)

	if result.Valid {
		t.Fatal("expected expired tenant to fail validation")
	}

	if result.Error != tenancy.ErrTenantExpired {
		t.Fatalf(
			"expected ErrTenantExpired, got %v",
			result.Error,
		)
	}
}

func TestValidatorRejectsPendingTenant(t *testing.T) {
	validator := tenancy.NewValidator()

	tenant := &tenancy.Tenant{
		ID:     "tenant-1",
		Status: tenancy.StatusPending,
	}

	result := validator.Validate(tenant)

	if result.Valid {
		t.Fatal("expected pending tenant to fail validation")
	}

	if result.Error != tenancy.ErrTenantPending {
		t.Fatalf(
			"expected ErrTenantPending, got %v",
			result.Error,
		)
	}
}

func TestValidatorAllowsPendingWhenConfigured(t *testing.T) {
	validator := tenancy.NewValidator()
	validator.AllowPending = true

	tenant := &tenancy.Tenant{
		ID:     "tenant-1",
		Status: tenancy.StatusPending,
	}

	result := validator.Validate(tenant)

	if !result.Valid {
		t.Fatalf(
			"expected pending tenant to be allowed, got %v",
			result.Error,
		)
	}
}

func TestValidateContextSuccess(t *testing.T) {
	validator := tenancy.NewValidator()

	ctx := &tenancy.TenantContext{
		TenantID:   "tenant-1",
		ResolvedAt: time.Now().UTC(),
	}

	if err := validator.ValidateContext(ctx); err != nil {
		t.Fatalf(
			"unexpected error: %v",
			err,
		)
	}
}

func TestValidateContextNil(t *testing.T) {
	validator := tenancy.NewValidator()

	err := validator.ValidateContext(nil)

	if err != tenancy.ErrMissingTenantContext {
		t.Fatalf(
			"expected ErrMissingTenantContext, got %v",
			err,
		)
	}
}

func TestValidateOwnershipSuccess(t *testing.T) {
	validator := tenancy.NewValidator()

	ctx := &tenancy.TenantContext{
		TenantID:   "tenant-1",
		ResolvedAt: time.Now().UTC(),
	}

	err := validator.ValidateOwnership(
		ctx,
		"tenant-1",
	)

	if err != nil {
		t.Fatalf(
			"unexpected error: %v",
			err,
		)
	}
}

func TestValidateOwnershipCrossTenant(t *testing.T) {
	validator := tenancy.NewValidator()

	ctx := &tenancy.TenantContext{
		TenantID:   "tenant-1",
		ResolvedAt: time.Now().UTC(),
	}

	err := validator.ValidateOwnership(
		ctx,
		"tenant-2",
	)

	if err != tenancy.ErrCrossTenantAccess {
		t.Fatalf(
			"expected ErrCrossTenantAccess, got %v",
			err,
		)
	}
}

func TestValidateExpirationSuccess(t *testing.T) {
	validator := tenancy.NewValidator()

	expiresAt := time.Now().UTC().Add(30 * 24 * time.Hour)

	tenant := &tenancy.Tenant{
		ID:        "tenant-1",
		ExpiresAt: &expiresAt,
	}

	err := validator.ValidateExpiration(
		tenant,
	)

	if err != nil {
		t.Fatalf(
			"unexpected error: %v",
			err,
		)
	}
}

func TestValidateExpirationExpired(t *testing.T) {
	validator := tenancy.NewValidator()

	expiresAt := time.Now().UTC().Add(-30 * 24 * time.Hour)

	tenant := &tenancy.Tenant{
		ID:        "tenant-1",
		ExpiresAt: &expiresAt,
	}

	err := validator.ValidateExpiration(
		tenant,
	)

	if err != tenancy.ErrTenantExpired {
		t.Fatalf(
			"expected ErrTenantExpired, got %v",
			err,
		)
	}
}
