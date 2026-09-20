<?php

declare(strict_types=1);

namespace HostMK\Security\Tests;

use HostMK\Security\AuditService;
use HostMK\Security\Tests\Fakes\FakeClient;
use InvalidArgumentException;
use PHPUnit\Framework\TestCase;

final class AuditServiceTest extends TestCase
{
    private FakeClient $client;

    private AuditService $service;

    protected function setUp(): void
    {
        $this->client = new FakeClient();

        $this->service = new AuditService(
            $this->client
        );
    }

    public function testListCallsApi(): void
    {
        $this->service->list();

        self::assertSame(
            'GET',
            $this->client->method
        );

        self::assertSame(
            '/audit',
            $this->client->endpoint
        );
    }

    public function testListWithFiltersCallsApi(): void
    {
        $this->service->list([
            'action' => 'login',
            'userId' => '123',
        ]);

        self::assertSame(
            'GET',
            $this->client->method
        );

        self::assertStringContainsString(
            '/audit?',
            $this->client->endpoint
        );
    }

    public function testGetThrowsWhenAuditIdEmpty(): void
    {
        $this->expectException(
            InvalidArgumentException::class
        );

        $this->service->get('');
    }

    public function testGetCallsApi(): void
    {
        $this->service->get(
            'audit-1'
        );

        self::assertSame(
            'GET',
            $this->client->method
        );

        self::assertSame(
            '/audit/audit-1',
            $this->client->endpoint
        );
    }

    public function testCreateThrowsWhenPayloadEmpty(): void
    {
        $this->expectException(
            InvalidArgumentException::class
        );

        $this->service->create([]);
    }

    public function testCreateCallsApi(): void
    {
        $payload = [
            'action' => 'login',
        ];

        $this->service->create(
            $payload
        );

        self::assertSame(
            'POST',
            $this->client->method
        );

        self::assertSame(
            '/audit',
            $this->client->endpoint
        );

        self::assertSame(
            $payload,
            $this->client->payload
        );
    }

    public function testDeleteThrowsWhenAuditIdEmpty(): void
    {
        $this->expectException(
            InvalidArgumentException::class
        );

        $this->service->delete('');
    }

    public function testDeleteCallsApi(): void
    {
        $this->service->delete(
            'audit-1'
        );

        self::assertSame(
            'DELETE',
            $this->client->method
        );

        self::assertSame(
            '/audit/audit-1',
            $this->client->endpoint
        );
    }

    public function testSearchCallsApi(): void
    {
        $this->service->search([
            'action' => 'login',
        ]);

        self::assertSame(
            'GET',
            $this->client->method
        );

        self::assertStringContainsString(
            '/audit/search',
            $this->client->endpoint
        );
    }

    public function testStatsCallsApi(): void
    {
        $this->service->stats();

        self::assertSame(
            'GET',
            $this->client->method
        );

        self::assertSame(
            '/audit/stats',
            $this->client->endpoint
        );
    }

    public function testExportCallsApi(): void
    {
        $payload = [
            'format' => 'csv',
        ];

        $this->service->export(
            $payload
        );

        self::assertSame(
            'POST',
            $this->client->method
        );

        self::assertSame(
            '/audit/export',
            $this->client->endpoint
        );

        self::assertSame(
            $payload,
            $this->client->payload
        );
    }

    public function testUserActivityThrowsWhenUserIdEmpty(): void
    {
        $this->expectException(
            InvalidArgumentException::class
        );

        $this->service->userActivity('');
    }

    public function testUserActivityCallsApi(): void
    {
        $this->service->userActivity(
            'user-1'
        );

        self::assertSame(
            'GET',
            $this->client->method
        );

        self::assertSame(
            '/audit/users/user-1',
            $this->client->endpoint
        );
    }

    public function testTenantActivityThrowsWhenTenantIdEmpty(): void
    {
        $this->expectException(
            InvalidArgumentException::class
        );

        $this->service->tenantActivity('');
    }

    public function testTenantActivityCallsApi(): void
    {
        $this->service->tenantActivity(
            'tenant-1'
        );

        self::assertSame(
            'GET',
            $this->client->method
        );

        self::assertSame(
            '/audit/tenants/tenant-1',
            $this->client->endpoint
        );
    }

    public function testRetentionThrowsWhenDaysInvalid(): void
    {
        $this->expectException(
            InvalidArgumentException::class
        );

        $this->service->retention(0);
    }

    public function testRetentionCallsApi(): void
    {
        $this->service->retention(
            30
        );

        self::assertSame(
            'POST',
            $this->client->method
        );

        self::assertSame(
            '/audit/retention',
            $this->client->endpoint
        );

        self::assertSame(
            [
                'days' => 30,
            ],
            $this->client->payload
        );
    }

    public function testPurgeThrowsWhenDateEmpty(): void
    {
        $this->expectException(
            InvalidArgumentException::class
        );

        $this->service->purge('');
    }

    public function testPurgeCallsApi(): void
    {
        $this->service->purge(
            '2026-01-01'
        );

        self::assertSame(
            'POST',
            $this->client->method
        );

        self::assertSame(
            '/audit/purge',
            $this->client->endpoint
        );

        self::assertSame(
            [
                'beforeDate' => '2026-01-01',
            ],
            $this->client->payload
        );
    }

    public function testCategoriesCallsApi(): void
    {
        $this->service->categories();

        self::assertSame(
            'GET',
            $this->client->method
        );

        self::assertSame(
            '/audit/categories',
            $this->client->endpoint
        );
    }

    public function testActionsCallsApi(): void
    {
        $this->service->actions();

        self::assertSame(
            'GET',
            $this->client->method
        );

        self::assertSame(
            '/audit/actions',
            $this->client->endpoint
        );
    }

    public function testHealthCallsApi(): void
    {
        $this->service->health();

        self::assertSame(
            'GET',
            $this->client->method
        );

        self::assertSame(
            '/audit/health',
            $this->client->endpoint
        );
    }
}