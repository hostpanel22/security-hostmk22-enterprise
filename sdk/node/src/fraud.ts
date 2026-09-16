import { HostMKClient } from "./client.js";

import type {
  ApiResponse,
  FraudAlert,
  FraudQuery,
  PaginatedResponse,
} from "./types.js";

export class FraudService {
  private readonly client: HostMKClient;

  constructor(
    client: HostMKClient,
  ) {
    this.client = client;
  }

  public async list(
    query?: FraudQuery,
  ): Promise<
    PaginatedResponse<FraudAlert>
  > {
    const params =
      new URLSearchParams();

    if (query?.tenantId) {
      params.set(
        "tenantId",
        query.tenantId,
      );
    }

    if (query?.userId) {
      params.set(
        "userId",
        query.userId,
      );
    }

    if (
      query?.minRiskScore !==
      undefined
    ) {
      params.set(
        "minRiskScore",
        String(
          query.minRiskScore,
        ),
      );
    }

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

    const queryString =
      params.toString();

    return this.client.get<
      PaginatedResponse<FraudAlert>
    >(
      queryString.length > 0
        ? `/fraud?${queryString}`
        : "/fraud",
    );
  }

  public async get(
    id: string,
  ): Promise<
    ApiResponse<FraudAlert>
  > {
    if (!id.trim()) {
      throw new Error(
        "fraud alert id is required",
      );
    }

    return this.client.get<
      ApiResponse<FraudAlert>
    >(
      `/fraud/${encodeURIComponent(
        id,
      )}`,
    );
  }

  public async create(
    alert: Omit<
      FraudAlert,
      "id" | "createdAt"
    >,
  ): Promise<
    ApiResponse<FraudAlert>
  > {
    return this.client.post<
      ApiResponse<FraudAlert>
    >(
      "/fraud",
      alert,
    );
  }

  public async delete(
    id: string,
  ): Promise<
    ApiResponse<boolean>
  > {
    if (!id.trim()) {
      throw new Error(
        "fraud alert id is required",
      );
    }

    return this.client.delete<
      ApiResponse<boolean>
    >(
      `/fraud/${encodeURIComponent(
        id,
      )}`,
    );
  }
}