<?php

declare(strict_types=1);

namespace HostMK\Laravel\Services;

use HostMK\Laravel\HostMKClient;

/** Base class for domain services sharing the HostMK HTTP client. */
abstract class AbstractService
{
    public function __construct(protected readonly HostMKClient $client)
    {
    }
}
