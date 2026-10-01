package com.hostmk.sdk.models.user;

import com.fasterxml.jackson.annotation.JsonProperty;

import java.util.List;

/** User creation payload. */
public record CreateUserRequest(String email, String name, String password, String status,
                                @JsonProperty("role_ids") List<String> roleIds) {
}
