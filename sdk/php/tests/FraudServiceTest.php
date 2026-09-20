<?php

declare(strict_types=1);

namespace HostMK\Security\Tests;

use HostMK\Security\FraudService;
use HostMK\Security\Tests\Fakes\FakeClient;
use InvalidArgumentException;
use PHPUnit\Framework\TestCase;

final class FraudServiceTest extends TestCase
{
    private FakeClient $client;

    private FraudService $service;

    protected function setUp(): void
    {
        $this->client = new FakeClient();

        $this->service = new FraudService(
            $this->client
        );
    }

    public function testCheckThrowsWhenPayloadEmpty(): void
    {
        $this->expectException(
            InvalidArgumentException::class
        );

        $this->expectExceptionMessage(
            'Fraud check payload is required.'
        );

        $this->service->check([]);
    }

    public function testCheckCallsApi(): void
    {
        $payload = [
            'userId' => 'user-1',
            'ip' => '127.0.0.1',
        ];

        $response = $this->service->check(
            $payload
        );

        $this->assertSame(
            'POST',
            $this->client->method
        );

        $this->assertSame(
            '/fraud/check',
            $this->client->endpoint
        );

        $this->assertSame(
            $payload,
            $this->client->payload
        );

        $this->assertTrue(
            $response['success']
        );
    }

    public function testScoreThrowsWhenPayloadEmpty(): void
    {
        $this->expectException(
            InvalidArgumentException::class
        );

        $this->expectExceptionMessage(
            'Fraud score payload is required.'
        );

        $this->service->score([]);
    }

    public function testScoreCallsApi(): void
    {
        $payload = [
            'userId' => 'user-1',
            'amount' => 100,
        ];

        $response = $this->service->score(
            $payload
        );

        $this->assertSame(
            '/fraud/score',
            $this->client->endpoint
        );

        $this->assertSame(
            'POST',
            $this->client->method
        );

        $this->assertSame(
            $payload,
            $this->client->payload
        );

        $this->assertTrue(
            $response['success']
        );
    }

    public function testReportThrowsWhenPayloadEmpty(): void
    {
        $this->expectException(
            InvalidArgumentException::class
        );

        $this->expectExceptionMessage(
            'Fraud report payload is required.'
        );

        $this->service->report([]);
    }

    public function testReportCallsApi(): void
    {
        $payload = [
            'type' => 'chargeback',
            'userId' => 'user-1',
        ];

        $response = $this->service->report(
            $payload
        );

        $this->assertSame(
            '/fraud/report',
            $this->client->endpoint
        );

        $this->assertSame(
            'POST',
            $this->client->method
        );

        $this->assertSame(
            $payload,
            $this->client->payload
        );

        $this->assertTrue(
            $response['success']
        );
    }

    public function testReportsCallsApi(): void
    {
        $response = $this->service->reports();

        $this->assertSame(
            'GET',
            $this->client->method
        );

        $this->assertSame(
            '/fraud/reports',
            $this->client->endpoint
        );

        $this->assertTrue(
            $response['success']
        );
    }

    public function testReportsWithFiltersCallsApi(): void
    {
        $response = $this->service->reports(
            [
                'status' => 'open',
                'page' => 2,
            ]
        );

        $this->assertSame(
            '/fraud/reports?status=open&page=2',
            $this->client->endpoint
        );

        $this->assertTrue(
            $response['success']
        );
    }

    public function testGetReportThrowsWhenReportIdEmpty(): void
    {
        $this->expectException(
            InvalidArgumentException::class
        );

        $this->expectExceptionMessage(
            'Report ID is required.'
        );

        $this->service->getReport('');
    }

    public function testGetReportCallsApi(): void
    {
        $response = $this->service->getReport(
            'report-123'
        );

        $this->assertSame(
            '/fraud/reports/report-123',
            $this->client->endpoint
        );

        $this->assertTrue(
            $response['success']
        );
    }

    public function testStatisticsCallsApi(): void
    {
        $response = $this->service->statistics(
            [
                'days' => 30,
            ]
        );

        $this->assertSame(
            'GET',
            $this->client->method
        );

        $this->assertSame(
            '/fraud/statistics?days=30',
            $this->client->endpoint
        );

        $this->assertTrue(
            $response['success']
        );
    }

    public function testUserRiskThrowsWhenUserIdEmpty(): void
    {
        $this->expectException(
            InvalidArgumentException::class
        );

        $this->expectExceptionMessage(
            'User ID is required.'
        );

        $this->service->userRisk('');
    }

    public function testUserRiskCallsApi(): void
    {
        $response = $this->service->userRisk(
            'user-123'
        );

        $this->assertSame(
            '/fraud/users/user-123/risk',
            $this->client->endpoint
        );

        $this->assertTrue(
            $response['success']
        );
    }

    public function testIpRiskThrowsWhenIpEmpty(): void
    {
        $this->expectException(
            InvalidArgumentException::class
        );

        $this->expectExceptionMessage(
            'IP address is required.'
        );

        $this->service->ipRisk('');
    }

    public function testIpRiskCallsApi(): void
    {
        $response = $this->service->ipRisk(
            '127.0.0.1'
        );

        $this->assertSame(
            '/fraud/ip/127.0.0.1',
            $this->client->endpoint
        );

        $this->assertTrue(
            $response['success']
        );
    }

    public function testDeviceRiskThrowsWhenDeviceIdEmpty(): void
    {
        $this->expectException(
            InvalidArgumentException::class
        );

        $this->expectExceptionMessage(
            'Device ID is required.'
        );

        $this->service->deviceRisk('');
    }

    public function testDeviceRiskCallsApi(): void
    {
        $response = $this->service->deviceRisk(
            'device-123'
        );

        $this->assertSame(
            '/fraud/devices/device-123',
            $this->client->endpoint
        );

        $this->assertTrue(
            $response['success']
        );
    }

    public function testWhitelistThrowsWhenPayloadEmpty(): void
    {
        $this->expectException(
            InvalidArgumentException::class
        );

        $this->expectExceptionMessage(
            'Whitelist payload is required.'
        );

        $this->service->whitelist([]);
    }

    public function testWhitelistCallsApi(): void
    {
        $payload = [
            'ip' => '127.0.0.1',
        ];

        $response = $this->service->whitelist(
            $payload
        );

        $this->assertSame(
            '/fraud/whitelist',
            $this->client->endpoint
        );

        $this->assertSame(
            $payload,
            $this->client->payload
        );

        $this->assertTrue(
            $response['success']
        );
    }

    public function testBlacklistThrowsWhenPayloadEmpty(): void
    {
        $this->expectException(
            InvalidArgumentException::class
        );

        $this->expectExceptionMessage(
            'Blacklist payload is required.'
        );

        $this->service->blacklist([]);
    }

    public function testBlacklistCallsApi(): void
    {
        $payload = [
            'ip' => '10.0.0.1',
        ];

        $response = $this->service->blacklist(
            $payload
        );

        $this->assertSame(
            '/fraud/blacklist',
            $this->client->endpoint
        );

        $this->assertSame(
            $payload,
            $this->client->payload
        );

        $this->assertTrue(
            $response['success']
        );
    }

    public function testGetWhitelistCallsApi(): void
    {
        $response = $this->service->getWhitelist();

        $this->assertSame(
            '/fraud/whitelist',
            $this->client->endpoint
        );

        $this->assertSame(
            'GET',
            $this->client->method
        );

        $this->assertTrue(
            $response['success']
        );
    }

    public function testGetBlacklistCallsApi(): void
    {
        $response = $this->service->getBlacklist();

        $this->assertSame(
            '/fraud/blacklist',
            $this->client->endpoint
        );

        $this->assertSame(
            'GET',
            $this->client->method
        );

        $this->assertTrue(
            $response['success']
        );
    }

    public function testHealthCallsApi(): void
    {
        $response = $this->service->health();

        $this->assertSame(
            '/fraud/health',
            $this->client->endpoint
        );

        $this->assertSame(
            'GET',
            $this->client->method
        );

        $this->assertTrue(
            $response['success']
        );
    }
}