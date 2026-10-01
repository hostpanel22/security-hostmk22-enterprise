package com.hostmk.sdk.models.auth;

import com.fasterxml.jackson.annotation.JsonProperty;

/** Initial platform bootstrap payload. */
public record BootstrapRequest(
        @JsonProperty("tenant_name") String tenantName,
        @JsonProperty("tenant_slug") String tenantSlug,
        String name,
        String email,
        String password
) {
}
