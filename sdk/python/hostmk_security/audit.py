from __future__ import annotations

from typing import Any, Dict, Optional

from .client import HostMKClient


class AuditService:
    """
    HostMK Audit Service.
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
        return self._client.get(
            f"/audit/events?page={page}&limit={limit}"
        )

    def get_event(
        self,
        event_id: str,
    ) -> Dict[str, Any]:
        if not event_id:
            raise RuntimeError(
                "Event ID is required."
            )

        return self._client.get(
            f"/audit/events/{event_id}"
        )

    def search(
        self,
        filters: Optional[Dict[str, Any]] = None,
    ) -> Dict[str, Any]:
        return self._client.post(
            "/audit/search",
            filters or {},
        )

    def export(
        self,
        filters: Optional[Dict[str, Any]] = None,
    ) -> Dict[str, Any]:
        return self._client.post(
            "/audit/export",
            filters or {},
        )

    def statistics(self) -> Dict[str, Any]:
        return self._client.get(
            "/audit/statistics"
        )