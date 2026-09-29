export type ApiListResponse<T> =
  | T[]
  | { items?: T[]; data?: T[]; results?: T[]; users?: T[]; tenants?: T[]; roles?: T[]; permissions?: T[]; secrets?: T[]; sessions?: T[]; records?: T[]; value?: T[]; entries?: T[] }
  | { count?: number; items?: T[]; data?: T[]; results?: T[]; users?: T[]; tenants?: T[]; roles?: T[]; permissions?: T[]; secrets?: T[]; sessions?: T[]; records?: T[]; value?: T[]; entries?: T[] };

export interface PaginationMeta {
  page?: number;
  limit?: number;
  total?: number;
  pages?: number;
}

export type ApiRecord = Record<string, unknown>;

export interface DashboardMetric {
  label: string;
  value: string | number;
  hint?: string;
}
