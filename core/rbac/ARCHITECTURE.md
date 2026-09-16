# 🏗️ RBAC Architecture

## Overview

The RBAC (Role-Based Access Control) module is responsible for authorization inside the HostMK Security Framework.

It provides a flexible and extensible permission system that supports:

- Multi-Tenant SaaS platforms
- Enterprise applications
- APIs
- Internal administration panels
- Microservices
- White-label environments

The architecture is intentionally decoupled from authentication systems such as JWT, API Keys, OAuth2, SAML, LDAP, or Session Management.

Authentication determines **who the user is**.

RBAC determines **what the user can do**.

---

# Architecture Goals

The RBAC package is designed around the following principles:

## Security First

Authorization checks must be deterministic and predictable.

Access is denied by default.

## Multi-Tenant Ready

Every RBAC object can belong to a tenant.

Isolation between tenants is mandatory.

## Extensible

The module can be integrated with:

- JWT
- API Keys
- OAuth2
- OpenID Connect
- SAML
- LDAP
- Active Directory

## Testable

Business logic is isolated from storage and transport layers.

## High Performance

Authorization checks must execute in constant time whenever possible.

## Thread Safe

Concurrent access is supported through synchronization primitives.

---

# Package Layout

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
├── utils.go
├── README.md
└── ARCHITECTURE.md
```

---

# High Level Flow

```text
Client
  │
  ▼
Authentication Layer
(JWT / API Key / Session)
  │
  ▼
RBAC Middleware
  │
  ▼
Manager
  │
  ▼
Role Resolution
  │
  ▼
Permission Evaluation
  │
  ▼
Authorization Decision
  │
  ├── Allow
  │
  └── Deny
```

---

# Core Components

## Permission

Permissions are the smallest authorization unit.

Structure:

```text
Resource + Action
```

Examples:

```text
users:create
users:read
users:update
users:delete

billing:read
billing:update

secrets:manage
```

Example object:

```go
Permission{
    Resource: "users",
    Action: "create",
}
```

Generated key:

```text
users:create
```

---

## Role

Roles aggregate permissions.

Example:

```text
Admin
├── users:create
├── users:read
├── users:update
└── users:delete
```

Example:

```go
Role{
    Name: "admin",
}
```

Roles are assigned to users externally by the authentication layer.

---

## Policy

Policies describe authorization requirements.

They can be used by middleware or business logic.

Example:

```text
Policy:
Role = admin

OR

Permission = users:create
```

Policies are intentionally lightweight to keep evaluation fast.

---

## Store

The Store layer abstracts persistence.

Current implementation:

```text
MemoryStore
```

Future implementations:

```text
PostgreSQLStore
MySQLStore
RedisStore
MongoStore
```

Store responsibilities:

```text
SaveRole
GetRole
DeleteRole

SavePermission
GetPermission
DeletePermission
```

---

## Manager

Manager is the RBAC service layer.

Responsibilities:

```text
Role Management
Permission Management
Assignments
Authorization Queries
```

The Manager is the entry point used by applications.

Example:

```go
manager.AssignPermissionToRole(...)
```

Example:

```go
manager.RoleHasPermission(...)
```

---

## Middleware

Middleware protects HTTP endpoints.

Responsibilities:

```text
Extract Context
Validate Headers
Resolve Identity
Execute Authorization Checks
```

Example headers:

```http
X-Tenant-ID
X-User-ID
X-Role
```

---

# Data Model

## Permission

```text
Permission
├── ID
├── TenantID
├── Name
├── Resource
├── Action
├── Description
├── CreatedAt
└── UpdatedAt
```

---

## Role

```text
Role
├── ID
├── TenantID
├── Name
├── Description
├── Permissions[]
├── CreatedAt
└── UpdatedAt
```

---

## Policy

```text
Policy
├── ID
├── Name
├── Description
├── Roles[]
├── Permissions[]
├── CreatedAt
└── UpdatedAt
```

---

# Multi-Tenant Design

Each tenant owns its RBAC configuration.

Example:

```text
Tenant A
 ├── Admin
 ├── Support
 └── Customer

Tenant B
 ├── Owner
 ├── Staff
 └── Viewer
```

Permissions are isolated.

No cross-tenant authorization is allowed.

---

# Authorization Lifecycle

## Step 1

Authentication succeeds.

Example:

```text
JWT Valid
```

---

## Step 2

User identity is loaded.

Example:

```text
User ID
Tenant ID
Assigned Role
```

---

## Step 3

RBAC middleware extracts context.

Example:

```text
tenant-001
user-123
admin
```

---

## Step 4

Role is resolved.

Example:

```text
admin
```

---

## Step 5

Permission check executes.

Example:

```text
users:create
```

---

## Step 6

Decision returned.

```text
ALLOW
```

or

```text
DENY
```

---

# Security Principles

## Default Deny

If authorization cannot be determined:

```text
DENY
```

---

## Tenant Isolation

Authorization must remain tenant-scoped.

---

## Least Privilege

Users should receive only required permissions.

---

## Immutable Decisions

Authorization results should never be modified after evaluation.

---

## Explicit Permissions

Permissions are always defined explicitly.

No hidden privilege escalation.

---

# Thread Safety

MemoryStore uses:

```go
sync.RWMutex
```

Benefits:

```text
Safe Concurrent Reads
Safe Concurrent Writes
Race Condition Protection
```

---

# Error Handling Strategy

Examples:

```go
ErrAccessDenied
ErrUnauthorized
ErrRoleNotFound
ErrPermissionNotFound
```

Goals:

- Predictable behavior
- Easy debugging
- Consistent API

---

# Performance Considerations

Current targets:

```text
Permission Lookup  O(1)
Role Lookup        O(1)
Assignment Check   O(1)
```

MemoryStore uses Go maps for fast lookups.

---

# Future Roadmap

## RBAC v2

Role inheritance:

```text
Admin
 └── Manager
      └── User
```

---

## RBAC v3

Conditional permissions:

```text
Allow if owner
Allow if tenant matches
Allow during business hours
```

---

## RBAC v4

Attribute Based Access Control (ABAC)

```text
User Attributes
Resource Attributes
Environment Attributes
```

---

## RBAC v5

Distributed authorization cache.

```text
Redis
Memcached
```

---

# Integration With HostMK Security Framework

RBAC integrates with:

```text
JWT
API Keys
MFA
Sessions
Rate Limiting
Audit Logs
Threat Detection
Webhooks
IP Reputation
```

---

# Summary

The RBAC module provides a secure, tenant-aware, extensible authorization layer for the HostMK Security Framework.

Design priorities:

- Security
- Simplicity
- Extensibility
- Performance
- Maintainability
- Enterprise Readiness