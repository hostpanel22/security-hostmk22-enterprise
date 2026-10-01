<?php

declare(strict_types=1);

namespace HostMK\Laravel\Services;

/** Permission CRUD and access checks. */
final class PermissionsService extends AbstractService
{
    /** @return array{count?:int,permissions?:list<array<string,mixed>>} */
    public function list(): array
    {
        return $this->client->get('/api/v1/permissions');
    }

    /** @param array<string,mixed> $permission @return array<string,mixed> */
    public function create(array $permission): array
    {
        return $this->client->post('/api/v1/permissions', $permission);
    }

    /** @return array<string,mixed> */
    public function delete(string $id): array
    {
        return $this->client->delete('/api/v1/permissions', ['id' => $id]);
    }

    /** @return array{user_id?:string,permission?:string,allowed?:bool} */
    public function check(string $userId, string $permission): array
    {
        return $this->client->post('/api/v1/permissions/check', [
            'user_id' => $userId,
            'permission' => $permission,
        ]);
    }
}
