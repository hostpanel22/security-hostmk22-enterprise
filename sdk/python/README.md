# HostMK Security Framework - Python SDK

Official Python SDK for the HostMK Security Framework.

The SDK provides a simple interface to integrate HostMK Security services into Python applications.

## Features

- Authentication Management
- Audit Logging
- Fraud Detection
- Permission Management
- Multi-Tenant Management
- REST API Client
- API Key Authentication
- Python 3.12+ Support
- Typed Interfaces
- Unit Tested

---

## Installation

### From Source

```bash
git clone https://github.com/azezaz/HostMK-Security-Framework.git

cd HostMK-Security-Framework/sdk/python

pip install -e .
```

### Install Dependencies

```bash
pip install requests
```

---

## Quick Start

```python
from hostmk_security import (
    HostMKClient,
    AuthService,
)

client = HostMKClient(
    api_key="your-api-key",
    base_url="https://api.hostmk.com"
)

auth = AuthService(client)

response = auth.health()

print(response)
```

---

# Client

## Create Client

```python
from hostmk_security import HostMKClient

client = HostMKClient(
    api_key="your-api-key",
    base_url="https://api.hostmk.com"
)
```

## Configuration

```python
client = HostMKClient(
    api_key="your-api-key",
    base_url="https://api.hostmk.com",
    timeout=60,
    user_agent="MyApplication/1.0"
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
auth.login(
    email="admin@example.com",
    password="password"
)
```

## Validate Token

```python
auth.validate_token(
    token="jwt-token"
)
```

## Refresh Token

```python
auth.refresh_token(
    refresh_token="refresh-token"
)
```

## Logout

```python
auth.logout(
    token="jwt-token"
)
```

---

# Audit Service

```python
from hostmk_security import AuditService

audit = AuditService(client)
```

## List Events

```python
audit.list_events()
```

## Search Events

```python
audit.search_events("login")
```

## Event Details

```python
audit.get_event("event-id")
```

## User Activity

```python
audit.user_activity("user-id")
```

## Tenant Activity

```python
audit.tenant_activity("tenant-id")
```

## Security Events

```python
audit.security_events()
```

## Statistics

```python
audit.statistics()
```

## Compliance Report

```python
audit.compliance_report("gdpr")
```

---

# Fraud Service

```python
from hostmk_security import FraudService

fraud = FraudService(client)
```

## Check IP

```python
fraud.check_ip("1.1.1.1")
```

## Check Email

```python
fraud.check_email(
    "user@example.com"
)
```

## Check Device

```python
fraud.check_device(
    "device-id"
)
```

## Fraud Score

```python
fraud.score({
    "email": "user@example.com"
})
```

## Analyze Transaction

```python
fraud.analyze({
    "amount": 100,
    "currency": "USD"
})
```

## Blacklist

```python
fraud.blacklist(
    "ip",
    "1.1.1.1"
)
```

## Whitelist

```python
fraud.whitelist(
    "email",
    "trusted@example.com"
)
```

---

# Permission Service

```python
from hostmk_security import PermissionService

permissions = PermissionService(client)
```

## Create Permission

```python
permissions.create_permission(
    "orders.read"
)
```

## Get Permission

```python
permissions.get_permission(
    "permission-id"
)
```

## Assign Permission

```python
permissions.assign_permission(
    "user-id",
    "orders.read"
)
```

## Revoke Permission

```python
permissions.revoke_permission(
    "user-id",
    "orders.read"
)
```

## User Permissions

```python
permissions.user_permissions(
    "user-id"
)
```

## Role Permissions

```python
permissions.role_permissions(
    "role-id"
)
```

---

# Tenancy Service

```python
from hostmk_security import TenancyService

tenancy = TenancyService(client)
```

## Create Tenant

```python
tenancy.create_tenant(
    name="Acme Corp",
    slug="acme"
)
```

## List Tenants

```python
tenancy.list_tenants()
```

## Get Tenant

```python
tenancy.get_tenant(
    "tenant-id"
)
```

## Update Tenant

```python
tenancy.update_tenant(
    "tenant-id",
    {
        "name": "New Name"
    }
)
```

## Tenant Usage

```python
tenancy.tenant_usage(
    "tenant-id"
)
```

## Tenant Statistics

```python
tenancy.tenant_statistics(
    "tenant-id"
)
```

## Tenant Domains

```python
tenancy.tenant_domains(
    "tenant-id"
)
```

## Tenant Settings

```python
tenancy.tenant_settings(
    "tenant-id"
)
```

## Activate Tenant

```python
tenancy.activate_tenant(
    "tenant-id"
)
```

## Suspend Tenant

```python
tenancy.suspend_tenant(
    "tenant-id"
)
```

## Delete Tenant

```python
tenancy.delete_tenant(
    "tenant-id"
)
```

---

# Running Tests

Run all tests:

```bash
python -m unittest discover -s tests -p "test_*.py" -v
```

Run a specific test file:

```bash
python -m unittest tests.test_auth -v
```

```bash
python -m unittest tests.test_client -v
```

```bash
python -m unittest tests.test_audit -v
```

```bash
python -m unittest tests.test_fraud -v
```

```bash
python -m unittest tests.test_permissions -v
```

```bash
python -m unittest tests.test_tenancy -v
```

---

# Package Verification

```bash
python -c "import hostmk_security"
```

```bash
python -c "from hostmk_security import *"
```

---

# Requirements

- Python 3.12+
- requests >= 2.32

---

# License

MIT License

---

# Repository

https://github.com/azezaz/HostMK-Security-Framework

---

# Version

Current Version:

```text
1.0.0
```