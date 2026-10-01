package com.hostmk.sdk.exceptions;

/** Base unchecked exception for SDK failures. */
public class HostMKException extends RuntimeException {
    public HostMKException(String message) {
        super(message);
    }

    public HostMKException(String message, Throwable cause) {
        super(message, cause);
    }
}
