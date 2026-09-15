package secrets_test

import (
	"testing"

	"github.com/azezaz/HostMK-Security-Framework/core/secrets"
)

func TestGenerateKey(t *testing.T) {

	key, err := secrets.GenerateKey()

	if err != nil {
		t.Fatal(err)
	}

	if len(key) != secrets.AES256KeySize {
		t.Fatalf(
			"expected key size %d got %d",
			secrets.AES256KeySize,
			len(key),
		)
	}
}

func TestNewEncryptionService(t *testing.T) {

	key := make(
		[]byte,
		secrets.AES256KeySize,
	)

	service, err := secrets.NewEncryptionService(
		key,
	)

	if err != nil {
		t.Fatal(err)
	}

	if service == nil {
		t.Fatal(
			"service should not be nil",
		)
	}
}

func TestNewEncryptionServiceInvalidKey(t *testing.T) {

	_, err := secrets.NewEncryptionService(
		[]byte("short"),
	)

	if err != secrets.ErrInvalidEncryptionKey {
		t.Fatalf(
			"expected %v got %v",
			secrets.ErrInvalidEncryptionKey,
			err,
		)
	}
}

func TestEncryptDecrypt(t *testing.T) {

	key, err := secrets.GenerateKey()

	if err != nil {
		t.Fatal(err)
	}

	service, err := secrets.NewEncryptionService(
		key,
	)

	if err != nil {
		t.Fatal(err)
	}

	plaintext := "super-secret-value"

	ciphertext, err := service.Encrypt(
		plaintext,
	)

	if err != nil {
		t.Fatal(err)
	}

	if ciphertext == plaintext {
		t.Fatal(
			"ciphertext should differ from plaintext",
		)
	}

	decrypted, err := service.Decrypt(
		ciphertext,
	)

	if err != nil {
		t.Fatal(err)
	}

	if decrypted != plaintext {
		t.Fatalf(
			"expected %s got %s",
			plaintext,
			decrypted,
		)
	}
}

func TestEncryptNilService(t *testing.T) {

	var service *secrets.EncryptionService

	_, err := service.Encrypt(
		"secret",
	)

	if err != secrets.ErrEncryptionUnavailable {
		t.Fatalf(
			"expected %v got %v",
			secrets.ErrEncryptionUnavailable,
			err,
		)
	}
}

func TestDecryptNilService(t *testing.T) {

	var service *secrets.EncryptionService

	_, err := service.Decrypt(
		"cipher",
	)

	if err != secrets.ErrEncryptionUnavailable {
		t.Fatalf(
			"expected %v got %v",
			secrets.ErrEncryptionUnavailable,
			err,
		)
	}
}

func TestDecryptInvalidCiphertext(t *testing.T) {

	key, err := secrets.GenerateKey()

	if err != nil {
		t.Fatal(err)
	}

	service, err := secrets.NewEncryptionService(
		key,
	)

	if err != nil {
		t.Fatal(err)
	}

	_, err = service.Decrypt(
		"not-base64",
	)

	if err != secrets.ErrInvalidCiphertext {
		t.Fatalf(
			"expected %v got %v",
			secrets.ErrInvalidCiphertext,
			err,
		)
	}
}

func TestEncryptSecret(t *testing.T) {

	key, err := secrets.GenerateKey()

	if err != nil {
		t.Fatal(err)
	}

	service, err := secrets.NewEncryptionService(
		key,
	)

	if err != nil {
		t.Fatal(err)
	}

	secret := secrets.NewSecret(
		"tenant-1",
		"api-key",
		"plain-value",
		"api",
	)

	err = service.EncryptSecret(
		secret,
	)

	if err != nil {
		t.Fatal(err)
	}

	if !secret.IsEncrypted() {
		t.Fatal(
			"secret should be encrypted",
		)
	}

	if secret.Value == "plain-value" {
		t.Fatal(
			"value should be encrypted",
		)
	}
}

func TestDecryptSecret(t *testing.T) {

	key, err := secrets.GenerateKey()

	if err != nil {
		t.Fatal(err)
	}

	service, err := secrets.NewEncryptionService(
		key,
	)

	if err != nil {
		t.Fatal(err)
	}

	secret := secrets.NewSecret(
		"tenant-1",
		"api-key",
		"plain-value",
		"api",
	)

	err = service.EncryptSecret(
		secret,
	)

	if err != nil {
		t.Fatal(err)
	}

	err = service.DecryptSecret(
		secret,
	)

	if err != nil {
		t.Fatal(err)
	}

	if secret.IsEncrypted() {
		t.Fatal(
			"secret should be decrypted",
		)
	}

	if secret.Value != "plain-value" {
		t.Fatalf(
			"expected plain-value got %s",
			secret.Value,
		)
	}
}

func TestEncryptSecretNil(t *testing.T) {

	key, _ := secrets.GenerateKey()

	service, _ := secrets.NewEncryptionService(
		key,
	)

	err := service.EncryptSecret(
		nil,
	)

	if err != secrets.ErrInvalidSecret {
		t.Fatalf(
			"expected %v got %v",
			secrets.ErrInvalidSecret,
			err,
		)
	}
}

func TestDecryptSecretNil(t *testing.T) {

	key, _ := secrets.GenerateKey()

	service, _ := secrets.NewEncryptionService(
		key,
	)

	err := service.DecryptSecret(
		nil,
	)

	if err != secrets.ErrInvalidSecret {
		t.Fatalf(
			"expected %v got %v",
			secrets.ErrInvalidSecret,
			err,
		)
	}
}

func TestKeySize(t *testing.T) {

	key, _ := secrets.GenerateKey()

	service, _ := secrets.NewEncryptionService(
		key,
	)

	if service.KeySize() != secrets.AES256KeySize {
		t.Fatalf(
			"expected %d got %d",
			secrets.AES256KeySize,
			service.KeySize(),
		)
	}
}

func TestValidate(t *testing.T) {

	key, _ := secrets.GenerateKey()

	service, _ := secrets.NewEncryptionService(
		key,
	)

	if err := service.Validate(); err != nil {
		t.Fatal(err)
	}
}

func TestValidateNilService(t *testing.T) {

	var service *secrets.EncryptionService

	if err := service.Validate(); err != secrets.ErrEncryptionUnavailable {
		t.Fatalf(
			"expected %v got %v",
			secrets.ErrEncryptionUnavailable,
			err,
		)
	}
}
