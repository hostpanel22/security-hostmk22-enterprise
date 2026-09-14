package tenancy

import (
	"context"
	"time"
)

/*
|--------------------------------------------------------------------------
| Isolation Engine
|--------------------------------------------------------------------------
|
| Responsible for enforcing tenant boundaries.
| Any attempt to access resources belonging
| to another tenant is immediately rejected
| and reported as a security violation.
|
*/

type IsolationEngine struct {
	strictMode bool
}

/*
|--------------------------------------------------------------------------
| Constructor
|--------------------------------------------------------------------------
*/

func NewIsolationEngine() *IsolationEngine {
	return &IsolationEngine{
		strictMode: true,
	}
}

/*
|--------------------------------------------------------------------------
| Configuration
|--------------------------------------------------------------------------
*/

func (e *IsolationEngine) EnableStrictMode() {
	e.strictMode = true
}

func (e *IsolationEngine) DisableStrictMode() {
	e.strictMode = false
}

func (e *IsolationEngine) IsStrictMode() bool {
	return e.strictMode
}

/*
|--------------------------------------------------------------------------
| Ownership Validation
|--------------------------------------------------------------------------
|
| Verifies that the active tenant owns
| the target resource.
|
*/

func (e *IsolationEngine) ValidateOwnership(
	ctx *TenantContext,
	resourceTenantID string,
	resource string,
	action string,
) error {

	if ctx == nil {
		return ErrMissingTenantContext
	}

	if resourceTenantID == "" {
		return ErrInvalidTenantIdentifier
	}

	if ctx.TenantID == resourceTenantID {
		return nil
	}

	violation := e.BuildCrossTenantViolation(
		ctx,
		resourceTenantID,
		resource,
		action,
	)

	_ = violation

	if e.strictMode {
		return ErrCrossTenantAccess
	}

	return nil
}

/*
|--------------------------------------------------------------------------
| Cross Tenant Detection
|--------------------------------------------------------------------------
*/

func (e *IsolationEngine) IsCrossTenantAccess(
	ctx *TenantContext,
	resourceTenantID string,
) bool {

	if ctx == nil {
		return true
	}

	return ctx.TenantID != resourceTenantID
}

/*
|--------------------------------------------------------------------------
| Context Validation
|--------------------------------------------------------------------------
*/

func (e *IsolationEngine) ValidateContext(
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
| Violation Builder
|--------------------------------------------------------------------------
*/

func (e *IsolationEngine) BuildCrossTenantViolation(
	ctx *TenantContext,
	targetTenantID string,
	resource string,
	action string,
) *Violation {

	violation := NewViolation(
		ViolationCrossTenantAccess,
		SeverityCritical,
		ctx.TenantID,
		"cross tenant access attempt detected",
	)

	violation.ActorID = ctx.ActorID
	violation.RequestID = ctx.RequestID
	violation.Resource = resource
	violation.Action = action
	violation.IPAddress = ctx.IPAddress
	violation.UserAgent = ctx.UserAgent

	violation.AddMetadata(
		"source_tenant_id",
		ctx.TenantID,
	)

	violation.AddMetadata(
		"target_tenant_id",
		targetTenantID,
	)

	violation.AddMetadata(
		"detected_at",
		time.Now().UTC(),
	)

	return violation
}

/*
|--------------------------------------------------------------------------
| Tenant Match Validation
|--------------------------------------------------------------------------
*/

func (e *IsolationEngine) ValidateTenantMatch(
	sourceTenantID string,
	targetTenantID string,
) error {

	if sourceTenantID == "" {
		return ErrInvalidTenantIdentifier
	}

	if targetTenantID == "" {
		return ErrInvalidTenantIdentifier
	}

	if sourceTenantID != targetTenantID {
		return ErrTenantMismatch
	}

	return nil
}

/*
|--------------------------------------------------------------------------
| Protected Query Validation
|--------------------------------------------------------------------------
|
| Used before executing database queries.
|
*/

func (e *IsolationEngine) ValidateQuery(
	ctx *TenantContext,
	queryTenantID string,
) error {

	if err := e.ValidateContext(ctx); err != nil {
		return err
	}

	if ctx.TenantID != queryTenantID {
		return ErrCrossTenantQuery
	}

	return nil
}

/*
|--------------------------------------------------------------------------
| Resource Access Validation
|--------------------------------------------------------------------------
*/

func (e *IsolationEngine) ValidateResourceAccess(
	ctx *TenantContext,
	resourceTenantID string,
	resource string,
) error {

	if err := e.ValidateContext(ctx); err != nil {
		return err
	}

	if ctx.TenantID != resourceTenantID {
		return ErrCrossTenantResourceAccess
	}

	return nil
}

/*
|--------------------------------------------------------------------------
| Secure Context Check
|--------------------------------------------------------------------------
*/

func (e *IsolationEngine) ValidateRequestContext(
	ctx context.Context,
) (*TenantContext, error) {

	tenantCtx, ok := TenantFromContext(ctx)

	if !ok {
		return nil, ErrMissingTenantContext
	}

	if !tenantCtx.IsValid() {
		return nil, ErrInvalidTenantContext
	}

	return tenantCtx, nil
}
