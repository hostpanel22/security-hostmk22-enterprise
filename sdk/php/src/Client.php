<?php

declare(strict_types=1);

namespace HostMK\Security;

use HostMK\Security\Contracts\HttpClientInterface;
use HostMK\Security\Contracts\TransportInterface;
use HostMK\Security\Transport\CurlTransport;
use RuntimeException;

final class Client implements HttpClientInterface
{
    private string $apiKey;

    private string $baseUrl;

    private int $timeout;

    private string $userAgent;

    private TransportInterface $transport;

    public function __construct(
        string $apiKey,
        string $baseUrl,
        int $timeout = 30,
        string $userAgent = 'HostMK-PHP-SDK/1.0.0',
        ?TransportInterface $transport = null
    ) {
        $this->apiKey = trim(
            $apiKey
        );

        $this->baseUrl = rtrim(
            trim($baseUrl),
            '/'
        );

        $this->timeout = $timeout;

        $this->userAgent = trim(
            $userAgent
        );

        $this->transport = $transport
            ?? new CurlTransport();

        $this->validateConfiguration();
    }

    public function get(
        string $endpoint,
        array $headers = []
    ): array {
        return $this->request(
            'GET',
            $endpoint,
            null,
            $headers
        );
    }

    public function post(
        string $endpoint,
        array $payload = [],
        array $headers = []
    ): array {
        return $this->request(
            'POST',
            $endpoint,
            $payload,
            $headers
        );
    }

    public function put(
        string $endpoint,
        array $payload = [],
        array $headers = []
    ): array {
        return $this->request(
            'PUT',
            $endpoint,
            $payload,
            $headers
        );
    }

    public function delete(
        string $endpoint,
        array $headers = []
    ): array {
        return $this->request(
            'DELETE',
            $endpoint,
            null,
            $headers
        );
    }

    private function request(
        string $method,
        string $endpoint,
        ?array $payload,
        array $headers
    ): array {
        return $this->transport->send(
            strtoupper($method),
            $this->buildUrl(
                $endpoint
            ),
            $this->buildHeaders(
                $headers
            ),
            $payload,
            $this->timeout
        );
    }

    private function buildHeaders(
        array $headers = []
    ): array {
        return array_merge(
            [
                'Accept' => 'application/json',
                'Content-Type' => 'application/json',
                'X-API-Key' => $this->apiKey,
                'User-Agent' => $this->userAgent,
            ],
            $headers
        );
    }

    private function buildUrl(
        string $endpoint
    ): string {
        return $this->baseUrl .
            '/' .
            ltrim(
                $endpoint,
                '/'
            );
    }

    private function validateConfiguration(): void
    {
        if ($this->apiKey === '') {
            throw new RuntimeException(
                'HostMK API key is required.'
            );
        }

        if ($this->baseUrl === '') {
            throw new RuntimeException(
                'HostMK base URL is required.'
            );
        }

        if (
            !filter_var(
                $this->baseUrl,
                FILTER_VALIDATE_URL
            )
        ) {
            throw new RuntimeException(
                'Invalid HostMK base URL.'
            );
        }

        if ($this->timeout <= 0) {
            throw new RuntimeException(
                'Timeout must be greater than zero.'
            );
        }

        if ($this->userAgent === '') {
            throw new RuntimeException(
                'User agent is required.'
            );
        }
    }

    public function getApiKey(): string
    {
        return $this->apiKey;
    }

    public function getBaseUrl(): string
    {
        return $this->baseUrl;
    }

    public function getTimeout(): int
    {
        return $this->timeout;
    }

    public function getUserAgent(): string
    {
        return $this->userAgent;
    }

    public function getTransport(): TransportInterface
    {
        return $this->transport;
    }
}