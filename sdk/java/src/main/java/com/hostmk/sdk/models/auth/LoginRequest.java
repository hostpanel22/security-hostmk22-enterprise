package com.hostmk.sdk.models.auth;

import com.fasterxml.jackson.annotation.JsonProperty;

/** Credentials accepted by the tenant-aware login endpoint. */
public record LoginRequest(
        @JsonProperty("tenant_slug") String tenantSlug,
        String email,
        String password,
        @JsonProperty("mfa_code") String mfaCode
) {
}
