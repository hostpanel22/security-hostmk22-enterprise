<?php

declare(strict_types=1);

namespace HostMK\Laravel\Tests;

use HostMK\Laravel\Exceptions\ApiException;
use HostMK\Laravel\HostMKManager;
use Illuminate\Http\Client\Request;
use Illuminate\Support\Facades\Http;

final class HostMKClientTest extends TestCase
{
    public function test_it_sends_api_key_and_decodes_data_envelope(): void
    {
        Http::fake([
            'security.example.com/api/v1/users' => Http::response([
                'success' => true,
                'data' => ['count' => 1, 'users' => [['id' => 'u1']]],
            ]),
        ]);

        $result = $this->app->make(HostMKManager::class)->users()->list();

        self::assertSame(1, $result['count']);
        Http::assertSent(static fn (Request $request): bool =>
            $request->url() === 'https://security.example.com/api/v1/users'
            && $request->hasHeader('X-API-Key', 'test-api-key')
        );
    }

    public function test_it_raises_typed_api_exception(): void
    {
        Http::fake([
            '*' => Http::response(['message' => 'missing permission'], 403),
        ]);

        $this->expectException(ApiException::class);
        $this->expectExceptionMessage('missing permission');

        $this->app->make(HostMKManager::class)->users()->list();
    }
}
