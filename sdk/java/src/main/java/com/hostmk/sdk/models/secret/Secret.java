package com.hostmk.sdk.models.secret;

import com.fasterxml.jackson.annotation.JsonIgnoreProperties;
import com.fasterxml.jackson.annotation.JsonProperty;

import java.time.Instant;

/** Managed secret metadata. */
@JsonIgnoreProperties(ignoreUnknown = true)
public record Secret(
        String id,
        @JsonProperty("tenant_id") String tenantId,
        String name,
        String type,
        String value,
        @JsonProperty("masked_value") String maskedValue,
        String description,
        @JsonProperty("revoked_at") Instant revokedAt,
        @JsonProperty("created_at") Instant createdAt,
        @JsonProperty("updated_at") Instant updatedAt
) {
}
