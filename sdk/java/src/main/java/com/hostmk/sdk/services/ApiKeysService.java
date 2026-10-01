package com.hostmk.sdk.services;

import com.hostmk.sdk.http.HostMKHttpClient;
import com.hostmk.sdk.models.ListResult;
import com.hostmk.sdk.models.apikey.ApiKey;
import com.hostmk.sdk.models.apikey.CreateApiKeyRequest;

import java.util.Map;

/** API key lifecycle operations. */
public final class ApiKeysService {
    private final HostMKHttpClient client;

    public ApiKeysService(HostMKHttpClient client) {
        this.client = client;
    }

    public ListResult<ApiKey> list() {
        return client.list("/api/v1/api-keys", Map.of(), "api_keys", ApiKey.class);
    }

    public ApiKey create(CreateApiKeyRequest request) {
        return client.post("/api/v1/api-keys", Map.of(), request, ApiKey.class);
    }

    public ApiKey rotate(String id) {
        return client.post("/api/v1/api-keys", Map.of("id", id, "action", "rotate"), null, ApiKey.class);
    }

    public void revoke(String id) {
        client.delete("/api/v1/api-keys", Map.of("id", id));
    }
}
