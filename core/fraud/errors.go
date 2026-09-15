package fraud

import "errors"

var (

	// Generic

	ErrInvalidSignal = errors.New(
		"invalid fraud signal",
	)

	ErrInvalidRule = errors.New(
		"invalid fraud rule",
	)

	ErrInvalidScore = errors.New(
		"invalid fraud score",
	)

	ErrStoreUnavailable = errors.New(
		"fraud store unavailable",
	)

	ErrDetectorUnavailable = errors.New(
		"fraud detector unavailable",
	)

	// Validation

	ErrEmptySignalType = errors.New(
		"signal type cannot be empty",
	)

	ErrEmptyTenantID = errors.New(
		"tenant id cannot be empty",
	)

	ErrEmptyUserID = errors.New(
		"user id cannot be empty",
	)

	ErrEmptyRuleName = errors.New(
		"rule name cannot be empty",
	)

	ErrInvalidThreshold = errors.New(
		"invalid threshold",
	)

	ErrInvalidWeight = errors.New(
		"invalid rule weight",
	)

	// Storage

	ErrSaveFailed = errors.New(
		"failed to save fraud data",
	)

	ErrDeleteFailed = errors.New(
		"failed to delete fraud data",
	)

	ErrLoadFailed = errors.New(
		"failed to load fraud data",
	)

	// Detection

	ErrFraudDetected = errors.New(
		"fraud detected",
	)

	ErrRiskTooHigh = errors.New(
		"risk score too high",
	)
)
