package com.hostmk.sdk.auth;

/** Thread-safe mutable credentials shared by all services. */
public final class CredentialStore {
    private volatile String apiKey;
    private volatile String bearerToken;

    public CredentialStore(String apiKey, String bearerToken) {
        this.apiKey = normalize(apiKey);
        this.bearerToken = normalize(bearerToken);
    }

    public String apiKey() {
        return apiKey;
    }

    public String bearerToken() {
        return bearerToken;
    }

    /** Uses a bearer token and clears the API key, matching the Go SDK behavior. */
    public synchronized void useBearerToken(String token) {
        apiKey = "";
        bearerToken = normalize(token);
    }

    public synchronized void useApiKey(String key) {
        apiKey = normalize(key);
        bearerToken = "";
    }

    public synchronized void clear() {
        apiKey = "";
        bearerToken = "";
    }

    private static String normalize(String value) {
        return value == null ? "" : value.trim();
    }
}
