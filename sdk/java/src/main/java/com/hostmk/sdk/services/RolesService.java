package com.hostmk.sdk.services;

import com.hostmk.sdk.http.HostMKHttpClient;
import com.hostmk.sdk.models.ListResult;
import com.hostmk.sdk.models.rbac.Role;

import java.util.Map;

/** RBAC role CRUD operations. */
public final class RolesService {
    private final HostMKHttpClient client;

    public RolesService(HostMKHttpClient client) {
        this.client = client;
    }

    public ListResult<Role> list() {
        return client.list("/api/v1/rbac/roles", Map.of(), "roles", Role.class);
    }

    public Role create(Role role) {
        return client.post("/api/v1/rbac/roles", Map.of(), role, Role.class);
    }

    public Role update(Role role) {
        return client.put("/api/v1/rbac/roles", Map.of(), role, Role.class);
    }

    public void delete(String id) {
        client.delete("/api/v1/rbac/roles", Map.of("id", id));
    }
}
