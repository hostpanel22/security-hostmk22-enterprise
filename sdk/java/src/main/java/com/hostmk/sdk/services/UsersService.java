package com.hostmk.sdk.services;

import com.hostmk.sdk.http.HostMKHttpClient;
import com.hostmk.sdk.models.ListResult;
import com.hostmk.sdk.models.user.CreateUserRequest;
import com.hostmk.sdk.models.user.UpdateUserRequest;
import com.hostmk.sdk.models.user.User;
import com.hostmk.sdk.models.user.UserActionResult;

import java.util.Map;

/** User lifecycle and MFA actions. */
public final class UsersService {
    private final HostMKHttpClient client;

    public UsersService(HostMKHttpClient client) {
        this.client = client;
    }

    public ListResult<User> list() {
        return client.list("/api/v1/users", Map.of(), "users", User.class);
    }

    public User create(CreateUserRequest request) {
        return client.post("/api/v1/users", Map.of(), request, User.class);
    }

    public User update(String id, UpdateUserRequest request) {
        return client.post("/api/v1/users/action", Map.of("id", id), request, User.class);
    }

    public UserActionResult action(String id, String action) {
        return client.post(
                "/api/v1/users/action",
                Map.of("id", id, "action", action),
                null,
                UserActionResult.class
        );
    }
}
