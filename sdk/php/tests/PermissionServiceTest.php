<?php

declare(strict_types=1);

namespace HostMK\Security\Tests;

use HostMK\Security\PermissionService;
use HostMK\Security\Tests\Fakes\FakeClient;
use InvalidArgumentException;
use PHPUnit\Framework\TestCase;

final class PermissionServiceTest extends TestCase
{
    private FakeClient $client;

    private PermissionService $service;

    protected function setUp(): void
    {
        $this->client = new FakeClient();

        $this->service = new PermissionService(
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
            '/permissions',
            $this->client->endpoint
        );

        $this->assertTrue(
            $response['success']
        );
    }

    public function testGetUserPermissionsThrowsWhenUserIdEmpty(): void
    {
        $this->expectException(
            InvalidArgumentException::class
        );

        $this->expectExceptionMessage(
            'User ID is required.'
        );

        $this->service->getUserPermissions('');
    }

    public function testGetUserPermissionsCallsApi(): void
    {
        $response = $this->service->getUserPermissions(
            'user-123'
        );

        $this->assertSame(
            'GET',
            $this->client->method
        );

        $this->assertSame(
            '/permissions/user-123',
            $this->client->endpoint
        );

        $this->assertTrue(
            $response['success']
        );
    }

    public function testCheckThrowsWhenUserIdEmpty(): void
    {
        $this->expectException(
            InvalidArgumentException::class
        );

        $this->expectExceptionMessage(
            'User ID is required.'
        );

        $this->service->check(
            '',
            'audit.read'
        );
    }

    public function testCheckThrowsWhenPermissionEmpty(): void
    {
        $this->expectException(
            InvalidArgumentException::class
        );

        $this->expectExceptionMessage(
            'Permission is required.'
        );

        $this->service->check(
            'user-123',
            ''
        );
    }

    public function testCheckCallsApi(): void
    {
        $response = $this->service->check(
            'user-123',
            'audit.read'
        );

        $this->assertSame(
            'POST',
            $this->client->method
        );

        $this->assertSame(
            '/permissions/check',
            $this->client->endpoint
        );

        $this->assertSame(
            [
                'userId' => 'user-123',
                'permission' => 'audit.read',
            ],
            $this->client->payload
        );

        $this->assertTrue(
            $response['success']
        );
    }

    public function testAssignRoleThrowsWhenUserIdEmpty(): void
    {
        $this->expectException(
            InvalidArgumentException::class
        );

        $this->expectExceptionMessage(
            'User ID is required.'
        );

        $this->service->assignRole(
            '',
            'admin'
        );
    }

    public function testAssignRoleThrowsWhenRoleEmpty(): void
    {
        $this->expectException(
            InvalidArgumentException::class
        );

        $this->expectExceptionMessage(
            'Role is required.'
        );

        $this->service->assignRole(
            'user-123',
            ''
        );
    }

    public function testAssignRoleCallsApi(): void
    {
        $response = $this->service->assignRole(
            'user-123',
            'admin'
        );

        $this->assertSame(
            'POST',
            $this->client->method
        );

        $this->assertSame(
            '/permissions/assign-role',
            $this->client->endpoint
        );

        $this->assertSame(
            [
                'userId' => 'user-123',
                'role' => 'admin',
            ],
            $this->client->payload
        );

        $this->assertTrue(
            $response['success']
        );
    }

    public function testRemoveRoleThrowsWhenUserIdEmpty(): void
    {
        $this->expectException(
            InvalidArgumentException::class
        );

        $this->expectExceptionMessage(
            'User ID is required.'
        );

        $this->service->removeRole(
            '',
            'admin'
        );
    }

    public function testRemoveRoleThrowsWhenRoleEmpty(): void
    {
        $this->expectException(
            InvalidArgumentException::class
        );

        $this->expectExceptionMessage(
            'Role is required.'
        );

        $this->service->removeRole(
            'user-123',
            ''
        );
    }

    public function testRemoveRoleCallsApi(): void
    {
        $response = $this->service->removeRole(
            'user-123',
            'admin'
        );

        $this->assertSame(
            'POST',
            $this->client->method
        );

        $this->assertSame(
            '/permissions/remove-role',
            $this->client->endpoint
        );

        $this->assertSame(
            [
                'userId' => 'user-123',
                'role' => 'admin',
            ],
            $this->client->payload
        );

        $this->assertTrue(
            $response['success']
        );
    }
}