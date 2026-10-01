package com.hostmk.sdk.models.secret;

/** Secret creation payload. */
public record CreateSecretRequest(String name, String type, String value, String description) {
}
