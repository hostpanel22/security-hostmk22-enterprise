package com.hostmk.sdk.http;

import com.fasterxml.jackson.core.JsonProcessingException;
import com.fasterxml.jackson.core.type.TypeReference;
import com.fasterxml.jackson.databind.JavaType;
import com.fasterxml.jackson.databind.JsonNode;
import com.fasterxml.jackson.databind.ObjectMapper;
import com.fasterxml.jackson.datatype.jsr310.JavaTimeModule;
import com.hostmk.sdk.auth.CredentialStore;
import com.hostmk.sdk.config.ClientConfig;
import com.hostmk.sdk.exceptions.ApiException;
import com.hostmk.sdk.exceptions.HostMKException;
import com.hostmk.sdk.models.ListResult;

import java.net.URI;
import java.net.URLEncoder;
import java.nio.charset.StandardCharsets;
import java.util.ArrayList;
import java.util.LinkedHashMap;
import java.util.List;
import java.util.Map;

/** Shared JSON HTTP client used by every domain service. */
public final class HostMKHttpClient {
    private final ClientConfig config;
    private final CredentialStore credentials;
    private final HttpTransport transport;
    private final ObjectMapper mapper;

    public HostMKHttpClient(ClientConfig config, CredentialStore credentials, HttpTransport transport) {
        this.config = config;
        this.credentials = credentials;
        this.transport = transport;
        this.mapper = new ObjectMapper().registerModule(new JavaTimeModule());
    }

    public CredentialStore credentials() {
        return credentials;
    }

    public <T> T get(String path, Map<String, ?> query, Class<T> type) {
        return request("GET", path, query, null, mapper.constructType(type));
    }

    public <T> T post(String path, Map<String, ?> query, Object body, Class<T> type) {
        return request("POST", path, query, body, mapper.constructType(type));
    }

    public <T> T put(String path, Map<String, ?> query, Object body, Class<T> type) {
        return request("PUT", path, query, body, mapper.constructType(type));
    }

    public void post(String path, Map<String, ?> query, Object body) {
        requestNode("POST", path, query, body);
    }

    public void delete(String path, Map<String, ?> query) {
        requestNode("DELETE", path, query, null);
    }

    public <T> ListResult<T> list(String path, Map<String, ?> query, String key, Class<T> itemType) {
        JsonNode data = requestNode("GET", path, query, null);
        JsonNode items = data;
        int count = 0;
        if (data != null && data.isObject()) {
            JsonNode countNode = data.get("count");
            count = countNode == null ? 0 : countNode.asInt();
            JsonNode keyed = data.get(key);
            if (keyed != null) {
                items = keyed;
            }
        }
        List<T> result = new ArrayList<>();
        if (items != null && items.isArray()) {
            JavaType type = mapper.getTypeFactory().constructCollectionType(List.class, itemType);
            result = mapper.convertValue(items, type);
        }
        return new ListResult<>(count == 0 ? result.size() : count, result);
    }

    public <T> T request(String method, String path, Map<String, ?> query, Object body, TypeReference<T> type) {
        return request(method, path, query, body, mapper.getTypeFactory().constructType(type));
    }

    private <T> T request(String method, String path, Map<String, ?> query, Object body, JavaType type) {
        JsonNode data = requestNode(method, path, query, body);
        if (type == null) {
            return null;
        }
        return mapper.convertValue(data, type);
    }

    private JsonNode requestNode(String method, String path, Map<String, ?> query, Object body) {
        String encodedBody = null;
        if (body != null) {
            try {
                encodedBody = mapper.writeValueAsString(body);
            } catch (JsonProcessingException exception) {
                throw new HostMKException("Unable to encode request body", exception);
            }
        }
        HttpResult result = transport.send(method, buildUri(path, query), headers(body != null), encodedBody, config.timeout());
        JsonNode root = null;
        if (result.body() != null && !result.body().isBlank()) {
            try {
                root = mapper.readTree(result.body());
            } catch (JsonProcessingException exception) {
                throw new HostMKException("Invalid JSON response received", exception);
            }
        }
        if (result.statusCode() < 200 || result.statusCode() >= 300) {
            String message = "HostMK API returned HTTP " + result.statusCode();
            if (root != null) {
                if (root.hasNonNull("message")) {
                    message = root.get("message").asText();
                } else if (root.hasNonNull("error")) {
                    message = root.get("error").asText();
                }
            }
            throw new ApiException(result.statusCode(), message, result.body());
        }
        if (root == null) {
            return mapper.nullNode();
        }
        JsonNode data = root.get("data");
        return data == null ? root : data;
    }

    private Map<String, String> headers(boolean hasBody) {
        Map<String, String> headers = new LinkedHashMap<>();
        headers.put("Accept", "application/json");
        headers.put("User-Agent", config.userAgent());
        if (hasBody) {
            headers.put("Content-Type", "application/json");
        }
        if (!credentials.apiKey().isBlank()) {
            headers.put("X-API-Key", credentials.apiKey());
        } else if (!credentials.bearerToken().isBlank()) {
            headers.put("Authorization", "Bearer " + credentials.bearerToken());
        }
        return headers;
    }

    private URI buildUri(String path, Map<String, ?> query) {
        if (path == null || !path.startsWith("/") || path.startsWith("//")) {
            throw new HostMKException("API path must be absolute");
        }
        String base = config.baseUrl().toString();
        String requestPath = path;
        if (base.endsWith("/api/v1") && (path.equals("/api/v1") || path.startsWith("/api/v1/"))) {
            requestPath = path.substring("/api/v1".length());
        }
        StringBuilder target = new StringBuilder(base).append(requestPath);
        if (query != null && !query.isEmpty()) {
            boolean first = true;
            for (Map.Entry<String, ?> entry : query.entrySet()) {
                if (entry.getValue() == null) {
                    continue;
                }
                target.append(first ? '?' : '&');
                first = false;
                target.append(encode(entry.getKey())).append('=').append(encode(String.valueOf(entry.getValue())));
            }
        }
        return URI.create(target.toString());
    }

    private static String encode(String value) {
        return URLEncoder.encode(value, StandardCharsets.UTF_8).replace("+", "%20");
    }
}
