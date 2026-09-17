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
  TenancyAPI,
} from "../src/tenancy.js";

describe(
  "TenancyAPI",
  () => {
    let client: HostMKClient;

    let tenancy: TenancyAPI;

    beforeEach(
      () => {
        client =
          new HostMKClient({
            apiKey:
              "test-api-key",
            baseUrl:
              "https://api.hostmk.io",
          });

        tenancy =
          new TenancyAPI(
            client,
          );
      },
    );

    it(
      "should list tenants without filters",
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
          await tenancy.list();

        expect(
          getSpy,
        ).toHaveBeenCalledWith(
          "/tenants",
        );

        expect(
          result.success,
        ).toBe(true);
      },
    );

    it(
      "should list tenants with filters",
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

        await tenancy.list({
          page: 1,
          limit: 20,
          search: "host",
          active: true,
        });

        expect(
          getSpy,
        ).toHaveBeenCalledWith(
          "/tenants?page=1&limit=20&search=host&active=true",
        );
      },
    );

    it(
      "should get tenant by id",
      async () => {
        const getSpy =
          vi.spyOn(
            client,
            "get",
          );

        getSpy.mockResolvedValue({
          success: true,
          data: {
            id: "tenant-1",
          },
        });

        const result =
          await tenancy.get(
            "tenant-1",
          );

        expect(
          getSpy,
        ).toHaveBeenCalledWith(
          "/tenants/tenant-1",
        );

        expect(
          result.success,
        ).toBe(true);
      },
    );

    it(
      "should create tenant",
      async () => {
        const postSpy =
          vi.spyOn(
            client,
            "post",
          );

        postSpy.mockResolvedValue({
          success: true,
          data: {
            id: "tenant-1",
          },
        });

        const payload = {
          name:
            "Tenant One",
          slug:
            "tenant-one",
        };

        const result =
          await tenancy.create(
            payload,
          );

        expect(
          postSpy,
        ).toHaveBeenCalledWith(
          "/tenants",
          payload,
        );

        expect(
          result.success,
        ).toBe(true);
      },
    );

    it(
      "should update tenant",
      async () => {
        const putSpy =
          vi.spyOn(
            client,
            "put",
          );

        putSpy.mockResolvedValue({
          success: true,
          data: {
            id: "tenant-1",
          },
        });

        const payload = {
          name:
            "Updated Tenant",
        };

        const result =
          await tenancy.update(
            "tenant-1",
            payload,
          );

        expect(
          putSpy,
        ).toHaveBeenCalledWith(
          "/tenants/tenant-1",
          payload,
        );

        expect(
          result.success,
        ).toBe(true);
      },
    );

    it(
      "should delete tenant",
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
          await tenancy.delete(
            "tenant-1",
          );

        expect(
          deleteSpy,
        ).toHaveBeenCalledWith(
          "/tenants/tenant-1",
        );

        expect(
          result.success,
        ).toBe(true);
      },
    );

    it(
      "should activate tenant",
      async () => {
        const postSpy =
          vi.spyOn(
            client,
            "post",
          );

        postSpy.mockResolvedValue({
          success: true,
          data: {
            id: "tenant-1",
            active: true,
          },
        });

        const result =
          await tenancy.activate(
            "tenant-1",
          );

        expect(
          postSpy,
        ).toHaveBeenCalledWith(
          "/tenants/tenant-1/activate",
        );

        expect(
          result.success,
        ).toBe(true);
      },
    );

    it(
      "should deactivate tenant",
      async () => {
        const postSpy =
          vi.spyOn(
            client,
            "post",
          );

        postSpy.mockResolvedValue({
          success: true,
          data: {
            id: "tenant-1",
            active: false,
          },
        });

        const result =
          await tenancy.deactivate(
            "tenant-1",
          );

        expect(
          postSpy,
        ).toHaveBeenCalledWith(
          "/tenants/tenant-1/deactivate",
        );

        expect(
          result.success,
        ).toBe(true);
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

        await tenancy.list({
          page: 2,
        });

        expect(
          getSpy,
        ).toHaveBeenCalledWith(
          "/tenants?page=2",
        );
      },
    );

    it(
      "should build query with active filter only",
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

        await tenancy.list({
          active: false,
        });

        expect(
          getSpy,
        ).toHaveBeenCalledWith(
          "/tenants?active=false",
        );
      },
    );
  },
);