package com.hostmk.sdk.models.monitoring;

import com.fasterxml.jackson.annotation.JsonIgnoreProperties;

/** Server version information. */
@JsonIgnoreProperties(ignoreUnknown = true)
public record VersionInfo(String name, String version, String env) {
}
