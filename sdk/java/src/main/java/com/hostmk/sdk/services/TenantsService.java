package com.hostmk.sdk.services;

import com.hostmk.sdk.http.HostMKHttpClient;
import com.hostmk.sdk.models.ListResult;
import com.hostmk.sdk.models.tenant.Tenant;

import java.util.Map;

/** Tenant creation and status management. */
public final class TenantsService {
    private final HostMKHttpClient client;

    public TenantsService(HostMKHttpClient client) {
        this.client = client;
    }

    public ListResult<Tenant> list() {
        return client.list("/api/v1/tenants", Map.of(), "tenants", Tenant.class);
    }

    public Tenant create(Tenant tenant) {
        return client.post("/api/v1/tenants", Map.of(), tenant, Tenant.class);
    }

    public Tenant update(String id, Tenant tenant) {
        return client.post("/api/v1/tenants/action", Map.of("id", id), tenant, Tenant.class);
    }

    public Tenant activate(String id) {
        return status(id, "activate");
    }

    public Tenant suspend(String id) {
        return status(id, "suspend");
    }

    private Tenant status(String id, String action) {
        return client.post("/api/v1/tenants/action", Map.of("id", id, "action", action), null, Tenant.class);
    }
}
