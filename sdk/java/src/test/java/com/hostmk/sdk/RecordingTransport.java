package com.hostmk.sdk;

import com.hostmk.sdk.http.HttpResult;
import com.hostmk.sdk.http.HttpTransport;

import java.net.URI;
import java.time.Duration;
import java.util.ArrayList;
import java.util.List;
import java.util.Map;
import java.util.function.Function;

final class RecordingTransport implements HttpTransport {
    record Request(String method, URI uri, Map<String, String> headers, String body) {
    }

    private final List<Request> requests = new ArrayList<>();
    private final Function<Request, HttpResult> responder;

    RecordingTransport(Function<Request, HttpResult> responder) {
        this.responder = responder;
    }

    @Override
    public HttpResult send(String method, URI uri, Map<String, String> headers, String body, Duration timeout) {
        Request request = new Request(method, uri, Map.copyOf(headers), body);
        requests.add(request);
        return responder.apply(request);
    }

    List<Request> requests() {
        return List.copyOf(requests);
    }

    Request last() {
        return requests.get(requests.size() - 1);
    }

    static HttpResult ok(String data) {
        return new HttpResult(200, "{\"success\":true,\"data\":" + data + "}", Map.of());
    }
}
