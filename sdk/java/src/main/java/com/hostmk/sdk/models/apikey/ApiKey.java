package com.hostmk.sdk.models.apikey;

import com.fasterxml.jackson.annotation.JsonIgnoreProperties;
import com.fasterxml.jackson.annotation.JsonProperty;

import java.time.Instant;
import java.util.List;

/** API key metadata; plaintext key is only populated after create or rotate. */
@JsonIgnoreProperties(ignoreUnknown = true)
public record ApiKey(
        String id,
        @JsonProperty("tenant_id") String tenantId,
        String name,
        String prefix,
        String key,
        List<String> scopes,
        @JsonProperty("rate_limit") Integer rateLimit,
        @JsonProperty("expires_at") Instant expiresAt,
        @JsonProperty("revoked_at") Instant revokedAt,
        @JsonProperty("last_used_at") Instant lastUsedAt,
        @JsonProperty("created_at") Instant createdAt,
        @JsonProperty("updated_at") Instant updatedAt
) {
}
