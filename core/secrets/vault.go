package secrets

import (
	"context"
	"strings"
	"sync"
)

type Vault struct {
	mu sync.RWMutex

	store Store

	encryption *EncryptionService
}

func NewVault(
	store Store,
	encryption *EncryptionService,
) (*Vault, error) {

	if store == nil {
		return nil, ErrStoreUnavailable
	}

	return &Vault{
		store:      store,
		encryption: encryption,
	}, nil
}

func (v *Vault) CreateSecret(
	ctx context.Context,
	secret *Secret,
) error {

	if v == nil {
		return ErrVaultUnavailable
	}

	if secret == nil {
		return ErrInvalidSecret
	}

	if err := secret.Validate(); err != nil {
		return err
	}

	_, err := v.store.GetByName(
		ctx,
		secret.Name,
	)

	if err == nil {
		return ErrSecretAlreadyExists
	}

	copySecret := secret.Clone()

	if v.encryption != nil &&
		!copySecret.Encrypted {

		if err := v.encryption.EncryptSecret(
			copySecret,
		); err != nil {
			return err
		}
	}

	return v.store.Save(
		ctx,
		copySecret,
	)
}

func (v *Vault) GetSecretByID(
	ctx context.Context,
	id string,
) (*Secret, error) {

	if v == nil {
		return nil, ErrVaultUnavailable
	}

	secret, err := v.store.GetByID(
		ctx,
		id,
	)

	if err != nil {
		return nil, err
	}

	copySecret := secret.Clone()

	if v.encryption != nil &&
		copySecret.Encrypted {

		if err := v.encryption.DecryptSecret(
			copySecret,
		); err != nil {
			return nil, err
		}
	}

	return copySecret, nil
}

func (v *Vault) GetSecretByName(
	ctx context.Context,
	name string,
) (*Secret, error) {

	if v == nil {
		return nil, ErrVaultUnavailable
	}

	secret, err := v.store.GetByName(
		ctx,
		name,
	)

	if err != nil {
		return nil, err
	}

	copySecret := secret.Clone()

	if v.encryption != nil &&
		copySecret.Encrypted {

		if err := v.encryption.DecryptSecret(
			copySecret,
		); err != nil {
			return nil, err
		}
	}

	return copySecret, nil
}

func (v *Vault) RotateSecret(
	ctx context.Context,
	id string,
	newValue string,
) error {

	if v == nil {
		return ErrVaultUnavailable
	}

	secret, err := v.store.GetByID(
		ctx,
		id,
	)

	if err != nil {
		return err
	}

	if err := secret.Rotate(
		newValue,
	); err != nil {
		return err
	}

	if v.encryption != nil {

		if err := v.encryption.EncryptSecret(
			secret,
		); err != nil {
			return err
		}
	}

	return v.store.Save(
		ctx,
		secret,
	)
}

func (v *Vault) DeleteSecret(
	ctx context.Context,
	id string,
) error {

	if v == nil {
		return ErrVaultUnavailable
	}

	return v.store.Delete(
		ctx,
		id,
	)
}

func (v *Vault) ListSecrets(
	ctx context.Context,
	tenantID string,
) ([]*Secret, error) {

	if v == nil {
		return nil, ErrVaultUnavailable
	}

	secrets, err := v.store.List(
		ctx,
	)

	if err != nil {
		return nil, err
	}

	result := make(
		[]*Secret,
		0,
	)

	for _, secret := range secrets {

		if secret == nil {
			continue
		}

		if tenantID != "" &&
			strings.TrimSpace(
				secret.TenantID,
			) != strings.TrimSpace(
				tenantID,
			) {
			continue
		}

		result = append(
			result,
			secret.Clone(),
		)
	}

	return result, nil
}

func (v *Vault) SecretCount(
	ctx context.Context,
	tenantID string,
) (int, error) {

	secrets, err := v.ListSecrets(
		ctx,
		tenantID,
	)

	if err != nil {
		return 0, err
	}

	return len(
		secrets,
	), nil
}

func (v *Vault) EncryptSecret(
	secret *Secret,
) error {

	if v == nil {
		return ErrVaultUnavailable
	}

	if v.encryption == nil {
		return ErrEncryptionUnavailable
	}

	return v.encryption.EncryptSecret(
		secret,
	)
}

func (v *Vault) DecryptSecret(
	secret *Secret,
) error {

	if v == nil {
		return ErrVaultUnavailable
	}

	if v.encryption == nil {
		return ErrEncryptionUnavailable
	}

	return v.encryption.DecryptSecret(
		secret,
	)
}

func (v *Vault) Clear(
	ctx context.Context,
) error {

	if v == nil {
		return ErrVaultUnavailable
	}

	return v.store.Clear(
		ctx,
	)
}

func (v *Vault) Store() Store {

	if v == nil {
		return nil
	}

	return v.store
}
