import {
  beforeEach,
  describe,
  expect,
  it,
  vi,
} from "vitest";

import {
  HostMKClient,
  HostMKError,
} from "../src/client.js";

describe(
  "HostMKClient",
  () => {
    beforeEach(() => {
      vi.restoreAllMocks();
    });

    it(
      "should create client",
      () => {
        const client =
          new HostMKClient({
            apiKey:
              "test-api-key",
            baseUrl:
              "https://api.hostmk.com",
          });

        expect(
          client,
        ).toBeInstanceOf(
          HostMKClient,
        );
      },
    );

    it(
      "should throw when apiKey is missing",
      () => {
        expect(
          () =>
            new HostMKClient({
              apiKey: "",
              baseUrl:
                "https://api.hostmk.com",
            }),
        ).toThrow(
          "apiKey is required",
        );
      },
    );

    it(
      "should throw when baseUrl is missing",
      () => {
        expect(
          () =>
            new HostMKClient({
              apiKey:
                "test-api-key",
              baseUrl: "",
            }),
        ).toThrow(
          "baseUrl is required",
        );
      },
    );

    it(
      "should perform GET request",
      async () => {
        vi.spyOn(
          global,
          "fetch",
        ).mockResolvedValue(
          new Response(
            JSON.stringify({
              success: true,
            }),
            {
              status: 200,
              headers: {
                "content-type":
                  "application/json",
              },
            },
          ),
        );

        const client =
          new HostMKClient({
            apiKey:
              "test-api-key",
            baseUrl:
              "https://api.hostmk.com",
          });

        const result =
          await client.get<{
            success: boolean;
          }>(
            "/health",
          );

        expect(
          result.success,
        ).toBe(true);

        expect(
          fetch,
        ).toHaveBeenCalledTimes(
          1,
        );
      },
    );

    it(
      "should perform POST request",
      async () => {
        vi.spyOn(
          global,
          "fetch",
        ).mockResolvedValue(
          new Response(
            JSON.stringify({
              id: "1",
            }),
            {
              status: 200,
              headers: {
                "content-type":
                  "application/json",
              },
            },
          ),
        );

        const client =
          new HostMKClient({
            apiKey:
              "test-api-key",
            baseUrl:
              "https://api.hostmk.com",
          });

        const result =
          await client.post<{
            id: string;
          }>(
            "/users",
            {
              name:
                "admin",
            },
          );

        expect(
          result.id,
        ).toBe("1");

        expect(
          fetch,
        ).toHaveBeenCalledTimes(
          1,
        );
      },
    );

    it(
      "should perform PUT request",
      async () => {
        vi.spyOn(
          global,
          "fetch",
        ).mockResolvedValue(
          new Response(
            JSON.stringify({
              updated: true,
            }),
            {
              status: 200,
              headers: {
                "content-type":
                  "application/json",
              },
            },
          ),
        );

        const client =
          new HostMKClient({
            apiKey:
              "test-api-key",
            baseUrl:
              "https://api.hostmk.com",
          });

        const result =
          await client.put<{
            updated: boolean;
          }>(
            "/users/1",
            {
              name:
                "new-name",
            },
          );

        expect(
          result.updated,
        ).toBe(true);
      },
    );

    it(
      "should perform DELETE request",
      async () => {
        vi.spyOn(
          global,
          "fetch",
        ).mockResolvedValue(
          new Response(
            JSON.stringify({
              deleted: true,
            }),
            {
              status: 200,
              headers: {
                "content-type":
                  "application/json",
              },
            },
          ),
        );

        const client =
          new HostMKClient({
            apiKey:
              "test-api-key",
            baseUrl:
              "https://api.hostmk.com",
          });

        const result =
          await client.delete<{
            deleted: boolean;
          }>(
            "/users/1",
          );

        expect(
          result.deleted,
        ).toBe(true);
      },
    );

    it(
      "should throw HostMKError on api error",
      async () => {
        vi.spyOn(
          global,
          "fetch",
        ).mockResolvedValue(
          new Response(
            JSON.stringify({
              error:
                "unauthorized",
            }),
            {
              status: 401,
              statusText:
                "Unauthorized",
              headers: {
                "content-type":
                  "application/json",
              },
            },
          ),
        );

        const client =
          new HostMKClient({
            apiKey:
              "invalid-key",
            baseUrl:
              "https://api.hostmk.com",
          });

        await expect(
          client.get(
            "/secure",
          ),
        ).rejects.toBeInstanceOf(
          HostMKError,
        );
      },
    );

    it(
      "should handle text responses",
      async () => {
        vi.spyOn(
          global,
          "fetch",
        ).mockResolvedValue(
          new Response(
            "OK",
            {
              status: 200,
              headers: {
                "content-type":
                  "text/plain",
              },
            },
          ),
        );

        const client =
          new HostMKClient({
            apiKey:
              "test-api-key",
            baseUrl:
              "https://api.hostmk.com",
          });

        const result =
          await client.get<string>(
            "/ping",
          );

        expect(
          result,
        ).toBe("OK");
      },
    );

    it(
      "should trim trailing slash from baseUrl",
      async () => {
        const fetchMock =
          vi.spyOn(
            global,
            "fetch",
          ).mockResolvedValue(
            new Response(
              JSON.stringify({
                success: true,
              }),
              {
                status: 200,
                headers: {
                  "content-type":
                    "application/json",
                },
              },
            ),
          );

        const client =
          new HostMKClient({
            apiKey:
              "test-api-key",
            baseUrl:
              "https://api.hostmk.com/",
          });

        await client.get(
          "/health",
        );

        expect(
          fetchMock,
        ).toHaveBeenCalledWith(
          expect.stringContaining(
            "https://api.hostmk.com/health",
          ),
          expect.any(
            Object,
          ),
        );
      },
    );
  },
);