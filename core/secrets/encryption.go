package secrets

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"io"
)

const (
	AES256KeySize = 32
	NonceSize     = 12
)

type EncryptionService struct {
	key []byte
}

func NewEncryptionService(
	key []byte,
) (*EncryptionService, error) {

	if len(key) != AES256KeySize {
		return nil, ErrInvalidEncryptionKey
	}

	service := &EncryptionService{
		key: make([]byte, AES256KeySize),
	}

	copy(
		service.key,
		key,
	)

	return service, nil
}

func GenerateKey() ([]byte, error) {

	key := make(
		[]byte,
		AES256KeySize,
	)

	if _, err := io.ReadFull(
		rand.Reader,
		key,
	); err != nil {

		return nil, ErrKeyGenerationFailed
	}

	return key, nil
}

func (e *EncryptionService) Encrypt(
	plaintext string,
) (string, error) {

	if e == nil {
		return "", ErrEncryptionUnavailable
	}

	block, err := aes.NewCipher(
		e.key,
	)

	if err != nil {
		return "", ErrEncryptionFailed
	}

	gcm, err := cipher.NewGCM(
		block,
	)

	if err != nil {
		return "", ErrEncryptionFailed
	}

	nonce := make(
		[]byte,
		gcm.NonceSize(),
	)

	if _, err := io.ReadFull(
		rand.Reader,
		nonce,
	); err != nil {

		return "", ErrEncryptionFailed
	}

	ciphertext := gcm.Seal(
		nonce,
		nonce,
		[]byte(plaintext),
		nil,
	)

	return base64.StdEncoding.EncodeToString(
		ciphertext,
	), nil
}

func (e *EncryptionService) Decrypt(
	ciphertext string,
) (string, error) {

	if e == nil {
		return "", ErrEncryptionUnavailable
	}

	decoded, err := base64.StdEncoding.DecodeString(
		ciphertext,
	)

	if err != nil {
		return "", ErrInvalidCiphertext
	}

	block, err := aes.NewCipher(
		e.key,
	)

	if err != nil {
		return "", ErrDecryptionFailed
	}

	gcm, err := cipher.NewGCM(
		block,
	)

	if err != nil {
		return "", ErrDecryptionFailed
	}

	nonceSize := gcm.NonceSize()

	if len(decoded) < nonceSize {
		return "", ErrInvalidCiphertext
	}

	nonce := decoded[:nonceSize]
	payload := decoded[nonceSize:]

	plaintext, err := gcm.Open(
		nil,
		nonce,
		payload,
		nil,
	)

	if err != nil {
		return "", ErrDecryptionFailed
	}

	return string(
		plaintext,
	), nil
}

func (e *EncryptionService) EncryptSecret(
	secret *Secret,
) error {

	if secret == nil {
		return ErrInvalidSecret
	}

	if secret.Encrypted {
		return nil
	}

	encrypted, err := e.Encrypt(
		secret.Value,
	)

	if err != nil {
		return err
	}

	secret.Value = encrypted
	secret.MarkEncrypted()

	return nil
}

func (e *EncryptionService) DecryptSecret(
	secret *Secret,
) error {

	if secret == nil {
		return ErrInvalidSecret
	}

	if !secret.Encrypted {
		return nil
	}

	decrypted, err := e.Decrypt(
		secret.Value,
	)

	if err != nil {
		return err
	}

	secret.Value = decrypted
	secret.MarkDecrypted()

	return nil
}

func (e *EncryptionService) KeySize() int {

	if e == nil {
		return 0
	}

	return len(
		e.key,
	)
}

func (e *EncryptionService) Validate() error {

	if e == nil {
		return ErrEncryptionUnavailable
	}

	if len(e.key) != AES256KeySize {
		return ErrInvalidEncryptionKey
	}

	return nil
}
