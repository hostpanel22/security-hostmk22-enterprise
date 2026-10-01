package com.hostmk.sdk.models.tenant;

import com.fasterxml.jackson.annotation.JsonIgnoreProperties;
import com.fasterxml.jackson.annotation.JsonProperty;

import java.time.Instant;
import java.util.Map;

/** HostMK tenant. */
@JsonIgnoreProperties(ignoreUnknown = true)
public record Tenant(
        String id,
        String name,
        String slug,
        String domain,
        String status,
        @JsonProperty("owner_id") String ownerId,
        String plan,
        String timezone,
        String country,
        Map<String, Object> quotas,
        Map<String, Object> branding,
        Map<String, Object> metadata,
        @JsonProperty("created_at") Instant createdAt,
        @JsonProperty("updated_at") Instant updatedAt,
        @JsonProperty("suspended_at") Instant suspendedAt,
        @JsonProperty("expires_at") Instant expiresAt
) {
}
