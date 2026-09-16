package rbac

const (

	// Default Roles

	DefaultAdminRole = "admin"

	DefaultUserRole = "user"

	DefaultViewerRole = "viewer"

	DefaultOperatorRole = "operator"

	DefaultAuditorRole = "auditor"

	// Permission Actions

	PermissionCreate = "create"

	PermissionRead = "read"

	PermissionUpdate = "update"

	PermissionDelete = "delete"

	PermissionList = "list"

	PermissionManage = "manage"

	PermissionExecute = "execute"

	// Resources

	ResourceUsers = "users"

	ResourceRoles = "roles"

	ResourcePermissions = "permissions"

	ResourcePolicies = "policies"

	ResourceTenants = "tenants"

	ResourceAuditLogs = "audit_logs"

	ResourceSecrets = "secrets"

	ResourceAPIKeys = "api_keys"

	ResourceSessions = "sessions"

	ResourceBilling = "billing"

	ResourceSettings = "settings"

	// Policy Effects

	EffectAllow = "allow"

	EffectDeny = "deny"

	// Context Keys

	ContextTenantID = "tenant_id"

	ContextUserID = "user_id"

	ContextRole = "role"

	// Headers

	HeaderTenantID = "X-Tenant-ID"

	HeaderUserID = "X-User-ID"

	HeaderRole = "X-Role"
)
