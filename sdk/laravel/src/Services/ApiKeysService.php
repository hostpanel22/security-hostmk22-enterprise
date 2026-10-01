<?php

declare(strict_types=1);

namespace HostMK\Laravel\Services;

/** API key lifecycle operations. */
final class ApiKeysService extends AbstractService
{
    /** @return array{count?:int,api_keys?:list<array<string,mixed>>} */
    public function list(): array
    {
        return $this->client->get('/api/v1/api-keys');
    }

    /** @param array<string,mixed> $payload @return array<string,mixed> */
    public function create(array $payload): array
    {
        return $this->client->post('/api/v1/api-keys', $payload);
    }

    /** @return array<string,mixed> */
    public function rotate(string $id): array
    {
        return $this->client->post('/api/v1/api-keys', [], ['id' => $id, 'action' => 'rotate']);
    }

    /** @return array<string,mixed> */
    public function revoke(string $id): array
    {
        return $this->client->delete('/api/v1/api-keys', ['id' => $id]);
    }
}
