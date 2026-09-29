import { apiGet, type ApiResponse } from './api';

export interface MonitoringItem {
  name: string;
  status: string;
  detail: string;
}

interface HealthData {
  status?: string;
  components?: Array<{ name?: string; status?: string; message?: string }>;
}

interface VersionData {
  name?: string;
  version?: string;
  env?: string;
}

export async function fetchMonitoringStatus(): Promise<MonitoringItem[]> {
  const [healthResponse, versionResponse] = await Promise.all([
    apiGet<ApiResponse<HealthData>>('/health'),
    apiGet<ApiResponse<VersionData>>('/api/v1/version'),
  ]);

  const health = healthResponse.data;
  const version = versionResponse.data;

  return [
    ...(health?.components ?? []).map((component) => ({
      name: component.name ?? 'Unnamed component',
      status: component.status ?? health?.status ?? 'unknown',
      detail: component.message ?? 'No status detail returned by the API',
    })),
    {
      name: version?.name ?? 'HostMK API',
      status: health?.status ?? 'unknown',
      detail: [version?.version, version?.env].filter(Boolean).join(' · ') || 'No version detail returned by the API',
    },
  ];
}
