<?php

declare(strict_types=1);

namespace HostMK\Laravel\Http\Middleware;

use Closure;
use HostMK\Laravel\HostMKManager;
use Illuminate\Http\Request;

/** Propagates an incoming bearer token to the request-scoped HostMK client. */
final class UseHostMKBearerToken
{
    public function __construct(private readonly HostMKManager $hostmk)
    {
    }

    public function handle(Request $request, Closure $next): mixed
    {
        $token = $request->bearerToken();
        if (is_string($token) && trim($token) !== '') {
            $this->hostmk->withToken($token);
        }

        return $next($request);
    }
}
