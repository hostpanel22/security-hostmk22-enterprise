package com.hostmk.sdk.models.audit;

import java.util.Map;

/** Audit event creation payload. */
public record CreateAuditEventRequest(
        String action,
        String resource,
        String severity,
        String status,
        String message,
        Map<String, Object> metadata
) {
}
