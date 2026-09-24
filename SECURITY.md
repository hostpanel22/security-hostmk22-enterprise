# Security Policy

## Overview

HostMK Security Framework takes security seriously and expects all reports to be handled privately and responsibly. This policy covers supported versions, vulnerability reporting, and coordinated disclosure expectations.

## Supported versions

The project currently supports the latest release line and the actively maintained branch. Legacy versions are not guaranteed to receive security patches.

| Version / branch | Status |
| --- | --- |
| latest release | ✅ supported |
| active main branch | ✅ supported |
| previous release line | ⚠️ limited support |
| unsupported or archived versions | ❌ not supported |

## Reporting a vulnerability

If you identify a vulnerability, please report it privately through one of the supported channels and do not disclose it publicly before we have had a chance to assess and remediate it.

### Do not

- open a public GitHub issue for the vulnerability
- publish exploit details or proof-of-concept code
- disclose the issue in public forums or social channels
- attempt unauthorized access to user data or production systems

### Please include

- title and summary
- affected component or module
- reproduction steps or proof of concept
- impact assessment
- affected version or commit
- suggested remediation, if known

### Contact channels

```text
security@hostmk.com
contact@hostmk.com
```

You may also use the private security reporting flow in GitHub if enabled for the repository and tag the issue as a confidential security report.

## Responsible disclosure

We ask researchers and contributors to:

1. report privately and not publicly disclose the issue before a fix is available
2. avoid modifying or deleting user data during validation
3. minimize impact on services and systems under test
4. act in good faith and provide sufficient detail for triage
5. cooperate with a reasonable remediation timeline

## Remediation expectations

The maintainers will review reports as quickly as possible and will coordinate remediation, release updates, and any required public disclosure in a responsible manner. In some cases, publication of a security advisory may occur after the fix is released.

## Security best practices

This project relies on the following principles:

- least privilege
- tenant isolation
- encryption for sensitive data
- secure configuration and secret handling
- code review for security-sensitive changes
- continuous validation through tests and CI

## Scope

This policy applies to the open-source project as well as any repository-managed security-related processes. Enterprise-only features placed under the ee/ directory remain subject to the separate enterprise licensing and support model defined by that code base.

## Acknowledgment

We appreciate responsible security research and collaboration with the community. Security reports that respect the policy above are welcome and will be handled with confidentiality and care.

---

# 🔍 Vulnerability Response Process

Once a report is received:

### Step 1

Acknowledgement

Target:

```text
Within 72 hours
```

### Step 2

Validation

The vulnerability is reproduced and assessed.

### Step 3

Risk Classification

Severity categories:

- Critical
- High
- Medium
- Low

### Step 4

Fix Development

A patch is created and tested.

### Step 5

Release

Security update is published.

### Step 6

Disclosure

A security advisory may be published.

---

# 📊 Severity Classification

## Critical

Examples:

- Remote code execution
- Authentication bypass
- Tenant isolation bypass
- Secret disclosure

Target Response:

```text
24-72 hours
```

---

## High

Examples:

- Privilege escalation
- Sensitive data exposure
- Access control bypass

Target Response:

```text
3-7 days
```

---

## Medium

Examples:

- Limited information disclosure
- Security misconfiguration

Target Response:

```text
7-14 days
```

---

## Low

Examples:

- Minor hardening issues
- Non-exploitable weaknesses

Target Response:

```text
Best effort
```

---

# 🔑 Cryptography Policy

HostMK Security Framework follows modern cryptographic standards.

Approved algorithms include:

- AES-256-GCM
- SHA-256
- SHA-512
- HMAC-SHA256
- HMAC-SHA512

Weak algorithms should not be used:

- MD5
- SHA1
- DES
- RC4

---

# 🏢 Enterprise Security

Enterprise Edition modules may include:

- SSO Integration
- Advanced Audit Reporting
- Compliance Controls
- Enterprise Monitoring
- Advanced Security Analytics

Enterprise security modules are governed by:

```text
ee/LICENSE
```

and are not covered by the Apache 2.0 license.

---

# 🧪 Security Testing

Before every release:

```bash
go fmt ./...
```

```bash
go test ./...
```

```bash
go vet ./...
```

```bash
go build ./...
```

Recommended additional tools:

```bash
gosec ./...
```

```bash
govulncheck ./...
```

---

# 🔐 Security Best Practices

Applications using HostMK Security Framework should:

- Use HTTPS everywhere
- Encrypt secrets
- Rotate credentials regularly
- Enable MFA
- Apply least privilege
- Audit access logs
- Monitor suspicious activity
- Keep dependencies updated
- Validate all inputs
- Protect APIs with authentication

---

# ⚖ Legal Notice

HostMK Security Framework is provided:

```text
AS IS
```

without warranty of any kind.

Users remain responsible for:

- Infrastructure security
- Deployment security
- Key management
- Regulatory compliance
- Operational security

---

# 📧 Security Contact

Security Team

```text
security@hostmk.com
```

General Contact

```text
contact@hostmk.com
```

---

# 🔐 Security First

Secure by Design.
Defense in Depth.
Enterprise Ready.
Built for Modern SaaS.