<?php

declare(strict_types=1);

namespace HostMK\Laravel\Services;

/** Tenant and platform settings. */
final class SettingsService extends AbstractService
{
    /** @return array{count?:int,settings?:list<array<string,mixed>>} */
    public function list(): array
    {
        return $this->client->get('/api/v1/settings');
    }

    /** @param array<string,mixed> $setting @return array<string,mixed> */
    public function save(array $setting): array
    {
        return $this->client->put('/api/v1/settings', $setting);
    }
}
