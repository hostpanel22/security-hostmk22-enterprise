import unittest

from hostmk_security.auth import AuthService


class FakeClient:
    def __init__(self):
        self.method = None
        self.endpoint = None
        self.payload = None

    def get(self, endpoint, headers=None):
        self.method = "GET"
        self.endpoint = endpoint

        return {
            "status": 200,
            "data": {
                "success": True
            }
        }

    def post(
        self,
        endpoint,
        payload=None,
        headers=None,
    ):
        self.method = "POST"
        self.endpoint = endpoint
        self.payload = payload

        return {
            "status": 200,
            "data": {
                "success": True
            }
        }

    def delete(
        self,
        endpoint,
        headers=None,
    ):
        self.method = "DELETE"
        self.endpoint = endpoint

        return {
            "status": 200,
            "data": {
                "success": True
            }
        }


class AuthServiceTest(
    unittest.TestCase
):
    def setUp(self):
        self.client = FakeClient()
        self.service = AuthService(
            self.client
        )

    def test_login(self):
        response = self.service.login(
            "admin@test.com",
            "password",
        )

        self.assertEqual(
            "POST",
            self.client.method,
        )

        self.assertEqual(
            "/auth/login",
            self.client.endpoint,
        )

        self.assertTrue(
            response["data"]["success"]
        )

    def test_login_email_required(self):
        with self.assertRaises(
            ValueError
        ):
            self.service.login(
                "",
                "password",
            )

    def test_login_password_required(self):
        with self.assertRaises(
            ValueError
        ):
            self.service.login(
                "admin@test.com",
                "",
            )

    def test_validate_token(self):
        self.service.validate_token(
            "token123"
        )

        self.assertEqual(
            "/auth/validate",
            self.client.endpoint,
        )

    def test_refresh_token(self):
        self.service.refresh_token(
            "refresh123"
        )

        self.assertEqual(
            "/auth/refresh",
            self.client.endpoint,
        )

    def test_logout(self):
        self.service.logout(
            "token123"
        )

        self.assertEqual(
            "/auth/logout",
            self.client.endpoint,
        )

    def test_sessions(self):
        self.service.sessions(
            "user-1"
        )

        self.assertEqual(
            "/auth/users/user-1/sessions",
            self.client.endpoint,
        )

    def test_revoke_session(self):
        self.service.revoke_session(
            "session-1"
        )

        self.assertEqual(
            "DELETE",
            self.client.method,
        )

        self.assertEqual(
            "/auth/sessions/session-1",
            self.client.endpoint,
        )

    def test_health(self):
        self.service.health()

        self.assertEqual(
            "/auth/health",
            self.client.endpoint,
        )


if __name__ == "__main__":
    unittest.main()