<?php

declare(strict_types=1);

namespace HostMK\Laravel;

use HostMK\Laravel\Auth\CredentialStore;
use Illuminate\Contracts\Foundation\Application;
use Illuminate\Http\Client\Factory;
use Illuminate\Support\ServiceProvider;

/** Registers request-scoped HostMK SDK services and publishable configuration. */
final class HostMKServiceProvider extends ServiceProvider
{
    public function register(): void
    {
        $this->mergeConfigFrom(__DIR__ . '/../config/hostmk.php', 'hostmk');

        $this->app->scoped(CredentialStore::class, static function (Application $app): CredentialStore {
            return new CredentialStore(
                (string) $app['config']->get('hostmk.api_key', ''),
                (string) $app['config']->get('hostmk.token', '')
            );
        });

        $this->app->scoped(HostMKClient::class, static function (Application $app): HostMKClient {
            return new HostMKClient(
                $app->make(Factory::class),
                $app->make(CredentialStore::class),
                (string) $app['config']->get('hostmk.base_url'),
                (int) $app['config']->get('hostmk.timeout', 30),
                (string) $app['config']->get('hostmk.user_agent', 'hostmk-laravel-sdk/1.0')
            );
        });

        $this->app->scoped(HostMKManager::class, static fn (Application $app): HostMKManager =>
            new HostMKManager($app->make(HostMKClient::class))
        );
        $this->app->alias(HostMKManager::class, 'hostmk');
    }

    public function boot(): void
    {
        $this->publishes([
            __DIR__ . '/../config/hostmk.php' => $this->app->configPath('hostmk.php'),
        ], 'hostmk-config');
    }
}
