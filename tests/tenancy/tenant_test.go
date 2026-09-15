package tenancy_test

import (
	"testing"
	"time"

	"github.com/azezaz/HostMK-Security-Framework/core/tenancy"
)

func TestTenantIsActive(t *testing.T) {
	tenant := &tenancy.Tenant{
		ID:     "tenant-1",
		Name:   "Tenant One",
		Slug:   "tenant-one",
		Status: tenancy.StatusActive,
	}

	if !tenant.IsActive() {
		t.Fatal("expected tenant to be active")
	}
}

func TestTenantIsPending(t *testing.T) {
	tenant := &tenancy.Tenant{
		Status: tenancy.StatusPending,
	}

	if !tenant.IsPending() {
		t.Fatal("expected tenant to be pending")
	}
}

func TestTenantIsSuspended(t *testing.T) {
	tenant := &tenancy.Tenant{
		Status: tenancy.StatusSuspended,
	}

	if !tenant.IsSuspended() {
		t.Fatal("expected tenant to be suspended")
	}
}

func TestTenantIsDisabled(t *testing.T) {
	tenant := &tenancy.Tenant{
		Status: tenancy.StatusDisabled,
	}

	if !tenant.IsDisabled() {
		t.Fatal("expected tenant to be disabled")
	}
}

func TestTenantIsDeleted(t *testing.T) {
	tenant := &tenancy.Tenant{
		Status: tenancy.StatusDeleted,
	}

	if !tenant.IsDeleted() {
		t.Fatal("expected tenant to be deleted")
	}
}

func TestTenantIsExpiredFromStatus(t *testing.T) {
	tenant := &tenancy.Tenant{
		Status: tenancy.StatusExpired,
	}

	if !tenant.IsExpired() {
		t.Fatal("expected tenant to be expired")
	}
}

func TestTenantIsExpiredFromDate(t *testing.T) {
	expiredAt := time.Now().UTC().Add(-24 * time.Hour)

	tenant := &tenancy.Tenant{
		Status:    tenancy.StatusActive,
		ExpiresAt: &expiredAt,
	}

	if !tenant.IsExpired() {
		t.Fatal("expected tenant to be expired from expiration date")
	}
}

func TestTenantCanAccessPlatform(t *testing.T) {
	tenant := &tenancy.Tenant{
		ID:     "tenant-1",
		Name:   "Tenant One",
		Slug:   "tenant-one",
		Status: tenancy.StatusActive,
	}

	if !tenant.CanAccessPlatform() {
		t.Fatal("expected tenant access to be granted")
	}
}

func TestTenantCannotAccessWhenSuspended(t *testing.T) {
	tenant := &tenancy.Tenant{
		Status: tenancy.StatusSuspended,
	}

	if tenant.CanAccessPlatform() {
		t.Fatal("expected suspended tenant access to be denied")
	}
}

func TestTenantCannotAccessWhenExpired(t *testing.T) {
	expiredAt := time.Now().UTC().Add(-24 * time.Hour)

	tenant := &tenancy.Tenant{
		Status:    tenancy.StatusActive,
		ExpiresAt: &expiredAt,
	}

	if tenant.CanAccessPlatform() {
		t.Fatal("expected expired tenant access to be denied")
	}
}

func TestTenantHasCustomDomain(t *testing.T) {
	tenant := &tenancy.Tenant{
		Domain: "panel.client.com",
	}

	if !tenant.HasCustomDomain() {
		t.Fatal("expected custom domain")
	}
}

func TestTenantSuspend(t *testing.T) {
	tenant := &tenancy.Tenant{
		Status: tenancy.StatusActive,
	}

	tenant.Suspend()

	if tenant.Status != tenancy.StatusSuspended {
		t.Fatal("expected tenant to be suspended")
	}

	if tenant.SuspendedAt == nil {
		t.Fatal("expected suspended_at to be set")
	}
}

func TestTenantActivate(t *testing.T) {
	now := time.Now().UTC()

	tenant := &tenancy.Tenant{
		Status:      tenancy.StatusSuspended,
		SuspendedAt: &now,
	}

	tenant.Activate()

	if tenant.Status != tenancy.StatusActive {
		t.Fatal("expected tenant to be active")
	}

	if tenant.SuspendedAt != nil {
		t.Fatal("expected suspended_at to be cleared")
	}
}

func TestTenantValidateSuccess(t *testing.T) {
	tenant := &tenancy.Tenant{
		ID:     "tenant-1",
		Name:   "Tenant One",
		Slug:   "tenant-one",
		Status: tenancy.StatusActive,
	}

	if err := tenant.Validate(); err != nil {
		t.Fatalf("unexpected validation error: %v", err)
	}
}

func TestTenantValidateFailure(t *testing.T) {
	tenant := &tenancy.Tenant{}

	if err := tenant.Validate(); err == nil {
		t.Fatal("expected validation error")
	}
}

func TestTenantToMap(t *testing.T) {
	tenant := &tenancy.Tenant{
		ID:     "tenant-1",
		Name:   "Tenant One",
		Slug:   "tenant-one",
		Status: tenancy.StatusActive,
	}

	data := tenant.ToMap()

	if data["id"] != "tenant-1" {
		t.Fatal("expected id to match")
	}

	if data["name"] != "Tenant One" {
		t.Fatal("expected name to match")
	}

	if data["slug"] != "tenant-one" {
		t.Fatal("expected slug to match")
	}
}
