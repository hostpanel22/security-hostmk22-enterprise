import { apiDelete, apiGet, apiPost, normalizeList, type ApiResponse } from './api';

export interface ApiKeyRecord {
  id?: string;
  name?: string;
  prefix?: string;
  scopes?: string[];
  rate_limit?: number | string;
  last_used_at?: string;
  expires_at?: string;
  revoked_at?: string;
  tenant_id?: string;
  key?: string;
  created_at?: string;
  updated_at?: string;
}

export interface ApiKeyWrite {
  name: string;
  scopes: string[];
  rate_limit?: number;
  expires_at?: string;
}

export async function fetchApiKeys(): Promise<ApiKeyRecord[]> {
  const payload = await apiGet<unknown>('/api/v1/api-keys');
  return normalizeList<ApiKeyRecord>(payload);
}

export async function createApiKey(input: ApiKeyWrite): Promise<ApiResponse<ApiKeyRecord>> {
  return apiPost<ApiResponse<ApiKeyRecord>>('/api/v1/api-keys', input);
}

export async function rotateApiKey(id: string): Promise<ApiResponse<ApiKeyRecord>> {
  return apiPost<ApiResponse<ApiKeyRecord>>('/api/v1/api-keys', undefined, { params: { action: 'rotate', id } });
}

export async function revokeApiKey(id: string): Promise<void> {
  await apiDelete('/api/v1/api-keys', { params: { id } });
}
