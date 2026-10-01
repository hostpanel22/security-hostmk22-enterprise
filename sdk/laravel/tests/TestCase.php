<?php

declare(strict_types=1);

namespace HostMK\Laravel\Tests;

use HostMK\Laravel\HostMKServiceProvider;
use Orchestra\Testbench\TestCase as Orchestra;

abstract class TestCase extends Orchestra
{
    /** @return list<class-string> */
    protected function getPackageProviders($app): array
    {
        return [HostMKServiceProvider::class];
    }

    protected function getEnvironmentSetUp($app): void
    {
        $app['config']->set('hostmk.base_url', 'https://security.example.com');
        $app['config']->set('hostmk.api_key', 'test-api-key');
        $app['config']->set('hostmk.token', '');
        $app['config']->set('hostmk.timeout', 10);
    }
}
