<?php

declare(strict_types=1);

namespace HostMK\Security\Tests;

use HostMK\Security\AuthService;
use HostMK\Security\Tests\Fakes\FakeClient;
use InvalidArgumentException;
use PHPUnit\Framework\TestCase;

final class AuthServiceTest extends TestCase
{
    public function testLoginThrowsWhenEmailEmpty(): void
    {
        $service = new AuthService(
            new FakeClient()
        );

        $this->expectException(
            InvalidArgumentException::class
        );

        $this->expectExceptionMessage(
            'Email is required.'
        );

        $service->login(
            '',
            'password'
        );
    }

    public function testLoginThrowsWhenEmailInvalid(): void
    {
        $service = new AuthService(
            new FakeClient()
        );

        $this->expectException(
            InvalidArgumentException::class
        );

        $this->expectExceptionMessage(
            'Invalid email address.'
        );

        $service->login(
            'invalid-email',
            'password'
        );
    }

    public function testLoginThrowsWhenPasswordEmpty(): void
    {
        $service = new AuthService(
            new FakeClient()
        );

        $this->expectException(
            InvalidArgumentException::class
        );

        $this->expectExceptionMessage(
            'Password is required.'
        );

        $service->login(
            'john@example.com',
            ''
        );
    }

    public function testLoginCallsApi(): void
    {
        $client = new FakeClient();

        $service = new AuthService(
            $client
        );

        $response = $service->login(
            'john@example.com',
            'secret'
        );

        $this->assertSame(
            '/auth/login',
            $client->endpoint
        );

        $this->assertSame(
            [
                'email' => 'john@example.com',
                'password' => 'secret',
            ],
            $client->payload
        );

        $this->assertTrue(
            $response['success']
        );
    }

    public function testValidateTokenThrowsWhenEmpty(): void
    {
        $service = new AuthService(
            new FakeClient()
        );

        $this->expectException(
            InvalidArgumentException::class
        );

        $service->validateToken('');
    }

    public function testValidateTokenCallsApi(): void
    {
        $client = new FakeClient();

        $service = new AuthService(
            $client
        );

        $service->validateToken(
            'token-123'
        );

        $this->assertSame(
            '/auth/validate',
            $client->endpoint
        );

        $this->assertSame(
            [
                'token' => 'token-123',
            ],
            $client->payload
        );
    }

    public function testRefreshTokenThrowsWhenEmpty(): void
    {
        $service = new AuthService(
            new FakeClient()
        );

        $this->expectException(
            InvalidArgumentException::class
        );

        $service->refreshToken('');
    }

    public function testRefreshTokenCallsApi(): void
    {
        $client = new FakeClient();

        $service = new AuthService(
            $client
        );

        $service->refreshToken(
            'refresh-123'
        );

        $this->assertSame(
            '/auth/refresh',
            $client->endpoint
        );

        $this->assertSame(
            [
                'refreshToken' => 'refresh-123',
            ],
            $client->payload
        );
    }

    public function testRevokeTokenThrowsWhenEmpty(): void
    {
        $service = new AuthService(
            new FakeClient()
        );

        $this->expectException(
            InvalidArgumentException::class
        );

        $service->revokeToken('');
    }

    public function testRevokeTokenCallsApi(): void
    {
        $client = new FakeClient();

        $service = new AuthService(
            $client
        );

        $service->revokeToken(
            'token-123'
        );

        $this->assertSame(
            '/auth/revoke',
            $client->endpoint
        );

        $this->assertSame(
            [
                'token' => 'token-123',
            ],
            $client->payload
        );
    }

    public function testLogoutThrowsWhenEmpty(): void
    {
        $service = new AuthService(
            new FakeClient()
        );

        $this->expectException(
            InvalidArgumentException::class
        );

        $service->logout('');
    }

    public function testLogoutCallsApi(): void
    {
        $client = new FakeClient();

        $service = new AuthService(
            $client
        );

        $service->logout(
            'token-123'
        );

        $this->assertSame(
            '/auth/logout',
            $client->endpoint
        );

        $this->assertSame(
            [
                'token' => 'token-123',
            ],
            $client->payload
        );
    }

    public function testMeCallsApi(): void
    {
        $client = new FakeClient();

        $service = new AuthService(
            $client
        );

        $service->me();

        $this->assertSame(
            '/auth/me',
            $client->endpoint
        );
    }

    public function testGetPermissionsCallsApi(): void
    {
        $client = new FakeClient();

        $service = new AuthService(
            $client
        );

        $service->getPermissions();

        $this->assertSame(
            '/auth/permissions',
            $client->endpoint
        );
    }

    public function testGetSessionsCallsApi(): void
    {
        $client = new FakeClient();

        $service = new AuthService(
            $client
        );

        $service->getSessions();

        $this->assertSame(
            '/auth/sessions',
            $client->endpoint
        );
    }

    public function testRevokeSessionThrowsWhenEmpty(): void
    {
        $service = new AuthService(
            new FakeClient()
        );

        $this->expectException(
            InvalidArgumentException::class
        );

        $service->revokeSession('');
    }

    public function testRevokeSessionCallsApi(): void
    {
        $client = new FakeClient();

        $service = new AuthService(
            $client
        );

        $service->revokeSession(
            'session-123'
        );

        $this->assertSame(
            '/auth/sessions/revoke',
            $client->endpoint
        );

        $this->assertSame(
            [
                'sessionId' => 'session-123',
            ],
            $client->payload
        );
    }
}