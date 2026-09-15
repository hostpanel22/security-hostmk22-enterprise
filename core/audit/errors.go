package audit

import "errors"

var (
	// Generic errors
	ErrNilEvent  = errors.New("audit: event is nil")
	ErrNilStore  = errors.New("audit: store is nil")
	ErrNilLogger = errors.New("audit: logger is nil")
	ErrInvalidID = errors.New("audit: invalid id")

	// Event validation
	ErrInvalidEventID       = errors.New("audit: invalid event id")
	ErrInvalidEventAction   = errors.New("audit: invalid event action")
	ErrInvalidEventCategory = errors.New("audit: invalid event category")
	ErrInvalidEventSeverity = errors.New("audit: invalid event severity")
	ErrInvalidEventStatus   = errors.New("audit: invalid event status")

	// Store operations
	ErrStoreUnavailable = errors.New("audit: store unavailable")
	ErrStoreWriteFailed = errors.New("audit: store write failed")
	ErrStoreReadFailed  = errors.New("audit: store read failed")

	// Export operations
	ErrExportFailed      = errors.New("audit: export failed")
	ErrUnsupportedFormat = errors.New("audit: unsupported export format")

	// Retention operations
	ErrInvalidRetentionPeriod = errors.New("audit: invalid retention period")
	ErrRetentionExecution     = errors.New("audit: retention execution failed")

	// Filter operations
	ErrInvalidFilter = errors.New("audit: invalid filter")

	// Security events
	ErrSecurityViolation = errors.New("audit: security violation detected")
	ErrUnauthorizedRead  = errors.New("audit: unauthorized audit read")
	ErrUnauthorizedWrite = errors.New("audit: unauthorized audit write")
)
