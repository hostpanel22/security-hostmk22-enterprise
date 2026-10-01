<?php

declare(strict_types=1);

namespace HostMK\Laravel\Facades;

use HostMK\Laravel\HostMKManager;
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
use Illuminate\Support\Facades\Facade;

/**
 * Laravel Facade for HostMK services.
 *
 * @method static AuthService auth()
 * @method static UsersService users()
 * @method static TenantsService tenants()
 * @method static SessionsService sessions()
 * @method static RolesService roles()
 * @method static PermissionsService permissions()
 * @method static ApiKeysService apiKeys()
 * @method static SecretsService secrets()
 * @method static AuditService audit()
 * @method static FraudService fraud()
 * @method static SettingsService settings()
 * @method static MonitoringService monitoring()
 * @method static HostMKManager withToken(string $token)
 * @method static HostMKManager withApiKey(string $apiKey)
 * @method static HostMKManager clearCredentials()
 */
final class HostMK extends Facade
{
    protected static function getFacadeAccessor(): string
    {
        return HostMKManager::class;
    }
}
