<?php

declare(strict_types=1);

namespace HostMK\Laravel;

use HostMK\Laravel\Auth\CredentialStore;
use HostMK\Laravel\Exceptions\ApiException;
use HostMK\Laravel\Exceptions\ConfigurationException;
use Illuminate\Http\Client\Factory;

/** Shared Laravel HTTP client for all HostMK modules. */
final class HostMKClient
{
    public function __construct(
        private readonly Factory $http,
        private readonly CredentialStore $credentials,
        private readonly string $baseUrl,
        private readonly int $timeout = 30,
        private readonly string $userAgent = 'hostmk-laravel-sdk/1.0'
    ) {
        $this->validateConfiguration();
    }

    public function credentials(): CredentialStore
    {
        return $this->credentials;
    }

    /** @param array<string,mixed> $query @return array<string,mixed>|list<mixed> */
    public function get(string $path, array $query = []): array
    {
        return $this->request('GET', $path, $query);
    }

    /** @param array<string,mixed> $payload @param array<string,mixed> $query @return array<string,mixed>|list<mixed> */
    public function post(string $path, array $payload = [], array $query = []): array
    {
        return $this->request('POST', $path, $query, $payload);
    }

    /** @param array<string,mixed> $payload @param array<string,mixed> $query @return array<string,mixed>|list<mixed> */
    public function put(string $path, array $payload = [], array $query = []): array
    {
        return $this->request('PUT', $path, $query, $payload);
    }

    /** @param array<string,mixed> $query @return array<string,mixed>|list<mixed> */
    public function delete(string $path, array $query = []): array
    {
        return $this->request('DELETE', $path, $query);
    }

    /**
     * @param array<string,mixed> $query
     * @param array<string,mixed>|null $payload
     * @return array<string,mixed>|list<mixed>
     */
    private function request(string $method, string $path, array $query, ?array $payload = null): array
    {
        $headers = [
            'Accept' => 'application/json',
            'User-Agent' => $this->userAgent,
        ];

        if ($this->credentials->apiKey() !== '') {
            $headers['X-API-Key'] = $this->credentials->apiKey();
        } elseif ($this->credentials->bearerToken() !== '') {
            $headers['Authorization'] = 'Bearer ' . $this->credentials->bearerToken();
        }

        $options = ['query' => $query];
        if ($payload !== null) {
            $options['json'] = $payload;
        }

        $response = $this->http
            ->withHeaders($headers)
            ->timeout($this->timeout)
            ->send($method, $this->buildUrl($path), $options);

        $decoded = $response->json();
        $decoded = is_array($decoded) ? $decoded : [];

        if ($response->failed()) {
            $message = $decoded['message'] ?? $decoded['error'] ?? ('HostMK API returned HTTP ' . $response->status());
            throw new ApiException((string) $message, $response->status(), $decoded);
        }

        $data = array_key_exists('data', $decoded) ? $decoded['data'] : $decoded;
        return is_array($data) ? $data : [];
    }

    private function buildUrl(string $path): string
    {
        if (!str_starts_with($path, '/') || str_starts_with($path, '//')) {
            throw new ConfigurationException('API path must be absolute.');
        }

        $base = rtrim($this->baseUrl, '/');
        if (str_ends_with($base, '/api/v1') && ($path === '/api/v1' || str_starts_with($path, '/api/v1/'))) {
            $path = substr($path, strlen('/api/v1'));
        }

        return $base . $path;
    }

    private function validateConfiguration(): void
    {
        if (filter_var($this->baseUrl, FILTER_VALIDATE_URL) === false) {
            throw new ConfigurationException('HOSTMK_BASE_URL must be a valid URL.');
        }
        if ($this->timeout <= 0) {
            throw new ConfigurationException('HOSTMK_TIMEOUT must be greater than zero.');
        }
        if (trim($this->userAgent) === '') {
            throw new ConfigurationException('HOSTMK_USER_AGENT is required.');
        }
    }
}
