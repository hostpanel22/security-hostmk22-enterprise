export type HttpMethod = 'GET' | 'POST' | 'PUT' | 'PATCH' | 'DELETE';

export interface ApiErrorShape {
  message?: string;
  errors?: Record<string, string[]>;
}

export interface AuthTokens {
  accessToken?: string | null;
  refreshToken?: string | null;
}

export interface CurrentUser {
  id?: string;
  email?: string;
  name?: string;
  roles?: string[];
}

export interface PermissionRecord {
  name: string;
  action?: string;
}
