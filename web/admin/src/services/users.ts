import { apiGet, apiPost, normalizeList } from './api';

export interface UserRecord {
  id?: string;
  email?: string;
  name?: string;
  status?: string;
  mfa_enabled?: boolean;
  last_login?: string;
  role_ids?: string[];
  tenant_id?: string;
  updated_at?: string;
  created_at?: string;
  email_verified_at?: string;
}

export interface UserWrite {
  email?: string;
  name?: string;
  password?: string;
  status?: string;
  role_ids?: string[];
}

export async function fetchUsers(): Promise<UserRecord[]> {
  const payload = await apiGet<unknown>('/api/v1/users');
  return normalizeList<UserRecord>(payload);
}

export async function createUser(user: Required<Pick<UserWrite, 'email' | 'name' | 'password' | 'status'>> & Pick<UserWrite, 'role_ids'>): Promise<UserRecord> {
  return apiPost<UserRecord>('/api/v1/users', user);
}

export async function updateUser(userId: string, user: Pick<UserWrite, 'name' | 'status' | 'role_ids'>): Promise<UserRecord> {
  return apiPost<UserRecord>('/api/v1/users/action', user, { params: { id: userId } });
}

export async function userAction(userId: string, action: 'suspend' | 'activate' | 'verify-email' | 'delete' | 'enable-mfa'): Promise<{ id?: string; secret?: string; otpauth_url?: string }> {
  return apiPost<{ id?: string; secret?: string; otpauth_url?: string }>('/api/v1/users/action', undefined, { params: { id: userId, action } });
}
