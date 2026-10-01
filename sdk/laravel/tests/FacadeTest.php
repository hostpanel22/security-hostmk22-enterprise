<?php

declare(strict_types=1);

namespace HostMK\Laravel\Tests;

use HostMK\Laravel\Facades\HostMK;
use HostMK\Laravel\Services\ApiKeysService;
use HostMK\Laravel\Services\AuditService;
use HostMK\Laravel\Services\AuthService;
use HostMK\Laravel\Services\FraudService;
use HostMK\Laravel\Services\MonitoringService;
use HostMK\Laravel\Services\PermissionsService;
use HostMK\Laravel\Services\RolesService;
use HostMK\Laravel\Services\SecretsService;
use HostMK\Laravel\Services\SessionsService;
use HostMK\Laravel\Services\SettingsService;
use HostMK\Laravel\Services\TenantsService;
use HostMK\Laravel\Services\UsersService;

final class FacadeTest extends TestCase
{
    public function test_facade_exposes_domain_services(): void
    {
        self::assertInstanceOf(AuthService::class, HostMK::auth());
        self::assertInstanceOf(UsersService::class, HostMK::users());
        self::assertInstanceOf(TenantsService::class, HostMK::tenants());
        self::assertInstanceOf(SessionsService::class, HostMK::sessions());
        self::assertInstanceOf(RolesService::class, HostMK::roles());
        self::assertInstanceOf(PermissionsService::class, HostMK::permissions());
        self::assertInstanceOf(ApiKeysService::class, HostMK::apiKeys());
        self::assertInstanceOf(SecretsService::class, HostMK::secrets());
        self::assertInstanceOf(AuditService::class, HostMK::audit());
        self::assertInstanceOf(FraudService::class, HostMK::fraud());
        self::assertInstanceOf(SettingsService::class, HostMK::settings());
        self::assertInstanceOf(MonitoringService::class, HostMK::monitoring());
    }

    public function test_configuration_can_be_published(): void
    {
        $this->artisan('vendor:publish', ['--tag' => 'hostmk-config'])
            ->assertExitCode(0);

        self::assertFileExists($this->app->configPath('hostmk.php'));
    }
}
