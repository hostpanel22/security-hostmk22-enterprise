# HostMK Security Framework

![Go](https://img.shields.io/badge/Go-1.26%2B-00ADD8?logo=go)
![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)
![Status](https://img.shields.io/badge/Status-Active-success)

HostMK Security Framework is a modular security toolkit for SaaS platforms, multi-tenant applications, APIs, and enterprise systems. The repository contains core security primitives, authorization building blocks, audit capabilities, fraud controls, and secret-management utilities designed for production use.

## Overview

This project focuses on the security controls most often required by modern platform teams:

- multi-tenant isolation and validation
- role and permission enforcement
- audit event capture and export
- fraud detection and violation scoring
- secret storage and rotation primitives
- HTTP API protections and middleware
- JWT-based authentication support

## Architecture

```mermaid
flowchart LR
    A[Applications] --> B[API Layer]
    B --> C[Auth & JWT]
    B --> D[Tenancy]
    B --> E[RBAC / Permissions]
    B --> F[Audit]
    B --> G[Fraud]
    B --> H[Secrets]
    D --> I[Database / Persistence]
    E --> I
    F --> I
    G --> I
    H --> I
```

## Repository layout

```text
HostMK-Security-Framework/
├── cmd/                     # application entrypoints
├── core/                    # OSS security modules
│   ├── audit/
│   ├── fraud/
│   ├── permissions/
│   ├── rbac/
│   ├── secrets/
│   └── tenancy/
├── internal/                # application implementation details
├── docs/                    # project documentation
├── ee/                      # enterprise-only components
├── sdk/                     # language SDKs
├── tests/                   # Go test suite
├── .github/                 # governance and CI workflows
├── LICENSE
├── NOTICE
├── SECURITY.md
├── CONTRIBUTING.md
├── SUPPORT.md
├── README.md
└── go.mod
```

## Features

| Area | Status | Notes |
| --- | --- | --- |
| Multi-tenancy | ✅ Active | Isolation, validation, tenant context, owner checks |
| Permissions | ✅ Active | Resource checks and policy enforcement |
| RBAC | ✅ Active | Role, permission, and policy abstractions |
| Audit | ✅ Active | Event capture, filters, retention, export |
| Fraud detection | ✅ Active | Signals, scoring, rules, violations |
| Secret management | ✅ Active | Encryption, storage, rotation, vault patterns |
| JWT / authentication | ✅ Active | Token generation and validation support |
| API protections | ✅ Active | Middleware, rate limiting, security headers, health endpoints |
| Enterprise extensions | ⚠️ Separate license | Located under ee/ |

## RBAC

The RBAC package is implemented in the repository under the core/rbac module. It includes:

- roles
- permissions
- policies
- access evaluation
- middleware integration
- tenant-aware authorization primitives

## Audit

The audit module provides structured event logging, filtering, retention, and export helpers for security-relevant activity.

## Fraud Detection

The fraud package includes rules, violations, scores, and middleware hooks to monitor suspicious activity and reduce risk exposure.

## Secrets Management

The secrets module focuses on encryption, secret storage, rotation models, and vault-oriented patterns for sensitive data management.

## Multi-tenancy

The tenancy layer enforces boundaries between tenant contexts, preventing unauthorized cross-tenant access and validating ownership and isolation rules.

## Installation

### Requirements

- Go 1.26+
- Git
- PostgreSQL-compatible database for runtime usage

### Clone

```bash
git clone https://github.com/azezaz/HostMK-Security-Framework.git
cd HostMK-Security-Framework
```

### Dependency setup

```bash
go mod tidy
```

### Build

```bash
go build ./...
```

### Test

```bash
go test ./...
```

## Configuration

The server configuration is loaded from environment variables and a local .env file when present. The default runtime expects values such as:

```env
APP_NAME=HostMK Security Framework
APP_ENV=development
HOST=0.0.0.0
PORT=8080
DATABASE_URL=postgres://user:password@localhost:5432/hostmk
JWT_SECRET=replace-with-secret
ENCRYPTION_KEY=replace-with-32-byte-key
```

The runtime configuration for JWT, database access, and encryption keys is validated in the application startup path.

## Database

The project uses a PostgreSQL-oriented persistence layer and includes migration support through the database package. Runtime configuration is provided through DATABASE_URL, and database migration is performed during startup in the server entrypoint.

## API

The application exposes HTTP endpoints under the /api/v1 namespace. Core routes include:

| Endpoint | Purpose |
| --- | --- |
| /health | health check |
| /ready | readiness check |
| /api/v1/version | version metadata |
| /api/v1/auth/bootstrap | bootstrap auth state |
| /api/v1/auth/login | login |
| /api/v1/auth/refresh | refresh token |
| /api/v1/users | user resources |
| /api/v1/tenants | tenant resources |
| /api/v1/rbac/roles | RBAC roles |
| /api/v1/permissions | permission management |
| /api/v1/secrets | secret resources |
| /api/v1/audit/events | audit events |
| /api/v1/fraud/rules | fraud rules |
| /api/v1/settings | settings |

For the full API contract and examples, see [docs/API.md](docs/API.md).

## Roadmap

### Current focus

- stabilization of core security primitives
- RBAC and permissions documentation alignment
- security policy and support process maturity
- SDK and documentation parity

### Planned enhancements

- advanced auth flows
- expanded rate limiting and webhook controls
- stronger operational telemetry and security automation
- broader enterprise feature packaging

## License

The core repository is licensed under the Apache License 2.0, except for files under the enterprise directory ee/, which are governed by the repository’s enterprise license terms.

- [LICENSE](LICENSE)
- [NOTICE](NOTICE)
- [ee/LICENSE](ee/LICENSE)

## Security

Please report vulnerabilities privately as described in [SECURITY.md](SECURITY.md). Do not open public issues for security problems.

## Support

For support requests, issues, and community questions, see [SUPPORT.md](SUPPORT.md).

## Contributing

Contributions are welcome. Please review [CONTRIBUTING.md](CONTRIBUTING.md) before opening a pull request.

## License and governance

This repository combines open-source and proprietary distribution boundaries:

- open-source components live under the main project tree and are Apache 2.0 licensed
- enterprise components are kept under the ee/ directory and are licensed separately
- legal notices and contributor rules are documented in the repository root

---

HostMK Security Framework © 2025
Secure by design. Multi-tenant by default. Enterprise ready.