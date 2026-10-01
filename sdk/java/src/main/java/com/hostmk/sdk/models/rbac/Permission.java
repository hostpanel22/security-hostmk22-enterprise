package com.hostmk.sdk.models.rbac;

import com.fasterxml.jackson.annotation.JsonIgnoreProperties;
import com.fasterxml.jackson.annotation.JsonProperty;

import java.time.Instant;

/** Tenant permission. */
@JsonIgnoreProperties(ignoreUnknown = true)
public record Permission(
        String id,
        @JsonProperty("tenant_id") String tenantId,
        String name,
        String resource,
        String action,
        String description,
        @JsonProperty("created_at") Instant createdAt,
        @JsonProperty("updated_at") Instant updatedAt
) {
}
