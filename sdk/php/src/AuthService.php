<?php

declare(strict_types=1);

namespace HostMK\Security;

use HostMK\Security\Contracts\HttpClientInterface;
use InvalidArgumentException;

final class AuthService
{
    private HttpClientInterface $client;

    public function __construct(
        HttpClientInterface $client
    ) {
        $this->client = $client;
    }

    public function login(
        string $email,
        string $password
    ): array {
        $email = trim($email);
        $password = trim($password);

        if ($email === '') {
            throw new InvalidArgumentException(
                'Email is required.'
            );
        }

        if (
            !filter_var(
                $email,
                FILTER_VALIDATE_EMAIL
            )
        ) {
            throw new InvalidArgumentException(
                'Invalid email address.'
            );
        }

        if ($password === '') {
            throw new InvalidArgumentException(
                'Password is required.'
            );
        }

        return $this->client->post(
            '/auth/login',
            [
                'email' => $email,
                'password' => $password,
            ]
        );
    }

    public function validateToken(
        string $token
    ): array {
        $token = trim($token);

        if ($token === '') {
            throw new InvalidArgumentException(
                'Token is required.'
            );
        }

        return $this->client->post(
            '/auth/validate',
            [
                'token' => $token,
            ]
        );
    }

    public function refreshToken(
        string $refreshToken
    ): array {
        $refreshToken = trim(
            $refreshToken
        );

        if ($refreshToken === '') {
            throw new InvalidArgumentException(
                'Refresh token is required.'
            );
        }

        return $this->client->post(
            '/auth/refresh',
            [
                'refreshToken' => $refreshToken,
            ]
        );
    }

    public function revokeToken(
        string $token
    ): array {
        $token = trim($token);

        if ($token === '') {
            throw new InvalidArgumentException(
                'Token is required.'
            );
        }

        return $this->client->post(
            '/auth/revoke',
            [
                'token' => $token,
            ]
        );
    }

    public function logout(
        string $token
    ): array {
        $token = trim($token);

        if ($token === '') {
            throw new InvalidArgumentException(
                'Token is required.'
            );
        }

        return $this->client->post(
            '/auth/logout',
            [
                'token' => $token,
            ]
        );
    }

    public function me(): array
    {
        return $this->client->get(
            '/auth/me'
        );
    }

    public function getPermissions(): array
    {
        return $this->client->get(
            '/auth/permissions'
        );
    }

    public function getSessions(): array
    {
        return $this->client->get(
            '/auth/sessions'
        );
    }

    public function revokeSession(
        string $sessionId
    ): array {
        $sessionId = trim(
            $sessionId
        );

        if ($sessionId === '') {
            throw new InvalidArgumentException(
                'Session ID is required.'
            );
        }

        return $this->client->post(
            '/auth/sessions/revoke',
            [
                'sessionId' => $sessionId,
            ]
        );
    }
}