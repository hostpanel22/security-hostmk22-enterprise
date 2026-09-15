package audit_test

import (
	"context"
	"testing"

	"github.com/azezaz/HostMK-Security-Framework/core/audit"
)

func TestNewAuditLogger(t *testing.T) {

	store := audit.NewMemoryStore()

	logger := audit.NewAuditLogger(store)

	if logger == nil {
		t.Fatal("logger should not be nil")
	}

	if !logger.IsEnabled() {
		t.Fatal("logger should be enabled")
	}
}

func TestLoggerEnableDisable(t *testing.T) {

	store := audit.NewMemoryStore()

	logger := audit.NewAuditLogger(store)

	logger.Disable()

	if logger.IsEnabled() {
		t.Fatal("logger should be disabled")
	}

	logger.Enable()

	if !logger.IsEnabled() {
		t.Fatal("logger should be enabled")
	}
}

func TestLoggerLog(t *testing.T) {

	store := audit.NewMemoryStore()

	logger := audit.NewAuditLogger(store)

	event := audit.NewEvent(
		audit.CategorySystem,
		"startup",
		audit.SeverityInfo,
		audit.StatusSuccess,
	)

	err := logger.Log(
		context.Background(),
		event,
	)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestLoggerLogNilEvent(t *testing.T) {

	store := audit.NewMemoryStore()

	logger := audit.NewAuditLogger(store)

	err := logger.Log(
		context.Background(),
		nil,
	)

	if err != audit.ErrNilEvent {
		t.Fatalf(
			"expected ErrNilEvent got %v",
			err,
		)
	}
}

func TestLoggerLogSuccess(t *testing.T) {

	store := audit.NewMemoryStore()

	logger := audit.NewAuditLogger(store)

	err := logger.LogSuccess(
		context.Background(),
		audit.CategorySystem,
		"startup",
		"application started",
	)

	if err != nil {
		t.Fatal(err)
	}
}

func TestLoggerLogFailure(t *testing.T) {

	store := audit.NewMemoryStore()

	logger := audit.NewAuditLogger(store)

	err := logger.LogFailure(
		context.Background(),
		audit.CategorySystem,
		"startup",
		"application failed",
	)

	if err != nil {
		t.Fatal(err)
	}
}

func TestLoggerLogSecurityViolation(t *testing.T) {

	store := audit.NewMemoryStore()

	logger := audit.NewAuditLogger(store)

	err := logger.LogSecurityViolation(
		context.Background(),
		audit.CategorySecurity,
		"intrusion",
		"unauthorized access detected",
	)

	if err != nil {
		t.Fatal(err)
	}
}

func TestLoggerLogAuthentication(t *testing.T) {

	store := audit.NewMemoryStore()

	logger := audit.NewAuditLogger(store)

	err := logger.LogAuthentication(
		context.Background(),
		"tenant-1",
		"user-1",
		true,
		"login success",
	)

	if err != nil {
		t.Fatal(err)
	}
}

func TestLoggerLogAuthorization(t *testing.T) {

	store := audit.NewMemoryStore()

	logger := audit.NewAuditLogger(store)

	err := logger.LogAuthorization(
		context.Background(),
		"tenant-1",
		"user-1",
		"orders",
		"read",
		true,
	)

	if err != nil {
		t.Fatal(err)
	}
}

func TestLoggerLogTenantViolation(t *testing.T) {

	store := audit.NewMemoryStore()

	logger := audit.NewAuditLogger(store)

	err := logger.LogTenantViolation(
		context.Background(),
		"tenant-1",
		"user-1",
		"orders",
		"cross tenant access",
	)

	if err != nil {
		t.Fatal(err)
	}
}
