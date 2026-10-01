<?php

declare(strict_types=1);

namespace HostMK\Laravel\Services;

/** Refresh-token session listing and revocation. */
final class SessionsService extends AbstractService
{
    /** @return array{count?:int,sessions?:list<array<string,mixed>>} */
    public function list(): array
    {
        return $this->client->get('/api/v1/sessions');
    }

    /** @return array<string,mixed> */
    public function revoke(string $id): array
    {
        return $this->client->post('/api/v1/sessions', [], ['id' => $id, 'action' => 'revoke']);
    }
}
