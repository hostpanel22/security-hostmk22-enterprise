package com.hostmk.sdk.models.user;

import com.fasterxml.jackson.annotation.JsonIgnoreProperties;
import com.fasterxml.jackson.annotation.JsonProperty;

import java.time.Instant;
import java.util.List;

/** HostMK user. */
@JsonIgnoreProperties(ignoreUnknown = true)
public record User(
        String id,
        @JsonProperty("tenant_id") String tenantId,
        String email,
        String name,
        String status,
        @JsonProperty("email_verified_at") Instant emailVerifiedAt,
        @JsonProperty("mfa_enabled") boolean mfaEnabled,
        @JsonProperty("role_ids") List<String> roleIds,
        @JsonProperty("created_at") Instant createdAt,
        @JsonProperty("updated_at") Instant updatedAt,
        @JsonProperty("last_login") Instant lastLogin
) {
}
