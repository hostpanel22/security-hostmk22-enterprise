<?php

declare(strict_types=1);

namespace HostMK\Laravel\Services;

/** User lifecycle and MFA operations. */
final class UsersService extends AbstractService
{
    /** @return array{count?:int,users?:list<array<string,mixed>>} */
    public function list(): array
    {
        return $this->client->get('/api/v1/users');
    }

    /** @param array<string,mixed> $payload @return array<string,mixed> */
    public function create(array $payload): array
    {
        return $this->client->post('/api/v1/users', $payload);
    }

    /** @param array<string,mixed> $payload @return array<string,mixed> */
    public function update(string $id, array $payload): array
    {
        return $this->client->post('/api/v1/users/action', $payload, ['id' => $id]);
    }

    /** @return array<string,mixed> */
    public function action(string $id, string $action): array
    {
        return $this->client->post('/api/v1/users/action', [], ['id' => $id, 'action' => $action]);
    }
}
