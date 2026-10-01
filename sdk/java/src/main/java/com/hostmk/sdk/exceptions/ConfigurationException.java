package com.hostmk.sdk.exceptions;

/** Raised when client configuration is invalid. */
public final class ConfigurationException extends HostMKException {
    public ConfigurationException(String message) {
        super(message);
    }
}
