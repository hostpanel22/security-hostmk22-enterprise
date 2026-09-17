import {
  beforeEach,
  describe,
  expect,
  it,
  vi,
} from "vitest";

import {
  HostMKClient,
} from "../src/client.js";

import {
  AuditService,
} from "../src/audit.js";

describe(
  "AuditService",
  () => {
    let client: HostMKClient;

    let audit: AuditService;

    beforeEach(
      () => {
        client =
          new HostMKClient({
            apiKey:
              "test-api-key",
            baseUrl:
              "https://api.hostmk.io",
          });

        audit =
          new AuditService(
            client,
          );
      },
    );

    it(
      "should list audit events without filters",
      async () => {
        const getSpy =
          vi.spyOn(
            client,
            "get",
          );

        getSpy.mockResolvedValue({
          success: true,
          data: [],
          meta: {
            page: 1,
            limit: 10,
            total: 0,
            totalPages: 0,
          },
        });

        const result =
          await audit.list();

        expect(
          getSpy,
        ).toHaveBeenCalledWith(
          "/audit",
        );

        expect(
          result.meta.page,
        ).toBe(1);
      },
    );

    it(
      "should list audit events with filters",
      async () => {
        const getSpy =
          vi.spyOn(
            client,
            "get",
          );

        getSpy.mockResolvedValue({
          success: true,
          data: [],
          meta: {
            page: 1,
            limit: 20,
            total: 0,
            totalPages: 0,
          },
        });

        await audit.list({
          tenantId:
            "tenant-1",
          actorId:
            "user-1",
          resource:
            "orders",
          action:
            "create",
          page: 1,
          limit: 20,
        });

        expect(
          getSpy,
        ).toHaveBeenCalledWith(
          "/audit?tenantId=tenant-1&actorId=user-1&resource=orders&action=create&page=1&limit=20",
        );
      },
    );

    it(
      "should get audit event",
      async () => {
        const getSpy =
          vi.spyOn(
            client,
            "get",
          );

        getSpy.mockResolvedValue({
          success: true,
          data: {
            id: "audit-1",
          },
        });

        const result =
          await audit.get(
            "audit-1",
          );

        expect(
          getSpy,
        ).toHaveBeenCalledWith(
          "/audit/audit-1",
        );

        expect(
          result.success,
        ).toBe(true);
      },
    );

    it(
      "should create audit event",
      async () => {
        const postSpy =
          vi.spyOn(
            client,
            "post",
          );

        postSpy.mockResolvedValue({
          success: true,
          data: {
            id: "audit-1",
          },
        });

        const payload = {
          tenantId:
            "tenant-1",
          actorId:
            "user-1",
          action:
            "create",
          resource:
            "orders",
          status:
            "success",
        };

        const result =
          await audit.create(
            payload,
          );

        expect(
          postSpy,
        ).toHaveBeenCalledWith(
          "/audit",
          payload,
        );

        expect(
          result.success,
        ).toBe(true);
      },
    );

    it(
      "should delete audit event",
      async () => {
        const deleteSpy =
          vi.spyOn(
            client,
            "delete",
          );

        deleteSpy.mockResolvedValue({
          success: true,
          data: true,
        });

        const result =
          await audit.delete(
            "audit-1",
          );

        expect(
          deleteSpy,
        ).toHaveBeenCalledWith(
          "/audit/audit-1",
        );

        expect(
          result.success,
        ).toBe(true);
      },
    );

    it(
      "should throw if audit id is empty on get",
      async () => {
        await expect(
          audit.get(
            "",
          ),
        ).rejects.toThrow(
          "audit id is required",
        );
      },
    );

    it(
      "should throw if audit id is blank on get",
      async () => {
        await expect(
          audit.get(
            "   ",
          ),
        ).rejects.toThrow(
          "audit id is required",
        );
      },
    );

    it(
      "should throw if audit id is empty on delete",
      async () => {
        await expect(
          audit.delete(
            "",
          ),
        ).rejects.toThrow(
          "audit id is required",
        );
      },
    );

    it(
      "should throw if audit id is blank on delete",
      async () => {
        await expect(
          audit.delete(
            "   ",
          ),
        ).rejects.toThrow(
          "audit id is required",
        );
      },
    );

    it(
      "should build query with page only",
      async () => {
        const getSpy =
          vi.spyOn(
            client,
            "get",
          );

        getSpy.mockResolvedValue({
          success: true,
          data: [],
          meta: {
            page: 2,
            limit: 10,
            total: 0,
            totalPages: 0,
          },
        });

        await audit.list({
          page: 2,
        });

        expect(
          getSpy,
        ).toHaveBeenCalledWith(
          "/audit?page=2",
        );
      },
    );
  },
);