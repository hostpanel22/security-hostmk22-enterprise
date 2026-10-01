package com.hostmk.sdk.models.fraud;

import com.fasterxml.jackson.annotation.JsonIgnoreProperties;
import com.fasterxml.jackson.annotation.JsonProperty;

import java.time.Instant;
import java.util.Map;

/** Fraud detection rule. */
@JsonIgnoreProperties(ignoreUnknown = true)
public record FraudRule(
        String id,
        @JsonProperty("tenant_id") String tenantId,
        String name,
        boolean enabled,
        String field,
        String operator,
        String threshold,
        @JsonProperty("window_seconds") Integer windowSeconds,
        String action,
        int score,
        Map<String, Object> metadata,
        @JsonProperty("created_at") Instant createdAt,
        @JsonProperty("updated_at") Instant updatedAt
) {
}
