# Contributing to HostMK Security Framework

Thanks for your interest in contributing to HostMK Security Framework. We welcome high-quality improvements to documentation, governance, security hardening, tests, and developer tooling.

This repository is intentionally split between open-source core components and enterprise-only code under ee/. The contribution process below applies to the open-source project unless a specific enterprise contribution path is explicitly authorized.

## Contribution goals

We value contributions that improve:

- security correctness
- auditability and observability
- multi-tenant safety
- RBAC and authorization design
- documentation clarity
- test coverage and build stability

## Repository scope

The open-source project is located in the main repository structure, while enterprise-only components remain under:

```text
ee/
```

Do not contribute proprietary or enterprise-only code unless explicitly permitted by the maintainers.

## Git workflow

Use a clear branch-based workflow with short-lived feature branches.

```bash
git checkout -b fix/tenant-validation
git commit -m "fix: validate tenant boundaries"
git push origin fix/tenant-validation
```

Recommended branch patterns:

- feature/*
- fix/*
- docs/*
- chore/*
- security/*

Pull requests should be small and focused and should document the purpose of the change clearly.

## Before opening a PR

Before submitting a contribution, please confirm the following:

1. the change is limited to the requested scope
2. documentation and comments are updated when needed
3. the code still builds and passes tests
4. the change does not expand scope beyond the issue or ticket
5. secrets or sensitive configuration are not committed

## Required validation

All contributions should pass the project checks before merge.

```bash
go fmt ./...
go build ./...
go test ./...
```

If the change affects security or authorization behavior, also review the relevant modules and tests for tenant, RBAC, permissions, or audit behavior.

## Go style and coding standards

Please keep code consistent with the repository’s Go conventions:

- prefer clear, explicit logic over clever shortcuts
- keep functions small and readable
- return explicit errors instead of silent failures
- validate security-sensitive inputs early
- avoid introducing broad side effects without justification
- keep concurrency and tenant boundaries easy to reason about

## Testing requirements

Contributions should include tests when new behavior is introduced or a bug is fixed.

Examples:

- unit tests for validation or policy behavior
- regression tests for tenant isolation issues
- audit or fraud event tests when changing event models

Test commands:

```bash
go test ./...
```

## Documentation requirements

Documentation must be updated when a contribution changes:

- public behavior
- configuration requirements
- module status
- API routes or contracts
- security guidance
- architecture or roadmap notes

Relevant documentation includes:

- [README.md](README.md)
- [docs/API.md](docs/API.md)
- [docs/architecture.md](docs/architecture.md)
- [SECURITY.md](SECURITY.md)
- [SUPPORT.md](SUPPORT.md)

## Security contributions

If the contribution is related to a vulnerability, cryptography, authorization, or compliance issue, please follow [SECURITY.md](SECURITY.md) and report the issue privately.

Do not disclose security issues publicly before coordination with maintainers.

## Commit guidance

Use clear and conventional commit messages.

Examples:

```text
feat: add tenant validation middleware
fix: prevent cross-tenant access in RBAC checks
docs: align README and support policy
security: harden secret handling validation
test: add regression coverage for fraud scoring
```

## Pull request checklist

Before opening or updating a PR, confirm:

- [ ] code is formatted
- [ ] go build ./... passes
- [ ] go test ./... passes
- [ ] docs are updated where necessary
- [ ] no secrets or credentials are included
- [ ] security-sensitive changes are reviewed

## Code of conduct

Contributors are expected to behave respectfully and professionally. Disruptive, abusive, or malicious behavior is not acceptable.

See [CODE_OF_CONDUCT.md](CODE_OF_CONDUCT.md) for the community expectations.

## Support and governance

For support requests, bug reports, and operational questions, please review [SUPPORT.md](SUPPORT.md). For responsible vulnerability disclosure, use [SECURITY.md](SECURITY.md).

## License

By contributing to this project, you agree that your work may be distributed under the repository’s applicable license terms. Open-source contributions are covered by the Apache 2.0 license for the non-enterprise code paths, while enterprise-only files remain governed by the enterprise license model.

Thank you for helping keep the framework secure, maintainable, and professional.