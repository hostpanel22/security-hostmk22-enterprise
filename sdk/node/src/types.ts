export interface ApiResponse<T> {
  success: boolean;
  data: T;
  message?: string;
}

export interface PaginationMeta {
  page: number;
  limit: number;
  total: number;
  totalPages: number;
}

export interface PaginatedResponse<T> {
  success: boolean;
  data: T[];
  meta: PaginationMeta;
}

export interface Tenant {
  id: string;
  name: string;
  slug: string;
  active: boolean;
  createdAt: string;
  updatedAt: string;
}

export interface Permission {
  id: string;
  tenantId: string;
  name: string;
  resource: string;
  action: string;
  description?: string;
  createdAt: string;
  updatedAt: string;
}

export interface Role {
  id: string;
  tenantId: string;
  name: string;
  description?: string;
  permissions: string[];
  parentRoles: string[];
  createdAt: string;
  updatedAt: string;
}

export interface Policy {
  id: string;
  tenantId: string;
  name: string;
  description?: string;
  effect: "allow" | "deny";
  resources: string[];
  actions: string[];
  createdAt: string;
  updatedAt: string;
}

export interface AuditEvent {
  id: string;
  tenantId: string;
  actorId: string;
  action: string;
  resource: string;
  status: string;
  ipAddress?: string;
  metadata?: Record<string, unknown>;
  createdAt: string;
}

export interface FraudAlert {
  id: string;
  tenantId: string;
  userId: string;
  riskScore: number;
  reason: string;
  metadata?: Record<string, unknown>;
  createdAt: string;
}

export interface LoginRequest {
  email: string;
  password: string;
}

export interface LoginResponse {
  token: string;
  refreshToken?: string;
  expiresAt: string;
}

export interface RefreshTokenRequest {
  refreshToken: string;
}

export interface ValidateTokenResponse {
  valid: boolean;
  userId?: string;
  tenantId?: string;
  expiresAt?: string;
}

export interface CreateRoleRequest {
  tenantId: string;
  name: string;
  description?: string;
}

export interface CreatePermissionRequest {
  tenantId: string;
  name: string;
  resource: string;
  action: string;
  description?: string;
}

export interface CreatePolicyRequest {
  tenantId: string;
  name: string;
  description?: string;
  effect: "allow" | "deny";
  resources: string[];
  actions: string[];
}

export interface AuditQuery {
  tenantId?: string;
  actorId?: string;
  resource?: string;
  action?: string;
  page?: number;
  limit?: number;
}

export interface FraudQuery {
  tenantId?: string;
  userId?: string;
  minRiskScore?: number;
  page?: number;
  limit?: number;
}