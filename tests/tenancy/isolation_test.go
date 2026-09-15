package tenancy_test

import (
	"testing"
	"time"

	"github.com/azezaz/HostMK-Security-Framework/core/tenancy"
)

func TestNewIsolationEngine(t *testing.T) {
	engine := tenancy.NewIsolationEngine()

	if engine == nil {
		t.Fatal("expected isolation engine instance")
	}

	if !engine.IsStrictMode() {
		t.Fatal("strict mode should be enabled by default")
	}
}

func TestEnableStrictMode(t *testing.T) {
	engine := tenancy.NewIsolationEngine()

	engine.DisableStrictMode()
	engine.EnableStrictMode()

	if !engine.IsStrictMode() {
		t.Fatal("strict mode should be enabled")
	}
}

func TestDisableStrictMode(t *testing.T) {
	engine := tenancy.NewIsolationEngine()

	engine.DisableStrictMode()

	if engine.IsStrictMode() {
		t.Fatal("strict mode should be disabled")
	}
}

func TestIsolationValidateOwnershipSuccess(t *testing.T) {
	engine := tenancy.NewIsolationEngine()

	ctx := &tenancy.TenantContext{
		TenantID:  "tenant-1",
		ActorID:   "user-1",
		RequestID: "req-1",
	}

	err := engine.ValidateOwnership(
		ctx,
		"tenant-1",
		"wallet",
		"read",
	)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestIsolationValidateOwnershipCrossTenant(t *testing.T) {
	engine := tenancy.NewIsolationEngine()

	ctx := &tenancy.TenantContext{
		TenantID:  "tenant-1",
		ActorID:   "user-1",
		RequestID: "req-1",
	}

	err := engine.ValidateOwnership(
		ctx,
		"tenant-2",
		"wallet",
		"read",
	)

	if err != tenancy.ErrCrossTenantAccess {
		t.Fatalf(
			"expected ErrCrossTenantAccess, got %v",
			err,
		)
	}
}

func TestValidateOwnershipNilContext(t *testing.T) {
	engine := tenancy.NewIsolationEngine()

	err := engine.ValidateOwnership(
		nil,
		"tenant-1",
		"wallet",
		"read",
	)

	if err != tenancy.ErrMissingTenantContext {
		t.Fatalf(
			"expected ErrMissingTenantContext, got %v",
			err,
		)
	}
}

func TestValidateTenantMatchSuccess(t *testing.T) {
	engine := tenancy.NewIsolationEngine()

	err := engine.ValidateTenantMatch(
		"tenant-1",
		"tenant-1",
	)

	if err != nil {
		t.Fatalf(
			"unexpected error: %v",
			err,
		)
	}
}

func TestValidateTenantMatchFailure(t *testing.T) {
	engine := tenancy.NewIsolationEngine()

	err := engine.ValidateTenantMatch(
		"tenant-1",
		"tenant-2",
	)

	if err != tenancy.ErrTenantMismatch {
		t.Fatalf(
			"expected ErrTenantMismatch, got %v",
			err,
		)
	}
}

func TestIsCrossTenantAccess(t *testing.T) {
	engine := tenancy.NewIsolationEngine()

	ctx := &tenancy.TenantContext{
		TenantID: "tenant-1",
	}

	if !engine.IsCrossTenantAccess(
		ctx,
		"tenant-2",
	) {
		t.Fatal("expected cross tenant detection")
	}
}

func TestIsNotCrossTenantAccess(t *testing.T) {
	engine := tenancy.NewIsolationEngine()

	ctx := &tenancy.TenantContext{
		TenantID: "tenant-1",
	}

	if engine.IsCrossTenantAccess(
		ctx,
		"tenant-1",
	) {
		t.Fatal("expected same tenant access")
	}
}

func TestIsolationValidateContextSuccess(t *testing.T) {
	engine := tenancy.NewIsolationEngine()

	ctx := &tenancy.TenantContext{
		TenantID:   "tenant-1",
		ResolvedAt: time.Now().UTC(),
	}

	err := engine.ValidateContext(
		ctx,
	)

	if err != nil {
		t.Fatalf(
			"unexpected error: %v",
			err,
		)
	}
}

func TestValidateContextFailure(t *testing.T) {
	engine := tenancy.NewIsolationEngine()

	err := engine.ValidateContext(
		nil,
	)

	if err != tenancy.ErrMissingTenantContext {
		t.Fatalf(
			"expected ErrMissingTenantContext, got %v",
			err,
		)
	}
}

func TestValidateQuerySuccess(t *testing.T) {
	engine := tenancy.NewIsolationEngine()

	ctx := &tenancy.TenantContext{
		TenantID:   "tenant-1",
		ResolvedAt: time.Now().UTC(),
	}

	err := engine.ValidateQuery(
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

func TestValidateQueryCrossTenant(t *testing.T) {
	engine := tenancy.NewIsolationEngine()

	ctx := &tenancy.TenantContext{
		TenantID:   "tenant-1",
		ResolvedAt: time.Now().UTC(),
	}

	err := engine.ValidateQuery(
		ctx,
		"tenant-2",
	)

	if err != tenancy.ErrCrossTenantQuery {
		t.Fatalf(
			"expected ErrCrossTenantQuery, got %v",
			err,
		)
	}
}

func TestValidateResourceAccessSuccess(t *testing.T) {
	engine := tenancy.NewIsolationEngine()

	ctx := &tenancy.TenantContext{
		TenantID:   "tenant-1",
		ResolvedAt: time.Now().UTC(),
	}

	err := engine.ValidateResourceAccess(
		ctx,
		"tenant-1",
		"service",
	)

	if err != nil {
		t.Fatalf(
			"unexpected error: %v",
			err,
		)
	}
}

func TestValidateResourceAccessFailure(t *testing.T) {
	engine := tenancy.NewIsolationEngine()

	ctx := &tenancy.TenantContext{
		TenantID:   "tenant-1",
		ResolvedAt: time.Now().UTC(),
	}

	err := engine.ValidateResourceAccess(
		ctx,
		"tenant-2",
		"service",
	)

	if err != tenancy.ErrCrossTenantResourceAccess {
		t.Fatalf(
			"expected ErrCrossTenantResourceAccess, got %v",
			err,
		)
	}
}

func TestBuildCrossTenantViolation(t *testing.T) {
	engine := tenancy.NewIsolationEngine()

	ctx := &tenancy.TenantContext{
		TenantID:  "tenant-1",
		ActorID:   "admin-1",
		RequestID: "req-1",
		IPAddress: "127.0.0.1",
		UserAgent: "Go-Test",
	}

	violation := engine.BuildCrossTenantViolation(
		ctx,
		"tenant-2",
		"wallet",
		"update",
	)

	if violation == nil {
		t.Fatal("expected violation")
	}

	if violation.TenantID != "tenant-1" {
		t.Fatal("invalid tenant id")
	}

	if violation.Resource != "wallet" {
		t.Fatal("invalid resource")
	}

	if violation.Action != "update" {
		t.Fatal("invalid action")
	}
}
