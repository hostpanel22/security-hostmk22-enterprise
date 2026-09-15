package secrets_test

import (
	"context"
	"testing"

	"github.com/azezaz/HostMK-Security-Framework/core/secrets"
)

func TestNewMemoryStore(t *testing.T) {

	store := secrets.NewMemoryStore()

	if store == nil {
		t.Fatal("store should not be nil")
	}

	if store.Count() != 0 {
		t.Fatal("new store should be empty")
	}
}

func TestStoreSave(t *testing.T) {

	ctx := context.Background()

	store := secrets.NewMemoryStore()

	secret := secrets.NewSecret(
		"tenant-1",
		"api-key",
		"value",
		"api_key",
	)

	err := store.Save(
		ctx,
		secret,
	)

	if err != nil {
		t.Fatalf(
			"unexpected error: %v",
			err,
		)
	}

	if store.Count() != 1 {
		t.Fatal("expected one secret")
	}
}

func TestStoreSaveNilSecret(t *testing.T) {

	ctx := context.Background()

	store := secrets.NewMemoryStore()

	err := store.Save(
		ctx,
		nil,
	)

	if err == nil {
		t.Fatal("expected error")
	}
}

func TestStoreGetByID(t *testing.T) {

	ctx := context.Background()

	store := secrets.NewMemoryStore()

	secret := secrets.NewSecret(
		"tenant-1",
		"db-password",
		"secret",
		"password",
	)

	_ = store.Save(
		ctx,
		secret,
	)

	found, err := store.GetByID(
		ctx,
		secret.ID,
	)

	if err != nil {
		t.Fatalf(
			"unexpected error: %v",
			err,
		)
	}

	if found.ID != secret.ID {
		t.Fatal("id mismatch")
	}
}

func TestStoreGetByIDNotFound(t *testing.T) {

	ctx := context.Background()

	store := secrets.NewMemoryStore()

	_, err := store.GetByID(
		ctx,
		"unknown",
	)

	if err == nil {
		t.Fatal("expected error")
	}
}

func TestStoreGetByName(t *testing.T) {

	ctx := context.Background()

	store := secrets.NewMemoryStore()

	secret := secrets.NewSecret(
		"tenant-1",
		"stripe-key",
		"sk_test",
		"api_key",
	)

	_ = store.Save(
		ctx,
		secret,
	)

	found, err := store.GetByName(
		ctx,
		"stripe-key",
	)

	if err != nil {
		t.Fatalf(
			"unexpected error: %v",
			err,
		)
	}

	if found.Name != "stripe-key" {
		t.Fatal("invalid secret")
	}
}

func TestStoreGetByNameNotFound(t *testing.T) {

	ctx := context.Background()

	store := secrets.NewMemoryStore()

	_, err := store.GetByName(
		ctx,
		"missing",
	)

	if err == nil {
		t.Fatal("expected error")
	}
}

func TestStoreList(t *testing.T) {

	ctx := context.Background()

	store := secrets.NewMemoryStore()

	_ = store.Save(
		ctx,
		secrets.NewSecret(
			"tenant-1",
			"secret-1",
			"value-1",
			"api_key",
		),
	)

	_ = store.Save(
		ctx,
		secrets.NewSecret(
			"tenant-1",
			"secret-2",
			"value-2",
			"password",
		),
	)

	list, err := store.List(
		ctx,
	)

	if err != nil {
		t.Fatalf(
			"unexpected error: %v",
			err,
		)
	}

	if len(list) != 2 {
		t.Fatal("expected two secrets")
	}
}

func TestStoreDelete(t *testing.T) {

	ctx := context.Background()

	store := secrets.NewMemoryStore()

	secret := secrets.NewSecret(
		"tenant-1",
		"delete-me",
		"value",
		"password",
	)

	_ = store.Save(
		ctx,
		secret,
	)

	err := store.Delete(
		ctx,
		secret.ID,
	)

	if err != nil {
		t.Fatalf(
			"unexpected error: %v",
			err,
		)
	}

	if store.Count() != 0 {
		t.Fatal("store should be empty")
	}
}

func TestStoreDeleteNotFound(t *testing.T) {

	ctx := context.Background()

	store := secrets.NewMemoryStore()

	err := store.Delete(
		ctx,
		"missing",
	)

	if err == nil {
		t.Fatal("expected error")
	}
}

func TestStoreClear(t *testing.T) {

	ctx := context.Background()

	store := secrets.NewMemoryStore()

	_ = store.Save(
		ctx,
		secrets.NewSecret(
			"tenant-1",
			"a",
			"b",
			"api_key",
		),
	)

	_ = store.Save(
		ctx,
		secrets.NewSecret(
			"tenant-1",
			"c",
			"d",
			"password",
		),
	)

	err := store.Clear(
		ctx,
	)

	if err != nil {
		t.Fatalf(
			"unexpected error: %v",
			err,
		)
	}

	if store.Count() != 0 {
		t.Fatal("store should be empty")
	}
}

func TestStoreReturnsClone(t *testing.T) {

	ctx := context.Background()

	store := secrets.NewMemoryStore()

	secret := secrets.NewSecret(
		"tenant-1",
		"clone-test",
		"value",
		"api_key",
	)

	_ = store.Save(
		ctx,
		secret,
	)

	found, err := store.GetByID(
		ctx,
		secret.ID,
	)

	if err != nil {
		t.Fatal(err)
	}

	found.Name = "modified"

	again, err := store.GetByID(
		ctx,
		secret.ID,
	)

	if err != nil {
		t.Fatal(err)
	}

	if again.Name != "clone-test" {
		t.Fatal("store leaked internal state")
	}
}

func TestStoreCount(t *testing.T) {

	ctx := context.Background()

	store := secrets.NewMemoryStore()

	if store.Count() != 0 {
		t.Fatal("expected zero")
	}

	_ = store.Save(
		ctx,
		secrets.NewSecret(
			"tenant-1",
			"count-test",
			"value",
			"api_key",
		),
	)

	if store.Count() != 1 {
		t.Fatal("expected one")
	}
}
