package com.hostmk.sdk.config;

import com.hostmk.sdk.exceptions.ConfigurationException;

import java.net.URI;
import java.time.Duration;

/** Immutable validated client configuration. */
public record ClientConfig(URI baseUrl, Duration timeout, String userAgent) {
    public ClientConfig {
        if (baseUrl == null || baseUrl.getHost() == null
                || !("http".equalsIgnoreCase(baseUrl.getScheme()) || "https".equalsIgnoreCase(baseUrl.getScheme()))
                || baseUrl.getUserInfo() != null || baseUrl.getQuery() != null || baseUrl.getFragment() != null) {
            throw new ConfigurationException("baseUrl must be an absolute HTTP(S) URL without credentials, query, or fragment");
        }
        if (timeout == null || timeout.isZero() || timeout.isNegative()) {
            throw new ConfigurationException("timeout must be greater than zero");
        }
        if (userAgent == null || userAgent.isBlank()) {
            throw new ConfigurationException("userAgent is required");
        }
        String normalized = baseUrl.toString().replaceFirst("/+$", "");
        baseUrl = URI.create(normalized);
        userAgent = userAgent.trim();
    }
}
