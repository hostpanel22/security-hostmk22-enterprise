<?php

declare(strict_types=1);

namespace HostMK\Laravel\Tests;

use HostMK\Laravel\HostMKManager;
use HostMK\Laravel\Http\Middleware\UseHostMKBearerToken;
use Illuminate\Http\Client\Request as ClientRequest;
use Illuminate\Http\Request;
use Illuminate\Support\Facades\Http;

final class MiddlewareTest extends TestCase
{
    public function test_middleware_propagates_incoming_bearer_token(): void
    {
        Http::fake(['*' => Http::response(['success' => true, 'data' => ['count' => 0, 'users' => []]])]);
        $manager = $this->app->make(HostMKManager::class);
        $middleware = new UseHostMKBearerToken($manager);
        $request = Request::create('/example', 'GET', server: ['HTTP_AUTHORIZATION' => 'Bearer request-token']);

        $middleware->handle($request, static function () use ($manager): string {
            $manager->users()->list();
            return 'ok';
        });

        Http::assertSent(static fn (ClientRequest $request): bool =>
            $request->hasHeader('Authorization', 'Bearer request-token')
            && !$request->hasHeader('X-API-Key')
        );
    }
}
