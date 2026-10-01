package com.hostmk.sdk.models.secret;

import com.fasterxml.jackson.annotation.JsonIgnoreProperties;
import com.fasterxml.jackson.annotation.JsonProperty;

import java.time.Instant;

/** Historical version of a managed secret. */
@JsonIgnoreProperties(ignoreUnknown = true)
public record SecretVersion(
        String id,
        @JsonProperty("secret_id") String secretId,
        @JsonProperty("masked_value") String maskedValue,
        @JsonProperty("created_at") Instant createdAt,
        @JsonProperty("created_by") String createdBy
) {
}
