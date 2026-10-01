package com.hostmk.sdk.http;

import java.net.URI;
import java.time.Duration;
import java.util.Map;

/** Pluggable transport contract used by the SDK and tests. */
@FunctionalInterface
public interface HttpTransport {
    HttpResult send(String method, URI uri, Map<String, String> headers, String body, Duration timeout);
}
