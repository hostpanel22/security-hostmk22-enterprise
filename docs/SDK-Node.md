# HostMK Security Framework Node.js SDK

## Overview

The HostMK Security Node.js SDK provides a modern and lightweight interface for integrating HostMK Security services into JavaScript and TypeScript applications.

Supported modules:

- Authentication
- Audit Logging
- Fraud Detection
- Permission Management
- Multi-Tenant Management

Supported runtimes:

```text
Node.js 18+
Node.js 20+
Node.js 22+
```

Current Version:

```text
v1.0.0
```

---

# Installation

## NPM

```bash
npm install @hostmk/security-sdk
```

## Yarn

```bash
yarn add @hostmk/security-sdk
```

## PNPM

```bash
pnpm add @hostmk/security-sdk
```

---

# Quick Start

```javascript
import { HostMKClient } from "@hostmk/security-sdk";

const client = new HostMKClient({
    apiKey: "YOUR_API_KEY",
    baseUrl: "https://api.hostmk.io/v1"
});
```

---

# Authentication Service

```javascript
import {
    HostMKClient,
    AuthService
} from "@hostmk/security-sdk";

const client = new HostMKClient({
    apiKey: process.env.HOSTMK_API_KEY,
    baseUrl: "https://api.hostmk.io/v1"
});

const auth = new AuthService(client);
```

## Login

```javascript
const result = await auth.login(
    "admin@example.com",
    "password123"
);
```

## Validate Token

```javascript
const result = await auth.validateToken(
    token
);
```

## Refresh Token

```javascript
const result = await auth.refreshToken(
    refreshToken
);
```

## Logout

```javascript
const result = await auth.logout();
```

## Sessions

```javascript
const result = await auth.sessions();
```

## Revoke Session

```javascript
const result = await auth.revokeSession(
    sessionId
);
```

## Health

```javascript
const result = await auth.health();
```

---

# Audit Service

```javascript
import { AuditService } from "@hostmk/security-sdk";

const audit = new AuditService(client);
```

## List Events

```javascript
const result = await audit.listEvents();
```

## Get Event

```javascript
const result = await audit.getEvent(
    "event-id"
);
```

## Search Events

```javascript
const result = await audit.searchEvents({
    action: "login"
});
```

## User Activity

```javascript
const result = await audit.userActivity(
    "user-id"
);
```

## Tenant Activity

```javascript
const result = await audit.tenantActivity(
    "tenant-id"
);
```

## Security Events

```javascript
const result = await audit.securityEvents();
```

## Security Events By Severity

```javascript
const result = await audit.securityEvents(
    "critical"
);
```

## Statistics

```javascript
const result = await audit.statistics();
```

## Compliance Report

```javascript
const result = await audit.complianceReport(
    "gdpr"
);
```

## Export Events

```javascript
const result = await audit.exportEvents();
```

---

# Fraud Service

```javascript
import { FraudService } from "@hostmk/security-sdk";

const fraud = new FraudService(client);
```

## Analyze

```javascript
const result = await fraud.analyze({
    email: "user@example.com",
    ip: "1.1.1.1"
});
```

## Score

```javascript
const result = await fraud.score({
    email: "user@example.com"
});
```

## Check Email

```javascript
const result = await fraud.checkEmail(
    "user@example.com"
);
```

## Check IP

```javascript
const result = await fraud.checkIp(
    "1.1.1.1"
);
```

## Check Device

```javascript
const result = await fraud.checkDevice(
    "device-id"
);
```

## Blacklist

```javascript
const result = await fraud.blacklist(
    "email",
    "blocked@example.com"
);
```

## Whitelist

```javascript
const result = await fraud.whitelist(
    "email",
    "trusted@example.com"
);
```

## Alerts

```javascript
const result = await fraud.alerts();
```

## Statistics

```javascript
const result = await fraud.statistics();
```

## Report

```javascript
const result = await fraud.report({
    type: "abuse"
});
```

---

# Permission Service

```javascript
import { PermissionService } from "@hostmk/security-sdk";

const permissions = new PermissionService(
    client
);
```

## List Permissions

```javascript
const result =
    await permissions.listPermissions();
```

## Create Permission

```javascript
const result =
    await permissions.createPermission(
        "users.create"
    );
```

## Get Permission

```javascript
const result =
    await permissions.getPermission(
        "permission-id"
    );
```

## Update Permission

```javascript
const result =
    await permissions.updatePermission(
        "permission-id",
        {
            name: "users.update"
        }
    );
```

## Delete Permission

```javascript
const result =
    await permissions.deletePermission(
        "permission-id"
    );
```

## Assign Permission

```javascript
const result =
    await permissions.assignPermission(
        "user-id",
        "users.create"
    );
```

## Revoke Permission

```javascript
const result =
    await permissions.revokePermission(
        "user-id",
        "users.create"
    );
```

## Bulk Assign

```javascript
const result =
    await permissions.bulkAssign(
        ["u1", "u2", "u3"],
        "users.create"
    );
```

## Bulk Revoke

```javascript
const result =
    await permissions.bulkRevoke(
        ["u1", "u2"],
        "users.create"
    );
```

## User Permissions

```javascript
const result =
    await permissions.userPermissions(
        "user-id"
    );
```

## Role Permissions

```javascript
const result =
    await permissions.rolePermissions(
        "role-id"
    );
```

---

# Tenancy Service

```javascript
import { TenancyService } from "@hostmk/security-sdk";

const tenancy =
    new TenancyService(client);
```

## List Tenants

```javascript
const result =
    await tenancy.listTenants();
```

## Create Tenant

```javascript
const result =
    await tenancy.createTenant(
        "Acme Inc",
        "acme"
    );
```

## Get Tenant

```javascript
const result =
    await tenancy.getTenant(
        "tenant-id"
    );
```

## Update Tenant

```javascript
const result =
    await tenancy.updateTenant(
        "tenant-id",
        {
            name: "Acme Corporation"
        }
    );
```

## Delete Tenant

```javascript
const result =
    await tenancy.deleteTenant(
        "tenant-id"
    );
```

## Suspend Tenant

```javascript
const result =
    await tenancy.suspendTenant(
        "tenant-id"
    );
```

## Activate Tenant

```javascript
const result =
    await tenancy.activateTenant(
        "tenant-id"
    );
```

## Usage

```javascript
const result =
    await tenancy.tenantUsage(
        "tenant-id"
    );
```

## Statistics

```javascript
const result =
    await tenancy.tenantStatistics(
        "tenant-id"
    );
```

## Domains

```javascript
const result =
    await tenancy.tenantDomains(
        "tenant-id"
    );
```

## Settings

```javascript
const result =
    await tenancy.tenantSettings(
        "tenant-id"
    );
```

## Update Settings

```javascript
const result =
    await tenancy.updateSettings(
        "tenant-id",
        {
            theme: "dark"
        }
    );
```

---

# TypeScript Support

```typescript
import {
    HostMKClient,
    AuthService,
    AuditService,
    FraudService,
    PermissionService,
    TenancyService
} from "@hostmk/security-sdk";
```

Example:

```typescript
const client =
    new HostMKClient({
        apiKey: process.env.HOSTMK_API_KEY!,
        baseUrl: process.env.HOSTMK_BASE_URL!
    });

const auth =
    new AuthService(client);

const response =
    await auth.login(
        "admin@example.com",
        "password"
    );
```

---

# Error Handling

```javascript
try {
    await auth.login(
        "admin@example.com",
        "password"
    );
} catch (error) {
    console.error(error.message);
}
```

---

# Environment Variables

```bash
HOSTMK_API_KEY=your_api_key
HOSTMK_BASE_URL=https://api.hostmk.io/v1
```

---

# Architecture

```text
HostMKClient
│
├── AuthService
├── AuditService
├── FraudService
├── PermissionService
└── TenancyService
```

---

# Security Best Practices

- Store API keys in environment variables.
- Never commit secrets to Git.
- Use HTTPS only.
- Rotate API keys regularly.
- Restrict API keys to required environments.

---

# Versioning

HostMK follows Semantic Versioning.

Example:

```text
1.0.0
1.1.0
2.0.0
```

---

# License

MIT License

Copyright (c) HostMK Security Framework