"use strict";
var __defProp = Object.defineProperty;
var __getOwnPropDesc = Object.getOwnPropertyDescriptor;
var __getOwnPropNames = Object.getOwnPropertyNames;
var __hasOwnProp = Object.prototype.hasOwnProperty;
var __export = (target, all) => {
  for (var name in all)
    __defProp(target, name, { get: all[name], enumerable: true });
};
var __copyProps = (to, from, except, desc) => {
  if (from && typeof from === "object" || typeof from === "function") {
    for (let key of __getOwnPropNames(from))
      if (!__hasOwnProp.call(to, key) && key !== except)
        __defProp(to, key, { get: () => from[key], enumerable: !(desc = __getOwnPropDesc(from, key)) || desc.enumerable });
  }
  return to;
};
var __toCommonJS = (mod) => __copyProps(__defProp({}, "__esModule", { value: true }), mod);

// src/index.ts
var index_exports = {};
__export(index_exports, {
  AuditService: () => AuditService,
  AuthAPI: () => AuthAPI,
  FraudService: () => FraudService,
  HostMKClient: () => HostMKClient,
  HostMKError: () => HostMKError,
  PermissionsAPI: () => PermissionsAPI,
  TenancyAPI: () => TenancyAPI
});
module.exports = __toCommonJS(index_exports);

// src/client.ts
var HostMKError = class extends Error {
  status;
  data;
  constructor(message, status, data) {
    super(message);
    this.name = "HostMKError";
    this.status = status;
    this.data = data;
  }
};
var HostMKClient = class {
  apiKey;
  baseUrl;
  timeout;
  constructor(options) {
    if (!options.apiKey) {
      throw new Error(
        "apiKey is required"
      );
    }
    if (!options.baseUrl) {
      throw new Error(
        "baseUrl is required"
      );
    }
    this.apiKey = options.apiKey;
    this.baseUrl = options.baseUrl.replace(
      /\/$/,
      ""
    );
    this.timeout = options.timeout ?? 3e4;
  }
  async get(path) {
    return this.request(
      path,
      {
        method: "GET"
      }
    );
  }
  async post(path, body) {
    return this.request(
      path,
      {
        method: "POST",
        body
      }
    );
  }
  async put(path, body) {
    return this.request(
      path,
      {
        method: "PUT",
        body
      }
    );
  }
  async delete(path) {
    return this.request(
      path,
      {
        method: "DELETE"
      }
    );
  }
  async request(path, options) {
    const controller = new AbortController();
    const timeout = setTimeout(
      () => controller.abort(),
      this.timeout
    );
    try {
      const response = await fetch(
        `${this.baseUrl}${path}`,
        {
          method: options.method ?? "GET",
          signal: controller.signal,
          headers: {
            Authorization: `Bearer ${this.apiKey}`,
            "Content-Type": "application/json",
            ...options.headers
          },
          body: options.body !== void 0 ? JSON.stringify(
            options.body
          ) : void 0
        }
      );
      let data;
      const contentType = response.headers.get(
        "content-type"
      ) ?? "";
      if (contentType.includes(
        "application/json"
      )) {
        data = await response.json();
      } else {
        data = await response.text();
      }
      if (!response.ok) {
        throw new HostMKError(
          response.statusText,
          response.status,
          data
        );
      }
      return data;
    } finally {
      clearTimeout(
        timeout
      );
    }
  }
};

// src/auth.ts
var AuthAPI = class {
  client;
  constructor(client) {
    this.client = client;
  }
  async login(payload) {
    return this.client.post(
      "/auth/login",
      payload
    );
  }
  async refreshToken(payload) {
    return this.client.post(
      "/auth/refresh",
      payload
    );
  }
  async validateToken(token) {
    return this.client.post(
      "/auth/validate",
      {
        token
      }
    );
  }
  async logout() {
    return this.client.post(
      "/auth/logout"
    );
  }
  async me() {
    return this.client.get(
      "/auth/me"
    );
  }
};

// src/permissions.ts
var PermissionsAPI = class {
  client;
  constructor(client) {
    this.client = client;
  }
  /*
   * Permissions
   */
  async createPermission(payload) {
    return this.client.post(
      "/rbac/permissions",
      payload
    );
  }
  async getPermission(name) {
    return this.client.get(
      `/rbac/permissions/${encodeURIComponent(
        name
      )}`
    );
  }
  async listPermissions() {
    return this.client.get(
      "/rbac/permissions"
    );
  }
  async deletePermission(name) {
    return this.client.delete(
      `/rbac/permissions/${encodeURIComponent(
        name
      )}`
    );
  }
  /*
   * Roles
   */
  async createRole(payload) {
    return this.client.post(
      "/rbac/roles",
      payload
    );
  }
  async getRole(name) {
    return this.client.get(
      `/rbac/roles/${encodeURIComponent(
        name
      )}`
    );
  }
  async listRoles() {
    return this.client.get(
      "/rbac/roles"
    );
  }
  async deleteRole(name) {
    return this.client.delete(
      `/rbac/roles/${encodeURIComponent(
        name
      )}`
    );
  }
  /*
   * Role ↔ Permission
   */
  async assignPermissionToRole(roleName, permissionName) {
    return this.client.post(
      `/rbac/roles/${encodeURIComponent(
        roleName
      )}/permissions`,
      {
        permissionName
      }
    );
  }
  async removePermissionFromRole(roleName, permissionName) {
    return this.client.delete(
      `/rbac/roles/${encodeURIComponent(
        roleName
      )}/permissions/${encodeURIComponent(
        permissionName
      )}`
    );
  }
  async roleHasPermission(roleName, permissionName) {
    return this.client.get(
      `/rbac/roles/${encodeURIComponent(
        roleName
      )}/permissions/${encodeURIComponent(
        permissionName
      )}`
    );
  }
  /*
   * Parent Roles
   */
  async assignParentRole(roleName, parentRole) {
    return this.client.post(
      `/rbac/roles/${encodeURIComponent(
        roleName
      )}/parents`,
      {
        parentRole
      }
    );
  }
  async removeParentRole(roleName, parentRole) {
    return this.client.delete(
      `/rbac/roles/${encodeURIComponent(
        roleName
      )}/parents/${encodeURIComponent(
        parentRole
      )}`
    );
  }
  /*
   * Policies
   */
  async createPolicy(payload) {
    return this.client.post(
      "/rbac/policies",
      payload
    );
  }
  async getPolicy(id) {
    return this.client.get(
      `/rbac/policies/${encodeURIComponent(
        id
      )}`
    );
  }
  async listPolicies() {
    return this.client.get(
      "/rbac/policies"
    );
  }
  async deletePolicy(id) {
    return this.client.delete(
      `/rbac/policies/${encodeURIComponent(
        id
      )}`
    );
  }
  /*
   * Access Check
   */
  async checkAccess(request) {
    return this.client.post(
      "/rbac/check-access",
      request
    );
  }
};

// src/audit.ts
var AuditService = class {
  client;
  constructor(client) {
    this.client = client;
  }
  async list(query) {
    const search = new URLSearchParams();
    if (query?.tenantId) {
      search.set(
        "tenantId",
        query.tenantId
      );
    }
    if (query?.actorId) {
      search.set(
        "actorId",
        query.actorId
      );
    }
    if (query?.resource) {
      search.set(
        "resource",
        query.resource
      );
    }
    if (query?.action) {
      search.set(
        "action",
        query.action
      );
    }
    if (query?.page !== void 0) {
      search.set(
        "page",
        String(query.page)
      );
    }
    if (query?.limit !== void 0) {
      search.set(
        "limit",
        String(query.limit)
      );
    }
    const qs = search.toString();
    return this.client.get(
      qs.length > 0 ? `/audit?${qs}` : "/audit"
    );
  }
  async get(id) {
    if (!id.trim()) {
      throw new Error(
        "audit id is required"
      );
    }
    return this.client.get(
      `/audit/${encodeURIComponent(
        id
      )}`
    );
  }
  async create(event) {
    return this.client.post(
      "/audit",
      event
    );
  }
  async delete(id) {
    if (!id.trim()) {
      throw new Error(
        "audit id is required"
      );
    }
    return this.client.delete(
      `/audit/${encodeURIComponent(
        id
      )}`
    );
  }
};

// src/fraud.ts
var FraudService = class {
  client;
  constructor(client) {
    this.client = client;
  }
  async list(query) {
    const params = new URLSearchParams();
    if (query?.tenantId) {
      params.set(
        "tenantId",
        query.tenantId
      );
    }
    if (query?.userId) {
      params.set(
        "userId",
        query.userId
      );
    }
    if (query?.minRiskScore !== void 0) {
      params.set(
        "minRiskScore",
        String(
          query.minRiskScore
        )
      );
    }
    if (query?.page !== void 0) {
      params.set(
        "page",
        String(query.page)
      );
    }
    if (query?.limit !== void 0) {
      params.set(
        "limit",
        String(query.limit)
      );
    }
    const queryString = params.toString();
    return this.client.get(
      queryString.length > 0 ? `/fraud?${queryString}` : "/fraud"
    );
  }
  async get(id) {
    if (!id.trim()) {
      throw new Error(
        "fraud alert id is required"
      );
    }
    return this.client.get(
      `/fraud/${encodeURIComponent(
        id
      )}`
    );
  }
  async create(alert) {
    return this.client.post(
      "/fraud",
      alert
    );
  }
  async delete(id) {
    if (!id.trim()) {
      throw new Error(
        "fraud alert id is required"
      );
    }
    return this.client.delete(
      `/fraud/${encodeURIComponent(
        id
      )}`
    );
  }
};

// src/tenancy.ts
var TenancyAPI = class {
  client;
  constructor(client) {
    this.client = client;
  }
  async list(query) {
    const params = new URLSearchParams();
    if (query?.page !== void 0) {
      params.set(
        "page",
        String(query.page)
      );
    }
    if (query?.limit !== void 0) {
      params.set(
        "limit",
        String(query.limit)
      );
    }
    if (query?.search) {
      params.set(
        "search",
        query.search
      );
    }
    if (query?.active !== void 0) {
      params.set(
        "active",
        String(query.active)
      );
    }
    const queryString = params.toString();
    return this.client.get(
      queryString.length > 0 ? `/tenants?${queryString}` : "/tenants"
    );
  }
  async get(tenantId) {
    return this.client.get(
      `/tenants/${tenantId}`
    );
  }
  async create(payload) {
    return this.client.post(
      "/tenants",
      payload
    );
  }
  async update(tenantId, payload) {
    return this.client.put(
      `/tenants/${tenantId}`,
      payload
    );
  }
  async delete(tenantId) {
    return this.client.delete(
      `/tenants/${tenantId}`
    );
  }
  async activate(tenantId) {
    return this.client.post(
      `/tenants/${tenantId}/activate`
    );
  }
  async deactivate(tenantId) {
    return this.client.post(
      `/tenants/${tenantId}/deactivate`
    );
  }
};
// Annotate the CommonJS export names for ESM import in node:
0 && (module.exports = {
  AuditService,
  AuthAPI,
  FraudService,
  HostMKClient,
  HostMKError,
  PermissionsAPI,
  TenancyAPI
});
