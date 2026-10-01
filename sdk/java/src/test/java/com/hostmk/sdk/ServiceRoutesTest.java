package com.hostmk.sdk;

import com.hostmk.sdk.models.apikey.CreateApiKeyRequest;
import com.hostmk.sdk.models.audit.CreateAuditEventRequest;
import com.hostmk.sdk.models.auth.BootstrapRequest;
import com.hostmk.sdk.models.auth.LoginRequest;
import com.hostmk.sdk.models.secret.CreateSecretRequest;
import com.hostmk.sdk.models.user.CreateUserRequest;
import org.junit.jupiter.api.Test;

import java.util.List;
import java.util.Map;
import java.util.Set;
import java.util.stream.Collectors;

import static org.junit.jupiter.api.Assertions.assertTrue;

class ServiceRoutesTest {
    @Test
    void everyRegisteredDomainUsesCanonicalRoutes() {
        RecordingTransport transport = new RecordingTransport(request -> {
            String path = request.uri().getPath();
            if (path.endsWith("/auth/login") || path.endsWith("/auth/refresh")) {
                return RecordingTransport.ok("{\"access_token\":\"access\",\"refresh_token\":\"refresh\",\"token_type\":\"Bearer\",\"expires_in\":900}");
            }
            if (path.endsWith("/users")) return RecordingTransport.ok("{\"count\":0,\"users\":[]}");
            if (path.endsWith("/tenants")) return RecordingTransport.ok("{\"count\":0,\"tenants\":[]}");
            if (path.endsWith("/sessions")) return RecordingTransport.ok("{\"count\":0,\"sessions\":[]}");
            if (path.endsWith("/rbac/roles")) return RecordingTransport.ok("{\"count\":0,\"roles\":[]}");
            if (path.endsWith("/permissions")) return RecordingTransport.ok("{\"count\":0,\"permissions\":[]}");
            if (path.endsWith("/api-keys")) return RecordingTransport.ok("{\"count\":0,\"api_keys\":[]}");
            if (path.endsWith("/secrets")) return RecordingTransport.ok("{\"count\":0,\"secrets\":[]}");
            if (path.endsWith("/audit/events")) return RecordingTransport.ok("{\"count\":0,\"events\":[]}");
            if (path.endsWith("/fraud/rules")) return RecordingTransport.ok("{\"count\":0,\"rules\":[]}");
            if (path.endsWith("/fraud/violations")) return RecordingTransport.ok("{\"count\":0,\"events\":[]}");
            if (path.endsWith("/settings")) return RecordingTransport.ok("{\"count\":0,\"settings\":[]}");
            return RecordingTransport.ok("{}");
        });
        HostMKClient client = HostMKClient.builder()
                .baseUrl("https://security.example.com")
                .token("access")
                .transport(transport)
                .build();

            client.auth().bootstrap(new BootstrapRequest("Acme", "acme", "Owner", "owner@example.com", "secret"));
            client.auth().login(new LoginRequest("acme", "owner@example.com", "secret", null));
            client.auth().refresh("refresh");
        client.users().list();
        client.users().create(new CreateUserRequest("a@example.com", "A", "secret", "active", List.of()));
            client.users().action("user-1", "disable");
        client.tenants().list();
            client.tenants().activate("tenant-1");
        client.sessions().list();
        client.roles().list();
        client.permissions().list();
            client.permissions().check("user-1", "users.read");
        client.apiKeys().list();
        client.apiKeys().create(new CreateApiKeyRequest("ci", List.of("audit.read"), null, null));
        client.secrets().list();
        client.secrets().create(new CreateSecretRequest("DB", "custom", "value", "database"));
        client.audit().list(25);
        client.audit().create(new CreateAuditEventRequest("login", "auth", "info", "success", null, Map.of()));
        client.fraud().listRules();
        client.fraud().analyze(Map.of("amount", 42));
        client.fraud().listViolations();
        client.settings().list();
        client.monitoring().health();
        client.monitoring().ready();
        client.monitoring().version();

        Set<String> paths = transport.requests().stream().map(request -> request.uri().getPath()).collect(Collectors.toSet());
        assertTrue(paths.containsAll(Set.of(
            "/api/v1/auth/bootstrap", "/api/v1/auth/login", "/api/v1/auth/refresh",
                "/api/v1/users", "/api/v1/tenants", "/api/v1/sessions", "/api/v1/rbac/roles",
            "/api/v1/users/action", "/api/v1/tenants/action", "/api/v1/permissions/check",
                "/api/v1/permissions", "/api/v1/api-keys", "/api/v1/secrets", "/api/v1/audit/events",
                "/api/v1/fraud/rules", "/api/v1/fraud/analyze", "/api/v1/fraud/violations",
                "/api/v1/settings", "/health", "/ready", "/api/v1/version"
        )));
    }
}
