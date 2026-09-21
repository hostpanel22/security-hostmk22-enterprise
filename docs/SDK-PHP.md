# HostMK Security Framework PHP SDK

## Overview

The HostMK Security PHP SDK provides a simple and secure interface for integrating HostMK Security services into PHP applications.

Supported services:

- Authentication
- Audit Logging
- Fraud Detection
- Permission Management
- Multi-Tenant Management

Requirements:

- PHP 8.2+
- ext-json

---

# Installation

Install with Composer:

```bash
composer require hostmk/security-sdk
```

Or clone the repository:

```bash
git clone https://github.com/azezaz/HostMK-Security-Framework.git
```

Generate autoload files:

```bash
composer install
```

---

# Basic Usage

```php
<?php

require_once 'vendor/autoload.php';

use HostMK\Security\Client;

$client = new Client(
    apiKey: 'YOUR_API_KEY',
    baseUrl: 'https://api.hostmk.io/v1'
);
```

---

# Authentication Service

```php
use HostMK\Security\AuthService;

$auth = new AuthService($client);
```

## Login

```php
$response = $auth->login(
    'admin@example.com',
    'password123'
);
```

## Validate Token

```php
$response = $auth->validateToken(
    $token
);
```

## Refresh Token

```php
$response = $auth->refreshToken(
    $refreshToken
);
```

## Logout

```php
$response = $auth->logout();
```

## Sessions

```php
$response = $auth->sessions();
```

---

# Audit Service

```php
use HostMK\Security\AuditService;

$audit = new AuditService($client);
```

## List Events

```php
$response = $audit->listEvents();
```

## Event Details

```php
$response = $audit->getEvent(
    'event-id'
);
```

## Search Events

```php
$response = $audit->searchEvents([
    'action' => 'login'
]);
```

## User Activity

```php
$response = $audit->userActivity(
    'user-id'
);
```

## Tenant Activity

```php
$response = $audit->tenantActivity(
    'tenant-id'
);
```

## Security Events

```php
$response = $audit->securityEvents();
```

## Statistics

```php
$response = $audit->statistics();
```

## Compliance Report

```php
$response = $audit->complianceReport(
    'gdpr'
);
```

---

# Fraud Service

```php
use HostMK\Security\FraudService;

$fraud = new FraudService($client);
```

## Analyze Request

```php
$response = $fraud->analyze([
    'email' => 'user@example.com',
    'ip' => '127.0.0.1'
]);
```

## Risk Score

```php
$response = $fraud->score([
    'email' => 'user@example.com'
]);
```

## Check Email

```php
$response = $fraud->checkEmail(
    'user@example.com'
);
```

## Check IP

```php
$response = $fraud->checkIp(
    '1.1.1.1'
);
```

## Check Device

```php
$response = $fraud->checkDevice(
    'device-id'
);
```

## Blacklist

```php
$response = $fraud->blacklist(
    'email',
    'spam@example.com'
);
```

## Whitelist

```php
$response = $fraud->whitelist(
    'email',
    'trusted@example.com'
);
```

## Alerts

```php
$response = $fraud->alerts();
```

---

# Permission Service

```php
use HostMK\Security\PermissionService;

$permissions = new PermissionService($client);
```

## List Permissions

```php
$response = $permissions->listPermissions();
```

## Create Permission

```php
$response = $permissions->createPermission(
    'users.create'
);
```

## Get Permission

```php
$response = $permissions->getPermission(
    'permission-id'
);
```

## Update Permission

```php
$response = $permissions->updatePermission(
    'permission-id',
    [
        'name' => 'users.update'
    ]
);
```

## Delete Permission

```php
$response = $permissions->deletePermission(
    'permission-id'
);
```

## Assign Permission

```php
$response = $permissions->assignPermission(
    'user-id',
    'users.create'
);
```

## Revoke Permission

```php
$response = $permissions->revokePermission(
    'user-id',
    'users.create'
);
```

## User Permissions

```php
$response = $permissions->userPermissions(
    'user-id'
);
```

## Role Permissions

```php
$response = $permissions->rolePermissions(
    'role-id'
);
```

---

# Tenancy Service

```php
use HostMK\Security\TenancyService;

$tenancy = new TenancyService($client);
```

## List Tenants

```php
$response = $tenancy->listTenants();
```

## Create Tenant

```php
$response = $tenancy->createTenant(
    'Acme Inc',
    'acme'
);
```

## Get Tenant

```php
$response = $tenancy->getTenant(
    'tenant-id'
);
```

## Update Tenant

```php
$response = $tenancy->updateTenant(
    'tenant-id',
    [
        'name' => 'Acme Corporation'
    ]
);
```

## Suspend Tenant

```php
$response = $tenancy->suspendTenant(
    'tenant-id'
);
```

## Activate Tenant

```php
$response = $tenancy->activateTenant(
    'tenant-id'
);
```

## Usage

```php
$response = $tenancy->tenantUsage(
    'tenant-id'
);
```

## Statistics

```php
$response = $tenancy->tenantStatistics(
    'tenant-id'
);
```

## Domains

```php
$response = $tenancy->tenantDomains(
    'tenant-id'
);
```

## Settings

```php
$response = $tenancy->tenantSettings(
    'tenant-id'
);
```

## Update Settings

```php
$response = $tenancy->updateSettings(
    'tenant-id',
    [
        'theme' => 'dark'
    ]
);
```

---

# Error Handling

```php
use HostMK\Security\Exceptions\ApiException;

try {
    $response = $auth->login(
        'admin@example.com',
        'password'
    );
} catch (ApiException $e) {
    echo $e->getMessage();
}
```

---

# Testing

Run PHPUnit:

```bash
vendor/bin/phpunit
```

Run coverage:

```bash
vendor/bin/phpunit --coverage-text
```

Run static analysis:

```bash
vendor/bin/phpstan analyse src tests
```

---

# Current Status

PHP SDK Version:

```text
v1.0.0
```

Validation Results:

```text
160 Tests Passed
327 Assertions Passed
PHPStan: No Errors
Production Ready
```

Coverage Summary:

```text
Methods Coverage: 92.56%
Lines Coverage: 84.98%
```

---

# Architecture

Main Components:

```text
Client
├── AuthService
├── AuditService
├── FraudService
├── PermissionService
└── TenancyService
```

Transport Layer:

```text
Contracts/
├── HttpClientInterface
├── TransportInterface
└── CurlAdapterInterface

Transport/
├── CurlTransport
└── NativeCurlAdapter
```

---

# License

MIT License

Copyright (c) HostMK Security Framework