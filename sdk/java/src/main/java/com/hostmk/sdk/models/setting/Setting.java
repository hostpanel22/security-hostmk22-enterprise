package com.hostmk.sdk.models.setting;

import com.fasterxml.jackson.annotation.JsonIgnoreProperties;
import com.fasterxml.jackson.annotation.JsonProperty;
import com.fasterxml.jackson.databind.JsonNode;

import java.time.Instant;

/** Tenant or platform setting. */
@JsonIgnoreProperties(ignoreUnknown = true)
public record Setting(
        String id,
        @JsonProperty("tenant_id") String tenantId,
        String key,
        JsonNode value,
        @JsonProperty("updated_at") Instant updatedAt
) {
}
