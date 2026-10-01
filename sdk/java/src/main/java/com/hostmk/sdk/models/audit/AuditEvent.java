package com.hostmk.sdk.models.audit;

import com.fasterxml.jackson.annotation.JsonIgnoreProperties;
import com.fasterxml.jackson.annotation.JsonProperty;

import java.time.Instant;
import java.util.Map;

/** Immutable audit event returned by HostMK. */
@JsonIgnoreProperties(ignoreUnknown = true)
public record AuditEvent(
        String id,
        @JsonProperty("tenant_id") String tenantId,
        @JsonProperty("user_id") String userId,
        String action,
        String resource,
        String severity,
        String status,
        @JsonProperty("ip_address") String ipAddress,
        @JsonProperty("user_agent") String userAgent,
        String message,
        Map<String, Object> metadata,
        @JsonProperty("created_at") Instant createdAt,
        @JsonProperty("updated_at") Instant updatedAt
) {
}
