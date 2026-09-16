import { HostMKClient } from "./client.js";
import type {
  ApiResponse,
  LoginRequest,
  LoginResponse,
  RefreshTokenRequest,
  ValidateTokenResponse,
} from "./types.js";

export class AuthAPI {
  private readonly client: HostMKClient;

  constructor(
    client: HostMKClient,
  ) {
    this.client = client;
  }

  public async login(
    payload: LoginRequest,
  ): Promise<ApiResponse<LoginResponse>> {
    return this.client.post<
      ApiResponse<LoginResponse>
    >(
      "/auth/login",
      payload,
    );
  }

  public async refreshToken(
    payload: RefreshTokenRequest,
  ): Promise<ApiResponse<LoginResponse>> {
    return this.client.post<
      ApiResponse<LoginResponse>
    >(
      "/auth/refresh",
      payload,
    );
  }

  public async validateToken(
    token: string,
  ): Promise<
    ApiResponse<ValidateTokenResponse>
  > {
    return this.client.post<
      ApiResponse<ValidateTokenResponse>
    >(
      "/auth/validate",
      {
        token,
      },
    );
  }

  public async logout(): Promise<
    ApiResponse<boolean>
  > {
    return this.client.post<
      ApiResponse<boolean>
    >(
      "/auth/logout",
    );
  }

  public async me(): Promise<
    ApiResponse<{
      id: string;
      email: string;
      tenantId: string;
      roles: string[];
    }>
  > {
    return this.client.get<
      ApiResponse<{
        id: string;
        email: string;
        tenantId: string;
        roles: string[];
      }>
    >(
      "/auth/me",
    );
  }
}