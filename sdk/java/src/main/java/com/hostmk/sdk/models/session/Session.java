package com.hostmk.sdk.models.session;

import com.fasterxml.jackson.annotation.JsonIgnoreProperties;
import com.fasterxml.jackson.annotation.JsonProperty;

import java.time.Instant;

/** User refresh-token session. */
@JsonIgnoreProperties(ignoreUnknown = true)
public record Session(
        String id,
        @JsonProperty("tenant_id") String tenantId,
        @JsonProperty("user_id") String userId,
        @JsonProperty("user_agent") String userAgent,
        @JsonProperty("ip_address") String ipAddress,
        @JsonProperty("expires_at") Instant expiresAt,
        @JsonProperty("revoked_at") Instant revokedAt,
        @JsonProperty("created_at") Instant createdAt
) {
}
