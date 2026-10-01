<?php

declare(strict_types=1);

namespace HostMK\Laravel\Services;

/** Public health, readiness, and version operations. */
final class MonitoringService extends AbstractService
{
    /** @return array<string,mixed> */
    public function health(): array
    {
        return $this->client->get('/health');
    }

    /** @return array<string,mixed> */
    public function ready(): array
    {
        return $this->client->get('/ready');
    }

    /** @return array<string,mixed> */
    public function version(): array
    {
        return $this->client->get('/api/v1/version');
    }
}
