from __future__ import annotations

from typing import Any, Dict

from .client import HostMKClient


class FraudService:
    def __init__(
        self,
        client: HostMKClient,
    ) -> None:
        self._client = client

    def evaluate(
        self,
        payload: Dict[str, Any],
    ) -> Dict[str, Any]:
        return self._client.post(
            "/fraud/evaluate",
            payload,
        )

    def score(
        self,
        payload: Dict[str, Any],
    ) -> Dict[str, Any]:
        return self._client.post(
            "/fraud/score",
            payload,
        )

    def blacklist(
        self,
        value: str,
    ) -> Dict[str, Any]:
        return self._client.post(
            "/fraud/blacklist",
            {
                "value": value,
            },
        )

    def whitelist(
        self,
        value: str,
    ) -> Dict[str, Any]:
        return self._client.post(
            "/fraud/whitelist",
            {
                "value": value,
            },
        )

    def health(self) -> Dict[str, Any]:
        return self._client.get(
            "/fraud/health"
        )