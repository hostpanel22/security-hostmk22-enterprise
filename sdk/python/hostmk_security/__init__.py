"""
HostMK Security Framework SDK
Official Python SDK.

Modules:
- Authentication
- Audit Logging
- Fraud Detection
- Permissions
- Multi-Tenancy
"""

from .client import HostMKClient
from .auth import AuthService
from .audit import AuditService
from .fraud import FraudService
from .permissions import PermissionService
from .tenancy import TenancyService

__title__ = "hostmk-security"
__description__ = "Official Python SDK for HostMK Security Framework"
__version__ = "1.0.0"
__author__ = "HostMK Team"
__license__ = "MIT"

__all__ = [
    "HostMKClient",
    "AuthService",
    "AuditService",
    "FraudService",
    "PermissionService",
    "TenancyService",
]