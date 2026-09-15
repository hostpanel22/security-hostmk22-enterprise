# 🔐 HostMK Security Framework

Enterprise-Grade Security Framework for SaaS, Multi-Tenant Applications, APIs, Platforms, and Cloud Services.

HostMK Security Framework is an open-source security toolkit designed to help developers integrate modern security controls into their applications while maintaining a clear separation between community and enterprise features.

---

## 🎯 Project Goals

HostMK Security Framework aims to provide:

- Multi-Tenant Security
- Role-Based Access Control (RBAC)
- Secret Management
- API Security
- Fraud Detection
- Audit Logging
- Security Middleware
- Authentication & Authorization
- Enterprise Security Extensions

The framework is designed for:

- SaaS Platforms
- Marketplace Applications
- Fintech Products
- Internal Business Platforms
- Cloud-Native Services
- White-Label Systems
- Microservice Architectures

---

# 🚀 Current Modules

## ✅ Tenancy

Multi-tenant isolation and validation.

Features:

- Tenant validation
- Tenant ownership checks
- Cross-tenant protection
- Tenant middleware support

Location:

```text
core/tenancy/
```

---

## ✅ Permissions

Permission management system.

Features:

- Permission definitions
- Permission validation
- Permission checking
- Access verification

Location:

```text
core/permissions/
```

---

## ✅ Audit

Security audit system.

Features:

- Audit events
- Audit logging
- Security event tracking
- Activity monitoring

Location:

```text
core/audit/
```

---

## ✅ Fraud

Fraud detection primitives.

Features:

- Risk evaluation
- Fraud indicators
- Event scoring
- Security monitoring

Location:

```text
core/fraud/
```

---

## ✅ Secrets

Enterprise-grade secret management.

Features:

- AES-256 Encryption
- Secret Vault
- Secret Rotation
- Secret Middleware
- Secret Storage
- Multi-Tenant Secret Isolation

Location:

```text
core/secrets/
```

---

# 🏗 Architecture

```text
HostMK-Security-Framework/
│
├── core/
│   ├── audit/
│   ├── fraud/
│   ├── permissions/
│   ├── secrets/
│   └── tenancy/
│
├── tests/
│
├── docs/
│   ├── architecture.md
│   ├── roadmap.md
│   ├── contributing.md
│   └── security.md
│
├── ee/
│   ├── LICENSE
│   ├── README.md
│   ├── sso_integration/
│   ├── audit_logs/
│   ├── compliance/
│   └── advanced_security/
│
├── LICENSE
├── NOTICE
└── README.md
```

---

# 🔓 Open Source Edition

The Open Source Edition includes all modules located inside:

```text
core/
```

Licensed under:

```text
Apache License 2.0
```

You may:

- Use commercially
- Modify
- Distribute
- Fork
- Build products using the framework

subject to Apache 2.0 terms.

---

# 🔒 Enterprise Edition

The Enterprise Edition is located in:

```text
ee/
```

Enterprise modules are NOT covered by Apache 2.0.

Examples:

- SSO Integrations
- Advanced Audit Reporting
- Compliance Controls
- Enterprise Monitoring
- Premium Security Features

Enterprise components require a separate commercial license.

See:

```text
ee/LICENSE
```

for details.

---

# 🧪 Testing

Run all tests:

```bash
go test ./...
```

Run formatting:

```bash
go fmt ./...
```

Run static analysis:

```bash
go vet ./...
```

Build project:

```bash
go build ./...
```

---

# 📊 Project Status

| Module | Status |
|----------|----------|
| Tenancy | ✅ Stable |
| Permissions | ✅ Stable |
| Audit | ✅ Stable |
| Fraud | ✅ Stable |
| Secrets | ✅ Stable |
| RBAC | 🚧 Planned |
| JWT/Auth | 🚧 Planned |
| API Keys | 🚧 Planned |
| Rate Limiting | 🚧 Planned |
| MFA/TOTP | 🚧 Planned |
| Session Security | 🚧 Planned |
| Webhook Security | 🚧 Planned |
| Security Headers | 🚧 Planned |
| IP Reputation | 🚧 Planned |
| Threat Detection Engine | 🚧 Planned |

---

# 🛣 Roadmap

### Phase 1

- Tenancy
- Permissions
- Audit
- Fraud
- Secrets

### Phase 2

- RBAC
- JWT Authentication
- API Keys
- Security Middleware

### Phase 3

- MFA/TOTP
- Session Security
- Security Headers
- Webhook Security

### Phase 4

- Threat Detection Engine
- IP Reputation
- Enterprise Security Extensions

---

# 🤝 Contributing

Contributions are welcome.

Before submitting changes:

1. Format code

```bash
go fmt ./...
```

2. Run tests

```bash
go test ./...
```

3. Run vet

```bash
go vet ./...
```

4. Ensure all checks pass.

---

# ⚖ License

Open Source components:

```text
Apache License 2.0
```

Enterprise components:

```text
HostMK Enterprise License
```

See:

```text
LICENSE
ee/LICENSE
NOTICE
```

---

# 📧 Contact

Commercial Licensing:

```text
contact@hostmk.com
```

---

# 🔐 HostMK Security Framework

Secure by Design.
Multi-Tenant by Default.
Enterprise Ready.
Built for Modern SaaS.