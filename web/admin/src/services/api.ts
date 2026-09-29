export type ApiResponse<T> = {
  success: boolean;
  message?: string;
  data?: T;
  error?: string;
  timestamp?: string;
};

export interface ApiRequestOptions {
  params?: Record<string, string | number | boolean | null | undefined>;
  filters?: Record<string, string | number | boolean | null | undefined>;
  search?: string;
  page?: number;
  limit?: number;
  sort?: string;
  headers?: Record<string, string>;
  signal?: AbortSignal;
}

export interface StandardApiError extends Error {
  status?: number;
  statusText?: string;
  method?: string;
  url?: string;
  details?: unknown;
}

export interface PaginatedResult<T> {
  items: T[];
  total?: number;
  page?: number;
  limit?: number;
  count?: number;
}

const listCandidateKeys = ['items', 'results', 'records', 'users', 'tenants', 'roles', 'permissions', 'sessions', 'secrets', 'values', 'entries', 'events', 'rules', 'settings', 'api_keys'];

const resolveBaseUrl = () => import.meta.env.VITE_API_BASE_URL ?? '';

const normalizeAuthToken = (): string | null => {
  try {
    return localStorage.getItem('hostmk.accessToken');
  } catch {
    return null;
  }
};

const buildQueryString = (params: Record<string, unknown> = {}): string => {
  const searchParams = new URLSearchParams();

  for (const [key, value] of Object.entries(params)) {
    if (value === undefined || value === null || value === '') {
      continue;
    }

    if (Array.isArray(value)) {
      value.forEach((entry) => {
        if (entry !== undefined && entry !== null && entry !== '') {
          searchParams.append(key, String(entry));
        }
      });
      continue;
    }

    searchParams.set(key, String(value));
  }

  const query = searchParams.toString();
  return query ? `?${query}` : '';
};

const appendRequestMeta = (options: ApiRequestOptions = {}) => {
  const params: Record<string, unknown> = { ...options.params };

  if (options.search) {
    params.search = options.search;
  }

  if (options.page) {
    params.page = options.page;
  }

  if (options.limit) {
    params.limit = options.limit;
  }

  if (options.sort) {
    params.sort = options.sort;
  }

  if (options.filters) {
    Object.entries(options.filters).forEach(([key, value]) => {
      if (value !== undefined && value !== null && value !== '') {
        params[key] = value;
      }
    });
  }

  return params;
};

const buildUrl = (path: string, options: ApiRequestOptions = {}) => {
  const normalizedPath = path.startsWith('/') ? path : `/${path}`;
  const params = appendRequestMeta(options);
  return `${resolveBaseUrl()}${normalizedPath}${buildQueryString(params)}`;
};

const parseResponse = async <T>(response: Response): Promise<T> => {
  const contentType = response.headers.get('content-type') ?? '';
  const isJson = contentType.includes('application/json');

  if (response.status === 204 || response.status === 205) {
    return undefined as T;
  }

  if (!isJson) {
    const text = await response.text();
    return (text as unknown) as T;
  }

  const payload = await response.json();
  return payload as T;
};

const toStandardApiError = (response: Response, payload: unknown, method: string, url: string): StandardApiError => {
  const message = typeof payload === 'object' && payload !== null
    ? ((payload as Record<string, unknown>).message as string | undefined)
      ?? ((payload as Record<string, unknown>).error as string | undefined)
      ?? ((payload as Record<string, unknown>).detail as string | undefined)
      ?? `Request failed: ${response.status}`
    : `Request failed: ${response.status}`;

  const error = new Error(message) as StandardApiError;
  error.status = response.status;
  error.statusText = response.statusText;
  error.method = method;
  error.url = url;
  error.details = payload;
  return error;
};

export async function apiRequest<T>(method: 'GET' | 'POST' | 'PUT' | 'PATCH' | 'DELETE', path: string, body?: unknown, options: ApiRequestOptions = {}): Promise<T> {
  const token = normalizeAuthToken();
  const headers = new Headers({
    Accept: 'application/json',
    ...options.headers,
  });

  if (token) {
    headers.set('Authorization', `Bearer ${token}`);
  }

  if (body !== undefined && !(body instanceof FormData) && !headers.has('Content-Type')) {
    headers.set('Content-Type', 'application/json');
  }

  const url = buildUrl(path, options);
  const response = await fetch(url, {
    method,
    headers,
    body: body === undefined ? undefined : typeof body === 'string' ? body : JSON.stringify(body),
    signal: options.signal,
  });

  const payload = await parseResponse<T>(response);

  if (!response.ok) {
    throw toStandardApiError(response, payload, method, url);
  }

  return payload;
}

export async function apiGet<T>(path: string, options: ApiRequestOptions = {}): Promise<T> {
  return apiRequest<T>('GET', path, undefined, options);
}

export async function apiPost<T>(path: string, body?: unknown, options: ApiRequestOptions = {}): Promise<T> {
  return apiRequest<T>('POST', path, body, options);
}

export async function apiPut<T>(path: string, body?: unknown, options: ApiRequestOptions = {}): Promise<T> {
  return apiRequest<T>('PUT', path, body, options);
}

export async function apiPatch<T>(path: string, body?: unknown, options: ApiRequestOptions = {}): Promise<T> {
  return apiRequest<T>('PATCH', path, body, options);
}

export async function apiDelete<T>(path: string, options: ApiRequestOptions = {}): Promise<T> {
  return apiRequest<T>('DELETE', path, undefined, options);
}

export function normalizeList<T>(payload: unknown): T[] {
  if (Array.isArray(payload)) {
    return payload as T[];
  }

  if (payload && typeof payload === 'object') {
    const record = payload as Record<string, unknown>;

    if ('data' in record) {
      return normalizeList<T>(record.data);
    }

    for (const key of listCandidateKeys) {
      const value = record[key];
      if (Array.isArray(value)) {
        return value as T[];
      }
    }

    if (typeof record.count === 'number') {
      return [];
    }
  }

  return [];
}

export function normalizePaginated<T>(payload: unknown): PaginatedResult<T> {
  if (payload && typeof payload === 'object' && 'data' in payload) {
    return normalizePaginated<T>((payload as ApiResponse<unknown>).data);
  }

  if (payload && typeof payload === 'object') {
    const record = payload as Record<string, unknown>;
    const items = normalizeList<T>(payload);
    return {
      items,
      total: typeof record.total === 'number' ? record.total : typeof record.count === 'number' ? record.count : undefined,
      page: typeof record.page === 'number' ? record.page : undefined,
      limit: typeof record.limit === 'number' ? record.limit : undefined,
      count: typeof record.count === 'number' ? record.count : items.length,
    };
  }

  return { items: normalizeList<T>(payload) };
}
