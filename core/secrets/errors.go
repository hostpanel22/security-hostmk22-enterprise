package secrets

import "errors"

var (

	// Generic

	ErrInvalidSecret = errors.New(
		"invalid secret",
	)

	ErrSecretNotFound = errors.New(
		"secret not found",
	)

	ErrSecretAlreadyExists = errors.New(
		"secret already exists",
	)

	ErrVaultUnavailable = errors.New(
		"vault unavailable",
	)

	ErrStoreUnavailable = errors.New(
		"secret store unavailable",
	)

	ErrEncryptionUnavailable = errors.New(
		"encryption service unavailable",
	)

	// Validation

	ErrEmptySecretName = errors.New(
		"secret name cannot be empty",
	)

	ErrEmptySecretValue = errors.New(
		"secret value cannot be empty",
	)

	ErrEmptyTenantID = errors.New(
		"tenant id cannot be empty",
	)

	ErrEmptyUserID = errors.New(
		"user id cannot be empty",
	)

	ErrInvalidSecretType = errors.New(
		"invalid secret type",
	)

	ErrInvalidSecretVersion = errors.New(
		"invalid secret version",
	)

	// Encryption

	ErrEncryptionFailed = errors.New(
		"failed to encrypt secret",
	)

	ErrDecryptionFailed = errors.New(
		"failed to decrypt secret",
	)

	ErrInvalidEncryptionKey = errors.New(
		"invalid encryption key",
	)

	ErrKeyGenerationFailed = errors.New(
		"failed to generate encryption key",
	)

	ErrInvalidCiphertext = errors.New(
		"invalid ciphertext",
	)

	// Rotation

	ErrRotationFailed = errors.New(
		"secret rotation failed",
	)

	ErrRotationNotRequired = errors.New(
		"secret rotation not required",
	)

	ErrRotationScheduleInvalid = errors.New(
		"invalid rotation schedule",
	)

	// Access Control

	ErrAccessDenied = errors.New(
		"access denied",
	)

	ErrUnauthorizedAccess = errors.New(
		"unauthorized secret access",
	)

	ErrCrossTenantAccess = errors.New(
		"cross tenant access denied",
	)

	// Storage

	ErrSaveFailed = errors.New(
		"failed to save secret",
	)

	ErrLoadFailed = errors.New(
		"failed to load secret",
	)

	ErrDeleteFailed = errors.New(
		"failed to delete secret",
	)

	ErrUpdateFailed = errors.New(
		"failed to update secret",
	)

	// Audit

	ErrAuditFailed = errors.New(
		"failed to audit secret operation",
	)

	// Middleware

	ErrMissingTenantContext = errors.New(
		"missing tenant context",
	)

	ErrMissingUserContext = errors.New(
		"missing user context",
	)

	ErrMissingSecretContext = errors.New(
		"missing secret context",
	)
)
