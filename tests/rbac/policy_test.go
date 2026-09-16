package rbac_test

import (
	"testing"
	"time"

	"github.com/azezaz/HostMK-Security-Framework/core/rbac"
)

func TestNewPolicy(t *testing.T) {
	policy := rbac.NewPolicy(
		"tenant-1",
		"admin-policy",
		"administration policy",
		rbac.PolicyAllow,
	)

	if policy == nil {
		t.Fatal("policy is nil")
	}

	if policy.Name != "admin-policy" {
		t.Fatalf("expected admin-policy, got %s", policy.Name)
	}

	if policy.Effect != rbac.PolicyAllow {
		t.Fatal("invalid effect")
	}
}

func TestPolicyValidate(t *testing.T) {
	policy := rbac.NewPolicy(
		"tenant",
		"policy",
		"description",
		rbac.PolicyAllow,
	)

	if err := policy.Validate(); err != nil {
		t.Fatal(err)
	}
}

func TestPolicyValidateNil(t *testing.T) {
	var policy *rbac.Policy

	if err := policy.Validate(); err != rbac.ErrInvalidPolicy {
		t.Fatalf("expected ErrInvalidPolicy, got %v", err)
	}
}

func TestPolicyValidateEmptyID(t *testing.T) {
	policy := rbac.NewPolicy(
		"tenant",
		"policy",
		"description",
		rbac.PolicyAllow,
	)

	policy.ID = ""

	if err := policy.Validate(); err != rbac.ErrInvalidPolicy {
		t.Fatalf("expected ErrInvalidPolicy, got %v", err)
	}
}

func TestPolicyValidateEmptyName(t *testing.T) {
	policy := rbac.NewPolicy(
		"tenant",
		"policy",
		"description",
		rbac.PolicyAllow,
	)

	policy.Name = ""

	if err := policy.Validate(); err != rbac.ErrEmptyPolicyName {
		t.Fatalf("expected ErrEmptyPolicyName, got %v", err)
	}
}

func TestPolicyClone(t *testing.T) {
	policy := rbac.NewPolicy(
		"tenant",
		"policy",
		"description",
		rbac.PolicyAllow,
	)

	policy.AddResource("orders")
	policy.AddAction("read")

	clone := policy.Clone()

	if clone == policy {
		t.Fatal("clone should be different instance")
	}

	if clone.Name != policy.Name {
		t.Fatal("clone mismatch")
	}
}

func TestPolicyTouch(t *testing.T) {
	policy := rbac.NewPolicy(
		"tenant",
		"policy",
		"description",
		rbac.PolicyAllow,
	)

	old := policy.UpdatedAt

	time.Sleep(time.Millisecond)

	policy.Touch()

	if !policy.UpdatedAt.After(old) {
		t.Fatal("touch failed")
	}
}

func TestPolicyAllow(t *testing.T) {
	policy := rbac.NewPolicy(
		"tenant",
		"policy",
		"description",
		rbac.PolicyAllow,
	)

	if !policy.IsAllow() {
		t.Fatal("should allow")
	}

	if policy.IsDeny() {
		t.Fatal("should not deny")
	}
}

func TestPolicyDeny(t *testing.T) {
	policy := rbac.NewPolicy(
		"tenant",
		"policy",
		"description",
		rbac.PolicyDeny,
	)

	if !policy.IsDeny() {
		t.Fatal("should deny")
	}

	if policy.IsAllow() {
		t.Fatal("should not allow")
	}
}

func TestPolicyResource(t *testing.T) {
	policy := rbac.NewPolicy(
		"tenant",
		"policy",
		"description",
		rbac.PolicyAllow,
	)

	policy.AddResource("orders")

	if !policy.HasResource("orders") {
		t.Fatal("resource missing")
	}

	if policy.ResourceCount() != 1 {
		t.Fatal("invalid resource count")
	}

	policy.RemoveResource("orders")

	if policy.HasResource("orders") {
		t.Fatal("resource should be removed")
	}
}

func TestPolicyAction(t *testing.T) {
	policy := rbac.NewPolicy(
		"tenant",
		"policy",
		"description",
		rbac.PolicyAllow,
	)

	policy.AddAction("read")

	if !policy.HasAction("read") {
		t.Fatal("action missing")
	}

	if policy.ActionCount() != 1 {
		t.Fatal("invalid action count")
	}

	policy.RemoveAction("read")

	if policy.HasAction("read") {
		t.Fatal("action should be removed")
	}
}

func TestPolicyMatches(t *testing.T) {
	policy := rbac.NewPolicy(
		"tenant",
		"policy",
		"description",
		rbac.PolicyAllow,
	)

	policy.AddResource("orders")
	policy.AddAction("read")

	if !policy.Matches(
		"orders",
		"read",
	) {
		t.Fatal("match expected")
	}

	if policy.Matches(
		"users",
		"read",
	) {
		t.Fatal("unexpected match")
	}
}

func TestPolicyBelongsToTenant(t *testing.T) {
	policy := rbac.NewPolicy(
		"tenant-1",
		"policy",
		"description",
		rbac.PolicyAllow,
	)

	if !policy.BelongsToTenant("tenant-1") {
		t.Fatal("tenant mismatch")
	}
}

func TestPolicyUpdateDescription(t *testing.T) {
	policy := rbac.NewPolicy(
		"tenant",
		"policy",
		"description",
		rbac.PolicyAllow,
	)

	policy.UpdateDescription(
		"new description",
	)

	if policy.Description != "new description" {
		t.Fatal("description not updated")
	}
}

func TestPolicyClearResources(t *testing.T) {
	policy := rbac.NewPolicy(
		"tenant",
		"policy",
		"description",
		rbac.PolicyAllow,
	)

	policy.AddResource("orders")
	policy.AddResource("users")

	policy.ClearResources()

	if policy.ResourceCount() != 0 {
		t.Fatal("resources not cleared")
	}
}

func TestPolicyClearActions(t *testing.T) {
	policy := rbac.NewPolicy(
		"tenant",
		"policy",
		"description",
		rbac.PolicyAllow,
	)

	policy.AddAction("read")
	policy.AddAction("write")

	policy.ClearActions()

	if policy.ActionCount() != 0 {
		t.Fatal("actions not cleared")
	}
}
