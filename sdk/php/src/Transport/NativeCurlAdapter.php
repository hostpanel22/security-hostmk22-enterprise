<?php

declare(strict_types=1);

namespace HostMK\Security\Transport;

use HostMK\Security\Contracts\CurlAdapterInterface;
use JsonException;
use RuntimeException;

final class NativeCurlAdapter implements CurlAdapterInterface
{
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

        $curl = curl_init();

        if ($curl === false) {
            throw new RuntimeException(
                'Unable to initialize cURL.'
            );
        }

        $options = [
            CURLOPT_URL => $url,
            CURLOPT_RETURNTRANSFER => true,
            CURLOPT_FOLLOWLOCATION => true,
            CURLOPT_TIMEOUT => $timeout,
            CURLOPT_CONNECTTIMEOUT => min(
                $timeout,
                10
            ),
            CURLOPT_CUSTOMREQUEST => strtoupper(
                $method
            ),
            CURLOPT_HTTPHEADER => $headers,
        ];

        if ($payload !== null) {
            try {
                $options[
                    CURLOPT_POSTFIELDS
                ] = json_encode(
                    $payload,
                    JSON_THROW_ON_ERROR
                );
            } catch (JsonException $exception) {
                curl_close(
                    $curl
                );

                throw new RuntimeException(
                    'Unable to encode request payload.',
                    0,
                    $exception
                );
            }
        }

        if (
            !curl_setopt_array(
                $curl,
                $options
            )
        ) {
            curl_close(
                $curl
            );

            throw new RuntimeException(
                'Unable to configure cURL request.'
            );
        }

        $response = curl_exec(
            $curl
        );

        if ($response === false) {
            $error = curl_error(
                $curl
            );

            curl_close(
                $curl
            );

            throw new RuntimeException(
                'HTTP request failed: ' .
                $error
            );
        }

        $statusCode = (int) curl_getinfo(
            $curl,
            CURLINFO_HTTP_CODE
        );

        curl_close(
            $curl
        );

        return [
            'status' => $statusCode,
            'body' => (string) $response,
        ];
    }
}