package com.hostmk.sdk.models.auth;

import com.hostmk.sdk.models.tenant.Tenant;
import com.hostmk.sdk.models.user.User;

/** Tenant and owner created during first-run bootstrap. */
public record BootstrapResult(Tenant tenant, User user) {
}
