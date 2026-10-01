package com.hostmk.sdk.services;

import com.hostmk.sdk.http.HostMKHttpClient;
import com.hostmk.sdk.models.ListResult;
import com.hostmk.sdk.models.secret.CreateSecretRequest;
import com.hostmk.sdk.models.secret.Secret;
import com.hostmk.sdk.models.secret.SecretVersion;

import java.util.Map;

/** Managed secret lifecycle and version history. */
public final class SecretsService {
    private final HostMKHttpClient client;

    public SecretsService(HostMKHttpClient client) {
        this.client = client;
    }

    public ListResult<Secret> list() {
        return client.list("/api/v1/secrets", Map.of(), "secrets", Secret.class);
    }

    public Secret create(CreateSecretRequest request) {
        return client.post("/api/v1/secrets", Map.of(), request, Secret.class);
    }

    public Secret rotate(String id, String value) {
        return client.post(
                "/api/v1/secrets",
                Map.of("id", id, "action", "rotate", "value", value),
                null,
                Secret.class
        );
    }

    public ListResult<SecretVersion> history(String id) {
        return client.list("/api/v1/secrets", Map.of("id", id, "history", true), "", SecretVersion.class);
    }

    public void revoke(String id) {
        client.delete("/api/v1/secrets", Map.of("id", id));
    }
}
