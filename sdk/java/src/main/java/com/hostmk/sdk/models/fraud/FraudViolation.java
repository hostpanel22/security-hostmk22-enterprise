package com.hostmk.sdk.models.fraud;

import com.fasterxml.jackson.annotation.JsonIgnoreProperties;
import com.fasterxml.jackson.annotation.JsonProperty;

import java.time.Instant;
import java.util.Map;

/** Fraud violation emitted after rule evaluation. */
@JsonIgnoreProperties(ignoreUnknown = true)
public record FraudViolation(
        String id,
        @JsonProperty("tenant_id") String tenantId,
        @JsonProperty("user_id") String userId,
        @JsonProperty("signal_id") String signalId,
        @JsonProperty("rule_id") String ruleId,
        double score,
        String level,
        String reason,
        String status,
        Map<String, Object> metadata,
        @JsonProperty("created_at") Instant createdAt,
        @JsonProperty("updated_at") Instant updatedAt
) {
}
