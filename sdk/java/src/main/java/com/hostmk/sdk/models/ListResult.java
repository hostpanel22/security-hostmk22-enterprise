package com.hostmk.sdk.models;

import java.util.List;

/** Typed list result used by collection services. */
public record ListResult<T>(int count, List<T> items) {
    public ListResult {
        items = items == null ? List.of() : List.copyOf(items);
    }
}
