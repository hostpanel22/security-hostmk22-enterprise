package rbac

import "errors"

var (

	// Generic

	ErrRBACUnavailable = errors.New(
		"rbac unavailable",
	)

	ErrAccessDenied = errors.New(
		"access denied",
	)

	ErrUnauthorized = errors.New(
		"unauthorized",
	)

	// Role

	ErrInvalidRole = errors.New(
		"invalid role",
	)

	ErrRoleNotFound = errors.New(
		"role not found",
	)

	ErrRoleAlreadyExists = errors.New(
		"role already exists",
	)

	ErrEmptyRoleID = errors.New(
		"role id cannot be empty",
	)

	ErrEmptyRoleName = errors.New(
		"role name cannot be empty",
	)

	ErrRoleInheritanceLoop = errors.New(
		"role inheritance loop detected",
	)

	// Permission

	ErrInvalidPermission = errors.New(
		"invalid permission",
	)

	ErrPermissionNotFound = errors.New(
		"permission not found",
	)

	ErrPermissionAlreadyExists = errors.New(
		"permission already exists",
	)

	ErrEmptyPermissionName = errors.New(
		"permission name cannot be empty",
	)

	ErrEmptyPermissionResource = errors.New(
		"permission resource cannot be empty",
	)

	ErrEmptyPermissionAction = errors.New(
		"permission action cannot be empty",
	)

	// Assignment

	ErrPermissionAlreadyAssigned = errors.New(
		"permission already assigned",
	)

	ErrPermissionNotAssigned = errors.New(
		"permission not assigned",
	)

	ErrRoleAlreadyAssigned = errors.New(
		"role already assigned",
	)

	ErrRoleNotAssigned = errors.New(
		"role not assigned",
	)

	// Policy

	ErrInvalidPolicy = errors.New(
		"invalid policy",
	)

	ErrPolicyNotFound = errors.New(
		"policy not found",
	)

	ErrPolicyAlreadyExists = errors.New(
		"policy already exists",
	)

	ErrEmptyPolicyName = errors.New(
		"policy name cannot be empty",
	)

	// Store

	ErrStoreUnavailable = errors.New(
		"rbac store unavailable",
	)

	ErrSaveFailed = errors.New(
		"failed to save object",
	)

	ErrDeleteFailed = errors.New(
		"failed to delete object",
	)

	ErrUpdateFailed = errors.New(
		"failed to update object",
	)

	ErrLoadFailed = errors.New(
		"failed to load object",
	)

	// Context

	ErrMissingTenantContext = errors.New(
		"missing tenant context",
	)

	ErrMissingUserContext = errors.New(
		"missing user context",
	)

	ErrInvalidTenantID = errors.New(
		"invalid tenant id",
	)

	ErrInvalidUserID = errors.New(
		"invalid user id",
	)

	// Validation

	ErrValidationFailed = errors.New(
		"validation failed",
	)
)
