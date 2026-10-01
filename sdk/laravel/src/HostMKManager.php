<?php

declare(strict_types=1);

namespace HostMK\Laravel;

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

/** Facade target composing all HostMK domain services. */
final class HostMKManager
{
    private readonly AuthService $auth;
    private readonly UsersService $users;
    private readonly TenantsService $tenants;
    private readonly SessionsService $sessions;
    private readonly RolesService $roles;
    private readonly PermissionsService $permissions;
    private readonly ApiKeysService $apiKeys;
    private readonly SecretsService $secrets;
    private readonly AuditService $audit;
    private readonly FraudService $fraud;
    private readonly SettingsService $settings;
    private readonly MonitoringService $monitoring;

    public function __construct(private readonly HostMKClient $client)
    {
        $this->auth = new AuthService($client);
        $this->users = new UsersService($client);
        $this->tenants = new TenantsService($client);
        $this->sessions = new SessionsService($client);
        $this->roles = new RolesService($client);
        $this->permissions = new PermissionsService($client);
        $this->apiKeys = new ApiKeysService($client);
        $this->secrets = new SecretsService($client);
        $this->audit = new AuditService($client);
        $this->fraud = new FraudService($client);
        $this->settings = new SettingsService($client);
        $this->monitoring = new MonitoringService($client);
    }

    public function auth(): AuthService { return $this->auth; }
    public function users(): UsersService { return $this->users; }
    public function tenants(): TenantsService { return $this->tenants; }
    public function sessions(): SessionsService { return $this->sessions; }
    public function roles(): RolesService { return $this->roles; }
    public function permissions(): PermissionsService { return $this->permissions; }
    public function apiKeys(): ApiKeysService { return $this->apiKeys; }
    public function secrets(): SecretsService { return $this->secrets; }
    public function audit(): AuditService { return $this->audit; }
    public function fraud(): FraudService { return $this->fraud; }
    public function settings(): SettingsService { return $this->settings; }
    public function monitoring(): MonitoringService { return $this->monitoring; }

    public function withToken(string $token): self
    {
        $this->client->credentials()->useBearerToken($token);
        return $this;
    }

    public function withApiKey(string $apiKey): self
    {
        $this->client->credentials()->useApiKey($apiKey);
        return $this;
    }

    public function clearCredentials(): self
    {
        $this->client->credentials()->clear();
        return $this;
    }
}
