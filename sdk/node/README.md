# HostMK Security SDK

Official Node.js SDK for HostMK Security Framework.

## Features

- Authentication API
- RBAC Management
- Permissions Management
- Roles Management
- Policies Management
- Fraud Detection API
- Multi-Tenant Management
- TypeScript Support
- ESM + CommonJS Support
- Node.js 22+

---

## Installation

```bash
npm install @hostmk/security-sdk
```

Requirements:

```text
Node.js >= 22
```

---

## Quick Start

```ts
import {
  HostMKClient,
  AuthAPI,
} from "@hostmk/security-sdk";

const client =
  new HostMKClient({
    apiKey: "your-api-key",
    baseUrl:
      "https://api.hostmk.com",
  });

const auth =
  new AuthAPI(client);

const result =
  await auth.login({
    email:
      "admin@example.com",
    password:
      "password",
  });

console.log(result);
```

---

# Client

## Create Client

```ts
import {
  HostMKClient,
} from "@hostmk/security-sdk";

const client =
  new HostMKClient({
    apiKey:
      process.env.HOSTMK_API_KEY!,
    baseUrl:
      process.env.HOSTMK_API_URL!,
  });
```

---

# Authentication

```ts
import {
  AuthAPI,
} from "@hostmk/security-sdk";

const auth =
  new AuthAPI(client);
```

## Login

```ts
await auth.login({
  email:
    "admin@example.com",
  password:
    "password",
});
```

## Refresh Token

```ts
await auth.refreshToken({
  refreshToken:
    "refresh-token",
});
```

## Validate Token

```ts
await auth.validateToken(
  "jwt-token",
);
```

## Current User

```ts
await auth.me();
```

## Logout

```ts
await auth.logout();
```

---

# Permissions API

```ts
import {
  PermissionsAPI,
} from "@hostmk/security-sdk";

const permissions =
  new PermissionsAPI(
    client,
  );
```

## Create Permission

```ts
await permissions.createPermission(
  {
    tenantId:
      "tenant-1",
    name:
      "users.read",
    resource:
      "users",
    action:
      "read",
  },
);
```

## Get Permission

```ts
await permissions.getPermission(
  "users.read",
);
```

## List Permissions

```ts
await permissions.listPermissions();
```

## Delete Permission

```ts
await permissions.deletePermission(
  "users.read",
);
```

---

# Roles

## Create Role

```ts
await permissions.createRole({
  tenantId:
    "tenant-1",
  name:
    "admin",
});
```

## Get Role

```ts
await permissions.getRole(
  "admin",
);
```

## List Roles

```ts
await permissions.listRoles();
```

## Delete Role

```ts
await permissions.deleteRole(
  "admin",
);
```

---

# Role Permissions

## Assign Permission

```ts
await permissions.assignPermissionToRole(
  "admin",
  "users.read",
);
```

## Remove Permission

```ts
await permissions.removePermissionFromRole(
  "admin",
  "users.read",
);
```

## Check Permission

```ts
await permissions.roleHasPermission(
  "admin",
  "users.read",
);
```

---

# Parent Roles

## Assign Parent Role

```ts
await permissions.assignParentRole(
  "admin",
  "manager",
);
```

## Remove Parent Role

```ts
await permissions.removeParentRole(
  "admin",
  "manager",
);
```

---

# Policies

## Create Policy

```ts
await permissions.createPolicy({
  tenantId:
    "tenant-1",
  name:
    "allow-users",
  effect:
    "allow",
  resources: [
    "users",
  ],
  actions: [
    "read",
  ],
});
```

## Get Policy

```ts
await permissions.getPolicy(
  "policy-id",
);
```

## List Policies

```ts
await permissions.listPolicies();
```

## Delete Policy

```ts
await permissions.deletePolicy(
  "policy-id",
);
```

---

# Fraud API

```ts
import {
  FraudAPI,
} from "@hostmk/security-sdk";

const fraud =
  new FraudAPI(client);
```

## Create Alert

```ts
await fraud.createAlert({
  tenantId:
    "tenant-1",
  userId:
    "user-1",
  riskScore:
    95,
  reason:
    "suspicious login",
});
```

## Get Alert

```ts
await fraud.getAlert(
  "alert-id",
);
```

## List Alerts

```ts
await fraud.listAlerts();
```

## Update Risk Score

```ts
await fraud.updateRiskScore(
  "alert-id",
  {
    riskScore: 100,
  },
);
```

## Resolve Alert

```ts
await fraud.resolveAlert(
  "alert-id",
);
```

## Reopen Alert

```ts
await fraud.reopenAlert(
  "alert-id",
);
```

---

# TypeScript

The SDK is fully typed.

```ts
import type {
  Permission,
  Role,
  Policy,
  FraudAlert,
  Tenant,
} from "@hostmk/security-sdk";
```

---

# Build

```bash
npm install
npm run typecheck
npm run test
npm run build
```

---

# Test

```bash
npm run test
```

Current test coverage:

```text
client.test.ts
auth.test.ts
permissions.test.ts
fraud.test.ts
```

40 tests passing.

---

# License

MIT

---

# Repository

https://github.com/azezaz/HostMK-Security-Framework