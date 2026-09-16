import { describe, expect, it, vi } from "vitest";

import { PermissionsAPI } from "../src/permissions.js";
import { HostMKClient } from "../src/client.js";

describe(
  "PermissionsAPI",
  () => {
    it(
      "creates permission",
      async () => {
        const client = {
          post: vi.fn().mockResolvedValue({
            success: true,
            data: {
              id: "1",
              name: "users.read",
            },
          }),
        } as unknown as HostMKClient;

        const api =
          new PermissionsAPI(
            client,
          );

        const result =
          await api.createPermission(
            {
              tenantId:
                "tenant-1",
              name:
                "users.read",
              resource:
                "users",
              action:
                "read",
            },
          );

        expect(
          result.success,
        ).toBe(true);

        expect(
          client.post,
        ).toHaveBeenCalledWith(
          "/rbac/permissions",
          {
            tenantId:
              "tenant-1",
            name:
              "users.read",
            resource:
              "users",
            action:
              "read",
          },
        );
      },
    );

    it(
      "gets permission",
      async () => {
        const client = {
          get: vi.fn().mockResolvedValue({
            success: true,
            data: {
              id: "1",
              name:
                "users.read",
            },
          }),
        } as unknown as HostMKClient;

        const api =
          new PermissionsAPI(
            client,
          );

        await api.getPermission(
          "users.read",
        );

        expect(
          client.get,
        ).toHaveBeenCalledWith(
          "/rbac/permissions/users.read",
        );
      },
    );

    it(
      "lists permissions",
      async () => {
        const client = {
          get: vi.fn().mockResolvedValue({
            success: true,
            data: [],
            meta: {
              page: 1,
              limit: 10,
              total: 0,
              totalPages: 0,
            },
          }),
        } as unknown as HostMKClient;

        const api =
          new PermissionsAPI(
            client,
          );

        await api.listPermissions();

        expect(
          client.get,
        ).toHaveBeenCalledWith(
          "/rbac/permissions",
        );
      },
    );

    it(
      "deletes permission",
      async () => {
        const client = {
          delete:
            vi.fn().mockResolvedValue({
              success: true,
              data: true,
            }),
        } as unknown as HostMKClient;

        const api =
          new PermissionsAPI(
            client,
          );

        await api.deletePermission(
          "users.read",
        );

        expect(
          client.delete,
        ).toHaveBeenCalledWith(
          "/rbac/permissions/users.read",
        );
      },
    );

    it(
      "creates role",
      async () => {
        const client = {
          post: vi.fn().mockResolvedValue({
            success: true,
            data: {
              id: "1",
              name:
                "admin",
            },
          }),
        } as unknown as HostMKClient;

        const api =
          new PermissionsAPI(
            client,
          );

        await api.createRole({
          tenantId:
            "tenant-1",
          name: "admin",
        });

        expect(
          client.post,
        ).toHaveBeenCalledWith(
          "/rbac/roles",
          {
            tenantId:
              "tenant-1",
            name:
              "admin",
          },
        );
      },
    );

    it(
      "assigns permission to role",
      async () => {
        const client = {
          post: vi.fn().mockResolvedValue({
            success: true,
            data: true,
          }),
        } as unknown as HostMKClient;

        const api =
          new PermissionsAPI(
            client,
          );

        await api.assignPermissionToRole(
          "admin",
          "users.read",
        );

        expect(
          client.post,
        ).toHaveBeenCalledWith(
          "/rbac/roles/admin/permissions",
          {
            permissionName:
              "users.read",
          },
        );
      },
    );

    it(
      "checks role permission",
      async () => {
        const client = {
          get: vi.fn().mockResolvedValue({
            success: true,
            data: {
              allowed: true,
            },
          }),
        } as unknown as HostMKClient;

        const api =
          new PermissionsAPI(
            client,
          );

        const result =
          await api.roleHasPermission(
            "admin",
            "users.read",
          );

        expect(
          result.data.allowed,
        ).toBe(true);

        expect(
          client.get,
        ).toHaveBeenCalledWith(
          "/rbac/roles/admin/permissions/users.read",
        );
      },
    );

    it(
      "assigns parent role",
      async () => {
        const client = {
          post: vi.fn().mockResolvedValue({
            success: true,
            data: true,
          }),
        } as unknown as HostMKClient;

        const api =
          new PermissionsAPI(
            client,
          );

        await api.assignParentRole(
          "admin",
          "manager",
        );

        expect(
          client.post,
        ).toHaveBeenCalledWith(
          "/rbac/roles/admin/parents",
          {
            parentRole:
              "manager",
          },
        );
      },
    );

    it(
      "creates policy",
      async () => {
        const client = {
          post: vi.fn().mockResolvedValue({
            success: true,
            data: {
              id: "policy-1",
            },
          }),
        } as unknown as HostMKClient;

        const api =
          new PermissionsAPI(
            client,
          );

        await api.createPolicy({
          tenantId:
            "tenant-1",
          name:
            "allow-users",
          effect:
            "allow",
          resources: [
            "users",
          ],
          actions: [
            "read",
          ],
        });

        expect(
          client.post,
        ).toHaveBeenCalledWith(
          "/rbac/policies",
          {
            tenantId:
              "tenant-1",
            name:
              "allow-users",
            effect:
              "allow",
            resources: [
              "users",
            ],
            actions: [
              "read",
            ],
          },
        );
      },
    );

    it(
      "checks access",
      async () => {
        const client = {
          post: vi.fn().mockResolvedValue({
            success: true,
            data: {
              allowed: true,
            },
          }),
        } as unknown as HostMKClient;

        const api =
          new PermissionsAPI(
            client,
          );

        const result =
          await api.checkAccess(
            {
              roleName:
                "admin",
              permissionName:
                "users.read",
            },
          );

        expect(
          result.data.allowed,
        ).toBe(true);

        expect(
          client.post,
        ).toHaveBeenCalledWith(
          "/rbac/check-access",
          {
            roleName:
              "admin",
            permissionName:
              "users.read",
          },
        );
      },
    );
  },
);