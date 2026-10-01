<?php

declare(strict_types=1);

namespace HostMK\Laravel\Tests;

use HostMK\Laravel\HostMKManager;
use Illuminate\Http\Client\Request;
use Illuminate\Support\Facades\Http;

final class AuthServiceTest extends TestCase
{
    public function test_login_switches_shared_services_to_bearer_authentication(): void
    {
        Http::fake([
            'security.example.com/api/v1/auth/login' => Http::response([
                'success' => true,
                'data' => [
                    'access_token' => 'access-2',
                    'refresh_token' => 'refresh-2',
                    'token_type' => 'Bearer',
                    'expires_in' => 900,
                ],
            ]),
            'security.example.com/api/v1/users' => Http::response([
                'success' => true,
                'data' => ['count' => 0, 'users' => []],
            ]),
        ]);

        $manager = $this->app->make(HostMKManager::class);
        $manager->auth()->login('acme', 'owner@example.com', 'secret');
        $manager->users()->list();

        Http::assertSent(static fn (Request $request): bool =>
            str_ends_with($request->url(), '/api/v1/users')
            && $request->hasHeader('Authorization', 'Bearer access-2')
            && !$request->hasHeader('X-API-Key')
        );
    }
}
