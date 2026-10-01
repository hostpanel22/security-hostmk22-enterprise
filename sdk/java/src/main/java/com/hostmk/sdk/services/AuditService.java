package com.hostmk.sdk.services;

import com.hostmk.sdk.http.HostMKHttpClient;
import com.hostmk.sdk.models.ListResult;
import com.hostmk.sdk.models.audit.AuditEvent;
import com.hostmk.sdk.models.audit.CreateAuditEventRequest;

import java.util.Map;

/** Audit event listing and creation. */
public final class AuditService {
    private final HostMKHttpClient client;

    public AuditService(HostMKHttpClient client) {
        this.client = client;
    }

    public ListResult<AuditEvent> list() {
        return list(0);
    }

    public ListResult<AuditEvent> list(int limit) {
        Map<String, ?> query = limit > 0 ? Map.of("limit", limit) : Map.of();
        return client.list("/api/v1/audit/events", query, "events", AuditEvent.class);
    }

    public AuditEvent create(CreateAuditEventRequest request) {
        return client.post("/api/v1/audit/events", Map.of(), request, AuditEvent.class);
    }
}
