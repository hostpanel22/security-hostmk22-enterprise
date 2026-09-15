package audit_test

import (
	"testing"

	"github.com/azezaz/HostMK-Security-Framework/core/audit"
)

func TestNewEvent(t *testing.T) {

	event := audit.NewEvent(
		audit.CategoryAuthentication,
		"user_login",
		audit.SeverityLow,
		audit.StatusSuccess,
	)

	if event == nil {
		t.Fatal("expected event")
	}

	if event.ID == "" {
		t.Fatal("expected event id")
	}

	if event.Action != "user_login" {
		t.Fatal("invalid action")
	}

	if event.Category != audit.CategoryAuthentication {
		t.Fatal("invalid category")
	}

	if event.Status != audit.StatusSuccess {
		t.Fatal("invalid status")
	}

	if event.Severity != audit.SeverityLow {
		t.Fatal("invalid severity")
	}
}

func TestEventValidateSuccess(t *testing.T) {

	event := audit.NewEvent(
		audit.CategorySystem,
		"startup",
		audit.SeverityInfo,
		audit.StatusSuccess,
	)

	if err := event.Validate(); err != nil {
		t.Fatalf(
			"unexpected error: %v",
			err,
		)
	}
}

func TestEventValidateMissingID(t *testing.T) {

	event := audit.NewEvent(
		audit.CategorySystem,
		"startup",
		audit.SeverityInfo,
		audit.StatusSuccess,
	)

	event.ID = ""

	err := event.Validate()

	if err != audit.ErrInvalidEventID {
		t.Fatalf(
			"expected ErrInvalidEventID got %v",
			err,
		)
	}
}

func TestEventValidateMissingAction(t *testing.T) {

	event := audit.NewEvent(
		audit.CategorySystem,
		"startup",
		audit.SeverityInfo,
		audit.StatusSuccess,
	)

	event.Action = ""

	err := event.Validate()

	if err != audit.ErrInvalidEventAction {
		t.Fatalf(
			"expected ErrInvalidEventAction got %v",
			err,
		)
	}
}

func TestEventValidateMissingCategory(t *testing.T) {

	event := audit.NewEvent(
		audit.CategorySystem,
		"startup",
		audit.SeverityInfo,
		audit.StatusSuccess,
	)

	event.Category = ""

	err := event.Validate()

	if err != audit.ErrInvalidEventCategory {
		t.Fatalf(
			"expected ErrInvalidEventCategory got %v",
			err,
		)
	}
}

func TestEventValidateMissingSeverity(t *testing.T) {

	event := audit.NewEvent(
		audit.CategorySystem,
		"startup",
		audit.SeverityInfo,
		audit.StatusSuccess,
	)

	event.Severity = ""

	err := event.Validate()

	if err != audit.ErrInvalidEventSeverity {
		t.Fatalf(
			"expected ErrInvalidEventSeverity got %v",
			err,
		)
	}
}

func TestEventValidateMissingStatus(t *testing.T) {

	event := audit.NewEvent(
		audit.CategorySystem,
		"startup",
		audit.SeverityInfo,
		audit.StatusSuccess,
	)

	event.Status = ""

	err := event.Validate()

	if err != audit.ErrInvalidEventStatus {
		t.Fatalf(
			"expected ErrInvalidEventStatus got %v",
			err,
		)
	}
}

func TestAddTag(t *testing.T) {

	event := audit.NewEvent(
		audit.CategorySecurity,
		"violation",
		audit.SeverityCritical,
		audit.StatusBlocked,
	)

	event.AddTag("security")

	if len(event.Tags) != 1 {
		t.Fatal("expected one tag")
	}

	if event.Tags[0] != "security" {
		t.Fatal("invalid tag")
	}
}

func TestSetMetadata(t *testing.T) {

	event := audit.NewEvent(
		audit.CategoryAPI,
		"request",
		audit.SeverityLow,
		audit.StatusSuccess,
	)

	event.SetMetadata(
		"endpoint",
		"/api/v1/orders",
	)

	value, exists := event.Metadata["endpoint"]

	if !exists {
		t.Fatal("metadata missing")
	}

	if value != "/api/v1/orders" {
		t.Fatal("invalid metadata value")
	}
}

func TestClone(t *testing.T) {

	event := audit.NewEvent(
		audit.CategoryFraud,
		"duplicate_order",
		audit.SeverityHigh,
		audit.StatusBlocked,
	)

	event.AddTag("fraud")

	event.SetMetadata(
		"score",
		95,
	)

	clone := event.Clone()

	if clone == nil {
		t.Fatal("expected clone")
	}

	if clone.ID != event.ID {
		t.Fatal("clone id mismatch")
	}

	if len(clone.Tags) != len(event.Tags) {
		t.Fatal("clone tags mismatch")
	}

	if clone.Metadata["score"] != 95 {
		t.Fatal("clone metadata mismatch")
	}
}

func TestIsCritical(t *testing.T) {

	event := audit.NewEvent(
		audit.CategorySecurity,
		"intrusion",
		audit.SeverityCritical,
		audit.StatusBlocked,
	)

	if !event.IsCritical() {
		t.Fatal("expected critical event")
	}
}

func TestIsHighRisk(t *testing.T) {

	event := audit.NewEvent(
		audit.CategoryFraud,
		"bot_detected",
		audit.SeverityHigh,
		audit.StatusBlocked,
	)

	if !event.IsHighRisk() {
		t.Fatal("expected high risk event")
	}
}

func TestIsFailure(t *testing.T) {

	event := audit.NewEvent(
		audit.CategoryAuthentication,
		"login",
		audit.SeverityMedium,
		audit.StatusFailed,
	)

	if !event.IsFailure() {
		t.Fatal("expected failure")
	}
}

func TestIsSuccess(t *testing.T) {

	event := audit.NewEvent(
		audit.CategoryAuthentication,
		"login",
		audit.SeverityLow,
		audit.StatusSuccess,
	)

	if !event.IsSuccess() {
		t.Fatal("expected success")
	}
}
