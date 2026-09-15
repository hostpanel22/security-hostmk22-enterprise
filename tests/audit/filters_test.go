package audit_test

import (
	"testing"
	"time"

	"github.com/azezaz/HostMK-Security-Framework/core/audit"
)

func TestNewFilter(t *testing.T) {

	filter := audit.NewFilter()

	if !filter.IsEmpty() {
		t.Fatal("new filter should be empty")
	}
}

func TestFilterMatchNilEvent(t *testing.T) {

	filter := audit.NewFilter()

	if filter.Match(nil) {
		t.Fatal("nil event should not match")
	}
}

func TestFilterMatchByTenant(t *testing.T) {

	event := audit.NewEvent(
		audit.CategorySystem,
		"startup",
		audit.SeverityInfo,
		audit.StatusSuccess,
	)

	event.TenantID = "tenant-1"

	filter := audit.NewFilter()

	filter.ByTenant("tenant-1")

	if !filter.Match(event) {
		t.Fatal("expected tenant match")
	}
}

func TestFilterRejectDifferentTenant(t *testing.T) {

	event := audit.NewEvent(
		audit.CategorySystem,
		"startup",
		audit.SeverityInfo,
		audit.StatusSuccess,
	)

	event.TenantID = "tenant-1"

	filter := audit.NewFilter()

	filter.ByTenant("tenant-2")

	if filter.Match(event) {
		t.Fatal("expected tenant mismatch")
	}
}

func TestFilterMatchByActor(t *testing.T) {

	event := audit.NewEvent(
		audit.CategoryAuthentication,
		"login",
		audit.SeverityLow,
		audit.StatusSuccess,
	)

	event.ActorID = "user-1"

	filter := audit.NewFilter()

	filter.ByActor("user-1")

	if !filter.Match(event) {
		t.Fatal("expected actor match")
	}
}

func TestFilterMatchByCategory(t *testing.T) {

	event := audit.NewEvent(
		audit.CategorySecurity,
		"intrusion",
		audit.SeverityCritical,
		audit.StatusBlocked,
	)

	filter := audit.NewFilter()

	filter.ByCategory(
		audit.CategorySecurity,
	)

	if !filter.Match(event) {
		t.Fatal("expected category match")
	}
}

func TestFilterMatchBySeverity(t *testing.T) {

	event := audit.NewEvent(
		audit.CategorySecurity,
		"intrusion",
		audit.SeverityCritical,
		audit.StatusBlocked,
	)

	filter := audit.NewFilter()

	filter.BySeverity(
		audit.SeverityCritical,
	)

	if !filter.Match(event) {
		t.Fatal("expected severity match")
	}
}

func TestFilterMatchByStatus(t *testing.T) {

	event := audit.NewEvent(
		audit.CategorySystem,
		"failure",
		audit.SeverityMedium,
		audit.StatusFailed,
	)

	filter := audit.NewFilter()

	filter.ByStatus(
		audit.StatusFailed,
	)

	if !filter.Match(event) {
		t.Fatal("expected status match")
	}
}

func TestFilterMatchByResource(t *testing.T) {

	event := audit.NewEvent(
		audit.CategoryAuthorization,
		"read",
		audit.SeverityLow,
		audit.StatusSuccess,
	)

	event.Resource = "orders"

	filter := audit.NewFilter()

	filter.ByResource(
		"orders",
	)

	if !filter.Match(event) {
		t.Fatal("expected resource match")
	}
}

func TestFilterMatchByTag(t *testing.T) {

	event := audit.NewEvent(
		audit.CategorySecurity,
		"scan",
		audit.SeverityHigh,
		audit.StatusSuccess,
	)

	event.AddTag("security")

	filter := audit.NewFilter()

	filter.ByTag("security")

	if !filter.Match(event) {
		t.Fatal("expected tag match")
	}
}

func TestFilterRejectMissingTag(t *testing.T) {

	event := audit.NewEvent(
		audit.CategorySecurity,
		"scan",
		audit.SeverityHigh,
		audit.StatusSuccess,
	)

	filter := audit.NewFilter()

	filter.ByTag("security")

	if filter.Match(event) {
		t.Fatal("expected tag mismatch")
	}
}

func TestFilterBetweenDates(t *testing.T) {

	event := audit.NewEvent(
		audit.CategorySystem,
		"startup",
		audit.SeverityInfo,
		audit.StatusSuccess,
	)

	now := time.Now().UTC()

	event.Timestamp = now

	from := now.Add(-time.Hour)
	to := now.Add(time.Hour)

	filter := audit.NewFilter()

	filter.Between(
		from,
		to,
	)

	if !filter.Match(event) {
		t.Fatal("expected date range match")
	}
}

func TestFilterOutsideDateRange(t *testing.T) {

	event := audit.NewEvent(
		audit.CategorySystem,
		"startup",
		audit.SeverityInfo,
		audit.StatusSuccess,
	)

	now := time.Now().UTC()

	event.Timestamp = now

	from := now.Add(time.Hour)
	to := now.Add(2 * time.Hour)

	filter := audit.NewFilter()

	filter.Between(
		from,
		to,
	)

	if filter.Match(event) {
		t.Fatal("expected date range mismatch")
	}
}

func TestFilterMultipleCriteria(t *testing.T) {

	event := audit.NewEvent(
		audit.CategoryAuthentication,
		"user_login",
		audit.SeverityLow,
		audit.StatusSuccess,
	)

	event.TenantID = "tenant-1"
	event.ActorID = "user-1"

	filter := audit.NewFilter()

	filter.
		ByTenant("tenant-1").
		ByActor("user-1").
		ByCategory(audit.CategoryAuthentication).
		ByStatus(audit.StatusSuccess)

	if !filter.Match(event) {
		t.Fatal("expected all criteria to match")
	}
}

func TestFilterIsEmptyFalse(t *testing.T) {

	filter := audit.NewFilter()

	filter.ByTenant("tenant-1")

	if filter.IsEmpty() {
		t.Fatal("filter should not be empty")
	}
}

func TestFilterBuilderChain(t *testing.T) {

	filter := audit.NewFilter()

	result := filter.
		ByTenant("tenant-1").
		ByActor("user-1").
		ByTag("security")

	if result == nil {
		t.Fatal("builder should return filter")
	}
}
