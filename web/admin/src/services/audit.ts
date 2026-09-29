import { apiGet, apiPost, normalizeList } from './api';

export interface AuditEventRecord {
  id?: string;
  action?: string;
  resource?: string;
  severity?: string;
  status?: string;
  created_at?: string;
  tenant_id?: string;
  user_id?: string;
  ip_address?: string;
  user_agent?: string;
  message?: string;
  metadata?: Record<string, unknown>;
  updated_at?: string;
}

export interface AuditEventWrite {
  action: string;
  resource: string;
  severity: string;
  status: string;
  message?: string;
  metadata?: Record<string, unknown>;
}

export async function fetchAuditEvents(limit = 100): Promise<AuditEventRecord[]> {
  const payload = await apiGet<unknown>('/api/v1/audit/events', { limit });
  return normalizeList<AuditEventRecord>(payload);
}

export async function createAuditEvent(event: AuditEventWrite): Promise<void> {
  await apiPost('/api/v1/audit/events', event);
}
