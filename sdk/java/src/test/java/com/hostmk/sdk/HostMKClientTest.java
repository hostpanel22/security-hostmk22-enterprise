package com.hostmk.sdk;

import com.hostmk.sdk.exceptions.ApiException;
import com.hostmk.sdk.exceptions.ConfigurationException;
import com.hostmk.sdk.http.HttpResult;
import org.junit.jupiter.api.Test;

import java.util.Map;

import static org.junit.jupiter.api.Assertions.assertEquals;
import static org.junit.jupiter.api.Assertions.assertThrows;

class HostMKClientTest {
    @Test
    void sendsApiKeyAndAvoidsDuplicatingApiPrefix() {
        RecordingTransport transport = new RecordingTransport(request -> RecordingTransport.ok("{\"count\":0,\"users\":[]}"));
        HostMKClient client = HostMKClient.builder()
                .baseUrl("https://security.example.com/api/v1")
                .apiKey("secret-key")
                .transport(transport)
                .build();

        client.users().list();

        assertEquals("https://security.example.com/api/v1/users", transport.last().uri().toString());
        assertEquals("secret-key", transport.last().headers().get("X-API-Key"));
    }

    @Test
    void rejectsInvalidConfiguration() {
        assertThrows(ConfigurationException.class, () -> HostMKClient.builder().baseUrl("ftp://example.com").build());
        assertThrows(ConfigurationException.class, () -> HostMKClient.builder().build());
    }

    @Test
    void exposesApiErrors() {
        RecordingTransport transport = new RecordingTransport(request ->
                new HttpResult(403, "{\"message\":\"missing permission\"}", Map.of()));
        HostMKClient client = HostMKClient.builder()
                .baseUrl("https://security.example.com")
                .token("token")
                .transport(transport)
                .build();

        ApiException exception = assertThrows(ApiException.class, () -> client.users().list());
        assertEquals(403, exception.statusCode());
        assertEquals("missing permission", exception.getMessage());
    }
}
