package com.hostmk.sdk.models.auth;

import com.fasterxml.jackson.annotation.JsonIgnoreProperties;
import com.fasterxml.jackson.annotation.JsonProperty;

/** Opaque access and refresh credentials returned by HostMK. */
@JsonIgnoreProperties(ignoreUnknown = true)
public record TokenSet(
        @JsonProperty("access_token") String accessToken,
        @JsonProperty("refresh_token") String refreshToken,
        @JsonProperty("token_type") String tokenType,
        @JsonProperty("expires_in") int expiresIn
) {
}
