package fraud_test

import (
	"testing"
	"time"

	"github.com/azezaz/HostMK-Security-Framework/core/fraud"
)

func TestNewRule(t *testing.T) {
	rule := fraud.NewRule(
		"High Login Rate",
		"Detect excessive logins",
		"login",
		100,
		10,
	)

	if rule == nil {
		t.Fatal("expected rule")
	}

	if rule.ID == "" {
		t.Fatal("expected generated id")
	}

	if rule.Name != "High Login Rate" {
		t.Fatalf("unexpected name: %s", rule.Name)
	}

	if rule.SignalType != "login" {
		t.Fatalf("unexpected signal type: %s", rule.SignalType)
	}

	if rule.Threshold != 100 {
		t.Fatalf("unexpected threshold: %f", rule.Threshold)
	}

	if rule.Weight != 10 {
		t.Fatalf("unexpected weight: %d", rule.Weight)
	}

	if !rule.Enabled {
		t.Fatal("new rule should be enabled")
	}

	if rule.CreatedAt.IsZero() {
		t.Fatal("expected created_at")
	}
}

func TestRuleValidateSuccess(t *testing.T) {
	rule := fraud.NewRule(
		"Test",
		"Test rule",
		"login",
		10,
		5,
	)

	if err := rule.Validate(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestRuleValidateNil(t *testing.T) {
	var rule *fraud.Rule

	err := rule.Validate()

	if err != fraud.ErrInvalidRule {
		t.Fatalf(
			"expected %v got %v",
			fraud.ErrInvalidRule,
			err,
		)
	}
}

func TestRuleValidateEmptyName(t *testing.T) {
	rule := fraud.NewRule(
		"",
		"desc",
		"login",
		10,
		1,
	)

	err := rule.Validate()

	if err != fraud.ErrEmptyRuleName {
		t.Fatalf(
			"expected %v got %v",
			fraud.ErrEmptyRuleName,
			err,
		)
	}
}

func TestRuleValidateEmptySignalType(t *testing.T) {
	rule := fraud.NewRule(
		"rule",
		"desc",
		"",
		10,
		1,
	)

	err := rule.Validate()

	if err != fraud.ErrEmptySignalType {
		t.Fatalf(
			"expected %v got %v",
			fraud.ErrEmptySignalType,
			err,
		)
	}
}

func TestRuleValidateInvalidThreshold(t *testing.T) {
	rule := fraud.NewRule(
		"rule",
		"desc",
		"login",
		0,
		1,
	)

	err := rule.Validate()

	if err != fraud.ErrInvalidThreshold {
		t.Fatalf(
			"expected %v got %v",
			fraud.ErrInvalidThreshold,
			err,
		)
	}
}

func TestRuleValidateInvalidWeight(t *testing.T) {
	rule := fraud.NewRule(
		"rule",
		"desc",
		"login",
		10,
		0,
	)

	err := rule.Validate()

	if err != fraud.ErrInvalidWeight {
		t.Fatalf(
			"expected %v got %v",
			fraud.ErrInvalidWeight,
			err,
		)
	}
}

func TestRuleEnable(t *testing.T) {
	rule := fraud.NewRule(
		"rule",
		"desc",
		"login",
		10,
		1,
	)

	rule.Disable()

	if rule.Enabled {
		t.Fatal("rule should be disabled")
	}

	rule.Enable()

	if !rule.Enabled {
		t.Fatal("rule should be enabled")
	}
}

func TestRuleDisable(t *testing.T) {
	rule := fraud.NewRule(
		"rule",
		"desc",
		"login",
		10,
		1,
	)

	rule.Disable()

	if rule.Enabled {
		t.Fatal("rule should be disabled")
	}
}

func TestRuleTouch(t *testing.T) {
	rule := fraud.NewRule(
		"rule",
		"desc",
		"login",
		10,
		1,
	)

	old := rule.UpdatedAt

	time.Sleep(time.Millisecond)

	rule.Touch()

	if !rule.UpdatedAt.After(old) {
		t.Fatal("updated_at should change")
	}
}

func TestRuleMatchSuccess(t *testing.T) {
	rule := fraud.NewRule(
		"rule",
		"desc",
		"login",
		100,
		1,
	)

	signal := fraud.NewSignal(
		"tenant",
		"user",
		"login",
		150,
		"api",
	)

	if !rule.Match(signal) {
		t.Fatal("expected rule match")
	}
}

func TestRuleMatchBelowThreshold(t *testing.T) {
	rule := fraud.NewRule(
		"rule",
		"desc",
		"login",
		100,
		1,
	)

	signal := fraud.NewSignal(
		"tenant",
		"user",
		"login",
		50,
		"api",
	)

	if rule.Match(signal) {
		t.Fatal("match should fail")
	}
}

func TestRuleMatchWrongSignalType(t *testing.T) {
	rule := fraud.NewRule(
		"rule",
		"desc",
		"login",
		100,
		1,
	)

	signal := fraud.NewSignal(
		"tenant",
		"user",
		"payment",
		150,
		"api",
	)

	if rule.Match(signal) {
		t.Fatal("match should fail")
	}
}

func TestRuleMatchDisabled(t *testing.T) {
	rule := fraud.NewRule(
		"rule",
		"desc",
		"login",
		100,
		1,
	)

	rule.Disable()

	signal := fraud.NewSignal(
		"tenant",
		"user",
		"login",
		150,
		"api",
	)

	if rule.Match(signal) {
		t.Fatal("disabled rule must not match")
	}
}

func TestRuleMatchNilSignal(t *testing.T) {
	rule := fraud.NewRule(
		"rule",
		"desc",
		"login",
		100,
		1,
	)

	if rule.Match(nil) {
		t.Fatal("nil signal must not match")
	}
}

func TestRuleClone(t *testing.T) {
	rule := fraud.NewRule(
		"rule",
		"description",
		"login",
		100,
		10,
	)

	clone := rule.Clone()

	if clone == nil {
		t.Fatal("expected clone")
	}

	if clone == rule {
		t.Fatal("clone must be different pointer")
	}

	if clone.ID != rule.ID {
		t.Fatal("id mismatch")
	}

	if clone.Name != rule.Name {
		t.Fatal("name mismatch")
	}

	if clone.SignalType != rule.SignalType {
		t.Fatal("signal type mismatch")
	}

	if clone.Weight != rule.Weight {
		t.Fatal("weight mismatch")
	}
}

func TestRuleCloneNil(t *testing.T) {
	var rule *fraud.Rule

	clone := rule.Clone()

	if clone != nil {
		t.Fatal("expected nil clone")
	}
}
