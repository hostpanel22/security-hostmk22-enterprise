<?php

declare(strict_types=1);

namespace HostMK\Security;

use HostMK\Security\Contracts\HttpClientInterface;
use InvalidArgumentException;

final class PermissionService
{
    private HttpClientInterface $client;

    public function __construct(
        HttpClientInterface $client
    ) {
        $this->client = $client;
    }

    /**
     * Liste toutes les permissions disponibles.
     */
    public function list(): array
    {
        return $this->client->get(
            '/permissions'
        );
    }

    /**
     * Retourne les permissions d'un utilisateur.
     */
    public function getUserPermissions(
        string $userId
    ): array {
        $userId = trim($userId);

        if ($userId === '') {
            throw new InvalidArgumentException(
                'User ID is required.'
            );
        }

        return $this->client->get(
            '/permissions/' .
            urlencode($userId)
        );
    }

    /**
     * Vérifie si un utilisateur possède une permission.
     */
    public function check(
        string $userId,
        string $permission
    ): array {
        $userId = trim($userId);
        $permission = trim($permission);

        if ($userId === '') {
            throw new InvalidArgumentException(
                'User ID is required.'
            );
        }

        if ($permission === '') {
            throw new InvalidArgumentException(
                'Permission is required.'
            );
        }

        return $this->client->post(
            '/permissions/check',
            [
                'userId' => $userId,
                'permission' => $permission,
            ]
        );
    }

    /**
     * Assigne un rôle à un utilisateur.
     */
    public function assignRole(
        string $userId,
        string $role
    ): array {
        $userId = trim($userId);
        $role = trim($role);

        if ($userId === '') {
            throw new InvalidArgumentException(
                'User ID is required.'
            );
        }

        if ($role === '') {
            throw new InvalidArgumentException(
                'Role is required.'
            );
        }

        return $this->client->post(
            '/permissions/assign-role',
            [
                'userId' => $userId,
                'role' => $role,
            ]
        );
    }

    /**
     * Retire un rôle à un utilisateur.
     */
    public function removeRole(
        string $userId,
        string $role
    ): array {
        $userId = trim($userId);
        $role = trim($role);

        if ($userId === '') {
            throw new InvalidArgumentException(
                'User ID is required.'
            );
        }

        if ($role === '') {
            throw new InvalidArgumentException(
                'Role is required.'
            );
        }

        return $this->client->post(
            '/permissions/remove-role',
            [
                'userId' => $userId,
                'role' => $role,
            ]
        );
    }
}