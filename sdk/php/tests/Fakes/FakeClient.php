<?php

declare(strict_types=1);

namespace HostMK\Security\Tests\Fakes;

use HostMK\Security\Contracts\HttpClientInterface;

final class FakeClient implements HttpClientInterface
{
    public string $method = '';

    public string $endpoint = '';

    public array $payload = [];

    public array $headers = [];

    public array $response = [
        'success' => true,
    ];

    public function get(
        string $endpoint,
        array $headers = []
    ): array {
        $this->method = 'GET';
        $this->endpoint = $endpoint;
        $this->headers = $headers;

        return $this->response;
    }

    public function post(
        string $endpoint,
        array $payload = [],
        array $headers = []
    ): array {
        $this->method = 'POST';
        $this->endpoint = $endpoint;
        $this->payload = $payload;
        $this->headers = $headers;

        return $this->response;
    }

    public function put(
        string $endpoint,
        array $payload = [],
        array $headers = []
    ): array {
        $this->method = 'PUT';
        $this->endpoint = $endpoint;
        $this->payload = $payload;
        $this->headers = $headers;

        return $this->response;
    }

    public function delete(
        string $endpoint,
        array $headers = []
    ): array {
        $this->method = 'DELETE';
        $this->endpoint = $endpoint;
        $this->headers = $headers;

        return $this->response;
    }

    public function setResponse(
        array $response
    ): self {
        $this->response = $response;

        return $this;
    }

    public function reset(): void
    {
        $this->method = '';
        $this->endpoint = '';
        $this->payload = [];
        $this->headers = [];

        $this->response = [
            'success' => true,
        ];
    }

    public function wasCalled(): bool
    {
        return $this->method !== '';
    }

    public function wasCalledWithMethod(
        string $method
    ): bool {
        return strtoupper($method)
            === $this->method;
    }

    public function wasCalledWithEndpoint(
        string $endpoint
    ): bool {
        return $this->endpoint
            === $endpoint;
    }

    public function getLastPayload(): array
    {
        return $this->payload;
    }

    public function getLastHeaders(): array
    {
        return $this->headers;
    }
}