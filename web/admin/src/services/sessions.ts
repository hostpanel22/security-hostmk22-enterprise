import { apiGet, apiPost, normalizeList } from './api';

export interface SessionRecord {
  id?: string;
  user_id?: string;
  user_agent?: string;
  ip_address?: string;
  expires_at?: string;
  revoked_at?: string;
  tenant_id?: string;
  created_at?: string;
}

export async function fetchSessions(): Promise<SessionRecord[]> {
  const payload = await apiGet<unknown>('/api/v1/sessions');
  return normalizeList<SessionRecord>(payload);
}

export async function revokeSession(sessionId: string): Promise<void> {
  await apiPost('/api/v1/sessions', undefined, { params: { action: 'revoke', id: sessionId } });
}
