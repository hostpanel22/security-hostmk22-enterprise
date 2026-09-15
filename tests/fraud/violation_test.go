package fraud_test

import (
	"testing"

	"github.com/azezaz/HostMK-Security-Framework/core/fraud"
)

func TestNewViolation(t *testing.T) {

	violation := fraud.NewViolation(
		"tenant-1",
		"user-1",
		"signal-1",
		"rule-1",
		"suspicious activity",
		95,
		fraud.RiskCritical,
	)

	if violation == nil {
		t.Fatal("expected violation")
	}

	if violation.ID == "" {
		t.Fatal("expected id")
	}

	if violation.TenantID != "tenant-1" {
		t.Fatal("invalid tenant id")
	}

	if violation.UserID != "user-1" {
		t.Fatal("invalid user id")
	}

	if violation.SignalID != "signal-1" {
		t.Fatal("invalid signal id")
	}

	if violation.RuleID != "rule-1" {
		t.Fatal("invalid rule id")
	}

	if violation.Score != 95 {
		t.Fatal("invalid score")
	}

	if violation.Level != fraud.RiskCritical {
		t.Fatal("invalid level")
	}
}

func TestViolationValidateSuccess(t *testing.T) {

	violation := fraud.NewViolation(
		"tenant-1",
		"user-1",
		"signal-1",
		"rule-1",
		"suspicious activity",
		90,
		fraud.RiskHigh,
	)

	if err := violation.Validate(); err != nil {
		t.Fatalf(
			"unexpected error: %v",
			err,
		)
	}
}

func TestViolationValidateNil(t *testing.T) {

	var violation *fraud.Violation

	err := violation.Validate()

	if err == nil {
		t.Fatal("expected error")
	}
}

func TestViolationValidateEmptyTenant(t *testing.T) {

	violation := fraud.NewViolation(
		"",
		"user-1",
		"signal-1",
		"rule-1",
		"reason",
		50,
		fraud.RiskMedium,
	)

	err := violation.Validate()

	if err != fraud.ErrEmptyTenantID {
		t.Fatalf(
			"expected ErrEmptyTenantID got %v",
			err,
		)
	}
}

func TestViolationValidateEmptyUser(t *testing.T) {

	violation := fraud.NewViolation(
		"tenant-1",
		"",
		"signal-1",
		"rule-1",
		"reason",
		50,
		fraud.RiskMedium,
	)

	err := violation.Validate()

	if err != fraud.ErrEmptyUserID {
		t.Fatalf(
			"expected ErrEmptyUserID got %v",
			err,
		)
	}
}

func TestViolationValidateEmptyReason(t *testing.T) {

	violation := fraud.NewViolation(
		"tenant-1",
		"user-1",
		"signal-1",
		"rule-1",
		"",
		50,
		fraud.RiskMedium,
	)

	err := violation.Validate()

	if err != fraud.ErrInvalidSignal {
		t.Fatalf(
			"expected ErrInvalidSignal got %v",
			err,
		)
	}
}

func TestViolationIsCritical(t *testing.T) {

	violation := fraud.NewViolation(
		"tenant",
		"user",
		"signal",
		"rule",
		"critical",
		100,
		fraud.RiskCritical,
	)

	if !violation.IsCritical() {
		t.Fatal("expected critical violation")
	}
}

func TestViolationIsHigh(t *testing.T) {

	violation := fraud.NewViolation(
		"tenant",
		"user",
		"signal",
		"rule",
		"high",
		80,
		fraud.RiskHigh,
	)

	if !violation.IsHigh() {
		t.Fatal("expected high violation")
	}
}

func TestViolationRequiresBlockHigh(t *testing.T) {

	violation := fraud.NewViolation(
		"tenant",
		"user",
		"signal",
		"rule",
		"high",
		80,
		fraud.RiskHigh,
	)

	if !violation.RequiresBlock() {
		t.Fatal("expected block")
	}
}

func TestViolationRequiresBlockCritical(t *testing.T) {

	violation := fraud.NewViolation(
		"tenant",
		"user",
		"signal",
		"rule",
		"critical",
		100,
		fraud.RiskCritical,
	)

	if !violation.RequiresBlock() {
		t.Fatal("expected block")
	}
}

func TestViolationRequiresBlockLow(t *testing.T) {

	violation := fraud.NewViolation(
		"tenant",
		"user",
		"signal",
		"rule",
		"low",
		10,
		fraud.RiskLow,
	)

	if violation.RequiresBlock() {
		t.Fatal("did not expect block")
	}
}

func TestViolationRequiresBlockMedium(t *testing.T) {

	violation := fraud.NewViolation(
		"tenant",
		"user",
		"signal",
		"rule",
		"medium",
		45,
		fraud.RiskMedium,
	)

	if violation.RequiresBlock() {
		t.Fatal("did not expect block")
	}
}
