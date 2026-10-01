package com.hostmk.sdk.models.apikey;

import com.fasterxml.jackson.annotation.JsonProperty;

import java.time.Instant;
import java.util.List;

/** API key creation payload. */
public record CreateApiKeyRequest(
        String name,
        List<String> scopes,
        @JsonProperty("rate_limit") Integer rateLimit,
        @JsonProperty("expires_at") Instant expiresAt
) {
}
