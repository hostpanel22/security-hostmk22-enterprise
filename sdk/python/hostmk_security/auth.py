from __future__ import annotations

from typing import Any, Dict

from .client import HostMKClient


class AuthService:
    """
    HostMK Authentication Service.

    Provides authentication and session
    management operations against the
    HostMK Security API.
    """

    def __init__(
        self,
        client: HostMKClient,
    ) -> None:
        self._client = client

    def login(
        self,
        email: str,
        password: str,
    ) -> Dict[str, Any]:
        email = email.strip()
        password = password.strip()

        if not email:
            raise ValueError(
                "Email is required."
            )

        if not password:
            raise ValueError(
                "Password is required."
            )

        return self._client.post(
            "/auth/login",
            {
                "email": email,
                "password": password,
            },
        )

    def validate_token(
        self,
        token: str,
    ) -> Dict[str, Any]:
        token = token.strip()

        if not token:
            raise ValueError(
                "Token is required."
            )

        return self._client.post(
            "/auth/validate",
            {
                "token": token,
            },
        )

    def refresh_token(
        self,
        refresh_token: str,
    ) -> Dict[str, Any]:
        refresh_token = refresh_token.strip()

        if not refresh_token:
            raise ValueError(
                "Refresh token is required."
            )

        return self._client.post(
            "/auth/refresh",
            {
                "refreshToken": refresh_token,
            },
        )

    def logout(
        self,
        token: str,
    ) -> Dict[str, Any]:
        token = token.strip()

        if not token:
            raise ValueError(
                "Token is required."
            )

        return self._client.post(
            "/auth/logout",
            {
                "token": token,
            },
        )

    def sessions(
        self,
        user_id: str,
    ) -> Dict[str, Any]:
        user_id = user_id.strip()

        if not user_id:
            raise ValueError(
                "User ID is required."
            )

        return self._client.get(
            f"/auth/users/{user_id}/sessions"
        )

    def revoke_session(
        self,
        session_id: str,
    ) -> Dict[str, Any]:
        session_id = session_id.strip()

        if not session_id:
            raise ValueError(
                "Session ID is required."
            )

        return self._client.delete(
            f"/auth/sessions/{session_id}"
        )

    def health(
        self,
    ) -> Dict[str, Any]:
        return self._client.get(
            "/auth/health"
        )