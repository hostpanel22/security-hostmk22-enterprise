package tenancy

import "errors"

/*
|--------------------------------------------------------------------------
| Tenant Resolution Errors
|--------------------------------------------------------------------------
*/

var (
	ErrTenantNotFound = errors.New(
		"tenant not found",
	)

	ErrTenantResolutionFailed = errors.New(
		"tenant resolution failed",
	)

	ErrInvalidTenantIdentifier = errors.New(
		"invalid tenant identifier",
	)
)

/*
|--------------------------------------------------------------------------
| Tenant State Errors
|--------------------------------------------------------------------------
*/

var (
	ErrTenantInactive = errors.New(
		"tenant is inactive",
	)

	ErrTenantPending = errors.New(
		"tenant is pending activation",
	)

	ErrTenantSuspended = errors.New(
		"tenant is suspended",
	)

	ErrTenantExpired = errors.New(
		"tenant subscription expired",
	)

	ErrTenantDisabled = errors.New(
		"tenant is disabled",
	)

	ErrTenantDeleted = errors.New(
		"tenant is deleted",
	)
)

/*
|--------------------------------------------------------------------------
| Tenant Isolation Errors
|--------------------------------------------------------------------------
*/

var (
	ErrCrossTenantAccess = errors.New(
		"cross tenant access detected",
	)

	ErrCrossTenantQuery = errors.New(
		"cross tenant query detected",
	)

	ErrCrossTenantResourceAccess = errors.New(
		"cross tenant resource access detected",
	)

	ErrTenantMismatch = errors.New(
		"tenant mismatch detected",
	)
)

/*
|--------------------------------------------------------------------------
| Authentication & Authorization Errors
|--------------------------------------------------------------------------
*/

var (
	ErrUnauthorizedTenantAccess = errors.New(
		"unauthorized tenant access",
	)

	ErrForbiddenTenantOperation = errors.New(
		"forbidden tenant operation",
	)

	ErrInvalidTenantContext = errors.New(
		"invalid tenant context",
	)

	ErrMissingTenantContext = errors.New(
		"missing tenant context",
	)
)

/*
|--------------------------------------------------------------------------
| Domain Resolution Errors
|--------------------------------------------------------------------------
*/

var (
	ErrInvalidTenantDomain = errors.New(
		"invalid tenant domain",
	)

	ErrTenantDomainNotConfigured = errors.New(
		"tenant domain not configured",
	)

	ErrUnknownTenantDomain = errors.New(
		"unknown tenant domain",
	)
)

/*
|--------------------------------------------------------------------------
| API & Request Errors
|--------------------------------------------------------------------------
*/

var (
	ErrMissingTenantHeader = errors.New(
		"missing tenant header",
	)

	ErrInvalidTenantHeader = errors.New(
		"invalid tenant header",
	)

	ErrInvalidTenantToken = errors.New(
		"invalid tenant token",
	)

	ErrTenantValidationFailed = errors.New(
		"tenant validation failed",
	)
)

/*
|--------------------------------------------------------------------------
| Security Monitoring Errors
|--------------------------------------------------------------------------
*/

var (
	ErrTenantSecurityViolation = errors.New(
		"tenant security violation",
	)

	ErrSuspiciousTenantActivity = errors.New(
		"suspicious tenant activity detected",
	)

	ErrTenantRiskThresholdExceeded = errors.New(
		"tenant risk threshold exceeded",
	)
)
