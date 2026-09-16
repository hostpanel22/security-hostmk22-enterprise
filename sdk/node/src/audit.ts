import { HostMKClient } from "./client.js";

import type {
  ApiResponse,
  AuditEvent,
  AuditQuery,
  PaginatedResponse,
} from "./types.js";

export class AuditService {
  private readonly client: HostMKClient;

  constructor(
    client: HostMKClient,
  ) {
    this.client = client;
  }

  public async list(
    query?: AuditQuery,
  ): Promise<
    PaginatedResponse<AuditEvent>
  > {
    const search =
      new URLSearchParams();

    if (query?.tenantId) {
      search.set(
        "tenantId",
        query.tenantId,
      );
    }

    if (query?.actorId) {
      search.set(
        "actorId",
        query.actorId,
      );
    }

    if (query?.resource) {
      search.set(
        "resource",
        query.resource,
      );
    }

    if (query?.action) {
      search.set(
        "action",
        query.action,
      );
    }

    if (
      query?.page !== undefined
    ) {
      search.set(
        "page",
        String(query.page),
      );
    }

    if (
      query?.limit !== undefined
    ) {
      search.set(
        "limit",
        String(query.limit),
      );
    }

    const qs =
      search.toString();

    return this.client.get<
      PaginatedResponse<AuditEvent>
    >(
      qs.length > 0
        ? `/audit?${qs}`
        : "/audit",
    );
  }

  public async get(
    id: string,
  ): Promise<
    ApiResponse<AuditEvent>
  > {
    if (!id.trim()) {
      throw new Error(
        "audit id is required",
      );
    }

    return this.client.get<
      ApiResponse<AuditEvent>
    >(
      `/audit/${encodeURIComponent(
        id,
      )}`,
    );
  }

  public async create(
    event: Omit<
      AuditEvent,
      "id" | "createdAt"
    >,
  ): Promise<
    ApiResponse<AuditEvent>
  > {
    return this.client.post<
      ApiResponse<AuditEvent>
    >(
      "/audit",
      event,
    );
  }

  public async delete(
    id: string,
  ): Promise<
    ApiResponse<boolean>
  > {
    if (!id.trim()) {
      throw new Error(
        "audit id is required",
      );
    }

    return this.client.delete<
      ApiResponse<boolean>
    >(
      `/audit/${encodeURIComponent(
        id,
      )}`,
    );
  }
}