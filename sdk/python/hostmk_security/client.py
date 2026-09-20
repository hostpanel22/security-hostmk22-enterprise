from __future__ import annotations

from typing import Any, Dict, Optional

import requests


class HostMKClient:
    """
    HostMK Security SDK HTTP Client.

    Features:
    - API key authentication
    - Base URL validation
    - Configurable timeout
    - Configurable User-Agent
    - GET / POST / PUT / DELETE helpers
    - Automatic JSON handling
    - Consistent error management
    """

    DEFAULT_TIMEOUT = 30
    DEFAULT_USER_AGENT = "HostMK-Python-SDK/1.0.0"

    def __init__(
        self,
        api_key: str,
        base_url: str,
        timeout: int = DEFAULT_TIMEOUT,
        user_agent: str = DEFAULT_USER_AGENT,
    ) -> None:
        api_key = api_key.strip()
        base_url = base_url.strip().rstrip("/")
        user_agent = user_agent.strip()

        if not api_key:
            raise RuntimeError(
                "HostMK API key is required."
            )

        if not base_url:
            raise RuntimeError(
                "HostMK base URL is required."
            )

        if not (
            base_url.startswith("http://")
            or base_url.startswith("https://")
        ):
            raise RuntimeError(
                "Invalid HostMK base URL."
            )

        if timeout <= 0:
            raise RuntimeError(
                "Timeout must be greater than zero."
            )

        if not user_agent:
            raise RuntimeError(
                "User-Agent is required."
            )

        self._api_key = api_key
        self._base_url = base_url
        self._timeout = timeout
        self._user_agent = user_agent

    @property
    def api_key(self) -> str:
        return self._api_key

    @property
    def base_url(self) -> str:
        return self._base_url

    @property
    def timeout(self) -> int:
        return self._timeout

    @property
    def user_agent(self) -> str:
        return self._user_agent

    def get(
        self,
        endpoint: str,
        headers: Optional[Dict[str, str]] = None,
    ) -> Dict[str, Any]:
        return self._request(
            method="GET",
            endpoint=endpoint,
            headers=headers,
        )

    def post(
        self,
        endpoint: str,
        payload: Optional[Dict[str, Any]] = None,
        headers: Optional[Dict[str, str]] = None,
    ) -> Dict[str, Any]:
        return self._request(
            method="POST",
            endpoint=endpoint,
            payload=payload,
            headers=headers,
        )

    def put(
        self,
        endpoint: str,
        payload: Optional[Dict[str, Any]] = None,
        headers: Optional[Dict[str, str]] = None,
    ) -> Dict[str, Any]:
        return self._request(
            method="PUT",
            endpoint=endpoint,
            payload=payload,
            headers=headers,
        )

    def delete(
        self,
        endpoint: str,
        headers: Optional[Dict[str, str]] = None,
    ) -> Dict[str, Any]:
        return self._request(
            method="DELETE",
            endpoint=endpoint,
            headers=headers,
        )

    def _request(
        self,
        method: str,
        endpoint: str,
        payload: Optional[Dict[str, Any]] = None,
        headers: Optional[Dict[str, str]] = None,
    ) -> Dict[str, Any]:
        url = self._build_url(endpoint)

        request_headers = self._build_headers(
            headers or {}
        )

        try:
            response = requests.request(
                method=method,
                url=url,
                json=payload,
                headers=request_headers,
                timeout=self._timeout,
            )
        except requests.RequestException as exc:
            raise RuntimeError(
                f"HTTP request failed: {exc}"
            ) from exc

        try:
            data = response.json()
        except ValueError as exc:
            raise RuntimeError(
                "Invalid JSON response received."
            ) from exc

        if response.status_code >= 400:
            message = (
                f"HostMK API returned HTTP "
                f"{response.status_code}"
            )

            if isinstance(data, dict):
                api_message = data.get(
                    "message"
                )

                if isinstance(
                    api_message,
                    str,
                ) and api_message:
                    message = api_message

            raise RuntimeError(message)

        return {
            "status": response.status_code,
            "data": data,
        }

    def _build_url(
        self,
        endpoint: str,
    ) -> str:
        endpoint = endpoint.strip()

        if not endpoint:
            raise RuntimeError(
                "Endpoint is required."
            )

        return (
            self._base_url
            + "/"
            + endpoint.lstrip("/")
        )

    def _build_headers(
        self,
        headers: Dict[str, str],
    ) -> Dict[str, str]:
        default_headers = {
            "Accept": "application/json",
            "Content-Type": "application/json",
            "X-API-Key": self._api_key,
            "User-Agent": self._user_agent,
        }

        default_headers.update(headers)

        return default_headers