package com.hostmk.sdk.services;

import com.hostmk.sdk.http.HostMKHttpClient;
import com.hostmk.sdk.models.auth.BootstrapRequest;
import com.hostmk.sdk.models.auth.BootstrapResult;
import com.hostmk.sdk.models.auth.LoginRequest;
import com.hostmk.sdk.models.auth.TokenSet;

import java.util.Map;

/** Authentication and initial platform bootstrap operations. */
public final class AuthService {
    private final HostMKHttpClient client;

    public AuthService(HostMKHttpClient client) {
        this.client = client;
    }

    public BootstrapResult bootstrap(BootstrapRequest request) {
        return client.post("/api/v1/auth/bootstrap", Map.of(), request, BootstrapResult.class);
    }

    public TokenSet login(LoginRequest request) {
        TokenSet tokens = client.post("/api/v1/auth/login", Map.of(), request, TokenSet.class);
        client.credentials().useBearerToken(tokens.accessToken());
        return tokens;
    }

    public TokenSet refresh(String refreshToken) {
        TokenSet tokens = client.post(
                "/api/v1/auth/refresh",
                Map.of(),
                Map.of("refresh_token", refreshToken),
                TokenSet.class
        );
        client.credentials().useBearerToken(tokens.accessToken());
        return tokens;
    }

    public void clearCredentials() {
        client.credentials().clear();
    }
}
