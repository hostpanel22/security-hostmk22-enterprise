package permissions

import "errors"

var (

	// Generic

	ErrPermissionDenied = errors.New(
		"permission denied",
	)

	ErrRoleNotFound = errors.New(
		"role not found",
	)

	ErrPermissionNotFound = errors.New(
		"permission not found",
	)

	ErrInvalidRole = errors.New(
		"invalid role",
	)

	ErrInvalidPermission = errors.New(
		"invalid permission",
	)

	ErrStoreUnavailable = errors.New(
		"permission store unavailable",
	)

	ErrCacheUnavailable = errors.New(
		"permission cache unavailable",
	)

	// User

	ErrUserNotFound = errors.New(
		"user not found",
	)

	ErrUserHasNoRoles = errors.New(
		"user has no roles",
	)

	// Tenant

	ErrTenantNotFound = errors.New(
		"tenant not found",
	)

	ErrCrossTenantAccess = errors.New(
		"cross tenant access denied",
	)

	// Middleware

	ErrMissingTenantContext = errors.New(
		"missing tenant context",
	)

	ErrMissingUserContext = errors.New(
		"missing user context",
	)

	ErrMissingPermission = errors.New(
		"missing required permission",
	)

	// Validation

	ErrEmptyRoleName = errors.New(
		"role name cannot be empty",
	)

	ErrEmptyPermissionName = errors.New(
		"permission name cannot be empty",
	)

	ErrEmptyUserID = errors.New(
		"user id cannot be empty",
	)

	ErrEmptyTenantID = errors.New(
		"tenant id cannot be empty",
	)

	ErrDuplicateRole = errors.New(
		"role already exists",
	)

	ErrDuplicatePermission = errors.New(
		"permission already exists",
	)

	ErrDuplicateAssignment = errors.New(
		"permission assignment already exists",
	)

	// Storage

	ErrSaveFailed = errors.New(
		"failed to save permission data",
	)

	ErrDeleteFailed = errors.New(
		"failed to delete permission data",
	)

	ErrLoadFailed = errors.New(
		"failed to load permission data",
	)
)
