package com.hostmk.sdk.services;

import com.hostmk.sdk.http.HostMKHttpClient;
import com.hostmk.sdk.models.ListResult;
import com.hostmk.sdk.models.setting.Setting;

import java.util.Map;

/** Tenant and platform settings. */
public final class SettingsService {
    private final HostMKHttpClient client;

    public SettingsService(HostMKHttpClient client) {
        this.client = client;
    }

    public ListResult<Setting> list() {
        return client.list("/api/v1/settings", Map.of(), "settings", Setting.class);
    }

    public Setting save(Setting setting) {
        return client.put("/api/v1/settings", Map.of(), setting, Setting.class);
    }
}
