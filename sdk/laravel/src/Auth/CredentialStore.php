<?php

declare(strict_types=1);

namespace HostMK\Laravel\Auth;

/** Request-scoped credential state shared by every service. */
final class CredentialStore
{
    public function __construct(
        private string $apiKey = '',
        private string $bearerToken = ''
    ) {
        $this->apiKey = trim($this->apiKey);
        $this->bearerToken = trim($this->bearerToken);
    }

    public function apiKey(): string
    {
        return $this->apiKey;
    }

    public function bearerToken(): string
    {
        return $this->bearerToken;
    }

    public function useBearerToken(string $token): void
    {
        $this->apiKey = '';
        $this->bearerToken = trim($token);
    }

    public function useApiKey(string $apiKey): void
    {
        $this->apiKey = trim($apiKey);
        $this->bearerToken = '';
    }

    public function clear(): void
    {
        $this->apiKey = '';
        $this->bearerToken = '';
    }
}
