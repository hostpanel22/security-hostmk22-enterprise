import unittest
from unittest.mock import Mock

from hostmk_security.fraud import FraudService


class FraudServiceTest(unittest.TestCase):

    def setUp(self) -> None:
        self.client = Mock()
        self.service = FraudService(
            self.client
        )

    def test_analyze(self) -> None:
        payload = {
            "transaction_id": "tx_123"
        }

        expected = {
            "status": 200,
            "data": {
                "risk": "low"
            }
        }

        self.client.post.return_value = expected

        result = self.service.analyze(
            payload
        )

        self.assertEqual(
            expected,
            result,
        )

        self.client.post.assert_called_once_with(
            "fraud/analyze",
            payload,
        )

    def test_score(self) -> None:
        payload = {
            "amount": 100
        }

        expected = {
            "status": 200,
            "data": {
                "score": 15
            }
        }

        self.client.post.return_value = expected

        result = self.service.score(
            payload
        )

        self.assertEqual(
            expected,
            result,
        )

        self.client.post.assert_called_once_with(
            "fraud/score",
            payload,
        )

    def test_check_ip(self) -> None:
        expected = {
            "status": 200
        }

        self.client.get.return_value = expected

        result = self.service.check_ip(
            "1.1.1.1"
        )

        self.assertEqual(
            expected,
            result,
        )

        self.client.get.assert_called_once_with(
            "fraud/ip/1.1.1.1"
        )

    def test_check_ip_requires_value(
        self,
    ) -> None:
        with self.assertRaises(
            RuntimeError
        ):
            self.service.check_ip("")

    def test_check_email(self) -> None:
        expected = {
            "status": 200
        }

        self.client.get.return_value = expected

        result = self.service.check_email(
            "user@example.com"
        )

        self.assertEqual(
            expected,
            result,
        )

        self.client.get.assert_called_once_with(
            "fraud/email/user@example.com"
        )

    def test_check_email_requires_value(
        self,
    ) -> None:
        with self.assertRaises(
            RuntimeError
        ):
            self.service.check_email("")

    def test_check_device(self) -> None:
        expected = {
            "status": 200
        }

        self.client.get.return_value = expected

        result = self.service.check_device(
            "device_123"
        )

        self.assertEqual(
            expected,
            result,
        )

        self.client.get.assert_called_once_with(
            "fraud/device/device_123"
        )

    def test_check_device_requires_value(
        self,
    ) -> None:
        with self.assertRaises(
            RuntimeError
        ):
            self.service.check_device("")

    def test_report(self) -> None:
        payload = {
            "reason": "fraud"
        }

        expected = {
            "status": 201
        }

        self.client.post.return_value = expected

        result = self.service.report(
            payload
        )

        self.assertEqual(
            expected,
            result,
        )

        self.client.post.assert_called_once_with(
            "fraud/report",
            payload,
        )

    def test_blacklist(self) -> None:
        expected = {
            "status": 201
        }

        self.client.post.return_value = expected

        result = self.service.blacklist(
            "ip",
            "8.8.8.8"
        )

        self.assertEqual(
            expected,
            result,
        )

        self.client.post.assert_called_once_with(
            "fraud/blacklist",
            {
                "type": "ip",
                "value": "8.8.8.8",
            },
        )

    def test_blacklist_requires_type(
        self,
    ) -> None:
        with self.assertRaises(
            RuntimeError
        ):
            self.service.blacklist(
                "",
                "8.8.8.8"
            )

    def test_blacklist_requires_value(
        self,
    ) -> None:
        with self.assertRaises(
            RuntimeError
        ):
            self.service.blacklist(
                "ip",
                ""
            )

    def test_whitelist(self) -> None:
        expected = {
            "status": 201
        }

        self.client.post.return_value = expected

        result = self.service.whitelist(
            "email",
            "trusted@example.com"
        )

        self.assertEqual(
            expected,
            result,
        )

        self.client.post.assert_called_once_with(
            "fraud/whitelist",
            {
                "type": "email",
                "value": "trusted@example.com",
            },
        )

    def test_whitelist_requires_type(
        self,
    ) -> None:
        with self.assertRaises(
            RuntimeError
        ):
            self.service.whitelist(
                "",
                "value"
            )

    def test_whitelist_requires_value(
        self,
    ) -> None:
        with self.assertRaises(
            RuntimeError
        ):
            self.service.whitelist(
                "email",
                ""
            )

    def test_alerts(self) -> None:
        expected = {
            "status": 200
        }

        self.client.get.return_value = expected

        result = self.service.alerts()

        self.assertEqual(
            expected,
            result,
        )

        self.client.get.assert_called_once_with(
            "fraud/alerts?page=1&limit=25"
        )

    def test_statistics(self) -> None:
        expected = {
            "status": 200
        }

        self.client.get.return_value = expected

        result = self.service.statistics()

        self.assertEqual(
            expected,
            result,
        )

        self.client.get.assert_called_once_with(
            "fraud/statistics"
        )

    def test_health(self) -> None:
        expected = {
            "status": 200
        }

        self.client.get.return_value = expected

        result = self.service.health()

        self.assertEqual(
            expected,
            result,
        )

        self.client.get.assert_called_once_with(
            "fraud/health"
        )


if __name__ == "__main__":
    unittest.main()