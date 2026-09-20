from __future__ import annotations

from typing import Any, Dict, Optional

from .client import HostMKClient


class AuditService:
    """
    HostMK Audit Service.

    Provides access to audit logs, security events,
    compliance records and monitoring endpoints.
    """

    def __init__(
        self,
        client: HostMKClient,
    ) -> None:
        self._client = client

    def list_events(
        self,
        page: int = 1,
        limit: int = 25,
    ) -> Dict[str, Any]:
        """
        Retrieve audit events.
        """

        return self._client.get(
            f"audit/events?page={page}&limit={limit}"
        )

    def get_event(
        self,
        event_id: str,
    ) -> Dict[str, Any]:
        """
        Retrieve a single audit event.
        """

        event_id = event_id.strip()

        if not event_id:
            raise RuntimeError(
                "Event ID is required."
            )

        return self._client.get(
            f"audit/events/{event_id}"
        )

    def search_events(
        self,
        filters: Dict[str, Any],
    ) -> Dict[str, Any]:
        """
        Search audit events.
        """

        return self._client.post(
            "audit/events/search",
            filters,
        )

    def user_activity(
        self,
        user_id: str,
        page: int = 1,
        limit: int = 25,
    ) -> Dict[str, Any]:
        """
        Retrieve activity for a specific user.
        """

        user_id = user_id.strip()

        if not user_id:
            raise RuntimeError(
                "User ID is required."
            )

        return self._client.get(
            f"audit/users/{user_id}/activity"
            f"?page={page}&limit={limit}"
        )

    def tenant_activity(
        self,
        tenant_id: str,
        page: int = 1,
        limit: int = 25,
    ) -> Dict[str, Any]:
        """
        Retrieve tenant audit activity.
        """

        tenant_id = tenant_id.strip()

        if not tenant_id:
            raise RuntimeError(
                "Tenant ID is required."
            )

        return self._client.get(
            f"audit/tenants/{tenant_id}/activity"
            f"?page={page}&limit={limit}"
        )

    def export_events(
        self,
        filters: Optional[Dict[str, Any]] = None,
    ) -> Dict[str, Any]:
        """
        Export audit logs.
        """

        return self._client.post(
            "audit/export",
            filters or {},
        )

    def compliance_report(
        self,
        report_type: str,
    ) -> Dict[str, Any]:
        """
        Generate compliance report.
        """

        report_type = report_type.strip()

        if not report_type:
            raise RuntimeError(
                "Report type is required."
            )

        return self._client.get(
            f"audit/compliance/{report_type}"
        )

    def security_events(
        self,
        severity: Optional[str] = None,
    ) -> Dict[str, Any]:
        """
        Retrieve security events.
        """

        endpoint = "audit/security-events"

        if severity:
            endpoint += (
                f"?severity={severity}"
            )

        return self._client.get(
            endpoint
        )

    def statistics(self) -> Dict[str, Any]:
        """
        Retrieve audit statistics.
        """

        return self._client.get(
            "audit/statistics"
        )

    def health(self) -> Dict[str, Any]:
        """
        Audit service health endpoint.
        """

        return self._client.get(
            "audit/health"
        )