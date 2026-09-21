# HostMK Security Framework Go SDK

## Overview

The HostMK Security Framework Go SDK provides a simple and strongly typed interface for interacting with the HostMK Security API.

Features:

- Authentication
- Audit Logs
- Fraud Detection
- Permission Management
- Multi-Tenant Management
- REST API Client
- JSON Serialization
- Error Handling
- Context Support
- Production Ready

Supported Go Version:

```text
Go 1.22+
```

---

# Installation

Install the SDK:

```bash
go get github.com/azezaz/HostMK-Security-Framework/sdk/go
```

---

# Import

```go
import (
    hostmk "github.com/azezaz/HostMK-Security-Framework/sdk/go"
)
```

---

# Create Client

```go
client := hostmk.NewClient(
    "your-api-key",
    "https://api.hostmk.com/v1",
)
```

---

# Authentication

## Login

```go
auth := hostmk.NewAuthService(client)

response, err := auth.Login(
    "user@example.com",
    "password",
)

if err != nil {
    panic(err)
}

fmt.Println(response)
```

---

## Validate Token

```go
result, err := auth.ValidateToken(token)
```

---

## Refresh Token

```go
result, err := auth.RefreshToken(refreshToken)
```

---

## Logout

```go
result, err := auth.Logout()
```

---

# Audit Service

Create service:

```go
audit := hostmk.NewAuditService(client)
```

---

## List Events

```go
events, err := audit.ListEvents()
```

---

## Get Event

```go
event, err := audit.GetEvent("event-id")
```

---

## Search Events

```go
result, err := audit.SearchEvents("login")
```

---

## User Activity

```go
activity, err := audit.UserActivity("user-id")
```

---

## Tenant Activity

```go
activity, err := audit.TenantActivity("tenant-id")
```

---

## Security Events

```go
events, err := audit.SecurityEvents("")
```

Filter by severity:

```go
events, err := audit.SecurityEvents("critical")
```

---

## Statistics

```go
stats, err := audit.Statistics()
```

---

## Compliance Report

```go
report, err := audit.ComplianceReport("gdpr")
```

Supported:

- gdpr
- iso27001
- soc2
- hipaa

---

# Fraud Service

Create service:

```go
fraud := hostmk.NewFraudService(client)
```

---

## Analyze Transaction

```go
result, err := fraud.Analyze(map[string]interface{}{
    "transaction_id": "tx-123",
    "amount": 100,
})
```

---

## Risk Score

```go
result, err := fraud.Score(map[string]interface{}{
    "email": "user@example.com",
})
```

---

## Check IP

```go
result, err := fraud.CheckIP("8.8.8.8")
```

---

## Check Email

```go
result, err := fraud.CheckEmail(
    "user@example.com",
)
```

---

## Check Device

```go
result, err := fraud.CheckDevice(
    "device-id",
)
```

---

## Blacklist

```go
result, err := fraud.Blacklist(
    "ip",
    "1.2.3.4",
)
```

---

## Whitelist

```go
result, err := fraud.Whitelist(
    "ip",
    "8.8.8.8",
)
```

---

## Alerts

```go
alerts, err := fraud.Alerts()
```

---

## Statistics

```go
stats, err := fraud.Statistics()
```

---

# Permission Service

Create service:

```go
permissions := hostmk.NewPermissionService(client)
```

---

## List Permissions

```go
result, err := permissions.ListPermissions()
```

---

## Create Permission

```go
result, err := permissions.CreatePermission(
    "users.create",
)
```

---

## Get Permission

```go
result, err := permissions.GetPermission(
    "permission-id",
)
```

---

## Update Permission

```go
result, err := permissions.UpdatePermission(
    "permission-id",
    map[string]interface{}{
        "description": "Create users",
    },
)
```

---

## Delete Permission

```go
result, err := permissions.DeletePermission(
    "permission-id",
)
```

---

## User Permissions

```go
result, err := permissions.UserPermissions(
    "user-id",
)
```

---

## Role Permissions

```go
result, err := permissions.RolePermissions(
    "role-id",
)
```

---

## Assign Permission

```go
result, err := permissions.AssignPermission(
    "user-id",
    "users.create",
)
```

---

## Revoke Permission

```go
result, err := permissions.RevokePermission(
    "user-id",
    "users.create",
)
```

---

# Tenancy Service

Create service:

```go
tenancy := hostmk.NewTenancyService(client)
```

---

## List Tenants

```go
result, err := tenancy.ListTenants()
```

---

## Create Tenant

```go
result, err := tenancy.CreateTenant(
    "Acme Inc",
    "acme",
)
```

---

## Get Tenant

```go
result, err := tenancy.GetTenant(
    "tenant-id",
)
```

---

## Update Tenant

```go
result, err := tenancy.UpdateTenant(
    "tenant-id",
    map[string]interface{}{
        "name": "Updated Tenant",
    },
)
```

---

## Delete Tenant

```go
result, err := tenancy.DeleteTenant(
    "tenant-id",
)
```

---

## Tenant Settings

```go
settings, err := tenancy.TenantSettings(
    "tenant-id",
)
```

---

## Update Settings

```go
result, err := tenancy.UpdateSettings(
    "tenant-id",
    map[string]interface{}{
        "theme": "dark",
    },
)
```

---

## Tenant Statistics

```go
stats, err := tenancy.TenantStatistics(
    "tenant-id",
)
```

---

## Tenant Usage

```go
usage, err := tenancy.TenantUsage(
    "tenant-id",
)
```

---

## Activate Tenant

```go
result, err := tenancy.ActivateTenant(
    "tenant-id",
)
```

---

## Suspend Tenant

```go
result, err := tenancy.SuspendTenant(
    "tenant-id",
)
```

---

## Tenant Domains

```go
domains, err := tenancy.TenantDomains(
    "tenant-id",
)
```

---

# Health Check

```go
result, err := client.Health()
```

---

# Error Handling

Example:

```go
result, err := auth.Login(
    "",
    "",
)

if err != nil {
    fmt.Println(err.Error())
}
```

---

# Project Structure

```text
sdk/go/
├── auth.go
├── audit.go
├── fraud.go
├── permissions.go
├── tenancy.go
├── client.go
├── errors.go
├── models.go
└── README.md
```

---

# Best Practices

- Reuse a single client instance
- Store API keys in environment variables
- Use HTTPS only
- Rotate API keys regularly
- Handle all returned errors
- Use contexts for request cancellation
- Log security-related events

---

# Version

```text
v1.0.0
```

---

# License

MIT License

Copyright (c) HostMK Security Framework