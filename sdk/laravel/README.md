# HostMK Security Laravel SDK

Laravel 10+ package for the routes registered by HostMK Security Framework.

## Installation

```bash
composer require hostmk/security-laravel
php artisan vendor:publish --tag=hostmk-config
```

Configure the application:

```dotenv
HOSTMK_BASE_URL=https://security.example.com
HOSTMK_API_KEY=your-api-key
# Or use HOSTMK_TOKEN for bearer authentication.
```

The service provider and `HostMK` facade are discovered automatically by Composer.

## Usage

```php
use HostMK\Laravel\Facades\HostMK;

$users = HostMK::users()->list();
$events = HostMK::audit()->list(100);
$analysis = HostMK::fraud()->analyze(['amount' => 2500]);
```

Login adopts the returned bearer token for every request-scoped service:

```php
$tokens = HostMK::auth()->login('tenant-slug', 'owner@example.com', 'secret');
$sessions = HostMK::sessions()->list();
```

To forward an incoming bearer token, register `HostMK\Laravel\Http\Middleware\UseHostMKBearerToken` on the desired Laravel routes. The client is request-scoped and safe for Laravel Octane workers.

Supported modules: Authentication, Users, Sessions, RBAC, Permissions, API Keys, Secrets, Audit, Fraud, Monitoring, Settings, and Tenants.

## Tests

```bash
composer install
composer test
```
