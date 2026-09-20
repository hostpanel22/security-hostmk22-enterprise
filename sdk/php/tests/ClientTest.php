<?php

declare(strict_types=1);

namespace HostMK\Security\Tests;

use HostMK\Security\Client;
use HostMK\Security\Tests\Fakes\FakeTransport;
use PHPUnit\Framework\TestCase;
use RuntimeException;

final class ClientTest extends TestCase
{
    public function testCanCreateClient(): void
    {
        $client = new Client(
            'test-api-key',
            'https://api.hostmk.com'
        );

        $this->assertInstanceOf(
            Client::class,
            $client
        );
    }

    public function testThrowsWhenApiKeyIsEmpty(): void
    {
        $this->expectException(
            RuntimeException::class
        );

        $this->expectExceptionMessage(
            'HostMK API key is required.'
        );

        new Client(
            '',
            'https://api.hostmk.com'
        );
    }

    public function testThrowsWhenBaseUrlIsEmpty(): void
    {
        $this->expectException(
            RuntimeException::class
        );

        $this->expectExceptionMessage(
            'HostMK base URL is required.'
        );

        new Client(
            'test-key',
            ''
        );
    }

    public function testThrowsWhenBaseUrlInvalid(): void
    {
        $this->expectException(
            RuntimeException::class
        );

        $this->expectExceptionMessage(
            'Invalid HostMK base URL.'
        );

        new Client(
            'test-key',
            'invalid-url'
        );
    }

    public function testThrowsWhenTimeoutInvalid(): void
    {
        $this->expectException(
            RuntimeException::class
        );

        $this->expectExceptionMessage(
            'Timeout must be greater than zero.'
        );

        new Client(
            'test-key',
            'https://api.hostmk.com',
            0
        );
    }

    public function testThrowsWhenUserAgentEmpty(): void
    {
        $this->expectException(
            RuntimeException::class
        );

        $this->expectExceptionMessage(
            'User agent is required.'
        );

        new Client(
            'test-key',
            'https://api.hostmk.com',
            30,
            ''
        );
    }

    public function testReturnsApiKey(): void
    {
        $client = new Client(
            'secret-key',
            'https://api.hostmk.com'
        );

        $this->assertSame(
            'secret-key',
            $client->getApiKey()
        );
    }

    public function testReturnsBaseUrl(): void
    {
        $client = new Client(
            'secret-key',
            'https://api.hostmk.com/'
        );

        $this->assertSame(
            'https://api.hostmk.com',
            $client->getBaseUrl()
        );
    }

    public function testReturnsTimeout(): void
    {
        $client = new Client(
            'secret-key',
            'https://api.hostmk.com',
            60
        );

        $this->assertSame(
            60,
            $client->getTimeout()
        );
    }

    public function testReturnsUserAgent(): void
    {
        $client = new Client(
            'secret-key',
            'https://api.hostmk.com',
            30,
            'HostMK-Test/1.0'
        );

        $this->assertSame(
            'HostMK-Test/1.0',
            $client->getUserAgent()
        );
    }

    public function testTrimsApiKey(): void
    {
        $client = new Client(
            '  secret-key  ',
            'https://api.hostmk.com'
        );

        $this->assertSame(
            'secret-key',
            $client->getApiKey()
        );
    }

    public function testGetUsesTransport(): void
    {
        $transport = new FakeTransport();

        $client = new Client(
            'api-key',
            'https://api.hostmk.com',
            30,
            'HostMK-Test',
            $transport
        );

        $response = $client->get(
            '/users'
        );

        $this->assertTrue(
            $transport->wasCalled()
        );

        $this->assertSame(
            'GET',
            $transport->getLastMethod()
        );

        $this->assertSame(
            'https://api.hostmk.com/users',
            $transport->getLastUrl()
        );

        $this->assertSame(
            200,
            $response['status']
        );
    }

    public function testPostUsesTransport(): void
    {
        $transport = new FakeTransport();

        $client = new Client(
            'api-key',
            'https://api.hostmk.com',
            30,
            'HostMK-Test',
            $transport
        );

        $payload = [
            'email' => 'test@example.com',
        ];

        $client->post(
            '/users',
            $payload
        );

        $this->assertSame(
            'POST',
            $transport->getLastMethod()
        );

        $this->assertSame(
            $payload,
            $transport->getLastPayload()
        );
    }

    public function testPutUsesTransport(): void
    {
        $transport = new FakeTransport();

        $client = new Client(
            'api-key',
            'https://api.hostmk.com',
            30,
            'HostMK-Test',
            $transport
        );

        $payload = [
            'name' => 'John'
        ];

        $client->put(
            '/users/1',
            $payload
        );

        $this->assertSame(
            'PUT',
            $transport->getLastMethod()
        );

        $this->assertSame(
            $payload,
            $transport->getLastPayload()
        );
    }

    public function testDeleteUsesTransport(): void
    {
        $transport = new FakeTransport();

        $client = new Client(
            'api-key',
            'https://api.hostmk.com',
            30,
            'HostMK-Test',
            $transport
        );

        $client->delete(
            '/users/1'
        );

        $this->assertSame(
            'DELETE',
            $transport->getLastMethod()
        );
    }

    public function testHeadersAreGenerated(): void
    {
        $transport = new FakeTransport();

        $client = new Client(
            'api-key',
            'https://api.hostmk.com',
            30,
            'HostMK-Agent',
            $transport
        );

        $client->get('/users');

        $headers = $transport->getLastHeaders();

        $this->assertSame(
            'application/json',
            $headers['Accept']
        );

        $this->assertSame(
            'application/json',
            $headers['Content-Type']
        );

        $this->assertSame(
            'api-key',
            $headers['X-API-Key']
        );

        $this->assertSame(
            'HostMK-Agent',
            $headers['User-Agent']
        );
    }

    public function testCustomHeadersOverrideDefaults(): void
    {
        $transport = new FakeTransport();

        $client = new Client(
            'api-key',
            'https://api.hostmk.com',
            30,
            'HostMK-Agent',
            $transport
        );

        $client->get(
            '/users',
            [
                'X-Tenant' => 'tenant-1',
            ]
        );

        $headers = $transport->getLastHeaders();

        $this->assertSame(
            'tenant-1',
            $headers['X-Tenant']
        );
    }

    public function testTimeoutPassedToTransport(): void
    {
        $transport = new FakeTransport();

        $client = new Client(
            'api-key',
            'https://api.hostmk.com',
            120,
            'HostMK-Agent',
            $transport
        );

        $client->get('/users');

        $this->assertSame(
            120,
            $transport->getLastTimeout()
        );
    }

    public function testReturnsInjectedTransport(): void
    {
        $transport = new FakeTransport();

        $client = new Client(
            'api-key',
            'https://api.hostmk.com',
            30,
            'HostMK-Agent',
            $transport
        );

        $this->assertSame(
            $transport,
            $client->getTransport()
        );
    }
}