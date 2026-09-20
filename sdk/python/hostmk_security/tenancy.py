from __future__ import annotations

from typing import Any, Dict

from .client import HostMKClient


class TenancyService:
    """
    HostMK Multi-Tenant Service.

    Features:
    - List tenants
    - Get tenant
    - Create tenant
    - Update tenant
    - Delete tenant
    - Suspend tenant
    - Activate tenant
    - Tenant usage
    - Tenant statistics
    - Tenant domains
    - Tenant settings
    - Health endpoint
    """

    def __init__(
        self,
        client: HostMKClient,
    ) -> None:
        self._client = client

    def list_tenants(
        self,
    ) -> Dict[str, Any]:
        return self._client.get(
            "/tenancy/tenants"
        )

    def get_tenant(
        self,
        tenant_id: str,
    ) -> Dict[str, Any]:
        tenant_id = tenant_id.strip()

        if not tenant_id:
            raise RuntimeError(
                "Tenant ID is required."
            )

        return self._client.get(
            f"/tenancy/tenants/{tenant_id}"
        )

    def create_tenant(
        self,
        name: str,
        slug: str,
    ) -> Dict[str, Any]:
        name = name.strip()
        slug = slug.strip()

        if not name:
            raise RuntimeError(
                "Tenant name is required."
            )

        if not slug:
            raise RuntimeError(
                "Tenant slug is required."
            )

        return self._client.post(
            "/tenancy/tenants",
            {
                "name": name,
                "slug": slug,
            },
        )

    def update_tenant(
        self,
        tenant_id: str,
        payload: Dict[str, Any],
    ) -> Dict[str, Any]:
        tenant_id = tenant_id.strip()

        if not tenant_id:
            raise RuntimeError(
                "Tenant ID is required."
            )

        return self._client.put(
            f"/tenancy/tenants/{tenant_id}",
            payload,
        )

    def delete_tenant(
        self,
        tenant_id: str,
    ) -> Dict[str, Any]:
        tenant_id = tenant_id.strip()

        if not tenant_id:
            raise RuntimeError(
                "Tenant ID is required."
            )

        return self._client.delete(
            f"/tenancy/tenants/{tenant_id}"
        )

    def suspend_tenant(
        self,
        tenant_id: str,
    ) -> Dict[str, Any]:
        tenant_id = tenant_id.strip()

        if not tenant_id:
            raise RuntimeError(
                "Tenant ID is required."
            )

        return self._client.post(
            f"/tenancy/tenants/{tenant_id}/suspend",
            {},
        )

    def activate_tenant(
        self,
        tenant_id: str,
    ) -> Dict[str, Any]:
        tenant_id = tenant_id.strip()

        if not tenant_id:
            raise RuntimeError(
                "Tenant ID is required."
            )

        return self._client.post(
            f"/tenancy/tenants/{tenant_id}/activate",
            {},
        )

    def tenant_usage(
        self,
        tenant_id: str,
    ) -> Dict[str, Any]:
        tenant_id = tenant_id.strip()

        if not tenant_id:
            raise RuntimeError(
                "Tenant ID is required."
            )

        return self._client.get(
            f"/tenancy/tenants/{tenant_id}/usage"
        )

    def tenant_statistics(
        self,
        tenant_id: str,
    ) -> Dict[str, Any]:
        tenant_id = tenant_id.strip()

        if not tenant_id:
            raise RuntimeError(
                "Tenant ID is required."
            )

        return self._client.get(
            f"/tenancy/tenants/{tenant_id}/statistics"
        )

    def tenant_domains(
        self,
        tenant_id: str,
    ) -> Dict[str, Any]:
        tenant_id = tenant_id.strip()

        if not tenant_id:
            raise RuntimeError(
                "Tenant ID is required."
            )

        return self._client.get(
            f"/tenancy/tenants/{tenant_id}/domains"
        )

    def tenant_settings(
        self,
        tenant_id: str,
    ) -> Dict[str, Any]:
        tenant_id = tenant_id.strip()

        if not tenant_id:
            raise RuntimeError(
                "Tenant ID is required."
            )

        return self._client.get(
            f"/tenancy/tenants/{tenant_id}/settings"
        )

    def update_settings(
        self,
        tenant_id: str,
        settings: Dict[str, Any],
    ) -> Dict[str, Any]:
        tenant_id = tenant_id.strip()

        if not tenant_id:
            raise RuntimeError(
                "Tenant ID is required."
            )

        return self._client.put(
            f"/tenancy/tenants/{tenant_id}/settings",
            settings,
        )

    def health(
        self,
    ) -> Dict[str, Any]:
        return self._client.get(
            "/tenancy/health"
        )