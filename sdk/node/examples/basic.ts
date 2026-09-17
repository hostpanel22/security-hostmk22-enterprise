import {
  HostMKClient,
  AuthAPI,
  PermissionsAPI,
  FraudAPI,
  AuditAPI,
  TenancyAPI,
} from "../src/index.js";

async function main(): Promise<void> {
  const client = new HostMKClient({
    apiKey:
      process.env.HOSTMK_API_KEY ??
      "replace-with-api-key",

    baseUrl:
      process.env.HOSTMK_BASE_URL ??
      "http://localhost:8080",
  });

  const auth =
    new AuthAPI(client);

  const permissions =
    new PermissionsAPI(client);

  const fraud =
    new FraudAPI(client);

  const audit =
    new AuditAPI(client);

  const tenancy =
    new TenancyAPI(client);

  try {
    console.log(
      "HostMK SDK Example",
    );

    /*
     * Auth
     */

    const me =
      await auth.me();

    console.log(
      "Current user:",
      me.data,
    );

    /*
     * Permissions
     */

    const roles =
      await permissions.listRoles();

    console.log(
      "Roles:",
      roles.data.length,
    );

    const permissionsList =
      await permissions.listPermissions();

    console.log(
      "Permissions:",
      permissionsList.data.length,
    );

    /*
     * Fraud
     */

    const alerts =
      await fraud.listAlerts();

    console.log(
      "Fraud alerts:",
      alerts.data.length,
    );

    /*
     * Audit
     */

    const events =
      await audit.listEvents();

    console.log(
      "Audit events:",
      events.data.length,
    );

    /*
     * Tenancy
     */

    const tenants =
      await tenancy.listTenants();

    console.log(
      "Tenants:",
      tenants.data.length,
    );

    console.log(
      "Example completed successfully",
    );
  } catch (error) {
    console.error(
      "Example failed:",
      error,
    );

    process.exitCode = 1;
  }
}

main();