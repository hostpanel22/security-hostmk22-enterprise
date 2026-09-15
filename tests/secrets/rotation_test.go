package secrets_test

import (
	"testing"
	"time"

	"github.com/azezaz/HostMK-Security-Framework/core/secrets"
)

func TestNewRotationPolicy(t *testing.T) {

	policy := secrets.NewRotationPolicy(
		30,
	)

	expected := 30 * 24 * time.Hour

	if policy.Interval != expected {
		t.Fatalf(
			"expected %v got %v",
			expected,
			policy.Interval,
		)
	}
}

func TestNewRotationPolicyDefault(t *testing.T) {

	policy := secrets.NewRotationPolicy(
		0,
	)

	expected := time.Duration(
		secrets.DefaultRotationDays,
	) * 24 * time.Hour

	if policy.Interval != expected {
		t.Fatalf(
			"expected %v got %v",
			expected,
			policy.Interval,
		)
	}
}

func TestNewRotationManager(t *testing.T) {

	policy := secrets.NewRotationPolicy(
		60,
	)

	manager := secrets.NewRotationManager(
		policy,
	)

	if manager == nil {
		t.Fatal(
			"manager should not be nil",
		)
	}

	if manager.PolicyDays() != 60 {
		t.Fatalf(
			"expected 60 got %d",
			manager.PolicyDays(),
		)
	}
}

func TestRotationManagerValidatePolicy(t *testing.T) {

	manager := secrets.NewRotationManager(
		secrets.NewRotationPolicy(
			90,
		),
	)

	if err := manager.ValidatePolicy(); err != nil {
		t.Fatal(
			err,
		)
	}
}

func TestRotationManagerValidatePolicyNil(t *testing.T) {

	var manager *secrets.RotationManager

	if err := manager.ValidatePolicy(); err == nil {
		t.Fatal(
			"expected error",
		)
	}
}

func TestNextRotationDateFromCreatedAt(t *testing.T) {

	manager := secrets.NewRotationManager(
		secrets.NewRotationPolicy(
			30,
		),
	)

	secret := secrets.NewSecret(
		"tenant-1",
		"api-key",
		"value",
		"api",
	)

	secret.CreatedAt = time.Now().
		Add(
			-10 * 24 * time.Hour,
		)

	next := manager.NextRotationDate(
		secret,
	)

	expected := secret.CreatedAt.Add(
		30 * 24 * time.Hour,
	)

	if !next.Equal(
		expected,
	) {
		t.Fatal(
			"unexpected next rotation date",
		)
	}
}

func TestNextRotationDateFromRotatedAt(t *testing.T) {

	manager := secrets.NewRotationManager(
		secrets.NewRotationPolicy(
			30,
		),
	)

	secret := secrets.NewSecret(
		"tenant-1",
		"api-key",
		"value",
		"api",
	)

	rotated := time.Now().
		Add(
			-5 * 24 * time.Hour,
		)

	secret.RotatedAt = &rotated

	next := manager.NextRotationDate(
		secret,
	)

	expected := rotated.Add(
		30 * 24 * time.Hour,
	)

	if !next.Equal(
		expected,
	) {
		t.Fatal(
			"unexpected next rotation date",
		)
	}
}

func TestNeedsRotationFalse(t *testing.T) {

	manager := secrets.NewRotationManager(
		secrets.NewRotationPolicy(
			90,
		),
	)

	secret := secrets.NewSecret(
		"tenant-1",
		"api-key",
		"value",
		"api",
	)

	if manager.NeedsRotation(
		secret,
	) {
		t.Fatal(
			"rotation should not be required",
		)
	}
}

func TestNeedsRotationTrue(t *testing.T) {

	manager := secrets.NewRotationManager(
		secrets.NewRotationPolicy(
			30,
		),
	)

	secret := secrets.NewSecret(
		"tenant-1",
		"api-key",
		"value",
		"api",
	)

	secret.CreatedAt = time.Now().
		Add(
			-60 * 24 * time.Hour,
		)

	if !manager.NeedsRotation(
		secret,
	) {
		t.Fatal(
			"rotation should be required",
		)
	}
}

func TestRotateSecret(t *testing.T) {

	manager := secrets.NewRotationManager(
		secrets.NewRotationPolicy(
			30,
		),
	)

	secret := secrets.NewSecret(
		"tenant-1",
		"api-key",
		"old-value",
		"api",
	)

	secret.CreatedAt = time.Now().
		Add(
			-60 * 24 * time.Hour,
		)

	err := manager.Rotate(
		secret,
		"new-value",
	)

	if err != nil {
		t.Fatal(
			err,
		)
	}

	if secret.Value != "new-value" {
		t.Fatal(
			"rotation failed",
		)
	}

	if secret.Version != 2 {
		t.Fatalf(
			"expected version 2 got %d",
			secret.Version,
		)
	}
}

func TestRotateSecretNotRequired(t *testing.T) {

	manager := secrets.NewRotationManager(
		secrets.NewRotationPolicy(
			90,
		),
	)

	secret := secrets.NewSecret(
		"tenant-1",
		"api-key",
		"value",
		"api",
	)

	err := manager.Rotate(
		secret,
		"new-value",
	)

	if err != secrets.ErrRotationNotRequired {
		t.Fatalf(
			"expected ErrRotationNotRequired got %v",
			err,
		)
	}
}

func TestDaysUntilRotation(t *testing.T) {

	manager := secrets.NewRotationManager(
		secrets.NewRotationPolicy(
			30,
		),
	)

	secret := secrets.NewSecret(
		"tenant-1",
		"api-key",
		"value",
		"api",
	)

	secret.CreatedAt = time.Now().
		Add(
			-10 * 24 * time.Hour,
		)

	days := manager.DaysUntilRotation(
		secret,
	)

	if days <= 0 {
		t.Fatalf(
			"expected positive days got %d",
			days,
		)
	}
}

func TestIsExpired(t *testing.T) {

	manager := secrets.NewRotationManager(
		secrets.NewRotationPolicy(
			30,
		),
	)

	secret := secrets.NewSecret(
		"tenant-1",
		"api-key",
		"value",
		"api",
	)

	secret.CreatedAt = time.Now().
		Add(
			-60 * 24 * time.Hour,
		)

	if !manager.IsExpired(
		secret,
	) {
		t.Fatal(
			"secret should be expired",
		)
	}
}

func TestPolicyDays(t *testing.T) {

	manager := secrets.NewRotationManager(
		secrets.NewRotationPolicy(
			45,
		),
	)

	if manager.PolicyDays() != 45 {
		t.Fatalf(
			"expected 45 got %d",
			manager.PolicyDays(),
		)
	}
}
