<?php

declare(strict_types=1);

namespace HostMK\Security\Contracts;

interface HttpClientInterface
{
    public function get(
        string $endpoint,
        array $headers = []
    ): array;

    public function post(
        string $endpoint,
        array $payload = [],
        array $headers = []
    ): array;

    public function put(
        string $endpoint,
        array $payload = [],
        array $headers = []
    ): array;

    public function delete(
        string $endpoint,
        array $headers = []
    ): array;
}