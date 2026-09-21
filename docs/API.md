# HostMK Security Framework API

## Overview

The HostMK Security Framework API provides security services for modern SaaS applications, multi-tenant platforms, e-commerce systems, fintech products, and enterprise infrastructures.

Main modules:

- Authentication
- Audit Logs
- Fraud Detection
- Permissions & RBAC
- Multi-Tenant Management

Base URL:

```text
https://api.hostmk.com/v1
```

Authentication:

```http
X-API-Key: YOUR_API_KEY
Content-Type: application/json
Accept: application/json
```

---

# Authentication API

## Login

### Request

```http
POST /auth/login
```

Body:

```json
{
  "email": "user@example.com",
  "password": "password"
}
```

### Response

```json
{
  "success": true,
  "token": "jwt-token",
  "refresh_token": "refresh-token"
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
  "token": "new-jwt-token"
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

# Audit API

## List Events

### Request

```http
GET /audit/events
```

### Response

```json
{
  "data": [
    {
      "id": "evt_123",
      "action": "login",
      "user_id": "usr_1",
      "created_at": "2026-01-01T00:00:00Z"
    }
  ]
}
```

---

## Get Event

### Request

```http
GET /audit/events/{eventId}
```

### Response

```json
{
  "id": "evt_123",
  "action": "login"
}
```

---

## Search Events

### Request

```http
POST /audit/search
```

Body:

```json
{
  "query": "login"
}
```

### Response

```json
{
  "results": []
}
```

---

## User Activity

### Request

```http
GET /audit/users/{userId}/activity
```

### Response

```json
{
  "activities": []
}
```

---

## Tenant Activity

### Request

```http
GET /audit/tenants/{tenantId}/activity
```

### Response

```json
{
  "activities": []
}
```

---

## Security Events

### Request

```http
GET /audit/security-events
```

### Response

```json
{
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
  "total_events": 10000
}
```

---

## Compliance Report

### Request

```http
GET /audit/compliance/{type}
```

Supported types:

- gdpr
- iso27001
- soc2
- hipaa

### Response

```json
{
  "type": "gdpr",
  "status": "compliant"
}
```

---

# Fraud Detection API

## Analyze Transaction

### Request

```http
POST /fraud/analyze
```

Body:

```json
{
  "transaction_id": "tx_123",
  "amount": 100
}
```

### Response

```json
{
  "risk_score": 12
}
```

---

## Risk Score

### Request

```http
POST /fraud/score
```

### Response

```json
{
  "score": 25
}
```

---

## Check IP

### Request

```http
POST /fraud/check-ip
```

Body:

```json
{
  "ip": "8.8.8.8"
}
```

### Response

```json
{
  "trusted": true
}
```

---

## Check Email

### Request

```http
POST /fraud/check-email
```

Body:

```json
{
  "email": "user@example.com"
}
```

### Response

```json
{
  "risk": "low"
}
```

---

## Check Device

### Request

```http
POST /fraud/check-device
```

Body:

```json
{
  "device_id": "device123"
}
```

### Response

```json
{
  "risk": "medium"
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
  "type": "ip",
  "value": "1.2.3.4"
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

### Response

```json
{
  "success": true
}
```

---

## Alerts

### Request

```http
GET /fraud/alerts
```

### Response

```json
{
  "alerts": []
}
```

---

# Permission API

## List Permissions

### Request

```http
GET /permissions
```

### Response

```json
{
  "data": []
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
  "id": "perm_123"
}
```

---

## Get Permission

### Request

```http
GET /permissions/{permissionId}
```

### Response

```json
{
  "id": "perm_123"
}
```

---

## Update Permission

### Request

```http
PUT /permissions/{permissionId}
```

### Response

```json
{
  "success": true
}
```

---

## Delete Permission

### Request

```http
DELETE /permissions/{permissionId}
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
GET /permissions/users/{userId}
```

### Response

```json
{
  "permissions": []
}
```

---

## Role Permissions

### Request

```http
GET /permissions/roles/{roleId}
```

### Response

```json
{
  "permissions": []
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
  "user_id": "usr_1",
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

# Multi-Tenant API

## List Tenants

### Request

```http
GET /tenants
```

### Response

```json
{
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
  "name": "Acme",
  "slug": "acme"
}
```

### Response

```json
{
  "id": "tenant_123"
}
```

---

## Get Tenant

### Request

```http
GET /tenants/{tenantId}
```

### Response

```json
{
  "id": "tenant_123"
}
```

---

## Update Tenant

### Request

```http
PUT /tenants/{tenantId}
```

### Response

```json
{
  "success": true
}
```

---

## Delete Tenant

### Request

```http
DELETE /tenants/{tenantId}
```

### Response

```json
{
  "success": true
}
```

---

## Tenant Settings

### Request

```http
GET /tenants/{tenantId}/settings
```

### Response

```json
{
  "settings": {}
}
```

---

## Update Settings

### Request

```http
PUT /tenants/{tenantId}/settings
```

### Response

```json
{
  "success": true
}
```

---

## Tenant Statistics

### Request

```http
GET /tenants/{tenantId}/statistics
```

### Response

```json
{
  "users": 500,
  "projects": 120
}
```

---

## Tenant Usage

### Request

```http
GET /tenants/{tenantId}/usage
```

### Response

```json
{
  "storage_mb": 512
}
```

---

## Activate Tenant

### Request

```http
POST /tenants/{tenantId}/activate
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
POST /tenants/{tenantId}/suspend
```

### Response

```json
{
  "success": true
}
```

---

## Tenant Domains

### Request

```http
GET /tenants/{tenantId}/domains
```

### Response

```json
{
  "domains": [
    "client.example.com"
  ]
}
```

---

# Health Check

## Request

```http
GET /health
```

## Response

```json
{
  "status": "ok",
  "version": "1.0.0"
}
```

---

# Error Format

Example:

```json
{
  "success": false,
  "message": "Validation failed",
  "errors": {
    "email": [
      "Email is required"
    ]
  }
}
```

HTTP Codes:

| Code | Meaning |
|--------|----------|
| 200 | Success |
| 201 | Created |
| 400 | Bad Request |
| 401 | Unauthorized |
| 403 | Forbidden |
| 404 | Not Found |
| 422 | Validation Error |
| 429 | Rate Limited |
| 500 | Internal Server Error |

---

# SDK Support

Official SDKs:

- PHP SDK
- Python SDK
- Node.js SDK
- Go SDK

Version:

```text
v1.0.0
```