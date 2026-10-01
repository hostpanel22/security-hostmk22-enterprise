package com.hostmk.sdk;

import com.hostmk.sdk.models.auth.LoginRequest;
import org.junit.jupiter.api.Test;

import static org.junit.jupiter.api.Assertions.assertEquals;
import static org.junit.jupiter.api.Assertions.assertFalse;

class AuthServiceTest {
    @Test
    void loginSwitchesSharedClientFromApiKeyToBearer() {
        RecordingTransport transport = new RecordingTransport(request -> {
            if (request.uri().getPath().endsWith("/auth/login")) {
                return RecordingTransport.ok("{\"access_token\":\"access-2\",\"refresh_token\":\"refresh-2\",\"token_type\":\"Bearer\",\"expires_in\":900}");
            }
            return RecordingTransport.ok("{\"count\":0,\"users\":[]}");
        });
        HostMKClient client = HostMKClient.builder()
                .baseUrl("https://security.example.com")
                .apiKey("bootstrap-key")
                .transport(transport)
                .build();

        client.auth().login(new LoginRequest("acme", "owner@example.com", "secret", null));
        client.users().list();

        assertEquals("Bearer access-2", transport.last().headers().get("Authorization"));
        assertFalse(transport.last().headers().containsKey("X-API-Key"));
    }
}
