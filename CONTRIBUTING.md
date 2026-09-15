# Contributing to HostMK Security Framework

First of all, thank you for your interest in contributing to HostMK Security Framework. 🚀

This project aims to provide a modern, secure, multi-tenant security framework for SaaS platforms, APIs, cloud services, and enterprise applications.

We welcome contributions from developers, security researchers, DevOps engineers, and open-source enthusiasts.

---

# 🎯 Project Mission

HostMK Security Framework focuses on:

- Multi-Tenant Security
- Access Control
- Secret Management
- Authentication & Authorization
- Security Middleware
- Fraud Detection
- Audit Logging
- Threat Detection
- Enterprise Security Extensions

Our goal is to build a production-ready security framework that follows industry best practices and modern security standards.

---

# 🤝 Ways to Contribute

You can contribute in many ways:

## Code Contributions

Examples:

- Bug fixes
- Security improvements
- New modules
- Performance optimizations
- Middleware development
- Unit tests
- Documentation improvements

---

## Documentation

Examples:

- README improvements
- Security documentation
- API examples
- Architecture diagrams
- Tutorials

---

## Security Research

Examples:

- Vulnerability reports
- Security reviews
- Threat modeling
- Cryptography reviews
- Penetration testing feedback

---

## Testing

Examples:

- Unit testing
- Integration testing
- Performance testing
- Security testing
- Cross-platform validation

---

# 📋 Before You Start

Before making changes:

1. Read the documentation.
2. Review the project architecture.
3. Check open issues.
4. Verify that your contribution aligns with project goals.
5. Ensure changes follow the project's coding standards.

---

# 🔧 Development Environment

Requirements:

```text
Go 1.22+
Git
```

Clone the repository:

```bash
git clone https://github.com/your-org/HostMK-Security-Framework.git
```

Move into the project:

```bash
cd HostMK-Security-Framework
```

Install dependencies:

```bash
go mod tidy
```

---

# 🏗 Project Structure

```text
HostMK-Security-Framework/
│
├── core/
│   ├── audit/
│   ├── fraud/
│   ├── permissions/
│   ├── secrets/
│   ├── tenancy/
│   └── ...
│
├── tests/
│
├── docs/
│
├── ee/
│
├── LICENSE
├── NOTICE
└── README.md
```

---

# 🔓 Open Source vs Enterprise

## Open Source

Located in:

```text
core/
```

Licensed under:

```text
Apache License 2.0
```

Community contributions are welcome.

---

## Enterprise Edition

Located in:

```text
ee/
```

Governed by:

```text
ee/LICENSE
```

Enterprise Edition code is proprietary.

Contributors may not submit Enterprise Edition code unless explicitly authorized by HostMK.

---

# 💻 Coding Standards

All code should follow:

## General Rules

- Keep code simple.
- Keep code readable.
- Avoid unnecessary complexity.
- Prefer explicit behavior.
- Write self-documenting code.

---

## Naming

Use clear and descriptive names.

Good:

```go
ValidateTenant()
EncryptSecret()
HasPermission()
```

Avoid:

```go
Do()
Run()
Process()
```

---

## Error Handling

Always return explicit errors.

Example:

```go
if tenantID == "" {
    return ErrEmptyTenantID
}
```

Avoid silent failures.

---

## Security First

Security takes priority over convenience.

Examples:

- Validate inputs.
- Sanitize data.
- Avoid insecure defaults.
- Use strong cryptography.
- Follow least privilege principles.

---

# 🧪 Testing Requirements

Every new feature must include tests.

Required:

```text
Feature Code
+
Unit Tests
```

Example:

```text
core/rbac/role.go
tests/rbac/role_test.go
```

---

## Run Tests

Run all tests:

```bash
go test ./...
```

Run package tests:

```bash
go test ./tests/secrets
```

---

## Coverage

New functionality should include meaningful test coverage.

Contributors should test:

- Success paths
- Error paths
- Edge cases
- Security scenarios

---

# 🔍 Static Analysis

Before submitting a contribution:

Format code:

```bash
go fmt ./...
```

Run vet:

```bash
go vet ./...
```

Build:

```bash
go build ./...
```

Recommended:

```bash
gosec ./...
```

Recommended:

```bash
govulncheck ./...
```

---

# 🚨 Security Contributions

Security-related contributions are highly encouraged.

Examples:

- Authentication hardening
- Access control improvements
- Cryptography reviews
- Security middleware
- Threat detection
- Vulnerability fixes

---

# 🔒 Reporting Security Vulnerabilities

Do NOT open public GitHub issues for security vulnerabilities.

Instead contact:

```text
security@hostmk.com
```

or

```text
contact@hostmk.com
```

See:

```text
docs/security.md
```

for the complete security policy.

---

# 📝 Commit Message Guidelines

Use clear commit messages.

Examples:

```text
feat: add rbac role validation
```

```text
feat: implement api key manager
```

```text
fix: prevent cross tenant access
```

```text
fix: validate secret rotation policy
```

```text
docs: update security documentation
```

```text
test: add vault integration tests
```

---

# 🔄 Pull Request Guidelines

Before opening a Pull Request:

- Ensure tests pass.
- Ensure code builds successfully.
- Ensure formatting is correct.
- Ensure documentation is updated.
- Ensure no secrets are committed.

Checklist:

```text
[ ] Tests pass
[ ] Build passes
[ ] Documentation updated
[ ] No secrets committed
[ ] Security reviewed
```

---

# 📚 Documentation Requirements

Documentation should be updated when:

- New modules are added.
- Public APIs change.
- Security features change.
- Architecture changes.

Relevant files:

```text
README.md
docs/architecture.md
docs/security.md
docs/roadmap.md
```

---

# ⚖ License

By contributing to this project, you agree that your contributions may be distributed under the project's applicable licenses.

Open Source contributions are governed by:

```text
Apache License 2.0
```

Enterprise Edition components remain subject to:

```text
ee/LICENSE
```

---

# 🙏 Code of Conduct

Be respectful.

We do not tolerate:

- Harassment
- Hate speech
- Discrimination
- Malicious contributions
- Intentional security backdoors

Professional and constructive collaboration is expected.

---

# 🚀 Thank You

Every contribution helps improve the security and quality of HostMK Security Framework.

Together we can build a modern, secure, enterprise-grade security ecosystem for developers worldwide.

---

🔐 HostMK Security Framework

Secure by Design.
Defense in Depth.
Multi-Tenant by Default.
Enterprise Ready.