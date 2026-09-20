<?php

declare(strict_types=1);

namespace HostMK\Security;

use HostMK\Security\Contracts\HttpClientInterface;
use InvalidArgumentException;

final class AuditService
{
    private HttpClientInterface $client;

    public function __construct(
        HttpClientInterface $client
    ) {
        $this->client = $client;
    }

    public function list(
        array $query = []
    ): array {
        $queryString = http_build_query(
            array_filter(
                $query,
                static fn (mixed $value): bool => $value !== null
            )
        );

        $endpoint = '/audit';

        if ($queryString !== '') {
            $endpoint .= '?' . $queryString;
        }

        return $this->client->get(
            $endpoint
        );
    }

    public function get(
        string $auditId
    ): array {
        $auditId = trim($auditId);

        if ($auditId === '') {
            throw new InvalidArgumentException(
                'Audit ID is required.'
            );
        }

        return $this->client->get(
            '/audit/' . urlencode($auditId)
        );
    }

    public function create(
        array $event
    ): array {
        if ($event === []) {
            throw new InvalidArgumentException(
                'Audit event payload is required.'
            );
        }

        return $this->client->post(
            '/audit',
            $event
        );
    }

    public function delete(
        string $auditId
    ): array {
        $auditId = trim($auditId);

        if ($auditId === '') {
            throw new InvalidArgumentException(
                'Audit ID is required.'
            );
        }

        return $this->client->delete(
            '/audit/' . urlencode($auditId)
        );
    }

    public function search(
        array $filters = []
    ): array {
        $queryString = http_build_query(
            array_filter(
                $filters,
                static fn (mixed $value): bool => $value !== null
            )
        );

        $endpoint = '/audit/search';

        if ($queryString !== '') {
            $endpoint .= '?' . $queryString;
        }

        return $this->client->get(
            $endpoint
        );
    }

    public function stats(
        array $filters = []
    ): array {
        $queryString = http_build_query(
            array_filter(
                $filters,
                static fn (mixed $value): bool => $value !== null
            )
        );

        $endpoint = '/audit/stats';

        if ($queryString !== '') {
            $endpoint .= '?' . $queryString;
        }

        return $this->client->get(
            $endpoint
        );
    }

    public function export(
        array $filters = []
    ): array {
        return $this->client->post(
            '/audit/export',
            $filters
        );
    }

    public function userActivity(
        string $userId,
        array $filters = []
    ): array {
        $userId = trim($userId);

        if ($userId === '') {
            throw new InvalidArgumentException(
                'User ID is required.'
            );
        }

        $queryString = http_build_query(
            array_filter(
                $filters,
                static fn (mixed $value): bool => $value !== null
            )
        );

        $endpoint =
            '/audit/users/' .
            urlencode($userId);

        if ($queryString !== '') {
            $endpoint .= '?' . $queryString;
        }

        return $this->client->get(
            $endpoint
        );
    }

    public function tenantActivity(
        string $tenantId,
        array $filters = []
    ): array {
        $tenantId = trim($tenantId);

        if ($tenantId === '') {
            throw new InvalidArgumentException(
                'Tenant ID is required.'
            );
        }

        $queryString = http_build_query(
            array_filter(
                $filters,
                static fn (mixed $value): bool => $value !== null
            )
        );

        $endpoint =
            '/audit/tenants/' .
            urlencode($tenantId);

        if ($queryString !== '') {
            $endpoint .= '?' . $queryString;
        }

        return $this->client->get(
            $endpoint
        );
    }

    public function retention(
        int $days
    ): array {
        if ($days <= 0) {
            throw new InvalidArgumentException(
                'Retention days must be greater than zero.'
            );
        }

        return $this->client->post(
            '/audit/retention',
            [
                'days' => $days,
            ]
        );
    }

    public function purge(
        string $beforeDate
    ): array {
        $beforeDate = trim($beforeDate);

        if ($beforeDate === '') {
            throw new InvalidArgumentException(
                'Date is required.'
            );
        }

        return $this->client->post(
            '/audit/purge',
            [
                'beforeDate' => $beforeDate,
            ]
        );
    }

    public function categories(): array
    {
        return $this->client->get(
            '/audit/categories'
        );
    }

    public function actions(): array
    {
        return $this->client->get(
            '/audit/actions'
        );
    }

    public function health(): array
    {
        return $this->client->get(
            '/audit/health'
        );
    }
}