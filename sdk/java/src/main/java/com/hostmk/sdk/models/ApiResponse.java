package com.hostmk.sdk.models;

import com.fasterxml.jackson.annotation.JsonIgnoreProperties;
import com.fasterxml.jackson.annotation.JsonProperty;

/** Standard HostMK response envelope. */
@JsonIgnoreProperties(ignoreUnknown = true)
public record ApiResponse<T>(
        boolean success,
        String message,
        T data,
        String error,
        String timestamp,
        @JsonProperty("status") Integer status
) {
}
