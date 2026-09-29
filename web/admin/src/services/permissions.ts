import { apiDelete, apiGet, apiPost, normalizeList } from './api';

export interface PermissionRecord {
  id?: string;
  name?: string;
  resource?: string;
  action?: string;
  description?: string;
  tenant_id?: string;
  created_at?: string;
  updated_at?: string;
}

export async function createPermission(permission: Pick<PermissionRecord, 'name' | 'resource' | 'action' | 'description'>): Promise<void> {
  await apiPost('/api/v1/permissions', permission);
}

export async function deletePermission(permissionId: string): Promise<void> {
  await apiDelete('/api/v1/permissions', { params: { id: permissionId } });
}

export async function checkPermission(userId: string, permission: string): Promise<{ user_id: string; permission: string; allowed: boolean }> {
  const response = await apiPost<{ data?: { user_id: string; permission: string; allowed: boolean } }>('/api/v1/permissions/check', { user_id: userId, permission });
  if (!response.data) throw new Error('Permission check returned no result.');
  return response.data;
}

export async function fetchPermissions(): Promise<PermissionRecord[]> {
  const payload = await apiGet<unknown>('/api/v1/permissions');
  return normalizeList<PermissionRecord>(payload);
}
