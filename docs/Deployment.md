# HostMK Security Framework API

## Overview

The HostMK Security Framework API provides a centralized security layer for SaaS platforms, multi-tenant applications, internal tools, and enterprise systems.

The API exposes services for:

- Authentication
- Audit Logging
- Fraud Detection
- Permission Management
- Multi-Tenant Management

Base URL:

```text
https://api.hostmk.io/v1
```

All requests and responses use JSON.

---

# Authentication

## API Key Authentication

Every request must include a valid API key.

Example:

```http
X-API-Key: YOUR_API_KEY
```

Example request:

```http
GET /auth/health HTTP/1.1
Host: api.hostmk.io
X-API-Key: YOUR_API_KEY
Content-Type: application/json
```

---

# Response Format

Successful response:

```json
{
  "success": true,
  "data": {}
}
```

Error response:

```json
{
  "success": false,
  "message": "Validation error"
}
```

---

# Auth Service

## Health Check

### Request

```http
GET /auth/health
```

### Response

```json
{
  "success": true,
  "status": "ok"
}
```

---

## Login

### Request

```http
POST /auth/login
```

Body:

```json
{
  "email": "user@example.com",
  "password": "secret"
}
```

### Response

```json
{
  "success": true,
  "token": "jwt-token"
}
```

---

## Validate Token

### Request

```http
POST /auth/validate
```

Body:

```json
{
  "token": "jwt-token"
}
```

### Response

```json
{
  "success": true,
  "valid": true
}
```

---

## Refresh Token

### Request

```http
POST /auth/refresh
```

Body:

```json
{
  "refresh_token": "refresh-token"
}
```

### Response

```json
{
  "success": true,
  "token": "new-token"
}
```

---

## Logout

### Request

```http
POST /auth/logout
```

### Response

```json
{
  "success": true
}
```

---

# Audit Service

## List Events

### Request

```http
GET /audit/events
```

### Response

```json
{
  "success": true,
  "data": [
    {
      "id": "evt_001",
      "action": "login",
      "user_id": "usr_001"
    }
  ]
}
```

---

## Event Details

### Request

```http
GET /audit/events/{event_id}
```

### Response

```json
{
  "success": true,
  "data": {
    "id": "evt_001",
    "action": "login"
  }
}
```

---

## User Activity

### Request

```http
GET /audit/users/{user_id}
```

### Response

```json
{
  "success": true,
  "events": []
}
```

---

## Tenant Activity

### Request

```http
GET /audit/tenants/{tenant_id}
```

### Response

```json
{
  "success": true,
  "events": []
}
```

---

## Security Events

### Request

```http
GET /audit/security-events
```

Optional Query:

```text
?severity=critical
```

### Response

```json
{
  "success": true,
  "events": []
}
```

---

## Statistics

### Request

```http
GET /audit/statistics
```

### Response

```json
{
  "success": true,
  "total_events": 1000
}
```

---

## Compliance Report

### Request

```http
GET /audit/compliance/{type}
```

Supported:

- gdpr
- iso27001
- soc2

### Response

```json
{
  "success": true,
  "report": {}
}
```

---

# Fraud Service

## Analyze

### Request

```http
POST /fraud/analyze
```

Body:

```json
{
  "email": "user@example.com",
  "ip": "1.2.3.4"
}
```

### Response

```json
{
  "success": true,
  "risk_score": 12
}
```

---

## Score

### Request

```http
POST /fraud/score
```

### Response

```json
{
  "success": true,
  "score": 75
}
```

---

## Check Email

### Request

```http
GET /fraud/email/{email}
```

### Response

```json
{
  "success": true,
  "risk": "low"
}
```

---

## Check IP

### Request

```http
GET /fraud/ip/{ip}
```

### Response

```json
{
  "success": true,
  "risk": "medium"
}
```

---

## Check Device

### Request

```http
GET /fraud/device/{device_id}
```

### Response

```json
{
  "success": true,
  "risk": "high"
}
```

---

## Blacklist

### Request

```http
POST /fraud/blacklist
```

Body:

```json
{
  "type": "email",
  "value": "spam@example.com"
}
```

### Response

```json
{
  "success": true
}
```

---

## Whitelist

### Request

```http
POST /fraud/whitelist
```

Body:

```json
{
  "type": "email",
  "value": "trusted@example.com"
}
```

### Response

```json
{
  "success": true
}
```

---

# Permission Service

## List Permissions

### Request

```http
GET /permissions
```

### Response

```json
{
  "success": true,
  "permissions": []
}
```

---

## Create Permission

### Request

```http
POST /permissions
```

Body:

```json
{
  "name": "users.create"
}
```

### Response

```json
{
  "success": true
}
```

---

## Assign Permission

### Request

```http
POST /permissions/assign
```

Body:

```json
{
  "user_id": "usr_001",
  "permission": "users.create"
}
```

### Response

```json
{
  "success": true
}
```

---

## Revoke Permission

### Request

```http
POST /permissions/revoke
```

### Response

```json
{
  "success": true
}
```

---

## User Permissions

### Request

```http
GET /permissions/users/{user_id}
```

### Response

```json
{
  "success": true,
  "permissions": []
}
```

---

## Role Permissions

### Request

```http
GET /permissions/roles/{role_id}
```

### Response

```json
{
  "success": true,
  "permissions": []
}
```

---

# Tenancy Service

## List Tenants

### Request

```http
GET /tenants
```

### Response

```json
{
  "success": true,
  "data": []
}
```

---

## Create Tenant

### Request

```http
POST /tenants
```

Body:

```json
{
  "name": "Acme Inc",
  "slug": "acme"
}
```

### Response

```json
{
  "success": true,
  "tenant_id": "tenant_001"
}
```

---

## Get Tenant

### Request

```http
GET /tenants/{tenant_id}
```

### Response

```json
{
  "success": true,
  "data": {}
}
```

---

## Update Tenant

### Request

```http
PUT /tenants/{tenant_id}
```

### Response

```json
{
  "success": true
}
```

---

## Suspend Tenant

### Request

```http
POST /tenants/{tenant_id}/suspend
```

### Response

```json
{
  "success": true
}
```

---

## Activate Tenant

### Request

```http
POST /tenants/{tenant_id}/activate
```

### Response

```json
{
  "success": true
}
```

---

## Tenant Usage

### Request

```http
GET /tenants/{tenant_id}/usage
```

### Response

```json
{
  "success": true,
  "usage": {}
}
```

---

## Tenant Statistics

### Request

```http
GET /tenants/{tenant_id}/statistics
```

### Response

```json
{
  "success": true,
  "statistics": {}
}
```

---

## Tenant Domains

### Request

```http
GET /tenants/{tenant_id}/domains
```

### Response

```json
{
  "success": true,
  "domains": []
}
```

---

## Tenant Settings

### Request

```http
GET /tenants/{tenant_id}/settings
```

### Response

```json
{
  "success": true,
  "settings": {}
}
```

---

# SDK Support

Official SDKs:

- PHP SDK
- Python SDK
- Node.js SDK
- Go SDK

Documentation:

- SDK-PHP.md
- SDK-Python.md
- SDK-Node.md
- SDK-Go.md

---

# Versioning

Current Version:

```text
v1.0.0
```

The API follows Semantic Versioning (SemVer).

Examples:

```text
v1.0.0
v1.1.0
v2.0.0
```

---

# License

MIT License

Copyright (c) HostMK Security Framework