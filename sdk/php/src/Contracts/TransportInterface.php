<?php

declare(strict_types=1);

namespace HostMK\Security\Contracts;

interface TransportInterface
{
    /**
     * Execute an HTTP request.
     *
     * @param string $method HTTP method.
     * @param string $url Absolute request URL.
     * @param array<string, string> $headers HTTP headers.
     * @param array<mixed>|null $payload Request payload.
     * @param int $timeout Timeout in seconds.
     *
     * @return array{
     *     status: int,
     *     body: string
     * }
     */
    public function send(
        string $method,
        string $url,
        array $headers = [],
        ?array $payload = null,
        int $timeout = 30
    ): array;
}