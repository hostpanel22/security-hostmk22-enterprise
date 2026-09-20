<?php

declare(strict_types=1);

namespace HostMK\Security\Tests\Fakes;

use HostMK\Security\Contracts\TransportInterface;

final class FakeTransport implements TransportInterface
{
    public string $method = '';

    public string $url = '';

    public array $headers = [];

    public ?array $payload = null;

    public int $timeout = 0;

    public int $callCount = 0;

    public array $response = [
        'status' => 200,
        'data' => [
            'success' => true,
        ],
    ];

    public function send(
        string $method,
        string $url,
        array $headers = [],
        ?array $payload = null,
        int $timeout = 30
    ): array {
        $this->callCount++;

        $this->method = $method;
        $this->url = $url;
        $this->headers = $headers;
        $this->payload = $payload;
        $this->timeout = $timeout;

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
        $this->url = '';
        $this->headers = [];
        $this->payload = null;
        $this->timeout = 0;
        $this->callCount = 0;

        $this->response = [
            'status' => 200,
            'data' => [
                'success' => true,
            ],
        ];
    }

    public function wasCalled(): bool
    {
        return $this->callCount > 0;
    }

    public function wasCalledWithMethod(
        string $method
    ): bool {
        return strtoupper($method)
            === $this->method;
    }

    public function wasCalledWithUrl(
        string $url
    ): bool {
        return $this->url === $url;
    }

    public function getLastMethod(): string
    {
        return $this->method;
    }

    public function getLastUrl(): string
    {
        return $this->url;
    }

    public function getLastHeaders(): array
    {
        return $this->headers;
    }

    public function getLastPayload(): ?array
    {
        return $this->payload;
    }

    public function getLastTimeout(): int
    {
        return $this->timeout;
    }

    public function getCallCount(): int
    {
        return $this->callCount;
    }
}