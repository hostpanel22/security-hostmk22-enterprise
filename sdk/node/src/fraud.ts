import { HostMKClient } from "./client.js";
import type {
  ApiResponse,
  FraudAlert,
  FraudQuery,
  PaginatedResponse,
} from "./types.js";

export interface CreateFraudAlertRequest {
  tenantId: string;
  userId: string;
  riskScore: number;
  reason: string;
  metadata?: Record<string, unknown>;
}

export interface UpdateRiskScoreRequest {
  riskScore: number;
}

export class FraudAPI {
  private readonly client: HostMKClient;

  constructor(
    client: HostMKClient,
  ) {
    this.client = client;
  }

  /**
   * Create fraud alert
   */
  public async createAlert(
    payload: CreateFraudAlertRequest,
  ): Promise<ApiResponse<FraudAlert>> {
    return this.client.post<
      ApiResponse<FraudAlert>
    >(
      "/fraud/alerts",
      payload,
    );
  }

  /**
   * Get alert by ID
   */
  public async getAlert(
    id: string,
  ): Promise<ApiResponse<FraudAlert>> {
    return this.client.get<
      ApiResponse<FraudAlert>
    >(
      `/fraud/alerts/${encodeURIComponent(
        id,
      )}`,
    );
  }

  /**
   * List alerts
   */
  public async listAlerts(
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

    if (query?.page) {
      params.set(
        "page",
        String(query.page),
      );
    }

    if (query?.limit) {
      params.set(
        "limit",
        String(query.limit),
      );
    }

    const suffix =
      params.toString().length > 0
        ? `?${params.toString()}`
        : "";

    return this.client.get<
      PaginatedResponse<FraudAlert>
    >(
      `/fraud/alerts${suffix}`,
    );
  }

  /**
   * Delete alert
   */
  public async deleteAlert(
    id: string,
  ): Promise<ApiResponse<boolean>> {
    return this.client.delete<
      ApiResponse<boolean>
    >(
      `/fraud/alerts/${encodeURIComponent(
        id,
      )}`,
    );
  }

  /**
   * Update risk score
   */
  public async updateRiskScore(
    id: string,
    payload: UpdateRiskScoreRequest,
  ): Promise<ApiResponse<FraudAlert>> {
    return this.client.put<
      ApiResponse<FraudAlert>
    >(
      `/fraud/alerts/${encodeURIComponent(
        id,
      )}/risk-score`,
      payload,
    );
  }

  /**
   * Resolve alert
   */
  public async resolveAlert(
    id: string,
  ): Promise<ApiResponse<boolean>> {
    return this.client.post<
      ApiResponse<boolean>
    >(
      `/fraud/alerts/${encodeURIComponent(
        id,
      )}/resolve`,
    );
  }

  /**
   * Reopen alert
   */
  public async reopenAlert(
    id: string,
  ): Promise<ApiResponse<boolean>> {
    return this.client.post<
      ApiResponse<boolean>
    >(
      `/fraud/alerts/${encodeURIComponent(
        id,
      )}/reopen`,
    );
  }

  /**
   * Get user alerts
   */
  public async getUserAlerts(
    userId: string,
  ): Promise<
    PaginatedResponse<FraudAlert>
  > {
    return this.client.get<
      PaginatedResponse<FraudAlert>
    >(
      `/fraud/users/${encodeURIComponent(
        userId,
      )}/alerts`,
    );
  }

  /**
   * Get high-risk alerts
   */
  public async getHighRiskAlerts(
    threshold = 80,
  ): Promise<
    PaginatedResponse<FraudAlert>
  > {
    return this.client.get<
      PaginatedResponse<FraudAlert>
    >(
      `/fraud/high-risk?threshold=${threshold}`,
    );
  }
}