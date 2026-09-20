from __future__ import annotations

import unittest
from unittest.mock import Mock, patch

import requests

from hostmk_security.client import HostMKClient


class HostMKClientTest(unittest.TestCase):

    def test_can_create_client(self) -> None:
        client = HostMKClient(
            api_key="test-key",
            base_url="https://api.hostmk.io",
        )

        self.assertEqual(
            "test-key",
            client.api_key,
        )

    def test_api_key_required(self) -> None:
        with self.assertRaises(RuntimeError) as context:
            HostMKClient(
                api_key="",
                base_url="https://api.hostmk.io",
            )

        self.assertEqual(
            "HostMK API key is required.",
            str(context.exception),
        )

    def test_base_url_required(self) -> None:
        with self.assertRaises(RuntimeError) as context:
            HostMKClient(
                api_key="key",
                base_url="",
            )

        self.assertEqual(
            "HostMK base URL is required.",
            str(context.exception),
        )

    def test_invalid_base_url(self) -> None:
        with self.assertRaises(RuntimeError) as context:
            HostMKClient(
                api_key="key",
                base_url="invalid-url",
            )

        self.assertEqual(
            "Invalid HostMK base URL.",
            str(context.exception),
        )

    def test_timeout_must_be_positive(self) -> None:
        with self.assertRaises(RuntimeError) as context:
            HostMKClient(
                api_key="key",
                base_url="https://api.hostmk.io",
                timeout=0,
            )

        self.assertEqual(
            "Timeout must be greater than zero.",
            str(context.exception),
        )

    def test_user_agent_required(self) -> None:
        with self.assertRaises(RuntimeError) as context:
            HostMKClient(
                api_key="key",
                base_url="https://api.hostmk.io",
                user_agent="   ",
            )

        self.assertEqual(
            "User-Agent is required.",
            str(context.exception),
        )

    def test_returns_properties(self) -> None:
        client = HostMKClient(
            api_key="secret",
            base_url="https://api.hostmk.io/",
            timeout=60,
            user_agent="HostMK-Test/1.0",
        )

        self.assertEqual(
            "secret",
            client.api_key,
        )

        self.assertEqual(
            "https://api.hostmk.io",
            client.base_url,
        )

        self.assertEqual(
            60,
            client.timeout,
        )

        self.assertEqual(
            "HostMK-Test/1.0",
            client.user_agent,
        )

    def test_build_url(self) -> None:
        client = HostMKClient(
            api_key="key",
            base_url="https://api.hostmk.io",
        )

        result = client._build_url(
            "/v1/auth/login"
        )

        self.assertEqual(
            "https://api.hostmk.io/v1/auth/login",
            result,
        )

    def test_endpoint_required(self) -> None:
        client = HostMKClient(
            api_key="key",
            base_url="https://api.hostmk.io",
        )

        with self.assertRaises(RuntimeError) as context:
            client._build_url("")

        self.assertEqual(
            "Endpoint is required.",
            str(context.exception),
        )

    def test_build_headers(self) -> None:
        client = HostMKClient(
            api_key="api-key",
            base_url="https://api.hostmk.io",
        )

        headers = client._build_headers(
            {
                "X-Custom": "custom",
            }
        )

        self.assertEqual(
            "api-key",
            headers["X-API-Key"],
        )

        self.assertEqual(
            "custom",
            headers["X-Custom"],
        )

        self.assertEqual(
            "application/json",
            headers["Accept"],
        )

    @patch("hostmk_security.client.requests.request")
    def test_get_request(
        self,
        request_mock: Mock,
    ) -> None:
        response = Mock()

        response.status_code = 200

        response.json.return_value = {
            "success": True,
        }

        request_mock.return_value = response

        client = HostMKClient(
            api_key="key",
            base_url="https://api.hostmk.io",
        )

        result = client.get(
            "/health"
        )

        self.assertEqual(
            200,
            result["status"],
        )

        self.assertTrue(
            result["data"]["success"]
        )

    @patch("hostmk_security.client.requests.request")
    def test_post_request(
        self,
        request_mock: Mock,
    ) -> None:
        response = Mock()

        response.status_code = 201

        response.json.return_value = {
            "created": True,
        }

        request_mock.return_value = response

        client = HostMKClient(
            api_key="key",
            base_url="https://api.hostmk.io",
        )

        result = client.post(
            "/items",
            {
                "name": "test",
            },
        )

        self.assertEqual(
            201,
            result["status"],
        )

        self.assertTrue(
            result["data"]["created"]
        )

    @patch("hostmk_security.client.requests.request")
    def test_put_request(
        self,
        request_mock: Mock,
    ) -> None:
        response = Mock()

        response.status_code = 200

        response.json.return_value = {
            "updated": True,
        }

        request_mock.return_value = response

        client = HostMKClient(
            api_key="key",
            base_url="https://api.hostmk.io",
        )

        result = client.put(
            "/items/1",
            {
                "name": "updated",
            },
        )

        self.assertEqual(
            200,
            result["status"],
        )

    @patch("hostmk_security.client.requests.request")
    def test_delete_request(
        self,
        request_mock: Mock,
    ) -> None:
        response = Mock()

        response.status_code = 200

        response.json.return_value = {
            "deleted": True,
        }

        request_mock.return_value = response

        client = HostMKClient(
            api_key="key",
            base_url="https://api.hostmk.io",
        )

        result = client.delete(
            "/items/1"
        )

        self.assertEqual(
            200,
            result["status"],
        )

    @patch("hostmk_security.client.requests.request")
    def test_http_error_response(
        self,
        request_mock: Mock,
    ) -> None:
        response = Mock()

        response.status_code = 401

        response.json.return_value = {
            "message": "Unauthorized",
        }

        request_mock.return_value = response

        client = HostMKClient(
            api_key="key",
            base_url="https://api.hostmk.io",
        )

        with self.assertRaises(RuntimeError) as context:
            client.get("/private")

        self.assertEqual(
            "Unauthorized",
            str(context.exception),
        )

    @patch("hostmk_security.client.requests.request")
    def test_invalid_json_response(
        self,
        request_mock: Mock,
    ) -> None:
        response = Mock()

        response.status_code = 200

        response.json.side_effect = ValueError()

        request_mock.return_value = response

        client = HostMKClient(
            api_key="key",
            base_url="https://api.hostmk.io",
        )

        with self.assertRaises(RuntimeError) as context:
            client.get("/health")

        self.assertEqual(
            "Invalid JSON response received.",
            str(context.exception),
        )

    @patch("hostmk_security.client.requests.request")
    def test_request_exception(
        self,
        request_mock: Mock,
    ) -> None:
        request_mock.side_effect = (
            requests.RequestException(
                "Network error"
            )
        )

        client = HostMKClient(
            api_key="key",
            base_url="https://api.hostmk.io",
        )

        with self.assertRaises(RuntimeError):
            client.get("/health")


if __name__ == "__main__":
    unittest.main()