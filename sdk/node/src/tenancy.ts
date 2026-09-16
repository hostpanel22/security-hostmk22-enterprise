import { HostMKClient } from "./client.js";
import type {
  ApiResponse,
  PaginatedResponse,
  Tenant,
} from "./types.js";

export interface CreateTenantRequest {
  name: string;
  slug: string;
}

export interface UpdateTenantRequest {
  name?: string;
  slug?: string;
  active?: boolean;
}

export interface TenantQuery {
  page?: number;
  limit?: number;
  search?: string;
  active?: boolean;
}

export class TenancyAPI {
  private readonly client: HostMKClient;

  constructor(
    client: HostMKClient,
  ) {
    this.client = client;
  }

  public async list(
    query?: TenantQuery,
  ): Promise<
    PaginatedResponse<Tenant>
  > {
    const params =
      new URLSearchParams();

    if (
      query?.page !== undefined
    ) {
      params.set(
        "page",
        String(query.page),
      );
    }

    if (
      query?.limit !== undefined
    ) {
      params.set(
        "limit",
        String(query.limit),
      );
    }

    if (
      query?.search
    ) {
      params.set(
        "search",
        query.search,
      );
    }

    if (
      query?.active !== undefined
    ) {
      params.set(
        "active",
        String(query.active),
      );
    }

    const queryString =
      params.toString();

    return this.client.get<
      PaginatedResponse<Tenant>
    >(
      queryString.length > 0
        ? `/tenants?${queryString}`
        : "/tenants",
    );
  }

  public async get(
    tenantId: string,
  ): Promise<
    ApiResponse<Tenant>
  > {
    return this.client.get<
      ApiResponse<Tenant>
    >(
      `/tenants/${tenantId}`,
    );
  }

  public async create(
    payload: CreateTenantRequest,
  ): Promise<
    ApiResponse<Tenant>
  > {
    return this.client.post<
      ApiResponse<Tenant>
    >(
      "/tenants",
      payload,
    );
  }

  public async update(
    tenantId: string,
    payload: UpdateTenantRequest,
  ): Promise<
    ApiResponse<Tenant>
  > {
    return this.client.put<
      ApiResponse<Tenant>
    >(
      `/tenants/${tenantId}`,
      payload,
    );
  }

  public async delete(
    tenantId: string,
  ): Promise<
    ApiResponse<boolean>
  > {
    return this.client.delete<
      ApiResponse<boolean>
    >(
      `/tenants/${tenantId}`,
    );
  }

  public async activate(
    tenantId: string,
  ): Promise<
    ApiResponse<Tenant>
  > {
    return this.client.post<
      ApiResponse<Tenant>
    >(
      `/tenants/${tenantId}/activate`,
    );
  }

  public async deactivate(
    tenantId: string,
  ): Promise<
    ApiResponse<Tenant>
  > {
    return this.client.post<
      ApiResponse<Tenant>
    >(
      `/tenants/${tenantId}/deactivate`,
    );
  }
}