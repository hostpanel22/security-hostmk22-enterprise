package secrets_test

import (
	"testing"

	"github.com/azezaz/HostMK-Security-Framework/core/secrets"
)

func TestNewSecret(t *testing.T) {

	secret := secrets.NewSecret(
		"tenant-1",
		"api-key",
		"secret-value",
		"api_key",
	)

	if secret == nil {
		t.Fatal("secret should not be nil")
	}

	if secret.TenantID != "tenant-1" {
		t.Fatal("invalid tenant id")
	}

	if secret.Name != "api-key" {
		t.Fatal("invalid name")
	}

	if secret.Value != "secret-value" {
		t.Fatal("invalid value")
	}

	if secret.Version != 1 {
		t.Fatal("version should start at 1")
	}
}

func TestSecretValidateSuccess(t *testing.T) {

	secret := secrets.NewSecret(
		"tenant-1",
		"api-key",
		"value",
		"api_key",
	)

	if err := secret.Validate(); err != nil {
		t.Fatalf(
			"expected nil error got %v",
			err,
		)
	}
}

func TestSecretValidateTenantRequired(t *testing.T) {

	secret := secrets.NewSecret(
		"",
		"api-key",
		"value",
		"api_key",
	)

	if err := secret.Validate(); err == nil {
		t.Fatal("expected error")
	}
}

func TestSecretValidateNameRequired(t *testing.T) {

	secret := secrets.NewSecret(
		"tenant-1",
		"",
		"value",
		"api_key",
	)

	if err := secret.Validate(); err == nil {
		t.Fatal("expected error")
	}
}

func TestSecretValidateValueRequired(t *testing.T) {

	secret := secrets.NewSecret(
		"tenant-1",
		"api-key",
		"",
		"api_key",
	)

	if err := secret.Validate(); err == nil {
		t.Fatal("expected error")
	}
}

func TestSecretValidateVersion(t *testing.T) {

	secret := secrets.NewSecret(
		"tenant-1",
		"api-key",
		"value",
		"api_key",
	)

	secret.Version = 0

	if err := secret.Validate(); err == nil {
		t.Fatal("expected error")
	}
}

func TestSecretRotate(t *testing.T) {

	secret := secrets.NewSecret(
		"tenant-1",
		"api-key",
		"old-value",
		"api_key",
	)

	oldVersion := secret.Version

	err := secret.Rotate(
		"new-value",
	)

	if err != nil {
		t.Fatalf(
			"unexpected error %v",
			err,
		)
	}

	if secret.Value != "new-value" {
		t.Fatal("rotation failed")
	}

	if secret.Version != oldVersion+1 {
		t.Fatal("version not incremented")
	}

	if secret.RotatedAt == nil {
		t.Fatal("rotated timestamp missing")
	}
}

func TestSecretRotateRejectsEmptyValue(t *testing.T) {

	secret := secrets.NewSecret(
		"tenant-1",
		"api-key",
		"value",
		"api_key",
	)

	if err := secret.Rotate(""); err == nil {
		t.Fatal("expected error")
	}
}

func TestSecretEncryptionFlags(t *testing.T) {

	secret := secrets.NewSecret(
		"tenant-1",
		"api-key",
		"value",
		"api_key",
	)

	if secret.IsEncrypted() {
		t.Fatal("should not be encrypted")
	}

	secret.MarkEncrypted()

	if !secret.IsEncrypted() {
		t.Fatal("should be encrypted")
	}

	secret.MarkDecrypted()

	if secret.IsEncrypted() {
		t.Fatal("should be decrypted")
	}
}

func TestSecretClone(t *testing.T) {

	secret := secrets.NewSecret(
		"tenant-1",
		"api-key",
		"value",
		"api_key",
	)

	clone := secret.Clone()

	if clone == nil {
		t.Fatal("clone is nil")
	}

	clone.Name = "modified"

	if secret.Name == clone.Name {
		t.Fatal("clone should be independent")
	}
}

func TestBelongsToTenant(t *testing.T) {

	secret := secrets.NewSecret(
		"tenant-123",
		"api-key",
		"value",
		"api_key",
	)

	if !secret.BelongsToTenant(
		"tenant-123",
	) {
		t.Fatal("should belong")
	}

	if secret.BelongsToTenant(
		"tenant-999",
	) {
		t.Fatal("should not belong")
	}
}
