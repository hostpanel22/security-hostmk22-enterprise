<?php

declare(strict_types=1);

namespace HostMK\Security\Transport;

use HostMK\Security\Contracts\CurlAdapterInterface;
use HostMK\Security\Contracts\TransportInterface;
use JsonException;
use RuntimeException;

final class CurlTransport implements TransportInterface
{
    private bool $verifySsl;

    private int $maxRedirects;

    private CurlAdapterInterface $adapter;

    public function __construct(
        bool $verifySsl = true,
        int $maxRedirects = 5,
        ?CurlAdapterInterface $adapter = null
    ) {
        $this->verifySsl = $verifySsl;
        $this->maxRedirects = $maxRedirects;
        $this->adapter = $adapter
            ?? new NativeCurlAdapter();
    }

    public function send(
        string $method,
        string $url,
        array $headers = [],
        ?array $payload = null,
        int $timeout = 30
    ): array {
        if (trim($url) === '') {
            throw new RuntimeException(
                'Request URL is required.'
            );
        }

        if (
            !filter_var(
                $url,
                FILTER_VALIDATE_URL
            )
        ) {
            throw new RuntimeException(
                'Invalid request URL.'
            );
        }

        if ($timeout <= 0) {
            throw new RuntimeException(
                'Timeout must be greater than zero.'
            );
        }

        $response = $this->adapter->send(
            $method,
            $url,
            $this->normalizeHeaders(
                $headers
            ),
            $payload,
            $timeout
        );

        if (
            !array_key_exists(
                'status',
                $response
            )
        ) {
            throw new RuntimeException(
                'Transport response is missing status.'
            );
        }

        if (
            !array_key_exists(
                'body',
                $response
            )
        ) {
            throw new RuntimeException(
                'Transport response is missing body.'
            );
        }

        return [
            'status' => (int) $response['status'],
            'data' => $this->decodeResponse(
                (string) $response['body']
            ),
        ];
    }

    /**
     * @param array<string,string> $headers
     *
     * @return array<int,string>
     */
    private function normalizeHeaders(
        array $headers
    ): array {
        $normalized = [];

        foreach (
            $headers as $name => $value
        ) {
            $normalized[] = sprintf(
                '%s: %s',
                $name,
                $value
            );
        }

        return $normalized;
    }

    private function decodeResponse(
        string $response
    ): mixed {
        if (
            trim($response) === ''
        ) {
            return [];
        }

        try {
            return json_decode(
                $response,
                true,
                512,
                JSON_THROW_ON_ERROR
            );
        } catch (JsonException $exception) {
            throw new RuntimeException(
                'Invalid JSON response received.',
                0,
                $exception
            );
        }
    }

    public function isSslVerificationEnabled(): bool
    {
        return $this->verifySsl;
    }

    public function getMaxRedirects(): int
    {
        return $this->maxRedirects;
    }

    public function getAdapter(): CurlAdapterInterface
    {
        return $this->adapter;
    }
}