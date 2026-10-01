<?php

declare(strict_types=1);

namespace HostMK\Laravel\Services;

/** RBAC role CRUD operations. */
final class RolesService extends AbstractService
{
    /** @return array{count?:int,roles?:list<array<string,mixed>>} */
    public function list(): array
    {
        return $this->client->get('/api/v1/rbac/roles');
    }

    /** @param array<string,mixed> $role @return array<string,mixed> */
    public function create(array $role): array
    {
        return $this->client->post('/api/v1/rbac/roles', $role);
    }

    /** @param array<string,mixed> $role @return array<string,mixed> */
    public function update(array $role): array
    {
        return $this->client->put('/api/v1/rbac/roles', $role);
    }

    /** @return array<string,mixed> */
    public function delete(string $id): array
    {
        return $this->client->delete('/api/v1/rbac/roles', ['id' => $id]);
    }
}
