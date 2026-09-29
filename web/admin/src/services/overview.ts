import { fetchApiKeys, type ApiKeyRecord } from './apiKeys';
import { fetchAuditEvents, type AuditEventRecord } from './audit';
import { fetchFraudViolations, type FraudViolationRecord } from './fraud';
import { fetchSecrets, type SecretRecord } from './secrets';
import { fetchSessions, type SessionRecord } from './sessions';
import { fetchTenants, type TenantRecord } from './tenants';
import { fetchUsers, type UserRecord } from './users';

export interface OverviewCount {
  label: string;
  value: number;
  hint: string;
}

interface OverviewData {
  security: OverviewCount[];
  tenants: OverviewCount[];
  fraud: OverviewCount[];
  audit: OverviewCount[];
}

const isActiveApiKey = (key: ApiKeyRecord, now: number) =>
  !key.revoked_at && (!key.expires_at || new Date(key.expires_at).getTime() > now);

const isActiveSecret = (secret: SecretRecord) => !secret.revoked_at;

const isCurrentSession = (session: SessionRecord) => !session.revoked_at;

const countEventsSince = (events: AuditEventRecord[], from: number) =>
  events.filter((event) => event.created_at && new Date(event.created_at).getTime() >= from).length;

export async function fetchOverviewCounts(): Promise<OverviewData> {
  const [users, sessions, apiKeys, secrets, tenants, fraud, audit] = await Promise.all([
    fetchUsers(),
    fetchSessions(),
    fetchApiKeys(),
    fetchSecrets(),
    fetchTenants(),
    fetchFraudViolations(),
    fetchAuditEvents(),
  ]);

  const now = Date.now();
  const day = 24 * 60 * 60 * 1000;

  return {
    security: [
      { label: 'Active Users', value: users.filter((user: UserRecord) => user.status === 'active').length, hint: 'Derived from /api/v1/users status' },
      { label: 'Current Sessions', value: sessions.filter((session) => isCurrentSession(session)).length, hint: 'Derived from /api/v1/sessions revocation state' },
      { label: 'Active API Keys', value: apiKeys.filter((key) => isActiveApiKey(key, now)).length, hint: 'Derived from /api/v1/api-keys lifecycle fields' },
      { label: 'Active Secrets', value: secrets.filter((secret) => isActiveSecret(secret)).length, hint: 'Derived from /api/v1/secrets revocation state' },
    ],
    tenants: [
      { label: 'Total Tenants', value: tenants.length, hint: 'From /api/v1/tenants' },
      { label: 'Active Tenants', value: tenants.filter((tenant: TenantRecord) => tenant.status === 'active').length, hint: 'Derived from /api/v1/tenants status' },
      { label: 'Suspended Tenants', value: tenants.filter((tenant) => tenant.status === 'suspended').length, hint: 'Derived from /api/v1/tenants status' },
    ],
    fraud: [
      { label: 'Fraud Violations', value: fraud.length, hint: 'From /api/v1/fraud/violations' },
      { label: 'Critical Violations', value: fraud.filter((violation: FraudViolationRecord) => violation.level === 'critical').length, hint: 'Derived from violation level' },
      { label: 'High Risk Violations', value: fraud.filter((violation) => violation.level === 'high').length, hint: 'Derived from violation level' },
    ],
    audit: [
      { label: 'Events Today', value: countEventsSince(audit, now - day), hint: 'Derived from /api/v1/audit/events created_at' },
      { label: 'Events 7 Days', value: countEventsSince(audit, now - (7 * day)), hint: 'Derived from /api/v1/audit/events created_at' },
      { label: 'Events 30 Days', value: countEventsSince(audit, now - (30 * day)), hint: 'Derived from /api/v1/audit/events created_at' },
    ],
  };
}
