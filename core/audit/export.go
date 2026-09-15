package audit

import (
	"context"
	"encoding/csv"
	"encoding/json"
	"io"
	"time"
)

type ExportFormat string

const (
	FormatJSON   ExportFormat = "json"
	FormatCSV    ExportFormat = "csv"
	FormatNDJSON ExportFormat = "ndjson"
)

type ExportOptions struct {
	Format ExportFormat

	Filter Filter

	Pretty bool
}

type Exporter struct {
	store Store
}

func NewExporter(store Store) *Exporter {
	return &Exporter{
		store: store,
	}
}

func (e *Exporter) Export(
	ctx context.Context,
	writer io.Writer,
	options ExportOptions,
) error {

	if writer == nil {
		return ErrExportFailed
	}

	if e.store == nil {
		return ErrStoreUnavailable
	}

	events, err := e.store.List(
		ctx,
		options.Filter,
	)

	if err != nil {
		return err
	}

	switch options.Format {

	case FormatJSON:
		return e.exportJSON(
			writer,
			events,
			options.Pretty,
		)

	case FormatCSV:
		return e.exportCSV(
			writer,
			events,
		)

	case FormatNDJSON:
		return e.exportNDJSON(
			writer,
			events,
		)

	default:
		return ErrUnsupportedFormat
	}
}

func (e *Exporter) exportJSON(
	writer io.Writer,
	events []*Event,
	pretty bool,
) error {

	var (
		data []byte
		err  error
	)

	if pretty {

		data, err = json.MarshalIndent(
			events,
			"",
			"  ",
		)

	} else {

		data, err = json.Marshal(
			events,
		)
	}

	if err != nil {
		return err
	}

	_, err = writer.Write(
		data,
	)

	return err
}

func (e *Exporter) exportNDJSON(
	writer io.Writer,
	events []*Event,
) error {

	encoder := json.NewEncoder(
		writer,
	)

	for _, event := range events {

		if err := encoder.Encode(
			event,
		); err != nil {
			return err
		}
	}

	return nil
}

func (e *Exporter) exportCSV(
	writer io.Writer,
	events []*Event,
) error {

	csvWriter := csv.NewWriter(
		writer,
	)

	header := []string{
		"id",
		"timestamp",
		"category",
		"action",
		"status",
		"severity",
		"tenant_id",
		"actor_id",
		"actor_type",
		"resource",
		"resource_id",
		"request_id",
		"session_id",
		"ip_address",
		"user_agent",
		"message",
	}

	if err := csvWriter.Write(
		header,
	); err != nil {
		return err
	}

	for _, event := range events {

		row := []string{
			event.ID,
			event.Timestamp.Format(time.RFC3339),
			string(event.Category),
			event.Action,
			string(event.Status),
			string(event.Severity),
			event.TenantID,
			event.ActorID,
			event.ActorType,
			event.Resource,
			event.ResourceID,
			event.RequestID,
			event.SessionID,
			event.IPAddress,
			event.UserAgent,
			event.Message,
		}

		if err := csvWriter.Write(
			row,
		); err != nil {
			return err
		}
	}

	csvWriter.Flush()

	return csvWriter.Error()
}

func (e *Exporter) ExportJSON(
	ctx context.Context,
	writer io.Writer,
	filter Filter,
) error {

	return e.Export(
		ctx,
		writer,
		ExportOptions{
			Format: FormatJSON,
			Filter: filter,
			Pretty: true,
		},
	)
}

func (e *Exporter) ExportCSV(
	ctx context.Context,
	writer io.Writer,
	filter Filter,
) error {

	return e.Export(
		ctx,
		writer,
		ExportOptions{
			Format: FormatCSV,
			Filter: filter,
		},
	)
}

func (e *Exporter) ExportNDJSON(
	ctx context.Context,
	writer io.Writer,
	filter Filter,
) error {

	return e.Export(
		ctx,
		writer,
		ExportOptions{
			Format: FormatNDJSON,
			Filter: filter,
		},
	)
}
