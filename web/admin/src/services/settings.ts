import { apiGet, normalizeList } from './api';

export interface SettingRecord {
  id?: string;
  key?: string;
  value?: string | number | boolean | Record<string, unknown>;
}

export async function fetchSettings(): Promise<SettingRecord[]> {
  const payload = await apiGet<unknown>('/api/v1/settings');
  return normalizeList<SettingRecord>(payload);
}
