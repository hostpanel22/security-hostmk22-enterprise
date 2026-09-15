package fraud_test

import (
	"context"
	"testing"

	"github.com/azezaz/HostMK-Security-Framework/core/fraud"
)

func TestNewMemoryStore(t *testing.T) {
	store := fraud.NewMemoryStore()

	if store == nil {
		t.Fatal("store should not be nil")
	}

	if store.SignalCount() != 0 {
		t.Fatal("signal count should be zero")
	}

	if store.ViolationCount() != 0 {
		t.Fatal("violation count should be zero")
	}
}

func TestSaveSignal(t *testing.T) {
	store := fraud.NewMemoryStore()

	signal := fraud.NewSignal(
		"tenant-1",
		"user-1",
		"login",
		10,
		"test",
	)

	err := store.SaveSignal(
		context.Background(),
		signal,
	)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if store.SignalCount() != 1 {
		t.Fatal("signal should be stored")
	}
}

func TestSaveViolation(t *testing.T) {
	store := fraud.NewMemoryStore()

	violation := fraud.NewViolation(
		"tenant-1",
		"user-1",
		"signal-1",
		"rule-1",
		"test",
		80,
		fraud.RiskHigh,
	)

	err := store.SaveViolation(
		context.Background(),
		violation,
	)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if store.ViolationCount() != 1 {
		t.Fatal("violation should be stored")
	}
}

func TestListSignals(t *testing.T) {
	store := fraud.NewMemoryStore()

	signal := fraud.NewSignal(
		"tenant",
		"user",
		"login",
		10,
		"test",
	)

	_ = store.SaveSignal(
		context.Background(),
		signal,
	)

	signals, err := store.ListSignals(
		context.Background(),
	)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(signals) != 1 {
		t.Fatal("expected one signal")
	}
}

func TestListViolations(t *testing.T) {
	store := fraud.NewMemoryStore()

	violation := fraud.NewViolation(
		"tenant",
		"user",
		"signal",
		"rule",
		"reason",
		90,
		fraud.RiskCritical,
	)

	_ = store.SaveViolation(
		context.Background(),
		violation,
	)

	violations, err := store.ListViolations(
		context.Background(),
	)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(violations) != 1 {
		t.Fatal("expected one violation")
	}
}

func TestClearStore(t *testing.T) {
	store := fraud.NewMemoryStore()

	_ = store.SaveSignal(
		context.Background(),
		fraud.NewSignal(
			"tenant",
			"user",
			"login",
			10,
			"test",
		),
	)

	err := store.Clear(
		context.Background(),
	)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if store.SignalCount() != 0 {
		t.Fatal("signals should be cleared")
	}
}
