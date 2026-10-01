<?php

declare(strict_types=1);

namespace HostMK\Laravel\Services;

/** Fraud rules, analysis, and violation operations. */
final class FraudService extends AbstractService
{
    /** @return array{count?:int,rules?:list<array<string,mixed>>} */
    public function listRules(): array
    {
        return $this->client->get('/api/v1/fraud/rules');
    }

    /** @param array<string,mixed> $rule @return array<string,mixed> */
    public function createRule(array $rule): array
    {
        return $this->client->post('/api/v1/fraud/rules', $rule);
    }

    /** @param array<string,mixed> $rule @return array<string,mixed> */
    public function updateRule(array $rule): array
    {
        return $this->client->put('/api/v1/fraud/rules', $rule);
    }

    /** @return array<string,mixed> */
    public function deleteRule(string $id): array
    {
        return $this->client->delete('/api/v1/fraud/rules', ['id' => $id]);
    }

    /** @param array<string,mixed> $attributes @return array<string,mixed> */
    public function analyze(array $attributes): array
    {
        return $this->client->post('/api/v1/fraud/analyze', $attributes);
    }

    /** @return array{count?:int,events?:list<array<string,mixed>>} */
    public function listViolations(): array
    {
        return $this->client->get('/api/v1/fraud/violations');
    }
}
