<?php

declare(strict_types=1);

return [
    'base_url' => env('HOSTMK_BASE_URL', 'http://localhost:8080'),
    'api_key' => env('HOSTMK_API_KEY', ''),
    'token' => env('HOSTMK_TOKEN', ''),
    'timeout' => (int) env('HOSTMK_TIMEOUT', 30),
    'user_agent' => env('HOSTMK_USER_AGENT', 'hostmk-laravel-sdk/1.0'),
];
