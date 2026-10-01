<?php

declare(strict_types=1);

namespace HostMK\Laravel\Tests;

use HostMK\Laravel\HostMKManager;
use Illuminate\Http\Client\Request;
use Illuminate\Support\Facades\Http;

final class ServiceRoutesTest extends TestCase
{
    public function test_all_registered_domains_use_canonical_paths(): void
    {
        Http::fake(['*' => Http::response(['success' => true, 'data' => []])]);
        $hostmk = $this->app->make(HostMKManager::class);

        $hostmk->auth()->bootstrap([
            'tenant_name' => 'Acme',
            'tenant_slug' => 'acme',
            'name' => 'Owner',
            'email' => 'owner@example.com',
            'password' => 'secret',
        ]);
        $hostmk->auth()->login('acme', 'owner@example.com', 'secret');
        $hostmk->auth()->refresh('refresh-token');
        $hostmk->users()->list();
        $hostmk->users()->create(['email' => 'user@example.com']);
        $hostmk->users()->update('user-1', ['name' => 'Updated']);
        $hostmk->users()->action('user-1', 'disable');
        $hostmk->tenants()->list();
        $hostmk->tenants()->create(['name' => 'Acme', 'slug' => 'acme']);
        $hostmk->tenants()->update('tenant-1', ['name' => 'Updated']);
        $hostmk->tenants()->activate('tenant-1');
        $hostmk->tenants()->suspend('tenant-1');
        $hostmk->sessions()->list();
        $hostmk->sessions()->revoke('session-1');
        $hostmk->roles()->list();
        $hostmk->roles()->create(['name' => 'operator']);
        $hostmk->roles()->update(['id' => 'role-1', 'name' => 'operator']);
        $hostmk->roles()->delete('role-1');
        $hostmk->permissions()->list();
        $hostmk->permissions()->create(['name' => 'users.read']);
        $hostmk->permissions()->delete('permission-1');
        $hostmk->permissions()->check('user-1', 'users.read');
        $hostmk->apiKeys()->list();
        $hostmk->apiKeys()->create(['name' => 'ci']);
        $hostmk->apiKeys()->rotate('key-1');
        $hostmk->apiKeys()->revoke('key-1');
        $hostmk->secrets()->list();
        $hostmk->secrets()->create(['name' => 'DB', 'value' => 'secret']);
        $hostmk->secrets()->rotate('secret-1', 'next');
        $hostmk->secrets()->history('secret-1');
        $hostmk->secrets()->revoke('secret-1');
        $hostmk->audit()->list();
        $hostmk->audit()->create(['event_type' => 'login']);
        $hostmk->fraud()->listRules();
        $hostmk->fraud()->createRule(['name' => 'rule']);
        $hostmk->fraud()->updateRule(['id' => 'rule-1', 'name' => 'rule']);
        $hostmk->fraud()->deleteRule('rule-1');
        $hostmk->fraud()->analyze(['amount' => 42]);
        $hostmk->fraud()->listViolations();
        $hostmk->settings()->list();
        $hostmk->settings()->save(['key' => 'limit', 'value' => 10]);
        $hostmk->monitoring()->health();
        $hostmk->monitoring()->ready();
        $hostmk->monitoring()->version();

        $expected = [
            'POST /api/v1/auth/bootstrap', 'POST /api/v1/auth/login', 'POST /api/v1/auth/refresh',
            'GET /api/v1/users', 'POST /api/v1/users', 'POST /api/v1/users/action',
            'GET /api/v1/tenants', 'POST /api/v1/tenants', 'POST /api/v1/tenants/action',
            'GET /api/v1/sessions', 'POST /api/v1/sessions',
            'GET /api/v1/rbac/roles', 'POST /api/v1/rbac/roles',
            'PUT /api/v1/rbac/roles', 'DELETE /api/v1/rbac/roles',
            'GET /api/v1/permissions', 'POST /api/v1/permissions', 'DELETE /api/v1/permissions',
            'POST /api/v1/permissions/check', 'GET /api/v1/api-keys', 'POST /api/v1/api-keys',
            'DELETE /api/v1/api-keys', 'GET /api/v1/secrets', 'POST /api/v1/secrets',
            'DELETE /api/v1/secrets', 'GET /api/v1/audit/events', 'POST /api/v1/audit/events',
            'GET /api/v1/fraud/rules', 'POST /api/v1/fraud/rules', 'PUT /api/v1/fraud/rules',
            'DELETE /api/v1/fraud/rules', 'POST /api/v1/fraud/analyze',
            'GET /api/v1/fraud/violations', 'GET /api/v1/settings', 'PUT /api/v1/settings',
            'GET /health', 'GET /ready', 'GET /api/v1/version',
        ];

        $sent = [];
        Http::assertSent(function (Request $request) use (&$sent): bool {
            $sent[] = $request->method() . ' ' . (string) parse_url($request->url(), PHP_URL_PATH);
            return true;
        });
        foreach ($expected as $endpoint) {
            self::assertContains($endpoint, $sent);
        }
    }
}
