package com.hostmk.sdk.exceptions;

/** HTTP error returned by the HostMK API. */
public final class ApiException extends HostMKException {
    private final int statusCode;
    private final String responseBody;

    public ApiException(int statusCode, String message, String responseBody) {
        super(message);
        this.statusCode = statusCode;
        this.responseBody = responseBody;
    }

    public int statusCode() {
        return statusCode;
    }

    public String responseBody() {
        return responseBody;
    }
}
