from __future__ import annotations

from typing import Any, Dict

from .client import HostMKClient


class PermissionService:
    def __init__(
        self,
        client: HostMKClient,
    ) -> None:
        self._client = client

    def list_permissions(
        self,
    ) -> Dict[str, Any]:
        return self._client.get(
            "/permissions"
        )

    def list_roles(
        self,
    ) -> Dict[str, Any]:
        return self._client.get(
            "/permissions/roles"
        )

    def assign(
        self,
        user_id: str,
        permission: str,
    ) -> Dict[str, Any]:
        return self._client.post(
            "/permissions/assign",
            {
                "user_id": user_id,
                "permission": permission,
            },
        )

    def revoke(
        self,
        user_id: str,
        permission: str,
    ) -> Dict[str, Any]:
        return self._client.post(
            "/permissions/revoke",
            {
                "user_id": user_id,
                "permission": permission,
            },
        )

    def check(
        self,
        user_id: str,
        permission: str,
    ) -> Dict[str, Any]:
        return self._client.post(
            "/permissions/check",
            {
                "user_id": user_id,
                "permission": permission,
            },
        )