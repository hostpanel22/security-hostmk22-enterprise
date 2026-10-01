package com.hostmk.sdk.models.user;

import com.fasterxml.jackson.annotation.JsonProperty;

/** Result returned by user actions such as MFA enrollment. */
public record UserActionResult(String id, String secret, @JsonProperty("otpauth_url") String otpAuthUrl) {
}
