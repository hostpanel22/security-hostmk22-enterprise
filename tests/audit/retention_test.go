package audit_test

import (
	"context"
	"testing"
	"time"

	"github.com/azezaz/HostMK-Security-Framework/core/audit"
)

func TestNewRetentionManager(t *testing.T) {

	store := audit.NewMemoryStore()

	manager := audit.NewRetentionManager(
		store,
	)

	if manager == nil {
		t.Fatal("manager should not be nil")
	}

	if len(manager.Policies()) != 0 {
		t.Fatal("expected no policies")
	}
}

func TestAddPolicy(t *testing.T) {

	store := audit.NewMemoryStore()

	manager := audit.NewRetentionManager(
		store,
	)

	err := manager.AddPolicy(
		audit.RetentionPolicy{
			Name:            "30-days",
			RetentionPeriod: 30 * 24 * time.Hour,
			Enabled:         true,
		},
	)

	if err != nil {
		t.Fatal(err)
	}

	if len(manager.Policies()) != 1 {
		t.Fatal("expected one policy")
	}
}

func TestAddPolicyInvalidPeriod(t *testing.T) {

	store := audit.NewMemoryStore()

	manager := audit.NewRetentionManager(
		store,
	)

	err := manager.AddPolicy(
		audit.RetentionPolicy{
			Name:            "invalid",
			RetentionPeriod: 0,
			Enabled:         true,
		},
	)

	if err != audit.ErrInvalidRetentionPeriod {
		t.Fatalf(
			"expected ErrInvalidRetentionPeriod got %v",
			err,
		)
	}
}

func TestPoliciesReturnsCopy(t *testing.T) {

	store := audit.NewMemoryStore()

	manager := audit.NewRetentionManager(
		store,
	)

	_ = manager.AddPolicy(
		audit.RetentionPolicy{
			Name:            "copy-test",
			RetentionPeriod: 24 * time.Hour,
			Enabled:         true,
		},
	)

	policies := manager.Policies()

	policies[0].Name = "modified"

	original := manager.Policies()

	if original[0].Name != "copy-test" {
		t.Fatal("policies should be immutable copy")
	}
}

func TestExecuteWithoutStore(t *testing.T) {

	manager := audit.NewRetentionManager(
		nil,
	)

	_, err := manager.Execute(
		context.Background(),
	)

	if err != audit.ErrStoreUnavailable {
		t.Fatalf(
			"expected ErrStoreUnavailable got %v",
			err,
		)
	}
}

func TestExecuteNoPolicies(t *testing.T) {

	store := audit.NewMemoryStore()

	manager := audit.NewRetentionManager(
		store,
	)

	results, err := manager.Execute(
		context.Background(),
	)

	if err != nil {
		t.Fatal(err)
	}

	if len(results) != 0 {
		t.Fatal("expected zero results")
	}
}

func TestExecuteDisabledPolicy(t *testing.T) {

	store := audit.NewMemoryStore()

	manager := audit.NewRetentionManager(
		store,
	)

	_ = manager.AddPolicy(
		audit.RetentionPolicy{
			Name:            "disabled",
			RetentionPeriod: 24 * time.Hour,
			Enabled:         false,
		},
	)

	results, err := manager.Execute(
		context.Background(),
	)

	if err != nil {
		t.Fatal(err)
	}

	if len(results) != 0 {
		t.Fatal("disabled policy should not execute")
	}
}

func TestExecutePolicy(t *testing.T) {

	store := audit.NewMemoryStore()

	manager := audit.NewRetentionManager(
		store,
	)

	_ = manager.AddPolicy(
		audit.RetentionPolicy{
			Name:            "30-days",
			RetentionPeriod: 30 * 24 * time.Hour,
			Enabled:         true,
		},
	)

	results, err := manager.Execute(
		context.Background(),
	)

	if err != nil {
		t.Fatal(err)
	}

	if len(results) != 1 {
		t.Fatalf(
			"expected 1 result got %d",
			len(results),
		)
	}

	if results[0].PolicyName != "30-days" {
		t.Fatal("unexpected policy name")
	}
}

func TestDefaultPolicies(t *testing.T) {

	policies := audit.DefaultPolicies()

	if len(policies) != 3 {
		t.Fatalf(
			"expected 3 policies got %d",
			len(policies),
		)
	}

	for _, policy := range policies {

		if !policy.Enabled {
			t.Fatal("default policy should be enabled")
		}

		if policy.RetentionPeriod <= 0 {
			t.Fatal("invalid retention period")
		}
	}
}

func TestLoadDefaults(t *testing.T) {

	store := audit.NewMemoryStore()

	manager := audit.NewRetentionManager(
		store,
	)

	err := manager.LoadDefaults()

	if err != nil {
		t.Fatal(err)
	}

	policies := manager.Policies()

	if len(policies) != 3 {
		t.Fatalf(
			"expected 3 policies got %d",
			len(policies),
		)
	}
}

func TestLoadDefaultsThenExecute(t *testing.T) {

	store := audit.NewMemoryStore()

	manager := audit.NewRetentionManager(
		store,
	)

	if err := manager.LoadDefaults(); err != nil {
		t.Fatal(err)
	}

	results, err := manager.Execute(
		context.Background(),
	)

	if err != nil {
		t.Fatal(err)
	}

	if len(results) != 3 {
		t.Fatalf(
			"expected 3 results got %d",
			len(results),
		)
	}
}

func TestRetentionResultTimestamp(t *testing.T) {

	store := audit.NewMemoryStore()

	manager := audit.NewRetentionManager(
		store,
	)

	_ = manager.AddPolicy(
		audit.RetentionPolicy{
			Name:            "timestamp-test",
			RetentionPeriod: time.Hour,
			Enabled:         true,
		},
	)

	results, err := manager.Execute(
		context.Background(),
	)

	if err != nil {
		t.Fatal(err)
	}

	if len(results) != 1 {
		t.Fatal("expected one result")
	}

	if results[0].ExecutedAt.IsZero() {
		t.Fatal("execution timestamp missing")
	}
}
