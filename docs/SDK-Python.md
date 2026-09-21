# HostMK Security Framework Python SDK

## Overview

The HostMK Security Python SDK provides a clean and lightweight interface for integrating HostMK Security services into Python applications.

Supported modules:

- Authentication
- Audit Logging
- Fraud Detection
- Permission Management
- Multi-Tenant Management

Supported Python versions:

```text
Python 3.9+
Python 3.10+
Python 3.11+
Python 3.12+
```

Current SDK Version:

```text
1.0.0
```

---

# Installation

## From Source

```bash
git clone https://github.com/azezaz/HostMK-Security-Framework.git
cd sdk/python
pip install -e .
```

## Install Dependencies

```bash
pip install requests
```

---

# Package Structure

```text
hostmk_security/
├── __init__.py
├── client.py
├── auth.py
├── audit.py
├── fraud.py
├── permissions.py
└── tenancy.py

tests/
├── test_client.py
├── test_auth.py
├── test_audit.py
├── test_fraud.py
├── test_permissions.py
└── test_tenancy.py
```

---

# Quick Start

```python
from hostmk_security import HostMKClient

client = HostMKClient(
    api_key="YOUR_API_KEY",
    base_url="https://api.hostmk.io/v1"
)
```

---

# Authentication Service

```python
from hostmk_security import AuthService

auth = AuthService(client)
```

## Login

```python
response = auth.login(
    "admin@example.com",
    "password123"
)
```

## Validate Token

```python
response = auth.validate_token(
    "jwt-token"
)
```

## Refresh Token

```python
response = auth.refresh_token(
    "refresh-token"
)
```

## Logout

```python
response = auth.logout()
```

## Sessions

```python
response = auth.sessions()
```

## Revoke Session

```python
response = auth.revoke_session(
    "session-id"
)
```

## Health Check

```python
response = auth.health()
```

---

# Audit Service

```python
from hostmk_security import AuditService

audit = AuditService(client)
```

## List Events

```python
response = audit.list_events()
```

## Event Details

```python
response = audit.get_event(
    "event-id"
)
```

## Search Events

```python
response = audit.search_events(
    {
        "action": "login"
    }
)
```

## User Activity

```python
response = audit.user_activity(
    "user-id"
)
```

## Tenant Activity

```python
response = audit.tenant_activity(
    "tenant-id"
)
```

## Security Events

```python
response = audit.security_events()
```

## Security Events By Severity

```python
response = audit.security_events(
    severity="critical"
)
```

## Statistics

```python
response = audit.statistics()
```

## Compliance Report

```python
response = audit.compliance_report(
    "gdpr"
)
```

## Export Events

```python
response = audit.export_events()
```

---

# Fraud Service

```python
from hostmk_security import FraudService

fraud = FraudService(client)
```

## Analyze

```python
response = fraud.analyze(
    {
        "email": "user@example.com",
        "ip": "1.1.1.1"
    }
)
```

## Score

```python
response = fraud.score(
    {
        "email": "user@example.com"
    }
)
```

## Check Email

```python
response = fraud.check_email(
    "user@example.com"
)
```

## Check IP

```python
response = fraud.check_ip(
    "1.1.1.1"
)
```

## Check Device

```python
response = fraud.check_device(
    "device-id"
)
```

## Blacklist

```python
response = fraud.blacklist(
    "email",
    "blocked@example.com"
)
```

## Whitelist

```python
response = fraud.whitelist(
    "email",
    "trusted@example.com"
)
```

## Alerts

```python
response = fraud.alerts()
```

## Statistics

```python
response = fraud.statistics()
```

## Report Fraud

```python
response = fraud.report(
    {
        "type": "abuse",
        "source": "example"
    }
)
```

## Health

```python
response = fraud.health()
```

---

# Permission Service

```python
from hostmk_security import PermissionService

permissions = PermissionService(client)
```

## List Permissions

```python
response = permissions.list_permissions()
```

## Create Permission

```python
response = permissions.create_permission(
    "users.create"
)
```

## Get Permission

```python
response = permissions.get_permission(
    "permission-id"
)
```

## Update Permission

```python
response = permissions.update_permission(
    "permission-id",
    {
        "name": "users.update"
    }
)
```

## Delete Permission

```python
response = permissions.delete_permission(
    "permission-id"
)
```

## Assign Permission

```python
response = permissions.assign_permission(
    "user-id",
    "users.create"
)
```

## Revoke Permission

```python
response = permissions.revoke_permission(
    "user-id",
    "users.create"
)
```

## Bulk Assign

```python
response = permissions.bulk_assign(
    ["u1", "u2", "u3"],
    "users.create"
)
```

## Bulk Revoke

```python
response = permissions.bulk_revoke(
    ["u1", "u2"],
    "users.create"
)
```

## User Permissions

```python
response = permissions.user_permissions(
    "user-id"
)
```

## Role Permissions

```python
response = permissions.role_permissions(
    "role-id"
)
```

## Health

```python
response = permissions.health()
```

---

# Tenancy Service

```python
from hostmk_security import TenancyService

tenancy = TenancyService(client)
```

## List Tenants

```python
response = tenancy.list_tenants()
```

## Create Tenant

```python
response = tenancy.create_tenant(
    "Acme Inc",
    "acme"
)
```

## Get Tenant

```python
response = tenancy.get_tenant(
    "tenant-id"
)
```

## Update Tenant

```python
response = tenancy.update_tenant(
    "tenant-id",
    {
        "name": "Acme Corporation"
    }
)
```

## Delete Tenant

```python
response = tenancy.delete_tenant(
    "tenant-id"
)
```

## Suspend Tenant

```python
response = tenancy.suspend_tenant(
    "tenant-id"
)
```

## Activate Tenant

```python
response = tenancy.activate_tenant(
    "tenant-id"
)
```

## Tenant Usage

```python
response = tenancy.tenant_usage(
    "tenant-id"
)
```

## Tenant Statistics

```python
response = tenancy.tenant_statistics(
    "tenant-id"
)
```

## Tenant Domains

```python
response = tenancy.tenant_domains(
    "tenant-id"
)
```

## Tenant Settings

```python
response = tenancy.tenant_settings(
    "tenant-id"
)
```

## Update Settings

```python
response = tenancy.update_settings(
    "tenant-id",
    {
        "theme": "dark"
    }
)
```

## Health

```python
response = tenancy.health()
```

---

# Error Handling

```python
try:
    auth.login(
        "admin@example.com",
        "password"
    )
except RuntimeError as exc:
    print(exc)
```

---

# Running Tests

Run all tests:

```bash
py -3.12 -m unittest discover -s tests -p "test_*.py" -v
```

Run individual suites:

```bash
py -3.12 -m unittest tests.test_auth -v
```

```bash
py -3.12 -m unittest tests.test_audit -v
```

```bash
py -3.12 -m unittest tests.test_fraud -v
```

```bash
py -3.12 -m unittest tests.test_permissions -v
```

```bash
py -3.12 -m unittest tests.test_tenancy -v
```

```bash
py -3.12 -m unittest tests.test_client -v
```

---

# Validate Source Files

```bash
py -3.12 -m py_compile hostmk_security\client.py
```

```bash
py -3.12 -m py_compile hostmk_security\auth.py
```

```bash
py -3.12 -m py_compile hostmk_security\audit.py
```

```bash
py -3.12 -m py_compile hostmk_security\fraud.py
```

```bash
py -3.12 -m py_compile hostmk_security\permissions.py
```

```bash
py -3.12 -m py_compile hostmk_security\tenancy.py
```

```bash
py -3.12 -m py_compile hostmk_security\__init__.py
```

---

# Current Validation Status

Latest validation:

```text
108 Tests Passed
0 Failures
0 Errors
0 Skipped
```

Modules validated:

```text
client.py
auth.py
audit.py
fraud.py
permissions.py
tenancy.py
__init__.py
```

Imports validated:

```python
from hostmk_security import (
    HostMKClient,
    AuthService,
    AuditService,
    FraudService,
    PermissionService,
    TenancyService,
)
```

Version check:

```python
import hostmk_security

print(hostmk_security.__version__)
```

Output:

```text
1.0.0
```

---

# License

MIT License

Copyright (c) HostMK Security Framework