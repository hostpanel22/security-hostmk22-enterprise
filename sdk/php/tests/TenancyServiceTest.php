<?php

declare(strict_types=1);

namespace HostMK\Security\Tests;

use HostMK\Security\TenancyService;
use HostMK\Security\Tests\Fakes\FakeClient;
use InvalidArgumentException;
use PHPUnit\Framework\TestCase;

final class TenancyServiceTest extends TestCase
{
    private FakeClient $client;

    private TenancyService $service;

    protected function setUp(): void
    {
        $this->client = new FakeClient();

        $this->service = new TenancyService(
            $this->client
        );
    }

    public function testListCallsApi(): void
    {
        $response = $this->service->list();

        $this->assertSame(
            'GET',
            $this->client->method
        );

        $this->assertSame(
            '/tenants',
            $this->client->endpoint
        );

        $this->assertTrue(
            $response['success']
        );
    }

    public function testListWithFiltersCallsApi(): void
    {
        $response = $this->service->list(
            [
                'page' => 2,
                'status' => 'active',
            ]
        );

        $this->assertSame(
            'GET',
            $this->client->method
        );

        $this->assertSame(
            '/tenants?page=2&status=active',
            $this->client->endpoint
        );

        $this->assertTrue(
            $response['success']
        );
    }

    public function testGetThrowsWhenTenantIdEmpty(): void
    {
        $this->expectException(
            InvalidArgumentException::class
        );

        $this->expectExceptionMessage(
            'Tenant ID is required.'
        );

        $this->service->get('');
    }

    public function testGetCallsApi(): void
    {
        $response = $this->service->get(
            'tenant-123'
        );

        $this->assertSame(
            'GET',
            $this->client->method
        );

        $this->assertSame(
            '/tenants/tenant-123',
            $this->client->endpoint
        );

        $this->assertTrue(
            $response['success']
        );
    }

    public function testCreateThrowsWhenNameEmpty(): void
    {
        $this->expectException(
            InvalidArgumentException::class
        );

        $this->expectExceptionMessage(
            'Tenant name is required.'
        );

        $this->service->create(
            '',
            'tenant-slug'
        );
    }

    public function testCreateThrowsWhenSlugEmpty(): void
    {
        $this->expectException(
            InvalidArgumentException::class
        );

        $this->expectExceptionMessage(
            'Tenant slug is required.'
        );

        $this->service->create(
            'Tenant Name',
            ''
        );
    }

    public function testCreateCallsApi(): void
    {
        $response = $this->service->create(
            'Tenant Name',
            'tenant-slug'
        );

        $this->assertSame(
            'POST',
            $this->client->method
        );

        $this->assertSame(
            '/tenants',
            $this->client->endpoint
        );

        $this->assertSame(
            [
                'name' => 'Tenant Name',
                'slug' => 'tenant-slug',
            ],
            $this->client->payload
        );

        $this->assertTrue(
            $response['success']
        );
    }

    public function testUpdateThrowsWhenTenantIdEmpty(): void
    {
        $this->expectException(
            InvalidArgumentException::class
        );

        $this->expectExceptionMessage(
            'Tenant ID is required.'
        );

        $this->service->update(
            '',
            [
                'name' => 'Updated',
            ]
        );
    }

    public function testUpdateThrowsWhenPayloadEmpty(): void
    {
        $this->expectException(
            InvalidArgumentException::class
        );

        $this->expectExceptionMessage(
            'Update payload is required.'
        );

        $this->service->update(
            'tenant-123',
            []
        );
    }

    public function testUpdateCallsApi(): void
    {
        $payload = [
            'name' => 'Updated Tenant',
        ];

        $response = $this->service->update(
            'tenant-123',
            $payload
        );

        $this->assertSame(
            'PUT',
            $this->client->method
        );

        $this->assertSame(
            '/tenants/tenant-123',
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

    public function testDeleteThrowsWhenTenantIdEmpty(): void
    {
        $this->expectException(
            InvalidArgumentException::class
        );

        $this->expectExceptionMessage(
            'Tenant ID is required.'
        );

        $this->service->delete('');
    }

    public function testDeleteCallsApi(): void
    {
        $response = $this->service->delete(
            'tenant-123'
        );

        $this->assertSame(
            'DELETE',
            $this->client->method
        );

        $this->assertSame(
            '/tenants/tenant-123',
            $this->client->endpoint
        );

        $this->assertTrue(
            $response['success']
        );
    }

    public function testActivateThrowsWhenTenantIdEmpty(): void
    {
        $this->expectException(
            InvalidArgumentException::class
        );

        $this->expectExceptionMessage(
            'Tenant ID is required.'
        );

        $this->service->activate('');
    }

    public function testActivateCallsApi(): void
    {
        $response = $this->service->activate(
            'tenant-123'
        );

        $this->assertSame(
            'POST',
            $this->client->method
        );

        $this->assertSame(
            '/tenants/tenant-123/activate',
            $this->client->endpoint
        );

        $this->assertTrue(
            $response['success']
        );
    }

    public function testDeactivateThrowsWhenTenantIdEmpty(): void
    {
        $this->expectException(
            InvalidArgumentException::class
        );

        $this->expectExceptionMessage(
            'Tenant ID is required.'
        );

        $this->service->deactivate('');
    }

    public function testDeactivateCallsApi(): void
    {
        $response = $this->service->deactivate(
            'tenant-123'
        );

        $this->assertSame(
            'POST',
            $this->client->method
        );

        $this->assertSame(
            '/tenants/tenant-123/deactivate',
            $this->client->endpoint
        );

        $this->assertTrue(
            $response['success']
        );
    }
}