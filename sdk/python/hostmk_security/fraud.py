from __future__ import annotations

from typing import Any, Dict, Optional

from .client import HostMKClient


class FraudService:
    """
    HostMK Fraud Detection Service.

    Provides fraud analysis, risk scoring,
    transaction screening and threat monitoring.
    """

    def __init__(
        self,
        client: HostMKClient,
    ) -> None:
        self._client = client

    def analyze(
        self,
        payload: Dict[str, Any],
    ) -> Dict[str, Any]:
        """
        Analyze a transaction or event.
        """

        return self._client.post(
            "fraud/analyze",
            payload,
        )

    def score(
        self,
        payload: Dict[str, Any],
    ) -> Dict[str, Any]:
        """
        Calculate fraud risk score.
        """

        return self._client.post(
            "fraud/score",
            payload,
        )

    def check_ip(
        self,
        ip_address: str,
    ) -> Dict[str, Any]:
        """
        Check IP reputation.
        """

        ip_address = ip_address.strip()

        if not ip_address:
            raise RuntimeError(
                "IP address is required."
            )

        return self._client.get(
            f"fraud/ip/{ip_address}"
        )

    def check_email(
        self,
        email: str,
    ) -> Dict[str, Any]:
        """
        Check email reputation.
        """

        email = email.strip()

        if not email:
            raise RuntimeError(
                "Email is required."
            )

        return self._client.get(
            f"fraud/email/{email}"
        )

    def check_device(
        self,
        device_id: str,
    ) -> Dict[str, Any]:
        """
        Check device reputation.
        """

        device_id = device_id.strip()

        if not device_id:
            raise RuntimeError(
                "Device ID is required."
            )

        return self._client.get(
            f"fraud/device/{device_id}"
        )

    def report(
        self,
        payload: Dict[str, Any],
    ) -> Dict[str, Any]:
        """
        Report fraudulent activity.
        """

        return self._client.post(
            "fraud/report",
            payload,
        )

    def blacklist(
        self,
        entity_type: str,
        value: str,
    ) -> Dict[str, Any]:
        """
        Add an entity to blacklist.
        """

        entity_type = entity_type.strip()
        value = value.strip()

        if not entity_type:
            raise RuntimeError(
                "Entity type is required."
            )

        if not value:
            raise RuntimeError(
                "Entity value is required."
            )

        return self._client.post(
            "fraud/blacklist",
            {
                "type": entity_type,
                "value": value,
            },
        )

    def whitelist(
        self,
        entity_type: str,
        value: str,
    ) -> Dict[str, Any]:
        """
        Add an entity to whitelist.
        """

        entity_type = entity_type.strip()
        value = value.strip()

        if not entity_type:
            raise RuntimeError(
                "Entity type is required."
            )

        if not value:
            raise RuntimeError(
                "Entity value is required."
            )

        return self._client.post(
            "fraud/whitelist",
            {
                "type": entity_type,
                "value": value,
            },
        )

    def alerts(
        self,
        page: int = 1,
        limit: int = 25,
    ) -> Dict[str, Any]:
        """
        Retrieve fraud alerts.
        """

        return self._client.get(
            f"fraud/alerts?page={page}&limit={limit}"
        )

    def statistics(self) -> Dict[str, Any]:
        """
        Retrieve fraud statistics.
        """

        return self._client.get(
            "fraud/statistics"
        )

    def health(self) -> Dict[str, Any]:
        """
        Fraud service health endpoint.
        """

        return self._client.get(
            "fraud/health"
        )