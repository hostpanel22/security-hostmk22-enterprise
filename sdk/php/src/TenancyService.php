<?php

declare(strict_types=1);

namespace HostMK\Security;

use HostMK\Security\Contracts\HttpClientInterface;
use InvalidArgumentException;

final class TenancyService
{
    private HttpClientInterface $client;

    public function __construct(
        HttpClientInterface $client
    ) {
        $this->client = $client;
    }

    /**
     * Liste les tenants.
     */
    public function list(
        array $query = []
    ): array {
        $queryString = http_build_query(
            array_filter(
                $query,
                static fn (mixed $value): bool => $value !== null
            )
        );

        $endpoint = '/tenants';

        if ($queryString !== '') {
            $endpoint .= '?' . $queryString;
        }

        return $this->client->get(
            $endpoint
        );
    }

    /**
     * Récupère un tenant.
     */
    public function get(
        string $tenantId
    ): array {
        $tenantId = trim($tenantId);

        if ($tenantId === '') {
            throw new InvalidArgumentException(
                'Tenant ID is required.'
            );
        }

        return $this->client->get(
            '/tenants/' .
            urlencode($tenantId)
        );
    }

    /**
     * Crée un tenant.
     */
    public function create(
        string $name,
        string $slug
    ): array {
        $name = trim($name);
        $slug = trim($slug);

        if ($name === '') {
            throw new InvalidArgumentException(
                'Tenant name is required.'
            );
        }

        if ($slug === '') {
            throw new InvalidArgumentException(
                'Tenant slug is required.'
            );
        }

        return $this->client->post(
            '/tenants',
            [
                'name' => $name,
                'slug' => $slug,
            ]
        );
    }

    /**
     * Met à jour un tenant.
     */
    public function update(
        string $tenantId,
        array $payload
    ): array {
        $tenantId = trim($tenantId);

        if ($tenantId === '') {
            throw new InvalidArgumentException(
                'Tenant ID is required.'
            );
        }

        if ($payload === []) {
            throw new InvalidArgumentException(
                'Update payload is required.'
            );
        }

        return $this->client->put(
            '/tenants/' .
            urlencode($tenantId),
            $payload
        );
    }

    /**
     * Supprime un tenant.
     */
    public function delete(
        string $tenantId
    ): array {
        $tenantId = trim($tenantId);

        if ($tenantId === '') {
            throw new InvalidArgumentException(
                'Tenant ID is required.'
            );
        }

        return $this->client->delete(
            '/tenants/' .
            urlencode($tenantId)
        );
    }

    /**
     * Active un tenant.
     */
    public function activate(
        string $tenantId
    ): array {
        $tenantId = trim($tenantId);

        if ($tenantId === '') {
            throw new InvalidArgumentException(
                'Tenant ID is required.'
            );
        }

        return $this->client->post(
            '/tenants/' .
            urlencode($tenantId) .
            '/activate'
        );
    }

    /**
     * Désactive un tenant.
     */
    public function deactivate(
        string $tenantId
    ): array {
        $tenantId = trim($tenantId);

        if ($tenantId === '') {
            throw new InvalidArgumentException(
                'Tenant ID is required.'
            );
        }

        return $this->client->post(
            '/tenants/' .
            urlencode($tenantId) .
            '/deactivate'
        );
    }
}