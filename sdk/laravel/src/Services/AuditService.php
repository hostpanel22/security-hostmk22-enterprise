<?php

declare(strict_types=1);

namespace HostMK\Laravel\Services;

/** Audit event listing and creation. */
final class AuditService extends AbstractService
{
    /** @return array{count?:int,events?:list<array<string,mixed>>} */
    public function list(?int $limit = null): array
    {
        return $this->client->get('/api/v1/audit/events', $limit !== null ? ['limit' => $limit] : []);
    }

    /** @param array<string,mixed> $event @return array<string,mixed> */
    public function create(array $event): array
    {
        return $this->client->post('/api/v1/audit/events', $event);
    }
}
