import { apiDelete, apiGet, apiPost, apiPut, normalizeList } from './api';

export interface RoleRecord {
  id?: string;
  name?: string;
  description?: string;
  permissions?: string[];
  tenant_id?: string;
  created_at?: string;
  updated_at?: string;
}

export async function createRole(role: Pick<RoleRecord, 'name' | 'description' | 'permissions'>): Promise<void> {
  await apiPost('/api/v1/rbac/roles', role);
}

export async function updateRole(role: RoleRecord & { id: string }): Promise<void> {
  await apiPut('/api/v1/rbac/roles', role);
}

export async function deleteRole(roleId: string): Promise<void> {
  await apiDelete('/api/v1/rbac/roles', { params: { id: roleId } });
}

export async function fetchRoles(): Promise<RoleRecord[]> {
  const payload = await apiGet<unknown>('/api/v1/rbac/roles');
  return normalizeList<RoleRecord>(payload);
}
