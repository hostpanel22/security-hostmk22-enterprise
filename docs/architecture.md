# HostMK Security Framework Architecture

## Overview

HostMK Security Framework is a modular security framework written in Go.

The framework provides reusable security components for:

- SaaS Platforms
- Multi-Tenant Applications
- API Platforms
- Internal Enterprise Services
- White-Label Systems
- Cloud-Native Applications

The architecture follows a security-first design focused on:

- Tenant Isolation
- Least Privilege
- Defense in Depth
- Auditability
- Extensibility
- Secure Defaults

---

# Architecture Principles

## Modular Design

Each security component is isolated in its own package.

Benefits:

- Independent testing
- Easier maintenance
- Reusable components
- Clear responsibilities

Example:

```text
core/
├── audit
├── fraud
├── permissions
├── secrets
└── tenancy
```

---

## Security by Default

Every module should:

- Validate inputs
- Fail securely
- Reject invalid states
- Avoid insecure defaults

---

## Defense in Depth

Security controls are implemented in multiple layers.

Example:

```text
Request
   |
   v
Middleware
   |
   v
Tenant Validation
   |
   v
Permission Validation
   |
   v
Business Logic
   |
   v
Audit Logging
```

---

## Tenant Isolation

Multi-tenant isolation is a core requirement.

Every operation must respect:

```text
Tenant A ≠ Tenant B
```

Cross-tenant access must always be denied unless explicitly authorized.

---

# Project Structure

```text
HostMK-Security-Framework/

cmd/
└── server/

core/

audit/
fraud/
permissions/
secrets/
tenancy/

docs/

architecture.md
roadmap.md
threat-model.md

examples/

sdk/

tests/

audit/
fraud/
permissions/
secrets/
tenancy/
```

---

# Core Modules

## Audit Module

Directory:

```text
core/audit
```

Purpose:

Provides security auditing and event tracking.

Responsibilities:

- Security event creation
- Event validation
- Audit records
- Change tracking
- Compliance support

Example Events:

- Login
- Logout
- Permission Change
- Secret Rotation
- Security Alert

---

## Tenancy Module

Directory:

```text
core/tenancy
```

Purpose:

Provides tenant isolation.

Responsibilities:

- Tenant creation
- Tenant validation
- Tenant ownership
- Cross-tenant protection

Security Goal:

Prevent tenant data leakage.

---

## Permissions Module

Directory:

```text
core/permissions
```

Purpose:

Provides authorization capabilities.

Responsibilities:

- Permission definition
- Permission validation
- Access control
- Authorization decisions

Current Model:

```text
User
   |
   v
Permission
   |
   v
Access Decision
```

Future:

```text
User
   |
   v
Role
   |
   v
Policy
   |
   v
Permission
```

---

## Fraud Module

Directory:

```text
core/fraud
```

Purpose:

Detect suspicious activity.

Responsibilities:

- Risk scoring
- Fraud indicators
- Suspicious behavior detection
- Threat evaluation

Example Indicators:

- Excessive requests
- Geographic anomalies
- Account abuse
- Automated activity

---

## Secrets Module

Directory:

```text
core/secrets
```

Purpose:

Secure management of sensitive information.

Responsibilities:

- Secret storage
- Encryption
- Rotation
- Secret access control
- Middleware integration

Protected Data:

- API Keys
- Database Passwords
- OAuth Credentials
- Tokens
- Encryption Keys

---

# Secrets Architecture

The Secrets module is composed of several layers.

```text
Secret
   |
   v
Encryption Service
   |
   v
Vault
   |
   v
Store
```

---

## Secret Model

```text
core/secrets/secret.go
```

Responsibilities:

- Secret definition
- Validation
- Versioning
- Rotation metadata

---

## Encryption Service

```text
core/secrets/encryption.go
```

Responsibilities:

- AES-256 Encryption
- AES-256 Decryption
- Key validation
- Secret encryption

Algorithm:

```text
AES-256-GCM
```

Benefits:

- Authenticated encryption
- Integrity verification
- Confidentiality

---

## Secret Store

```text
core/secrets/store.go
```

Responsibilities:

- Persist secrets
- Retrieve secrets
- Delete secrets
- List secrets

Current Implementation:

```text
MemoryStore
```

Future Implementations:

```text
PostgreSQLStore
RedisStore
VaultStore
CloudStore
```

---

## Secret Vault

```text
core/secrets/vault.go
```

Responsibilities:

- Secret lifecycle management
- Encryption orchestration
- Secret retrieval
- Secret rotation

Operations:

```text
Create
Read
Rotate
Delete
List
Count
```

---

## Secret Rotation

```text
core/secrets/rotation.go
```

Responsibilities:

- Rotation policies
- Expiration checks
- Rotation scheduling

Default Policy:

```text
90 days
```

---

## Secret Middleware

```text
core/secrets/middleware.go
```

Responsibilities:

- Tenant extraction
- User extraction
- Context injection
- Request validation

Headers:

```text
X-Tenant-ID
X-User-ID
```

---

# Request Lifecycle

## Standard Request

```text
Client
   |
   v
HTTP Request
   |
   v
Security Middleware
   |
   v
Tenant Validation
   |
   v
Permission Validation
   |
   v
Business Logic
   |
   v
Audit Event
   |
   v
Response
```

---

## Secret Access Request

```text
Client
   |
   v
Middleware
   |
   v
Vault
   |
   v
Store
   |
   v
Encrypted Secret
   |
   v
Decrypt
   |
   v
Plain Secret
```

---

# Testing Strategy

Every module must include tests.

Current Test Directories:

```text
tests/

audit/
fraud/
permissions/
secrets/
tenancy/
```

Testing Categories:

- Unit Tests
- Validation Tests
- Error Tests
- Security Tests
- Edge Cases

Required Commands:

```bash
go test ./...
go vet ./...
go build ./...
```

All commands must pass before release.

---

# Security Layers

Current Layers:

```text
Layer 1
Tenant Isolation

Layer 2
Permissions

Layer 3
Secrets Protection

Layer 4
Fraud Detection

Layer 5
Audit Logging
```

---

# Planned Architecture Extensions

The next modules will extend the architecture.

## Identity Layer

```text
core/rbac
core/auth
core/apikeys
```

Features:

- RBAC
- JWT
- API Keys

---

## Protection Layer

```text
core/ratelimit
core/mfa
core/session
```

Features:

- Rate Limiting
- MFA
- Session Security

---

## Application Security Layer

```text
core/webhooks
core/headers
```

Features:

- Webhook Verification
- Security Headers

---

## Threat Intelligence Layer

```text
core/reputation
core/threat
```

Features:

- IP Reputation
- Threat Detection Engine

---

# Long-Term Architecture

```text
                  Applications
                         |
                         v
                 Security Middleware
                         |
      +------------------+------------------+
      |                  |                  |
      v                  v                  v

   Tenancy        Authentication        Permissions
      |                  |                  |
      +------------------+------------------+
                         |
                         v

                     Secrets

                         |
                         v

                     Auditing

                         |
                         v

                 Threat Detection

                         |
                         v

                     Analytics
```

---

# Architectural Goals

The framework is designed to provide:

- Secure by Default Components
- Strong Tenant Isolation
- Extensible Security Modules
- Enterprise Security Features
- Cloud-Native Compatibility
- High Test Coverage
- Simple Integration

The objective is to become a complete security foundation for modern SaaS and multi-tenant applications.