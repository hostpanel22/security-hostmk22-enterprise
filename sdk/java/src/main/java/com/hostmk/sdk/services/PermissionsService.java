package com.hostmk.sdk.services;

import com.hostmk.sdk.http.HostMKHttpClient;
import com.hostmk.sdk.models.ListResult;
import com.hostmk.sdk.models.rbac.Permission;
import com.hostmk.sdk.models.rbac.PermissionCheckResult;

import java.util.Map;

/** Permission CRUD and authorization checks. */
public final class PermissionsService {
    private final HostMKHttpClient client;

    public PermissionsService(HostMKHttpClient client) {
        this.client = client;
    }

    public ListResult<Permission> list() {
        return client.list("/api/v1/permissions", Map.of(), "permissions", Permission.class);
    }

    public Permission create(Permission permission) {
        return client.post("/api/v1/permissions", Map.of(), permission, Permission.class);
    }

    public void delete(String id) {
        client.delete("/api/v1/permissions", Map.of("id", id));
    }

    public PermissionCheckResult check(String userId, String permission) {
        return client.post(
                "/api/v1/permissions/check",
                Map.of(),
                Map.of("user_id", userId, "permission", permission),
                PermissionCheckResult.class
        );
    }
}
