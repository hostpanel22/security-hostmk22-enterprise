package fraud_test

import (
	"testing"

	"github.com/azezaz/HostMK-Security-Framework/core/fraud"
)

func TestNewSignal(t *testing.T) {
	signal := fraud.NewSignal(
		"tenant-1",
		"user-1",
		"LOGIN",
		75.5,
		"gateway",
	)

	if signal == nil {
		t.Fatal("expected signal")
	}

	if signal.TenantID != "tenant-1" {
		t.Fatalf("expected tenant-1 got %s", signal.TenantID)
	}

	if signal.UserID != "user-1" {
		t.Fatalf("expected user-1 got %s", signal.UserID)
	}

	if signal.Type != "login" {
		t.Fatalf("expected login got %s", signal.Type)
	}

	if signal.Value != 75.5 {
		t.Fatalf("expected 75.5 got %f", signal.Value)
	}

	if signal.Source != "gateway" {
		t.Fatalf("expected gateway got %s", signal.Source)
	}

	if signal.ID == "" {
		t.Fatal("expected generated id")
	}

	if signal.CreatedAt.IsZero() {
		t.Fatal("expected created at")
	}

	if signal.Metadata == nil {
		t.Fatal("expected metadata map")
	}
}

func TestSignalValidateSuccess(t *testing.T) {
	signal := fraud.NewSignal(
		"tenant-1",
		"user-1",
		"login",
		10,
		"api",
	)

	if err := signal.Validate(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestSignalValidateNil(t *testing.T) {
	var signal *fraud.Signal

	err := signal.Validate()

	if err != fraud.ErrInvalidSignal {
		t.Fatalf(
			"expected %v got %v",
			fraud.ErrInvalidSignal,
			err,
		)
	}
}

func TestSignalValidateEmptyTenant(t *testing.T) {
	signal := fraud.NewSignal(
		"",
		"user-1",
		"login",
		10,
		"api",
	)

	err := signal.Validate()

	if err != fraud.ErrEmptyTenantID {
		t.Fatalf(
			"expected %v got %v",
			fraud.ErrEmptyTenantID,
			err,
		)
	}
}

func TestSignalValidateEmptyUser(t *testing.T) {
	signal := fraud.NewSignal(
		"tenant-1",
		"",
		"login",
		10,
		"api",
	)

	err := signal.Validate()

	if err != fraud.ErrEmptyUserID {
		t.Fatalf(
			"expected %v got %v",
			fraud.ErrEmptyUserID,
			err,
		)
	}
}

func TestSignalValidateEmptyType(t *testing.T) {
	signal := fraud.NewSignal(
		"tenant-1",
		"user-1",
		"",
		10,
		"api",
	)

	err := signal.Validate()

	if err != fraud.ErrEmptySignalType {
		t.Fatalf(
			"expected %v got %v",
			fraud.ErrEmptySignalType,
			err,
		)
	}
}

func TestSignalAddMetadata(t *testing.T) {
	signal := fraud.NewSignal(
		"tenant-1",
		"user-1",
		"login",
		10,
		"api",
	)

	signal.AddMetadata(
		"ip",
		"127.0.0.1",
	)

	value := signal.GetMetadata("ip")

	if value != "127.0.0.1" {
		t.Fatalf(
			"expected 127.0.0.1 got %s",
			value,
		)
	}
}

func TestSignalGetMetadataMissing(t *testing.T) {
	signal := fraud.NewSignal(
		"tenant-1",
		"user-1",
		"login",
		10,
		"api",
	)

	value := signal.GetMetadata(
		"missing",
	)

	if value != "" {
		t.Fatalf(
			"expected empty string got %s",
			value,
		)
	}
}

func TestSignalGetMetadataNilMap(t *testing.T) {
	signal := fraud.NewSignal(
		"tenant-1",
		"user-1",
		"login",
		10,
		"api",
	)

	signal.Metadata = nil

	value := signal.GetMetadata(
		"anything",
	)

	if value != "" {
		t.Fatalf(
			"expected empty string got %s",
			value,
		)
	}
}

func TestSignalAddMetadataNilMap(t *testing.T) {
	signal := fraud.NewSignal(
		"tenant-1",
		"user-1",
		"login",
		10,
		"api",
	)

	signal.Metadata = nil

	signal.AddMetadata(
		"country",
		"FR",
	)

	if signal.Metadata == nil {
		t.Fatal("metadata should be initialized")
	}

	if signal.Metadata["country"] != "FR" {
		t.Fatal("metadata value not stored")
	}
}

func TestSignalClone(t *testing.T) {
	original := fraud.NewSignal(
		"tenant-1",
		"user-1",
		"login",
		10,
		"api",
	)

	original.AddMetadata(
		"ip",
		"127.0.0.1",
	)

	clone := original.Clone()

	if clone == nil {
		t.Fatal("expected clone")
	}

	if clone == original {
		t.Fatal("clone should be different pointer")
	}

	if clone.ID != original.ID {
		t.Fatal("id mismatch")
	}

	if clone.TenantID != original.TenantID {
		t.Fatal("tenant mismatch")
	}

	if clone.Metadata["ip"] != "127.0.0.1" {
		t.Fatal("metadata mismatch")
	}

	clone.Metadata["ip"] = "8.8.8.8"

	if original.Metadata["ip"] == "8.8.8.8" {
		t.Fatal("clone should not modify original")
	}
}

func TestSignalCloneNil(t *testing.T) {
	var signal *fraud.Signal

	clone := signal.Clone()

	if clone != nil {
		t.Fatal("expected nil clone")
	}
}
