# 🔐 RBAC Package

The RBAC package provides Role-Based Access Control capabilities for HostMK-Security-Framework.

It is designed for multi-tenant SaaS platforms, APIs, microservices, internal tools and enterprise applications.

---

# Features

✅ Role Management

✅ Permission Management

✅ Policy Management

✅ Authorization Decisions

✅ Multi-Tenant Support

✅ HTTP Middleware

✅ Thread-Safe Memory Store

✅ Extensible Architecture

✅ Production Ready Foundation

---

# Installation

```bash
go get github.com/azezaz/HostMK-Security-Framework
```

Import package:

```go
import "github.com/azezaz/HostMK-Security-Framework/core/rbac"
```

---

# Package Structure

```text
core/rbac/

├── constants.go
├── errors.go
├── manager.go
├── middleware.go
├── permission.go
├── policy.go
├── role.go
├── store.go
├── types.go
└── utils.go
```

---

# Core Concepts

The RBAC system is built around four entities:

```text
User
 │
 ▼
Role
 │
 ▼
Permission
 │
 ▼
Resource Action
```

Example:

```text
User
 └── Admin

Admin
 ├── users:create
 ├── users:read
 ├── users:update
 └── users:delete
```

---

# Creating Permissions

Permissions define which actions are allowed on specific resources.

Example:

```go
permission := rbac.NewPermission(
    "tenant-1",
    "users_create",
    "users",
    "create",
    "Create platform users",
)
```

Permission Key:

```go
permission.Key()
```

Returns:

```text
users:create
```

---

# Creating Roles

Roles group multiple permissions together.

Example:

```go
role := rbac.NewRole(
    "tenant-1",
    "admin",
    "Platform administrator",
)
```

Adding permissions:

```go
err := role.AddPermission(
    "users:create",
)
```

---

# Memory Store

The package includes a thread-safe in-memory store.

Create store:

```go
store := rbac.NewMemoryStore()
```

Store supports:

- SaveRole
- GetRole
- ListRoles
- DeleteRole
- SavePermission
- GetPermission
- ListPermissions
- DeletePermission
- Clear

---

# Manager

The Manager provides high-level RBAC operations.

Create manager:

```go
store := rbac.NewMemoryStore()

manager, err := rbac.NewManager(
    store,
)
```

Create role:

```go
err = manager.CreateRole(
    ctx,
    role,
)
```

Create permission:

```go
err = manager.CreatePermission(
    ctx,
    permission,
)
```

Assign permission:

```go
err = manager.AssignPermissionToRole(
    ctx,
    "admin",
    "users_create",
)
```

Check permission:

```go
allowed, err := manager.RoleHasPermission(
    ctx,
    "admin",
    "users:create",
)
```

---

# Middleware

The RBAC package provides HTTP middleware.

Headers:

```text
X-Tenant-ID
X-User-ID
X-Role
```

Example:

```go
middleware := rbac.RBACMiddleware(
    rbac.MiddlewareConfig{
        RequireTenant: true,
        RequireUser: true,
        RequireRole: true,
    },
)
```

Protect endpoint:

```go
http.Handle(
    "/admin",
    middleware(
        adminHandler,
    ),
)
```

---

# Context Helpers

Retrieve values from context.

Tenant:

```go
tenantID, ok := rbac.TenantIDFromContext(
    ctx,
)
```

User:

```go
userID, ok := rbac.UserIDFromContext(
    ctx,
)
```

Role:

```go
role, ok := rbac.RoleFromContext(
    ctx,
)
```

---

# Multi-Tenant Support

All RBAC entities can be associated with a tenant.

Example:

```text
tenant-a
tenant-b
tenant-c
```

Isolation is enforced at application level.

---

# Authorization Flow

```text
HTTP Request
      │
      ▼
RBAC Middleware
      │
      ▼
Context Extraction
      │
      ▼
Role Resolution
      │
      ▼
Permission Check
      │
      ▼
Allow / Deny
```

---

# Built-In Roles

Available constants:

```go
rbac.DefaultAdminRole
rbac.DefaultUserRole
rbac.DefaultViewerRole
rbac.DefaultOperatorRole
rbac.DefaultAuditorRole
```

---

# Built-In Actions

```go
rbac.PermissionCreate
rbac.PermissionRead
rbac.PermissionUpdate
rbac.PermissionDelete
rbac.PermissionList
rbac.PermissionManage
rbac.PermissionExecute
```

---

# Built-In Resources

```go
rbac.ResourceUsers
rbac.ResourceRoles
rbac.ResourcePermissions
rbac.ResourcePolicies
rbac.ResourceTenants
rbac.ResourceAuditLogs
rbac.ResourceSecrets
rbac.ResourceAPIKeys
rbac.ResourceSessions
rbac.ResourceBilling
rbac.ResourceSettings
```

---

# Error Handling

Common errors:

```go
rbac.ErrAccessDenied

rbac.ErrUnauthorized

rbac.ErrRoleNotFound

rbac.ErrPermissionNotFound

rbac.ErrRoleAlreadyExists

rbac.ErrPermissionAlreadyExists
```

Example:

```go
if err == rbac.ErrAccessDenied {
    // Handle authorization failure
}
```

---

# Testing

Run RBAC tests:

```bash
go test ./tests/rbac -v
```

Run all tests:

```bash
go test ./...
```

---

# Future Extensions

The package is designed to integrate with:

- JWT Authentication
- API Keys
- OAuth2
- OpenID Connect
- MFA / TOTP
- SAML SSO
- LDAP
- Active Directory
- Threat Detection Engine

---

# License

Licensed under Apache License 2.0.

Enterprise Edition modules located under `/ee` are governed by the HostMK Enterprise License.

---

# HostMK Security Framework

Secure. Multi-Tenant. Extensible.