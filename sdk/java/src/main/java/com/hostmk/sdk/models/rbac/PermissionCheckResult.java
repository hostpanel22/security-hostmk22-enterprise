package com.hostmk.sdk.models.rbac;

import com.fasterxml.jackson.annotation.JsonProperty;

/** Permission decision for a user. */
public record PermissionCheckResult(
        @JsonProperty("user_id") String userId,
        String permission,
        boolean allowed
) {
}
