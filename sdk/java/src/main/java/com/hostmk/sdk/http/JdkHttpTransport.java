package com.hostmk.sdk.http;

import com.hostmk.sdk.exceptions.HostMKException;

import java.io.IOException;
import java.net.URI;
import java.net.http.HttpClient;
import java.net.http.HttpRequest;
import java.net.http.HttpResponse;
import java.time.Duration;
import java.util.Map;

/** Java 17 {@link HttpClient}-based transport. */
public final class JdkHttpTransport implements HttpTransport {
    private final HttpClient client;

    public JdkHttpTransport() {
        this(HttpClient.newBuilder().followRedirects(HttpClient.Redirect.NORMAL).build());
    }

    public JdkHttpTransport(HttpClient client) {
        this.client = client;
    }

    @Override
    public HttpResult send(String method, URI uri, Map<String, String> headers, String body, Duration timeout) {
        HttpRequest.Builder builder = HttpRequest.newBuilder(uri).timeout(timeout);
        headers.forEach(builder::header);
        HttpRequest.BodyPublisher publisher = body == null
                ? HttpRequest.BodyPublishers.noBody()
                : HttpRequest.BodyPublishers.ofString(body);
        builder.method(method, publisher);
        try {
            HttpResponse<String> response = client.send(builder.build(), HttpResponse.BodyHandlers.ofString());
            return new HttpResult(response.statusCode(), response.body(), response.headers().map());
        } catch (InterruptedException exception) {
            Thread.currentThread().interrupt();
            throw new HostMKException("HTTP request interrupted", exception);
        } catch (IOException exception) {
            throw new HostMKException("HTTP request failed", exception);
        }
    }
}
