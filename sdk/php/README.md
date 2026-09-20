# HostMK Security SDK for PHP

Official PHP SDK for the HostMK Security Framework.

The SDK provides a simple and typed interface for interacting with HostMK Security services, including:

- Authentication
- Permission Management
- Audit Logging
- Fraud Detection
- Multi-Tenant Management

---

## Requirements

- PHP 8.2 or higher
- JSON Extension

---

## Installation

Using Composer:

```bash
composer require hostmk/security-sdk
```

---

## Configuration

Create an environment file:

```env
HOSTMK_API_KEY=your-api-key
HOSTMK_BASE_URL=https://api.hostmk.com
```

---

## Basic Usage

```php
<?php

require __DIR__ . '/vendor/autoload.php';

use HostMK\Security\Client;

$client = new Client(
    $_ENV['HOSTMK_API_KEY'],
    $_ENV['HOSTMK_BASE_URL']
);
```

---

## Authentication

```php
use HostMK\Security\AuthService;

$auth = new AuthService($client);

$response = $auth->validateToken(
    'jwt-token'
);
```

---

## Permissions

```php
use HostMK\Security\PermissionService;

$permissions = new PermissionService(
    $client
);

$response = $permissions->check(
    'user-id',
    'tenant-id',
    'users.read'
);
```

---

## Audit Logs

```php
use HostMK\Security\AuditService;

$audit = new AuditService($client);

$events = $audit->list();
```

Create an audit event:

```php
$audit->create([
    'tenantId' => 'tenant-1',
    'actorId' => 'user-1',
    'resource' => 'orders',
    'action' => 'create',
    'metadata' => []
]);
```

---

## Fraud Detection

```php
use HostMK\Security\FraudService;

$fraud = new FraudService(
    $client
);

$result = $fraud->evaluate([
    'amount' => 1000,
    'currency' => 'USD'
]);
```

---

## Tenancy

```php
use HostMK\Security\TenancyService;

$tenancy = new TenancyService(
    $client
);

$tenants = $tenancy->list();
```

Create a tenant:

```php
$tenant = $tenancy->create([
    'name' => 'Demo Tenant',
    'slug' => 'demo'
]);
```

---

## Testing

Run PHPUnit:

```bash
composer test
```

Run PHPStan:

```bash
composer analyse
```

Run all checks:

```bash
composer check
```

---

## Development

Install dependencies:

```bash
composer install
```

Validate composer configuration:

```bash
composer validate
```

---

## Project Structure

```text
src/
├── Client.php
├── AuthService.php
├── PermissionService.php
├── AuditService.php
├── FraudService.php
└── TenancyService.php

tests/
├── ClientTest.php
├── AuthServiceTest.php
├── PermissionServiceTest.php
├── AuditServiceTest.php
├── FraudServiceTest.php
└── TenancyServiceTest.php
```

---

## License

MIT License

Copyright (c) 2026 HostMK Security Framework