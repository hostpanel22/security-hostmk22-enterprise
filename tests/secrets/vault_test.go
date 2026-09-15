package secrets_test

import (
	"context"
	"testing"

	"github.com/azezaz/HostMK-Security-Framework/core/secrets"
)

func createVault(t *testing.T) *secrets.Vault {
	t.Helper()

	store := secrets.NewMemoryStore()

	key, err := secrets.GenerateKey()
	if err != nil {
		t.Fatalf("generate key: %v", err)
	}

	enc, err := secrets.NewEncryptionService(key)
	if err != nil {
		t.Fatalf("new encryption: %v", err)
	}

	vault, err := secrets.NewVault(
		store,
		enc,
	)
	if err != nil {
		t.Fatalf("new vault: %v", err)
	}

	return vault
}

func TestVaultCreateSecret(t *testing.T) {
	vault := createVault(t)

	secret := secrets.NewSecret(
		"tenant-1",
		"api-key",
		"secret-value",
		"api",
	)

	err := vault.CreateSecret(
		context.Background(),
		secret,
	)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestVaultGetSecretByID(t *testing.T) {
	vault := createVault(t)

	secret := secrets.NewSecret(
		"tenant-1",
		"api-key",
		"value",
		"api",
	)

	if err := vault.CreateSecret(
		context.Background(),
		secret,
	); err != nil {
		t.Fatal(err)
	}

	result, err := vault.GetSecretByID(
		context.Background(),
		secret.ID,
	)

	if err != nil {
		t.Fatal(err)
	}

	if result.Name != "api-key" {
		t.Fatalf("expected api-key got %s", result.Name)
	}
}

func TestVaultGetSecretByName(t *testing.T) {
	vault := createVault(t)

	secret := secrets.NewSecret(
		"tenant-1",
		"jwt-secret",
		"value",
		"jwt",
	)

	if err := vault.CreateSecret(
		context.Background(),
		secret,
	); err != nil {
		t.Fatal(err)
	}

	result, err := vault.GetSecretByName(
		context.Background(),
		"jwt-secret",
	)

	if err != nil {
		t.Fatal(err)
	}

	if result.Name != "jwt-secret" {
		t.Fatalf("unexpected name %s", result.Name)
	}
}

func TestVaultRotateSecret(t *testing.T) {
	vault := createVault(t)

	secret := secrets.NewSecret(
		"tenant-1",
		"token",
		"old",
		"api",
	)

	if err := vault.CreateSecret(
		context.Background(),
		secret,
	); err != nil {
		t.Fatal(err)
	}

	err := vault.RotateSecret(
		context.Background(),
		secret.ID,
		"new-value",
	)

	if err != nil {
		t.Fatal(err)
	}

	updated, err := vault.GetSecretByID(
		context.Background(),
		secret.ID,
	)

	if err != nil {
		t.Fatal(err)
	}

	if updated.Version != 2 {
		t.Fatalf(
			"expected version 2 got %d",
			updated.Version,
		)
	}
}

func TestVaultDeleteSecret(t *testing.T) {
	vault := createVault(t)

	secret := secrets.NewSecret(
		"tenant-1",
		"delete-me",
		"value",
		"api",
	)

	if err := vault.CreateSecret(
		context.Background(),
		secret,
	); err != nil {
		t.Fatal(err)
	}

	err := vault.DeleteSecret(
		context.Background(),
		secret.ID,
	)

	if err != nil {
		t.Fatal(err)
	}

	_, err = vault.GetSecretByID(
		context.Background(),
		secret.ID,
	)

	if err == nil {
		t.Fatal("expected error")
	}
}

func TestVaultListSecrets(t *testing.T) {
	vault := createVault(t)

	_ = vault.CreateSecret(
		context.Background(),
		secrets.NewSecret(
			"tenant-a",
			"one",
			"1",
			"api",
		),
	)

	_ = vault.CreateSecret(
		context.Background(),
		secrets.NewSecret(
			"tenant-a",
			"two",
			"2",
			"api",
		),
	)

	_ = vault.CreateSecret(
		context.Background(),
		secrets.NewSecret(
			"tenant-b",
			"three",
			"3",
			"api",
		),
	)

	list, err := vault.ListSecrets(
		context.Background(),
		"tenant-a",
	)

	if err != nil {
		t.Fatal(err)
	}

	if len(list) != 2 {
		t.Fatalf(
			"expected 2 secrets got %d",
			len(list),
		)
	}
}

func TestVaultSecretCount(t *testing.T) {
	vault := createVault(t)

	_ = vault.CreateSecret(
		context.Background(),
		secrets.NewSecret(
			"tenant-1",
			"a",
			"1",
			"api",
		),
	)

	_ = vault.CreateSecret(
		context.Background(),
		secrets.NewSecret(
			"tenant-1",
			"b",
			"2",
			"api",
		),
	)

	count, err := vault.SecretCount(
		context.Background(),
		"tenant-1",
	)

	if err != nil {
		t.Fatal(err)
	}

	if count != 2 {
		t.Fatalf(
			"expected 2 got %d",
			count,
		)
	}
}

func TestVaultClear(t *testing.T) {
	vault := createVault(t)

	_ = vault.CreateSecret(
		context.Background(),
		secrets.NewSecret(
			"tenant-1",
			"a",
			"1",
			"api",
		),
	)

	err := vault.Clear(
		context.Background(),
	)

	if err != nil {
		t.Fatal(err)
	}

	count, err := vault.SecretCount(
		context.Background(),
		"",
	)

	if err != nil {
		t.Fatal(err)
	}

	if count != 0 {
		t.Fatalf(
			"expected 0 got %d",
			count,
		)
	}
}
