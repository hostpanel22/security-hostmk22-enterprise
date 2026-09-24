package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"strings"
)

type AEAD struct {
	gcm cipher.AEAD
}

func NewAEAD(key string) (*AEAD, error) {
	raw, err := decodeKey(key)
	if err != nil {
		return nil, err
	}
	block, err := aes.NewCipher(raw)
	if err != nil {
		return nil, err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}
	return &AEAD{gcm: gcm}, nil
}

func decodeKey(key string) ([]byte, error) {
	key = strings.TrimSpace(key)
	if key == "" {
		return nil, fmt.Errorf("encryption key is required")
	}
	if decoded, err := hex.DecodeString(key); err == nil && (len(decoded) == 16 || len(decoded) == 24 || len(decoded) == 32) {
		return decoded, nil
	}
	sum := sha256.Sum256([]byte(key))
	return sum[:], nil
}

func (a *AEAD) Encrypt(plain []byte) ([]byte, error) {
	if a == nil || a.gcm == nil {
		return nil, fmt.Errorf("aead not initialized")
	}
	nonce := make([]byte, a.gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, err
	}
	return a.gcm.Seal(nonce, nonce, plain, nil), nil
}

func (a *AEAD) Decrypt(payload []byte) ([]byte, error) {
	if a == nil || a.gcm == nil {
		return nil, fmt.Errorf("aead not initialized")
	}
	nonceSize := a.gcm.NonceSize()
	if len(payload) < nonceSize {
		return nil, fmt.Errorf("ciphertext too short")
	}
	nonce, ciphertext := payload[:nonceSize], payload[nonceSize:]
	return a.gcm.Open(nil, nonce, ciphertext, nil)
}
