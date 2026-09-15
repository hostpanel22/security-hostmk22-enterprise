package audit_test

import (
	"bytes"
	"context"
	"strings"
	"testing"

	"github.com/azezaz/HostMK-Security-Framework/core/audit"
)

func TestNewExporter(t *testing.T) {

	store := audit.NewMemoryStore()

	exporter := audit.NewExporter(
		store,
	)

	if exporter == nil {
		t.Fatal("exporter should not be nil")
	}
}

func TestExportJSON(t *testing.T) {

	store := audit.NewMemoryStore()

	event := audit.NewEvent(
		audit.CategorySystem,
		"startup",
		audit.SeverityInfo,
		audit.StatusSuccess,
	)

	if err := store.Save(
		context.Background(),
		event,
	); err != nil {
		t.Fatal(err)
	}

	exporter := audit.NewExporter(
		store,
	)

	var buffer bytes.Buffer

	err := exporter.ExportJSON(
		context.Background(),
		&buffer,
		audit.NewFilter(),
	)

	if err != nil {
		t.Fatal(err)
	}

	if buffer.Len() == 0 {
		t.Fatal("expected json output")
	}

	if !strings.Contains(
		buffer.String(),
		"startup",
	) {
		t.Fatal("expected action in json")
	}
}

func TestExportCSV(t *testing.T) {

	store := audit.NewMemoryStore()

	event := audit.NewEvent(
		audit.CategorySystem,
		"startup",
		audit.SeverityInfo,
		audit.StatusSuccess,
	)

	if err := store.Save(
		context.Background(),
		event,
	); err != nil {
		t.Fatal(err)
	}

	exporter := audit.NewExporter(
		store,
	)

	var buffer bytes.Buffer

	err := exporter.ExportCSV(
		context.Background(),
		&buffer,
		audit.NewFilter(),
	)

	if err != nil {
		t.Fatal(err)
	}

	output := buffer.String()

	if !strings.Contains(
		output,
		"id",
	) {
		t.Fatal("csv header missing")
	}

	if !strings.Contains(
		output,
		"startup",
	) {
		t.Fatal("csv content missing")
	}
}

func TestExportNDJSON(t *testing.T) {

	store := audit.NewMemoryStore()

	event := audit.NewEvent(
		audit.CategoryAPI,
		"request",
		audit.SeverityLow,
		audit.StatusSuccess,
	)

	if err := store.Save(
		context.Background(),
		event,
	); err != nil {
		t.Fatal(err)
	}

	exporter := audit.NewExporter(
		store,
	)

	var buffer bytes.Buffer

	err := exporter.ExportNDJSON(
		context.Background(),
		&buffer,
		audit.NewFilter(),
	)

	if err != nil {
		t.Fatal(err)
	}

	if buffer.Len() == 0 {
		t.Fatal("expected ndjson output")
	}

	if !strings.Contains(
		buffer.String(),
		"request",
	) {
		t.Fatal("expected request action")
	}
}

func TestExportNilWriter(t *testing.T) {

	store := audit.NewMemoryStore()

	exporter := audit.NewExporter(
		store,
	)

	err := exporter.Export(
		context.Background(),
		nil,
		audit.ExportOptions{
			Format: audit.FormatJSON,
		},
	)

	if err != audit.ErrExportFailed {
		t.Fatalf(
			"expected ErrExportFailed got %v",
			err,
		)
	}
}

func TestExportNilStore(t *testing.T) {

	exporter := audit.NewExporter(
		nil,
	)

	var buffer bytes.Buffer

	err := exporter.Export(
		context.Background(),
		&buffer,
		audit.ExportOptions{
			Format: audit.FormatJSON,
		},
	)

	if err != audit.ErrStoreUnavailable {
		t.Fatalf(
			"expected ErrStoreUnavailable got %v",
			err,
		)
	}
}

func TestExportUnsupportedFormat(t *testing.T) {

	store := audit.NewMemoryStore()

	exporter := audit.NewExporter(
		store,
	)

	var buffer bytes.Buffer

	err := exporter.Export(
		context.Background(),
		&buffer,
		audit.ExportOptions{
			Format: audit.ExportFormat(
				"xml",
			),
		},
	)

	if err != audit.ErrUnsupportedFormat {
		t.Fatalf(
			"expected ErrUnsupportedFormat got %v",
			err,
		)
	}
}

func TestExportPrettyJSON(t *testing.T) {

	store := audit.NewMemoryStore()

	event := audit.NewEvent(
		audit.CategorySystem,
		"pretty",
		audit.SeverityInfo,
		audit.StatusSuccess,
	)

	if err := store.Save(
		context.Background(),
		event,
	); err != nil {
		t.Fatal(err)
	}

	exporter := audit.NewExporter(
		store,
	)

	var buffer bytes.Buffer

	err := exporter.Export(
		context.Background(),
		&buffer,
		audit.ExportOptions{
			Format: audit.FormatJSON,
			Pretty: true,
			Filter: audit.NewFilter(),
		},
	)

	if err != nil {
		t.Fatal(err)
	}

	if !strings.Contains(
		buffer.String(),
		"\n",
	) {
		t.Fatal("expected formatted json")
	}
}

func TestExportWithFilter(t *testing.T) {

	store := audit.NewMemoryStore()

	event := audit.NewEvent(
		audit.CategorySecurity,
		"intrusion",
		audit.SeverityCritical,
		audit.StatusBlocked,
	)

	event.TenantID = "tenant-1"

	if err := store.Save(
		context.Background(),
		event,
	); err != nil {
		t.Fatal(err)
	}

	exporter := audit.NewExporter(
		store,
	)

	filter := audit.NewFilter()

	filter.ByTenant(
		"tenant-1",
	)

	var buffer bytes.Buffer

	err := exporter.ExportJSON(
		context.Background(),
		&buffer,
		filter,
	)

	if err != nil {
		t.Fatal(err)
	}

	if !strings.Contains(
		buffer.String(),
		"intrusion",
	) {
		t.Fatal("expected filtered event")
	}
}
