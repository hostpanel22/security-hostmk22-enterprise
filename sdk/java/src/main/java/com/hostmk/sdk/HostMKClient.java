package com.hostmk.sdk;

import com.hostmk.sdk.auth.CredentialStore;
import com.hostmk.sdk.config.ClientConfig;
import com.hostmk.sdk.exceptions.ConfigurationException;
import com.hostmk.sdk.http.HostMKHttpClient;
import com.hostmk.sdk.http.HttpTransport;
import com.hostmk.sdk.http.JdkHttpTransport;
import com.hostmk.sdk.services.ApiKeysService;
import com.hostmk.sdk.services.AuditService;
import com.hostmk.sdk.services.AuthService;
import com.hostmk.sdk.services.FraudService;
import com.hostmk.sdk.services.MonitoringService;
import com.hostmk.sdk.services.PermissionsService;
import com.hostmk.sdk.services.RolesService;
import com.hostmk.sdk.services.SecretsService;
import com.hostmk.sdk.services.SessionsService;
import com.hostmk.sdk.services.SettingsService;
import com.hostmk.sdk.services.TenantsService;
import com.hostmk.sdk.services.UsersService;

import java.net.URI;
import java.time.Duration;

/**
 * Entry point for the HostMK Security Java SDK.
 *
 * <pre>{@code
 * HostMKClient client = HostMKClient.builder()
 *     .baseUrl("https://security.example.com")
 *     .token(accessToken)
 *     .build();
 * client.users().list();
 * }</pre>
 */
public final class HostMKClient {
    private final AuthService auth;
    private final UsersService users;
    private final TenantsService tenants;
    private final SessionsService sessions;
    private final RolesService roles;
    private final PermissionsService permissions;
    private final ApiKeysService apiKeys;
    private final SecretsService secrets;
    private final AuditService audit;
    private final FraudService fraud;
    private final SettingsService settings;
    private final MonitoringService monitoring;

    private HostMKClient(Builder builder) {
        if (builder.baseUrl == null || builder.baseUrl.isBlank()) {
            throw new ConfigurationException("baseUrl is required");
        }
        ClientConfig config = new ClientConfig(URI.create(builder.baseUrl.trim()), builder.timeout, builder.userAgent);
        CredentialStore credentials = new CredentialStore(builder.apiKey, builder.token);
        HostMKHttpClient http = new HostMKHttpClient(config, credentials, builder.transport);
        auth = new AuthService(http);
        users = new UsersService(http);
        tenants = new TenantsService(http);
        sessions = new SessionsService(http);
        roles = new RolesService(http);
        permissions = new PermissionsService(http);
        apiKeys = new ApiKeysService(http);
        secrets = new SecretsService(http);
        audit = new AuditService(http);
        fraud = new FraudService(http);
        settings = new SettingsService(http);
        monitoring = new MonitoringService(http);
    }

    public static Builder builder() {
        return new Builder();
    }

    public AuthService auth() { return auth; }
    public UsersService users() { return users; }
    public TenantsService tenants() { return tenants; }
    public SessionsService sessions() { return sessions; }
    public RolesService roles() { return roles; }
    public PermissionsService permissions() { return permissions; }
    public ApiKeysService apiKeys() { return apiKeys; }
    public SecretsService secrets() { return secrets; }
    public AuditService audit() { return audit; }
    public FraudService fraud() { return fraud; }
    public SettingsService settings() { return settings; }
    public MonitoringService monitoring() { return monitoring; }

    /** Fluent SDK configuration builder. */
    public static final class Builder {
        private String baseUrl;
        private String apiKey = "";
        private String token = "";
        private Duration timeout = Duration.ofSeconds(30);
        private String userAgent = "hostmk-java-sdk/1.0";
        private HttpTransport transport = new JdkHttpTransport();

        private Builder() {
        }

        public Builder baseUrl(String baseUrl) {
            this.baseUrl = baseUrl;
            return this;
        }

        public Builder apiKey(String apiKey) {
            this.apiKey = apiKey;
            return this;
        }

        public Builder token(String token) {
            this.token = token;
            return this;
        }

        public Builder timeout(Duration timeout) {
            this.timeout = timeout;
            return this;
        }

        public Builder userAgent(String userAgent) {
            this.userAgent = userAgent;
            return this;
        }

        /** Overrides the transport, primarily for controlled environments and tests. */
        public Builder transport(HttpTransport transport) {
            if (transport == null) {
                throw new ConfigurationException("transport is required");
            }
            this.transport = transport;
            return this;
        }

        public HostMKClient build() {
            return new HostMKClient(this);
        }
    }
}
