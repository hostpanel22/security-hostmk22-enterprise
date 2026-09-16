import { describe, expect, it, vi } from "vitest";

import { FraudAPI } from "../src/fraud.js";
import { HostMKClient } from "../src/client.js";

describe(
  "FraudAPI",
  () => {
    it(
      "creates alert",
      async () => {
        const client = {
          post: vi.fn().mockResolvedValue({
            success: true,
            data: {
              id: "alert-1",
            },
          }),
        } as unknown as HostMKClient;

        const api =
          new FraudAPI(client);

        const result =
          await api.createAlert({
            tenantId:
              "tenant-1",
            userId:
              "user-1",
            riskScore: 90,
            reason:
              "suspicious login",
          });

        expect(
          result.success,
        ).toBe(true);

        expect(
          client.post,
        ).toHaveBeenCalled();
      },
    );

    it(
      "gets alert",
      async () => {
        const client = {
          get: vi.fn().mockResolvedValue({
            success: true,
            data: {
              id: "alert-1",
            },
          }),
        } as unknown as HostMKClient;

        const api =
          new FraudAPI(client);

        const result =
          await api.getAlert(
            "alert-1",
          );

        expect(
          result.data.id,
        ).toBe("alert-1");
      },
    );

    it(
      "lists alerts",
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
          new FraudAPI(client);

        const result =
          await api.listAlerts();

        expect(
          result.meta.page,
        ).toBe(1);
      },
    );

    it(
      "deletes alert",
      async () => {
        const client = {
          delete:
            vi.fn().mockResolvedValue({
              success: true,
              data: true,
            }),
        } as unknown as HostMKClient;

        const api =
          new FraudAPI(client);

        const result =
          await api.deleteAlert(
            "alert-1",
          );

        expect(
          result.data,
        ).toBe(true);
      },
    );

    it(
      "updates risk score",
      async () => {
        const client = {
          put: vi.fn().mockResolvedValue({
            success: true,
            data: {
              id: "alert-1",
              riskScore: 99,
            },
          }),
        } as unknown as HostMKClient;

        const api =
          new FraudAPI(client);

        const result =
          await api.updateRiskScore(
            "alert-1",
            {
              riskScore: 99,
            },
          );

        expect(
          result.data.riskScore,
        ).toBe(99);
      },
    );

    it(
      "resolves alert",
      async () => {
        const client = {
          post: vi.fn().mockResolvedValue({
            success: true,
            data: true,
          }),
        } as unknown as HostMKClient;

        const api =
          new FraudAPI(client);

        const result =
          await api.resolveAlert(
            "alert-1",
          );

        expect(
          result.data,
        ).toBe(true);
      },
    );

    it(
      "reopens alert",
      async () => {
        const client = {
          post: vi.fn().mockResolvedValue({
            success: true,
            data: true,
          }),
        } as unknown as HostMKClient;

        const api =
          new FraudAPI(client);

        const result =
          await api.reopenAlert(
            "alert-1",
          );

        expect(
          result.data,
        ).toBe(true);
      },
    );

    it(
      "gets user alerts",
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
          new FraudAPI(client);

        const result =
          await api.getUserAlerts(
            "user-1",
          );

        expect(
          result.meta.page,
        ).toBe(1);
      },
    );

    it(
      "gets high risk alerts",
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
          new FraudAPI(client);

        const result =
          await api.getHighRiskAlerts();

        expect(
          result.meta.page,
        ).toBe(1);
      },
    );

    it(
      "gets high risk alerts with threshold",
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
          new FraudAPI(client);

        const result =
          await api.getHighRiskAlerts(
            95,
          );

        expect(
          result.meta.page,
        ).toBe(1);
      },
    );
  },
);