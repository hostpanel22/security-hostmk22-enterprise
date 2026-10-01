<?php

declare(strict_types=1);

namespace HostMK\Laravel\Services;

/** Managed secret lifecycle and history. */
final class SecretsService extends AbstractService
{
    /** @return array{count?:int,secrets?:list<array<string,mixed>>} */
    public function list(): array
    {
        return $this->client->get('/api/v1/secrets');
    }

    /** @param array<string,mixed> $payload @return array<string,mixed> */
    public function create(array $payload): array
    {
        return $this->client->post('/api/v1/secrets', $payload);
    }

    /** @return array<string,mixed> */
    public function rotate(string $id, string $value): array
    {
        return $this->client->post('/api/v1/secrets', [], [
            'id' => $id,
            'action' => 'rotate',
            'value' => $value,
        ]);
    }

    /** @return list<array<string,mixed>> */
    public function history(string $id): array
    {
        return $this->client->get('/api/v1/secrets', ['id' => $id, 'history' => 'true']);
    }

    /** @return array<string,mixed> */
    public function revoke(string $id): array
    {
        return $this->client->delete('/api/v1/secrets', ['id' => $id]);
    }
}
