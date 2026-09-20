import unittest

from hostmk_security.permissions import PermissionService


class FakeClient:
    def get(self, endpoint):
        return {
            "method": "GET",
            "endpoint": endpoint,
        }

    def post(self, endpoint, payload):
        return {
            "method": "POST",
            "endpoint": endpoint,
            "payload": payload,
        }

    def put(self, endpoint, payload):
        return {
            "method": "PUT",
            "endpoint": endpoint,
            "payload": payload,
        }

    def delete(self, endpoint):
        return {
            "method": "DELETE",
            "endpoint": endpoint,
        }


class PermissionServiceTest(unittest.TestCase):

    def setUp(self):
        self.service = PermissionService(
            FakeClient()
        )

    def test_list_permissions(self):
        result = self.service.list_permissions()

        self.assertEqual(
            result["endpoint"],
            "/permissions",
        )

    def test_get_permission(self):
        result = self.service.get_permission(
            "perm-1"
        )

        self.assertEqual(
            result["endpoint"],
            "/permissions/perm-1",
        )

    def test_get_permission_requires_id(self):
        with self.assertRaises(RuntimeError):
            self.service.get_permission("")

    def test_create_permission(self):
        result = self.service.create_permission(
            "users.read",
            "Read users",
        )

        self.assertEqual(
            result["payload"]["name"],
            "users.read",
        )

    def test_create_permission_requires_name(self):
        with self.assertRaises(RuntimeError):
            self.service.create_permission("")

    def test_update_permission(self):
        result = self.service.update_permission(
            "perm-1",
            {
                "description": "updated",
            },
        )

        self.assertEqual(
            result["method"],
            "PUT",
        )

    def test_update_permission_requires_id(self):
        with self.assertRaises(RuntimeError):
            self.service.update_permission(
                "",
                {},
            )

    def test_delete_permission(self):
        result = self.service.delete_permission(
            "perm-1"
        )

        self.assertEqual(
            result["method"],
            "DELETE",
        )

    def test_delete_permission_requires_id(self):
        with self.assertRaises(RuntimeError):
            self.service.delete_permission("")

    def test_assign_permission(self):
        result = self.service.assign_permission(
            "user-1",
            "users.read",
        )

        self.assertEqual(
            result["endpoint"],
            "/permissions/users/user-1/assign",
        )

    def test_assign_permission_requires_user(self):
        with self.assertRaises(RuntimeError):
            self.service.assign_permission(
                "",
                "users.read",
            )

    def test_assign_permission_requires_permission(self):
        with self.assertRaises(RuntimeError):
            self.service.assign_permission(
                "user-1",
                "",
            )

    def test_revoke_permission(self):
        result = self.service.revoke_permission(
            "user-1",
            "users.read",
        )

        self.assertEqual(
            result["endpoint"],
            "/permissions/users/user-1/revoke",
        )

    def test_role_permissions(self):
        result = self.service.role_permissions(
            "role-1"
        )

        self.assertEqual(
            result["endpoint"],
            "/permissions/roles/role-1",
        )

    def test_role_permissions_requires_id(self):
        with self.assertRaises(RuntimeError):
            self.service.role_permissions("")

    def test_user_permissions(self):
        result = self.service.user_permissions(
            "user-1"
        )

        self.assertEqual(
            result["endpoint"],
            "/permissions/users/user-1",
        )

    def test_user_permissions_requires_id(self):
        with self.assertRaises(RuntimeError):
            self.service.user_permissions("")

    def test_bulk_assign(self):
        result = self.service.bulk_assign(
            ["u1", "u2"],
            "users.read",
        )

        self.assertEqual(
            result["endpoint"],
            "/permissions/bulk-assign",
        )

    def test_bulk_assign_requires_users(self):
        with self.assertRaises(RuntimeError):
            self.service.bulk_assign(
                [],
                "users.read",
            )

    def test_bulk_assign_requires_permission(self):
        with self.assertRaises(RuntimeError):
            self.service.bulk_assign(
                ["u1"],
                "",
            )

    def test_bulk_revoke(self):
        result = self.service.bulk_revoke(
            ["u1", "u2"],
            "users.read",
        )

        self.assertEqual(
            result["endpoint"],
            "/permissions/bulk-revoke",
        )

    def test_bulk_revoke_requires_users(self):
        with self.assertRaises(RuntimeError):
            self.service.bulk_revoke(
                [],
                "users.read",
            )

    def test_bulk_revoke_requires_permission(self):
        with self.assertRaises(RuntimeError):
            self.service.bulk_revoke(
                ["u1"],
                "",
            )

    def test_health(self):
        result = self.service.health()

        self.assertEqual(
            result["endpoint"],
            "/permissions/health",
        )


if __name__ == "__main__":
    unittest.main()