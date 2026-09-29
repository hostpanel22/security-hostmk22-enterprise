import { apiDelete, apiGet, apiPost, normalizeList, type ApiResponse } from './api';

export interface SecretRecord {
  id?: string;
  name?: string;
  type?: string;
  description?: string;
  created_at?: string;
  revoked_at?: string;
  tenant_id?: string;
  masked_value?: string;
  updated_at?: string;
}

export interface SecretVersionRecord {
  id?: string;
  secret_id?: string;
  masked_value?: string;
  created_at?: string;
  created_by?: string;
}

export interface SecretWrite {
  name: string;
  type: string;
  value: string;
  description?: string;
}

export async function fetchSecrets(): Promise<SecretRecord[]> {
  const payload = await apiGet<unknown>('/api/v1/secrets');
  return normalizeList<SecretRecord>(payload);
}

export async function createSecret(input: SecretWrite): Promise<ApiResponse<SecretRecord>> {
  return apiPost<ApiResponse<SecretRecord>>('/api/v1/secrets', input);
}

export async function rotateSecret(id: string, value: string): Promise<ApiResponse<SecretRecord>> {
  return apiPost<ApiResponse<SecretRecord>>('/api/v1/secrets', undefined, { params: { action: 'rotate', id, value } });
}

export async function fetchSecretHistory(id: string): Promise<SecretVersionRecord[]> {
  const payload = await apiGet<unknown>('/api/v1/secrets', { params: { id, history: true } });
  return normalizeList<SecretVersionRecord>(payload);
}

export async function revokeSecret(id: string): Promise<void> {
  await apiDelete('/api/v1/secrets', { params: { id } });
}
