import {
  beforeEach,
  describe,
  expect,
  it,
  vi,
} from "vitest";

import { AuthAPI } from "../src/auth.js";
import { HostMKClient } from "../src/client.js";

describe(
  "AuthAPI",
  () => {
    let client: HostMKClient;

    let auth: AuthAPI;

    beforeEach(() => {
      vi.restoreAllMocks();

      client =
        new HostMKClient({
          apiKey:
            "test-api-key",
          baseUrl:
            "https://api.hostmk.com",
        });

      auth =
        new AuthAPI(
          client,
        );
    });

    it(
      "should login successfully",
      async () => {
        vi.spyOn(
          client,
          "post",
        ).mockResolvedValue({
          success: true,
          data: {
            token:
              "jwt-token",
            refreshToken:
              "refresh-token",
            expiresAt:
              "2027-01-01T00:00:00Z",
          },
        });

        const result =
          await auth.login({
            email:
              "admin@hostmk.com",
            password:
              "password123",
          });

        expect(
          result.success,
        ).toBe(true);

        expect(
          result.data.token,
        ).toBe(
          "jwt-token",
        );
      },
    );

    it(
      "should refresh token",
      async () => {
        vi.spyOn(
          client,
          "post",
        ).mockResolvedValue({
          success: true,
          data: {
            token:
              "new-token",
            refreshToken:
              "new-refresh-token",
            expiresAt:
              "2027-01-01T00:00:00Z",
          },
        });

        const result =
          await auth.refreshToken({
            refreshToken:
              "old-refresh-token",
          });

        expect(
          result.success,
        ).toBe(true);

        expect(
          result.data.token,
        ).toBe(
          "new-token",
        );
      },
    );

    it(
      "should validate token",
      async () => {
        vi.spyOn(
          client,
          "post",
        ).mockResolvedValue({
          success: true,
          data: {
            valid: true,
            userId:
              "user-1",
            tenantId:
              "tenant-1",
            expiresAt:
              "2027-01-01T00:00:00Z",
          },
        });

        const result =
          await auth.validateToken(
            "jwt-token",
          );

        expect(
          result.success,
        ).toBe(true);

        expect(
          result.data.valid,
        ).toBe(true);

        expect(
          result.data.userId,
        ).toBe(
          "user-1",
        );
      },
    );

    it(
      "should logout successfully",
      async () => {
        vi.spyOn(
          client,
          "post",
        ).mockResolvedValue({
          success: true,
          data: true,
        });

        const result =
          await auth.logout();

        expect(
          result.success,
        ).toBe(true);

        expect(
          result.data,
        ).toBe(true);
      },
    );

    it(
      "should get current user profile",
      async () => {
        vi.spyOn(
          client,
          "get",
        ).mockResolvedValue({
          success: true,
          data: {
            id: "user-1",
            email:
              "admin@hostmk.com",
            tenantId:
              "tenant-1",
            roles: [
              "admin",
            ],
          },
        });

        const result =
          await auth.me();

        expect(
          result.success,
        ).toBe(true);

        expect(
          result.data.id,
        ).toBe(
          "user-1",
        );

        expect(
          result.data.roles,
        ).toContain(
          "admin",
        );
      },
    );

    it(
      "should call login endpoint",
      async () => {
        const spy =
          vi.spyOn(
            client,
            "post",
          ).mockResolvedValue({
            success: true,
            data: {
              token:
                "jwt-token",
              expiresAt:
                "2027-01-01T00:00:00Z",
            },
          });

        await auth.login({
          email:
            "test@test.com",
          password:
            "secret",
        });

        expect(
          spy,
        ).toHaveBeenCalledWith(
          "/auth/login",
          {
            email:
              "test@test.com",
            password:
              "secret",
          },
        );
      },
    );

    it(
      "should call refresh endpoint",
      async () => {
        const spy =
          vi.spyOn(
            client,
            "post",
          ).mockResolvedValue({
            success: true,
            data: {},
          });

        await auth.refreshToken({
          refreshToken:
            "refresh-token",
        });

        expect(
          spy,
        ).toHaveBeenCalledWith(
          "/auth/refresh",
          {
            refreshToken:
              "refresh-token",
          },
        );
      },
    );

    it(
      "should call validate endpoint",
      async () => {
        const spy =
          vi.spyOn(
            client,
            "post",
          ).mockResolvedValue({
            success: true,
            data: {
              valid: true,
            },
          });

        await auth.validateToken(
          "jwt-token",
        );

        expect(
          spy,
        ).toHaveBeenCalledWith(
          "/auth/validate",
          {
            token:
              "jwt-token",
          },
        );
      },
    );

    it(
      "should call logout endpoint",
      async () => {
        const spy =
          vi.spyOn(
            client,
            "post",
          ).mockResolvedValue({
            success: true,
            data: true,
          });

        await auth.logout();

        expect(
          spy,
        ).toHaveBeenCalledWith(
          "/auth/logout",
        );
      },
    );

    it(
      "should call me endpoint",
      async () => {
        const spy =
          vi.spyOn(
            client,
            "get",
          ).mockResolvedValue({
            success: true,
            data: {},
          });

        await auth.me();

        expect(
          spy,
        ).toHaveBeenCalledWith(
          "/auth/me",
        );
      },
    );
  },
);