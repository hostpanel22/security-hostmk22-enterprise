package com.hostmk.sdk.models.monitoring;

import com.fasterxml.jackson.annotation.JsonIgnoreProperties;

import java.util.List;

/** Health/ready response. */
@JsonIgnoreProperties(ignoreUnknown = true)
public record HealthStatus(String status, String service, String version, String env, List<Component> components) {
    /** Health of one backing component. */
    @JsonIgnoreProperties(ignoreUnknown = true)
    public record Component(String name, String status, String message) {
    }
}
