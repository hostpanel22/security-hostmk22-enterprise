package com.hostmk.sdk.models.rbac;

import com.fasterxml.jackson.annotation.JsonIgnoreProperties;
import com.fasterxml.jackson.annotation.JsonProperty;

import java.time.Instant;
import java.util.List;

/** Tenant role. */
@JsonIgnoreProperties(ignoreUnknown = true)
public record Role(
        String id,
        @JsonProperty("tenant_id") String tenantId,
        String name,
        String description,
        List<String> permissions,
        @JsonProperty("created_at") Instant createdAt,
        @JsonProperty("updated_at") Instant updatedAt
) {
}
