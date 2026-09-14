package tenancy

import (
	"time"
)

/*
|--------------------------------------------------------------------------
| Validation Result
|--------------------------------------------------------------------------
*/

type ValidationResult struct {
	Valid bool `json:"valid"`

	Error error `json:"-"`

	Reason string `json:"reason,omitempty"`

	Violation *Violation `json:"violation,omitempty"`
}

/*
|--------------------------------------------------------------------------
| Tenant Validator
|--------------------------------------------------------------------------
|
| Responsible for validating tenant lifecycle,
| security state and platform access eligibility.
|
*/

type Validator struct {
	AllowPending bool
}

/*
|--------------------------------------------------------------------------
| Constructor
|--------------------------------------------------------------------------
*/

func NewValidator() *Validator {
	return &Validator{
		AllowPending: false,
	}
}

/*
|--------------------------------------------------------------------------
| Main Validation
|--------------------------------------------------------------------------
*/

/*
Validate checks whether a tenant
is allowed to access protected resources.
*/
func (v *Validator) Validate(
	tenant *Tenant,
) *ValidationResult {

	if tenant == nil {
		return &ValidationResult{
			Valid:  false,
			Error:  ErrTenantNotFound,
			Reason: "tenant does not exist",
		}
	}

	if tenant.IsDeleted() {

		return v.failed(
			tenant,
			ErrTenantDeleted,
			"tenant has been deleted",
			SeverityCritical,
		)
	}

	if tenant.IsDisabled() {

		return v.failed(
			tenant,
			ErrTenantDisabled,
			"tenant is disabled",
			SeverityHigh,
		)
	}

	if tenant.IsSuspended() {

		return v.failed(
			tenant,
			ErrTenantSuspended,
			"tenant is suspended",
			SeverityHigh,
		)
	}

	if tenant.IsExpired() {

		return v.failed(
			tenant,
			ErrTenantExpired,
			"tenant subscription expired",
			SeverityMedium,
		)
	}

	if tenant.IsPending() && !v.AllowPending {

		return v.failed(
			tenant,
			ErrTenantPending,
			"tenant activation pending",
			SeverityLow,
		)
	}

	if !tenant.IsActive() &&
		!tenant.IsPending() {

		return v.failed(
			tenant,
			ErrTenantInactive,
			"tenant inactive",
			SeverityMedium,
		)
	}

	return &ValidationResult{
		Valid: true,
	}
}

/*
|--------------------------------------------------------------------------
| Tenant Context Validation
|--------------------------------------------------------------------------
*/

/*
ValidateContext ensures that a request
contains a valid tenant context.
*/
func (v *Validator) ValidateContext(
	ctx *TenantContext,
) error {

	if ctx == nil {
		return ErrMissingTenantContext
	}

	if !ctx.IsValid() {
		return ErrInvalidTenantContext
	}

	return nil
}

/*
|--------------------------------------------------------------------------
| Tenant Ownership Validation
|--------------------------------------------------------------------------
*/

/*
ValidateOwnership verifies that a resource
belongs to the active tenant.
*/
func (v *Validator) ValidateOwnership(
	ctx *TenantContext,
	resourceTenantID string,
) error {

	if ctx == nil {
		return ErrMissingTenantContext
	}

	if resourceTenantID == "" {
		return ErrInvalidTenantIdentifier
	}

	if ctx.TenantID != resourceTenantID {
		return ErrCrossTenantAccess
	}

	return nil
}

/*
|--------------------------------------------------------------------------
| Expiration Validation
|--------------------------------------------------------------------------
*/

/*
ValidateExpiration checks expiration rules.
*/
func (v *Validator) ValidateExpiration(
	tenant *Tenant,
) error {

	if tenant == nil {
		return ErrTenantNotFound
	}

	if tenant.ExpiresAt == nil {
		return nil
	}

	if time.Now().UTC().After(
		*tenant.ExpiresAt,
	) {
		return ErrTenantExpired
	}

	return nil
}

/*
|--------------------------------------------------------------------------
| Internal Helpers
|--------------------------------------------------------------------------
*/

func (v *Validator) failed(
	tenant *Tenant,
	err error,
	message string,
	severity string,
) *ValidationResult {

	violation := NewViolation(
		ViolationSecurityPolicy,
		severity,
		tenant.ID,
		message,
	)

	return &ValidationResult{
		Valid:     false,
		Error:     err,
		Reason:    message,
		Violation: violation,
	}
}
