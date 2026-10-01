package com.hostmk.sdk.services;

import com.hostmk.sdk.http.HostMKHttpClient;
import com.hostmk.sdk.models.monitoring.HealthStatus;
import com.hostmk.sdk.models.monitoring.VersionInfo;

import java.util.Map;

/** Public health, readiness, and version endpoints. */
public final class MonitoringService {
    private final HostMKHttpClient client;

    public MonitoringService(HostMKHttpClient client) {
        this.client = client;
    }

    public HealthStatus health() {
        return client.get("/health", Map.of(), HealthStatus.class);
    }

    public HealthStatus ready() {
        return client.get("/ready", Map.of(), HealthStatus.class);
    }

    public VersionInfo version() {
        return client.get("/api/v1/version", Map.of(), VersionInfo.class);
    }
}
