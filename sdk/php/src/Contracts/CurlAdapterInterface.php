<?php

declare(strict_types=1);

namespace HostMK\Security\Contracts;

/**
 * Low-level cURL adapter contract.
 *
 * This interface isolates native cURL operations from the
 * transport layer to make HTTP requests fully testable.
 *
 * Production:
 *   - NativeCurlAdapter
 *
 * Tests:
 *   - FakeCurlAdapter
 */
interface CurlAdapterInterface
{
    /**
     * Execute an HTTP request.
     *
     * Example response:
     *
     * [
     *     'status' => 200,
     *     'body' => '{"success":true}'
     * ]
     *
     * @param string $method HTTP method
     * @param string $url Absolute request URL
     * @param array<string,string> $headers HTTP headers
     * @param array<mixed>|null $payload Request payload
     * @param int $timeout Timeout in seconds
     *
     * @return array{
     *     status:int,
     *     body:string
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