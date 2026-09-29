import { apiGet, normalizeList } from './api';

export interface TenantRecord {
  id?: string;
  name?: string;
  slug?: string;
  status?: string;
  plan?: string;
  quotas?: Record<string, unknown>;
  branding?: Record<string, unknown>;
}

export async function fetchTenants(): Promise<TenantRecord[]> {
  const payload = await apiGet<unknown>('/api/v1/tenants');
  return normalizeList<TenantRecord>(payload);
}
