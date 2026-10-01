package com.hostmk.sdk.models.user;

import com.fasterxml.jackson.annotation.JsonProperty;

import java.util.List;

/** Mutable user fields accepted by the action endpoint. */
public record UpdateUserRequest(String name, String status, @JsonProperty("role_ids") List<String> roleIds) {
}
