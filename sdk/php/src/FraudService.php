<?php

declare(strict_types=1);

namespace HostMK\Security;

use HostMK\Security\Contracts\HttpClientInterface;
use InvalidArgumentException;

final class FraudService
{
    private HttpClientInterface $client;

    public function __construct(
        HttpClientInterface $client
    ) {
        $this->client = $client;
    }

    public function check(
        array $payload
    ): array {
        if ($payload === []) {
            throw new InvalidArgumentException(
                'Fraud check payload is required.'
            );
        }

        return $this->client->post(
            '/fraud/check',
            $payload
        );
    }

    public function score(
        array $payload
    ): array {
        if ($payload === []) {
            throw new InvalidArgumentException(
                'Fraud score payload is required.'
            );
        }

        return $this->client->post(
            '/fraud/score',
            $payload
        );
    }

    public function report(
        array $payload
    ): array {
        if ($payload === []) {
            throw new InvalidArgumentException(
                'Fraud report payload is required.'
            );
        }

        return $this->client->post(
            '/fraud/report',
            $payload
        );
    }

    public function reports(
        array $query = []
    ): array {
        $queryString = http_build_query(
            array_filter(
                $query,
                static fn (
                    mixed $value
                ): bool => $value !== null
            )
        );

        $endpoint = '/fraud/reports';

        if ($queryString !== '') {
            $endpoint .= '?' . $queryString;
        }

        return $this->client->get(
            $endpoint
        );
    }

    public function getReport(
        string $reportId
    ): array {
        $reportId = trim(
            $reportId
        );

        if ($reportId === '') {
            throw new InvalidArgumentException(
                'Report ID is required.'
            );
        }

        return $this->client->get(
            '/fraud/reports/' .
            urlencode(
                $reportId
            )
        );
    }

    public function statistics(
        array $filters = []
    ): array {
        $queryString = http_build_query(
            array_filter(
                $filters,
                static fn (
                    mixed $value
                ): bool => $value !== null
            )
        );

        $endpoint = '/fraud/statistics';

        if ($queryString !== '') {
            $endpoint .= '?' . $queryString;
        }

        return $this->client->get(
            $endpoint
        );
    }

    public function userRisk(
        string $userId
    ): array {
        $userId = trim(
            $userId
        );

        if ($userId === '') {
            throw new InvalidArgumentException(
                'User ID is required.'
            );
        }

        return $this->client->get(
            '/fraud/users/' .
            urlencode(
                $userId
            ) .
            '/risk'
        );
    }

    public function ipRisk(
        string $ip
    ): array {
        $ip = trim($ip);

        if ($ip === '') {
            throw new InvalidArgumentException(
                'IP address is required.'
            );
        }

        return $this->client->get(
            '/fraud/ip/' .
            urlencode(
                $ip
            )
        );
    }

    public function deviceRisk(
        string $deviceId
    ): array {
        $deviceId = trim(
            $deviceId
        );

        if ($deviceId === '') {
            throw new InvalidArgumentException(
                'Device ID is required.'
            );
        }

        return $this->client->get(
            '/fraud/devices/' .
            urlencode(
                $deviceId
            )
        );
    }

    public function whitelist(
        array $payload
    ): array {
        if ($payload === []) {
            throw new InvalidArgumentException(
                'Whitelist payload is required.'
            );
        }

        return $this->client->post(
            '/fraud/whitelist',
            $payload
        );
    }

    public function blacklist(
        array $payload
    ): array {
        if ($payload === []) {
            throw new InvalidArgumentException(
                'Blacklist payload is required.'
            );
        }

        return $this->client->post(
            '/fraud/blacklist',
            $payload
        );
    }

    public function getWhitelist(): array
    {
        return $this->client->get(
            '/fraud/whitelist'
        );
    }

    public function getBlacklist(): array
    {
        return $this->client->get(
            '/fraud/blacklist'
        );
    }

    public function health(): array
    {
        return $this->client->get(
            '/fraud/health'
        );
    }
}