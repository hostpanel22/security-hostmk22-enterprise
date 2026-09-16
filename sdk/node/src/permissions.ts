import { HostMKClient } from "./client.js";
import type {
  ApiResponse,
  PaginatedResponse,
  Permission,
  Policy,
  Role,
  CreatePermissionRequest,
  CreatePolicyRequest,
  CreateRoleRequest,
} from "./types.js";

export interface AccessCheckRequest {
  roleName: string;
  permissionName: string;
}

export interface AccessCheckResponse {
  allowed: boolean;
}

export class PermissionsAPI {
  private readonly client: HostMKClient;

  constructor(
    client: HostMKClient,
  ) {
    this.client = client;
  }

  /*
   * Permissions
   */

  public async createPermission(
    payload: CreatePermissionRequest,
  ): Promise<ApiResponse<Permission>> {
    return this.client.post<
      ApiResponse<Permission>
    >(
      "/rbac/permissions",
      payload,
    );
  }

  public async getPermission(
    name: string,
  ): Promise<ApiResponse<Permission>> {
    return this.client.get<
      ApiResponse<Permission>
    >(
      `/rbac/permissions/${encodeURIComponent(
        name,
      )}`,
    );
  }

  public async listPermissions(): Promise<
    PaginatedResponse<Permission>
  > {
    return this.client.get<
      PaginatedResponse<Permission>
    >(
      "/rbac/permissions",
    );
  }

  public async deletePermission(
    name: string,
  ): Promise<ApiResponse<boolean>> {
    return this.client.delete<
      ApiResponse<boolean>
    >(
      `/rbac/permissions/${encodeURIComponent(
        name,
      )}`,
    );
  }

  /*
   * Roles
   */

  public async createRole(
    payload: CreateRoleRequest,
  ): Promise<ApiResponse<Role>> {
    return this.client.post<
      ApiResponse<Role>
    >(
      "/rbac/roles",
      payload,
    );
  }

  public async getRole(
    name: string,
  ): Promise<ApiResponse<Role>> {
    return this.client.get<
      ApiResponse<Role>
    >(
      `/rbac/roles/${encodeURIComponent(
        name,
      )}`,
    );
  }

  public async listRoles(): Promise<
    PaginatedResponse<Role>
  > {
    return this.client.get<
      PaginatedResponse<Role>
    >(
      "/rbac/roles",
    );
  }

  public async deleteRole(
    name: string,
  ): Promise<ApiResponse<boolean>> {
    return this.client.delete<
      ApiResponse<boolean>
    >(
      `/rbac/roles/${encodeURIComponent(
        name,
      )}`,
    );
  }

  /*
   * Role ↔ Permission
   */

  public async assignPermissionToRole(
    roleName: string,
    permissionName: string,
  ): Promise<ApiResponse<boolean>> {
    return this.client.post<
      ApiResponse<boolean>
    >(
      `/rbac/roles/${encodeURIComponent(
        roleName,
      )}/permissions`,
      {
        permissionName,
      },
    );
  }

  public async removePermissionFromRole(
    roleName: string,
    permissionName: string,
  ): Promise<ApiResponse<boolean>> {
    return this.client.delete<
      ApiResponse<boolean>
    >(
      `/rbac/roles/${encodeURIComponent(
        roleName,
      )}/permissions/${encodeURIComponent(
        permissionName,
      )}`,
    );
  }

  public async roleHasPermission(
    roleName: string,
    permissionName: string,
  ): Promise<
    ApiResponse<AccessCheckResponse>
  > {
    return this.client.get<
      ApiResponse<AccessCheckResponse>
    >(
      `/rbac/roles/${encodeURIComponent(
        roleName,
      )}/permissions/${encodeURIComponent(
        permissionName,
      )}`,
    );
  }

  /*
   * Parent Roles
   */

  public async assignParentRole(
    roleName: string,
    parentRole: string,
  ): Promise<ApiResponse<boolean>> {
    return this.client.post<
      ApiResponse<boolean>
    >(
      `/rbac/roles/${encodeURIComponent(
        roleName,
      )}/parents`,
      {
        parentRole,
      },
    );
  }

  public async removeParentRole(
    roleName: string,
    parentRole: string,
  ): Promise<ApiResponse<boolean>> {
    return this.client.delete<
      ApiResponse<boolean>
    >(
      `/rbac/roles/${encodeURIComponent(
        roleName,
      )}/parents/${encodeURIComponent(
        parentRole,
      )}`,
    );
  }

  /*
   * Policies
   */

  public async createPolicy(
    payload: CreatePolicyRequest,
  ): Promise<ApiResponse<Policy>> {
    return this.client.post<
      ApiResponse<Policy>
    >(
      "/rbac/policies",
      payload,
    );
  }

  public async getPolicy(
    id: string,
  ): Promise<ApiResponse<Policy>> {
    return this.client.get<
      ApiResponse<Policy>
    >(
      `/rbac/policies/${encodeURIComponent(
        id,
      )}`,
    );
  }

  public async listPolicies(): Promise<
    PaginatedResponse<Policy>
  > {
    return this.client.get<
      PaginatedResponse<Policy>
    >(
      "/rbac/policies",
    );
  }

  public async deletePolicy(
    id: string,
  ): Promise<ApiResponse<boolean>> {
    return this.client.delete<
      ApiResponse<boolean>
    >(
      `/rbac/policies/${encodeURIComponent(
        id,
      )}`,
    );
  }

  /*
   * Access Check
   */

  public async checkAccess(
    request: AccessCheckRequest,
  ): Promise<
    ApiResponse<AccessCheckResponse>
  > {
    return this.client.post<
      ApiResponse<AccessCheckResponse>
    >(
      "/rbac/check-access",
      request,
    );
  }
}