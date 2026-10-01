package com.hostmk.sdk.services;

import com.hostmk.sdk.http.HostMKHttpClient;
import com.hostmk.sdk.models.ListResult;
import com.hostmk.sdk.models.session.Session;

import java.util.Map;

/** Refresh-token session operations. */
public final class SessionsService {
    private final HostMKHttpClient client;

    public SessionsService(HostMKHttpClient client) {
        this.client = client;
    }

    public ListResult<Session> list() {
        return client.list("/api/v1/sessions", Map.of(), "sessions", Session.class);
    }

    public void revoke(String id) {
        client.post("/api/v1/sessions", Map.of("id", id, "action", "revoke"), null);
    }
}
