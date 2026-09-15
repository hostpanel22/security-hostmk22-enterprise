package audit_test

import (
	"context"
	"testing"
	"time"

	"github.com/azezaz/HostMK-Security-Framework/core/audit"
)

func TestMemoryStoreSave(t *testing.T) {

	store := audit.NewMemoryStore()

	event := audit.NewEvent(
		audit.CategorySystem,
		"application_started",
		audit.SeverityInfo,
		audit.StatusSuccess,
	)

	err := store.Save(
		context.Background(),
		event,
	)

	if err != nil {
		t.Fatalf(
			"unexpected error: %v",
			err,
		)
	}
}

func TestMemoryStoreSaveNilEvent(t *testing.T) {

	store := audit.NewMemoryStore()

	err := store.Save(
		context.Background(),
		nil,
	)

	if err != audit.ErrNilEvent {
		t.Fatalf(
			"expected ErrNilEvent got %v",
			err,
		)
	}
}

func TestMemoryStoreFindByID(t *testing.T) {

	store := audit.NewMemoryStore()

	event := audit.NewEvent(
		audit.CategoryAuthentication,
		"user_login",
		audit.SeverityLow,
		audit.StatusSuccess,
	)

	err := store.Save(
		context.Background(),
		event,
	)

	if err != nil {
		t.Fatal(err)
	}

	found, err := store.FindByID(
		context.Background(),
		event.ID,
	)

	if err != nil {
		t.Fatal(err)
	}

	if found == nil {
		t.Fatal("expected event")
	}

	if found.ID != event.ID {
		t.Fatal("event mismatch")
	}
}

func TestMemoryStoreFindInvalidID(t *testing.T) {

	store := audit.NewMemoryStore()

	_, err := store.FindByID(
		context.Background(),
		"",
	)

	if err != audit.ErrInvalidID {
		t.Fatalf(
			"expected ErrInvalidID got %v",
			err,
		)
	}
}

func TestMemoryStoreList(t *testing.T) {

	store := audit.NewMemoryStore()

	for i := 0; i < 5; i++ {

		event := audit.NewEvent(
			audit.CategoryAPI,
			"request",
			audit.SeverityLow,
			audit.StatusSuccess,
		)

		err := store.Save(
			context.Background(),
			event,
		)

		if err != nil {
			t.Fatal(err)
		}
	}

	events, err := store.List(
		context.Background(),
		audit.NewFilter(),
	)

	if err != nil {
		t.Fatal(err)
	}

	if len(events) != 5 {
		t.Fatalf(
			"expected 5 events got %d",
			len(events),
		)
	}
}

func TestMemoryStoreCount(t *testing.T) {

	store := audit.NewMemoryStore()

	for i := 0; i < 3; i++ {

		event := audit.NewEvent(
			audit.CategorySystem,
			"heartbeat",
			audit.SeverityInfo,
			audit.StatusSuccess,
		)

		if err := store.Save(
			context.Background(),
			event,
		); err != nil {
			t.Fatal(err)
		}
	}

	count, err := store.Count(
		context.Background(),
		audit.NewFilter(),
	)

	if err != nil {
		t.Fatal(err)
	}

	if count != 3 {
		t.Fatalf(
			"expected 3 got %d",
			count,
		)
	}
}

func TestMemoryStoreDelete(t *testing.T) {

	store := audit.NewMemoryStore()

	event := audit.NewEvent(
		audit.CategorySystem,
		"delete_test",
		audit.SeverityInfo,
		audit.StatusSuccess,
	)

	if err := store.Save(
		context.Background(),
		event,
	); err != nil {
		t.Fatal(err)
	}

	err := store.Delete(
		context.Background(),
		event.ID,
	)

	if err != nil {
		t.Fatal(err)
	}

	if store.Exists(event.ID) {
		t.Fatal("event should be deleted")
	}
}

func TestMemoryStoreDeleteOlderThan(t *testing.T) {

	store := audit.NewMemoryStore()

	oldEvent := audit.NewEvent(
		audit.CategorySystem,
		"old",
		audit.SeverityInfo,
		audit.StatusSuccess,
	)

	oldEvent.Timestamp = time.Now().
		Add(-48 * time.Hour)

	if err := store.Save(
		context.Background(),
		oldEvent,
	); err != nil {
		t.Fatal(err)
	}

	deleted, err := store.DeleteOlderThan(
		context.Background(),
		time.Now().Add(-24*time.Hour),
	)

	if err != nil {
		t.Fatal(err)
	}

	if deleted != 1 {
		t.Fatalf(
			"expected 1 deleted got %d",
			deleted,
		)
	}
}

func TestMemoryStoreExists(t *testing.T) {

	store := audit.NewMemoryStore()

	event := audit.NewEvent(
		audit.CategorySystem,
		"exists",
		audit.SeverityInfo,
		audit.StatusSuccess,
	)

	if err := store.Save(
		context.Background(),
		event,
	); err != nil {
		t.Fatal(err)
	}

	if !store.Exists(event.ID) {
		t.Fatal("expected event to exist")
	}
}

func TestMemoryStoreClear(t *testing.T) {

	store := audit.NewMemoryStore()

	event := audit.NewEvent(
		audit.CategorySystem,
		"clear",
		audit.SeverityInfo,
		audit.StatusSuccess,
	)

	if err := store.Save(
		context.Background(),
		event,
	); err != nil {
		t.Fatal(err)
	}

	store.Clear()

	if store.Size() != 0 {
		t.Fatal("store should be empty")
	}
}

func TestMemoryStoreSize(t *testing.T) {

	store := audit.NewMemoryStore()

	for i := 0; i < 10; i++ {

		event := audit.NewEvent(
			audit.CategorySystem,
			"size_test",
			audit.SeverityInfo,
			audit.StatusSuccess,
		)

		if err := store.Save(
			context.Background(),
			event,
		); err != nil {
			t.Fatal(err)
		}
	}

	if store.Size() != 10 {
		t.Fatalf(
			"expected 10 got %d",
			store.Size(),
		)
	}
}
