import unittest
from unittest.mock import Mock

from hostmk_security.audit import AuditService


class AuditServiceTest(unittest.TestCase):

    def setUp(self) -> None:
        self.client = Mock()
        self.service = AuditService(
            self.client
        )

    def test_list_events(self) -> None:
        expected = {
            "status": 200,
            "data": [],
        }

        self.client.get.return_value = expected

        result = self.service.list_events()

        self.assertEqual(
            expected,
            result,
        )

        self.client.get.assert_called_once_with(
            "audit/events?page=1&limit=25"
        )

    def test_get_event(self) -> None:
        expected = {
            "status": 200,
            "data": {
                "id": "evt_123"
            },
        }

        self.client.get.return_value = expected

        result = self.service.get_event(
            "evt_123"
        )

        self.assertEqual(
            expected,
            result,
        )

        self.client.get.assert_called_once_with(
            "audit/events/evt_123"
        )

    def test_get_event_requires_id(
        self,
    ) -> None:
        with self.assertRaises(
            RuntimeError
        ):
            self.service.get_event("")

    def test_search_events(
        self,
    ) -> None:
        filters = {
            "action": "login"
        }

        expected = {
            "status": 200,
            "data": [],
        }

        self.client.post.return_value = (
            expected
        )

        result = self.service.search_events(
            filters
        )

        self.assertEqual(
            expected,
            result,
        )

        self.client.post.assert_called_once_with(
            "audit/events/search",
            filters,
        )

    def test_user_activity(
        self,
    ) -> None:
        expected = {
            "status": 200,
            "data": [],
        }

        self.client.get.return_value = expected

        result = self.service.user_activity(
            "user_1"
        )

        self.assertEqual(
            expected,
            result,
        )

        self.client.get.assert_called_once_with(
            "audit/users/user_1/activity?page=1&limit=25"
        )

    def test_user_activity_requires_id(
        self,
    ) -> None:
        with self.assertRaises(
            RuntimeError
        ):
            self.service.user_activity(
                ""
            )

    def test_tenant_activity(
        self,
    ) -> None:
        expected = {
            "status": 200,
            "data": [],
        }

        self.client.get.return_value = expected

        result = self.service.tenant_activity(
            "tenant_1"
        )

        self.assertEqual(
            expected,
            result,
        )

        self.client.get.assert_called_once_with(
            "audit/tenants/tenant_1/activity?page=1&limit=25"
        )

    def test_tenant_activity_requires_id(
        self,
    ) -> None:
        with self.assertRaises(
            RuntimeError
        ):
            self.service.tenant_activity(
                ""
            )

    def test_export_events(
        self,
    ) -> None:
        filters = {
            "severity": "high"
        }

        expected = {
            "status": 200
        }

        self.client.post.return_value = (
            expected
        )

        result = self.service.export_events(
            filters
        )

        self.assertEqual(
            expected,
            result,
        )

        self.client.post.assert_called_once_with(
            "audit/export",
            filters,
        )

    def test_compliance_report(
        self,
    ) -> None:
        expected = {
            "status": 200
        }

        self.client.get.return_value = expected

        result = (
            self.service.compliance_report(
                "gdpr"
            )
        )

        self.assertEqual(
            expected,
            result,
        )

        self.client.get.assert_called_once_with(
            "audit/compliance/gdpr"
        )

    def test_compliance_report_requires_type(
        self,
    ) -> None:
        with self.assertRaises(
            RuntimeError
        ):
            self.service.compliance_report(
                ""
            )

    def test_security_events(
        self,
    ) -> None:
        expected = {
            "status": 200
        }

        self.client.get.return_value = expected

        result = (
            self.service.security_events()
        )

        self.assertEqual(
            expected,
            result,
        )

        self.client.get.assert_called_once_with(
            "audit/security-events"
        )

    def test_security_events_with_severity(
        self,
    ) -> None:
        expected = {
            "status": 200
        }

        self.client.get.return_value = expected

        result = (
            self.service.security_events(
                "critical"
            )
        )

        self.assertEqual(
            expected,
            result,
        )

        self.client.get.assert_called_once_with(
            "audit/security-events?severity=critical"
        )

    def test_statistics(
        self,
    ) -> None:
        expected = {
            "status": 200
        }

        self.client.get.return_value = expected

        result = (
            self.service.statistics()
        )

        self.assertEqual(
            expected,
            result,
        )

        self.client.get.assert_called_once_with(
            "audit/statistics"
        )

    def test_health(
        self,
    ) -> None:
        expected = {
            "status": 200
        }

        self.client.get.return_value = expected

        result = (
            self.service.health()
        )

        self.assertEqual(
            expected,
            result,
        )

        self.client.get.assert_called_once_with(
            "audit/health"
        )


if __name__ == "__main__":
    unittest.main()