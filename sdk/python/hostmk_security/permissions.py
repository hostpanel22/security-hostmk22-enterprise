from __future__ import annotations

from typing import Any, Dict, List

from .client import HostMKClient


class PermissionService:
    """
    HostMK Permission Service.

    Features:
    - List permissions
    - Get permission details
    - Create permissions
    - Update permissions
    - Delete permissions
    - Assign permissions
    - Revoke permissions
    - Role permissions
    - User permissions
    - Health endpoint
    """

    def __init__(self, client: HostMKClient) -> None:
        self._client = client

    def list_permissions(self) -> Dict[str, Any]:
        return self._client.get(
            "/permissions"
        )

    def get_permission(
        self,
        permission_id: str,
    ) -> Dict[str, Any]:
        permission_id = permission_id.strip()

        if not permission_id:
            raise RuntimeError(
                "Permission ID is required."
            )

        return self._client.get(
            f"/permissions/{permission_id}"
        )

    def create_permission(
        self,
        name: str,
        description: str = "",
    ) -> Dict[str, Any]:
        name = name.strip()

        if not name:
            raise RuntimeError(
                "Permission name is required."
            )

        return self._client.post(
            "/permissions",
            {
                "name": name,
                "description": description,
            },
        )

    def update_permission(
        self,
        permission_id: str,
        payload: Dict[str, Any],
    ) -> Dict[str, Any]:
        permission_id = permission_id.strip()

        if not permission_id:
            raise RuntimeError(
                "Permission ID is required."
            )

        return self._client.put(
            f"/permissions/{permission_id}",
            payload,
        )

    def delete_permission(
        self,
        permission_id: str,
    ) -> Dict[str, Any]:
        permission_id = permission_id.strip()

        if not permission_id:
            raise RuntimeError(
                "Permission ID is required."
            )

        return self._client.delete(
            f"/permissions/{permission_id}"
        )

    def assign_permission(
        self,
        user_id: str,
        permission: str,
    ) -> Dict[str, Any]:
        user_id = user_id.strip()
        permission = permission.strip()

        if not user_id:
            raise RuntimeError(
                "User ID is required."
            )

        if not permission:
            raise RuntimeError(
                "Permission is required."
            )

        return self._client.post(
            f"/permissions/users/{user_id}/assign",
            {
                "permission": permission,
            },
        )

    def revoke_permission(
        self,
        user_id: str,
        permission: str,
    ) -> Dict[str, Any]:
        user_id = user_id.strip()
        permission = permission.strip()

        if not user_id:
            raise RuntimeError(
                "User ID is required."
            )

        if not permission:
            raise RuntimeError(
                "Permission is required."
            )

        return self._client.post(
            f"/permissions/users/{user_id}/revoke",
            {
                "permission": permission,
            },
        )

    def role_permissions(
        self,
        role_id: str,
    ) -> Dict[str, Any]:
        role_id = role_id.strip()

        if not role_id:
            raise RuntimeError(
                "Role ID is required."
            )

        return self._client.get(
            f"/permissions/roles/{role_id}"
        )

    def user_permissions(
        self,
        user_id: str,
    ) -> Dict[str, Any]:
        user_id = user_id.strip()

        if not user_id:
            raise RuntimeError(
                "User ID is required."
            )

        return self._client.get(
            f"/permissions/users/{user_id}"
        )

    def bulk_assign(
        self,
        user_ids: List[str],
        permission: str,
    ) -> Dict[str, Any]:
        if not user_ids:
            raise RuntimeError(
                "User IDs are required."
            )

        permission = permission.strip()

        if not permission:
            raise RuntimeError(
                "Permission is required."
            )

        return self._client.post(
            "/permissions/bulk-assign",
            {
                "users": user_ids,
                "permission": permission,
            },
        )

    def bulk_revoke(
        self,
        user_ids: List[str],
        permission: str,
    ) -> Dict[str, Any]:
        if not user_ids:
            raise RuntimeError(
                "User IDs are required."
            )

        permission = permission.strip()

        if not permission:
            raise RuntimeError(
                "Permission is required."
            )

        return self._client.post(
            "/permissions/bulk-revoke",
            {
                "users": user_ids,
                "permission": permission,
            },
        )

    def health(self) -> Dict[str, Any]:
        return self._client.get(
            "/permissions/health"
        )