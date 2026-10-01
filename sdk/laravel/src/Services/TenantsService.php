<?php

declare(strict_types=1);

namespace HostMK\Laravel\Services;

/** Tenant creation, update, and status operations. */
final class TenantsService extends AbstractService
{
    /** @return array{count?:int,tenants?:list<array<string,mixed>>} */
    public function list(): array
    {
        return $this->client->get('/api/v1/tenants');
    }

    /** @param array<string,mixed> $tenant @return array<string,mixed> */
    public function create(array $tenant): array
    {
        return $this->client->post('/api/v1/tenants', $tenant);
    }

    /** @param array<string,mixed> $tenant @return array<string,mixed> */
    public function update(string $id, array $tenant): array
    {
        $tenant['id'] = $id;
        return $this->client->post('/api/v1/tenants/action', $tenant, ['id' => $id]);
    }

    /** @return array<string,mixed> */
    public function activate(string $id): array
    {
        return $this->setStatus($id, 'activate');
    }

    /** @return array<string,mixed> */
    public function suspend(string $id): array
    {
        return $this->setStatus($id, 'suspend');
    }

    /** @return array<string,mixed> */
    private function setStatus(string $id, string $action): array
    {
        return $this->client->post('/api/v1/tenants/action', [], ['id' => $id, 'action' => $action]);
    }
}
