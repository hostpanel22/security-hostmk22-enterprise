package com.hostmk.sdk.http;

import java.util.List;
import java.util.Map;

/** Raw response returned by an HTTP transport. */
public record HttpResult(int statusCode, String body, Map<String, List<String>> headers) {
}
