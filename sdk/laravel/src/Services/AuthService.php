<?php

declare(strict_types=1);

namespace HostMK\Laravel\Services;

/** Authentication, refresh, and platform bootstrap operations. */
final class AuthService extends AbstractService
{
    /** @param array{tenant_name:string,tenant_slug:string,name:string,email:string,password:string} $payload @return array<string,mixed> */
    public function bootstrap(array $payload): array
    {
        return $this->client->post('/api/v1/auth/bootstrap', $payload);
    }

    /** @return array<string,mixed> */
    public function login(string $tenantSlug, string $email, string $password, ?string $mfaCode = null): array
    {
        $payload = [
            'tenant_slug' => $tenantSlug,
            'email' => $email,
            'password' => $password,
        ];
        if ($mfaCode !== null && trim($mfaCode) !== '') {
            $payload['mfa_code'] = $mfaCode;
        }
        $tokens = $this->client->post('/api/v1/auth/login', $payload);
        $this->adoptToken($tokens);
        return $tokens;
    }

    /** @return array<string,mixed> */
    public function refresh(string $refreshToken): array
    {
        $tokens = $this->client->post('/api/v1/auth/refresh', ['refresh_token' => $refreshToken]);
        $this->adoptToken($tokens);
        return $tokens;
    }

    public function clearCredentials(): void
    {
        $this->client->credentials()->clear();
    }

    /** @param array<string,mixed> $tokens */
    private function adoptToken(array $tokens): void
    {
        $accessToken = $tokens['access_token'] ?? '';
        if (is_string($accessToken) && trim($accessToken) !== '') {
            $this->client->credentials()->useBearerToken($accessToken);
        }
    }
}
