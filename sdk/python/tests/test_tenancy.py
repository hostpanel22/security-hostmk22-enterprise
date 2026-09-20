import unittest

from hostmk_security.tenancy import TenancyService


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


class TenancyServiceTest(unittest.TestCase):

    def setUp(self):
        self.service = TenancyService(
            FakeClient()
        )

    def test_list_tenants(self):
        result = self.service.list_tenants()

        self.assertEqual(
            result["endpoint"],
            "/tenancy/tenants",
        )

    def test_get_tenant(self):
        result = self.service.get_tenant(
            "tenant-1"
        )

        self.assertEqual(
            result["endpoint"],
            "/tenancy/tenants/tenant-1",
        )

    def test_get_tenant_requires_id(self):
        with self.assertRaises(RuntimeError):
            self.service.get_tenant("")

    def test_create_tenant(self):
        result = self.service.create_tenant(
            "Tenant One",
            "tenant-one",
        )

        self.assertEqual(
            result["payload"]["name"],
            "Tenant One",
        )

        self.assertEqual(
            result["payload"]["slug"],
            "tenant-one",
        )

    def test_create_tenant_requires_name(self):
        with self.assertRaises(RuntimeError):
            self.service.create_tenant(
                "",
                "tenant-one",
            )

    def test_create_tenant_requires_slug(self):
        with self.assertRaises(RuntimeError):
            self.service.create_tenant(
                "Tenant One",
                "",
            )

    def test_update_tenant(self):
        result = self.service.update_tenant(
            "tenant-1",
            {
                "name": "Updated",
            },
        )

        self.assertEqual(
            result["method"],
            "PUT",
        )

    def test_update_tenant_requires_id(self):
        with self.assertRaises(RuntimeError):
            self.service.update_tenant(
                "",
                {},
            )

    def test_delete_tenant(self):
        result = self.service.delete_tenant(
            "tenant-1"
        )

        self.assertEqual(
            result["method"],
            "DELETE",
        )

    def test_delete_tenant_requires_id(self):
        with self.assertRaises(RuntimeError):
            self.service.delete_tenant("")

    def test_suspend_tenant(self):
        result = self.service.suspend_tenant(
            "tenant-1"
        )

        self.assertEqual(
            result["endpoint"],
            "/tenancy/tenants/tenant-1/suspend",
        )

    def test_suspend_tenant_requires_id(self):
        with self.assertRaises(RuntimeError):
            self.service.suspend_tenant("")

    def test_activate_tenant(self):
        result = self.service.activate_tenant(
            "tenant-1"
        )

        self.assertEqual(
            result["endpoint"],
            "/tenancy/tenants/tenant-1/activate",
        )

    def test_activate_tenant_requires_id(self):
        with self.assertRaises(RuntimeError):
            self.service.activate_tenant("")

    def test_tenant_usage(self):
        result = self.service.tenant_usage(
            "tenant-1"
        )

        self.assertEqual(
            result["endpoint"],
            "/tenancy/tenants/tenant-1/usage",
        )

    def test_tenant_usage_requires_id(self):
        with self.assertRaises(RuntimeError):
            self.service.tenant_usage("")

    def test_tenant_statistics(self):
        result = self.service.tenant_statistics(
            "tenant-1"
        )

        self.assertEqual(
            result["endpoint"],
            "/tenancy/tenants/tenant-1/statistics",
        )

    def test_tenant_statistics_requires_id(self):
        with self.assertRaises(RuntimeError):
            self.service.tenant_statistics("")

    def test_tenant_domains(self):
        result = self.service.tenant_domains(
            "tenant-1"
        )

        self.assertEqual(
            result["endpoint"],
            "/tenancy/tenants/tenant-1/domains",
        )

    def test_tenant_domains_requires_id(self):
        with self.assertRaises(RuntimeError):
            self.service.tenant_domains("")

    def test_tenant_settings(self):
        result = self.service.tenant_settings(
            "tenant-1"
        )

        self.assertEqual(
            result["endpoint"],
            "/tenancy/tenants/tenant-1/settings",
        )

    def test_tenant_settings_requires_id(self):
        with self.assertRaises(RuntimeError):
            self.service.tenant_settings("")

    def test_update_settings(self):
        result = self.service.update_settings(
            "tenant-1",
            {
                "theme": "dark",
            },
        )

        self.assertEqual(
            result["method"],
            "PUT",
        )

    def test_update_settings_requires_id(self):
        with self.assertRaises(RuntimeError):
            self.service.update_settings(
                "",
                {},
            )

    def test_health(self):
        result = self.service.health()

        self.assertEqual(
            result["endpoint"],
            "/tenancy/health",
        )


if __name__ == "__main__":
    unittest.main()