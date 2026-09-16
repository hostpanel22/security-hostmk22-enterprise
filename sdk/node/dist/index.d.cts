interface HostMKClientOptions {
    apiKey: string;
    baseUrl: string;
    timeout?: number;
}
interface RequestOptions {
    method?: string;
    body?: unknown;
    headers?: Record<string, string>;
}
declare class HostMKError extends Error {
    readonly status: number;
    readonly data: unknown;
    constructor(message: string, status: number, data?: unknown);
}
declare class HostMKClient {
    private readonly apiKey;
    private readonly baseUrl;
    private readonly timeout;
    constructor(options: HostMKClientOptions);
    get<T>(path: string): Promise<T>;
    post<T>(path: string, body?: unknown): Promise<T>;
    put<T>(path: string, body?: unknown): Promise<T>;
    delete<T>(path: string): Promise<T>;
    private request;
}

interface ApiResponse<T> {
    success: boolean;
    data: T;
    message?: string;
}
interface PaginationMeta {
    page: number;
    limit: number;
    total: number;
    totalPages: number;
}
interface PaginatedResponse<T> {
    success: boolean;
    data: T[];
    meta: PaginationMeta;
}
interface Tenant {
    id: string;
    name: string;
    slug: string;
    active: boolean;
    createdAt: string;
    updatedAt: string;
}
interface Permission {
    id: string;
    tenantId: string;
    name: string;
    resource: string;
    action: string;
    description?: string;
    createdAt: string;
    updatedAt: string;
}
interface Role {
    id: string;
    tenantId: string;
    name: string;
    description?: string;
    permissions: string[];
    parentRoles: string[];
    createdAt: string;
    updatedAt: string;
}
interface Policy {
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
interface AuditEvent {
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
interface FraudAlert {
    id: string;
    tenantId: string;
    userId: string;
    riskScore: number;
    reason: string;
    metadata?: Record<string, unknown>;
    createdAt: string;
}
interface LoginRequest {
    email: string;
    password: string;
}
interface LoginResponse {
    token: string;
    refreshToken?: string;
    expiresAt: string;
}
interface RefreshTokenRequest {
    refreshToken: string;
}
interface ValidateTokenResponse {
    valid: boolean;
    userId?: string;
    tenantId?: string;
    expiresAt?: string;
}
interface CreateRoleRequest {
    tenantId: string;
    name: string;
    description?: string;
}
interface CreatePermissionRequest {
    tenantId: string;
    name: string;
    resource: string;
    action: string;
    description?: string;
}
interface CreatePolicyRequest {
    tenantId: string;
    name: string;
    description?: string;
    effect: "allow" | "deny";
    resources: string[];
    actions: string[];
}
interface AuditQuery {
    tenantId?: string;
    actorId?: string;
    resource?: string;
    action?: string;
    page?: number;
    limit?: number;
}
interface FraudQuery {
    tenantId?: string;
    userId?: string;
    minRiskScore?: number;
    page?: number;
    limit?: number;
}

declare class AuthAPI {
    private readonly client;
    constructor(client: HostMKClient);
    login(payload: LoginRequest): Promise<ApiResponse<LoginResponse>>;
    refreshToken(payload: RefreshTokenRequest): Promise<ApiResponse<LoginResponse>>;
    validateToken(token: string): Promise<ApiResponse<ValidateTokenResponse>>;
    logout(): Promise<ApiResponse<boolean>>;
    me(): Promise<ApiResponse<{
        id: string;
        email: string;
        tenantId: string;
        roles: string[];
    }>>;
}

interface AccessCheckRequest {
    roleName: string;
    permissionName: string;
}
interface AccessCheckResponse {
    allowed: boolean;
}
declare class PermissionsAPI {
    private readonly client;
    constructor(client: HostMKClient);
    createPermission(payload: CreatePermissionRequest): Promise<ApiResponse<Permission>>;
    getPermission(name: string): Promise<ApiResponse<Permission>>;
    listPermissions(): Promise<PaginatedResponse<Permission>>;
    deletePermission(name: string): Promise<ApiResponse<boolean>>;
    createRole(payload: CreateRoleRequest): Promise<ApiResponse<Role>>;
    getRole(name: string): Promise<ApiResponse<Role>>;
    listRoles(): Promise<PaginatedResponse<Role>>;
    deleteRole(name: string): Promise<ApiResponse<boolean>>;
    assignPermissionToRole(roleName: string, permissionName: string): Promise<ApiResponse<boolean>>;
    removePermissionFromRole(roleName: string, permissionName: string): Promise<ApiResponse<boolean>>;
    roleHasPermission(roleName: string, permissionName: string): Promise<ApiResponse<AccessCheckResponse>>;
    assignParentRole(roleName: string, parentRole: string): Promise<ApiResponse<boolean>>;
    removeParentRole(roleName: string, parentRole: string): Promise<ApiResponse<boolean>>;
    createPolicy(payload: CreatePolicyRequest): Promise<ApiResponse<Policy>>;
    getPolicy(id: string): Promise<ApiResponse<Policy>>;
    listPolicies(): Promise<PaginatedResponse<Policy>>;
    deletePolicy(id: string): Promise<ApiResponse<boolean>>;
    checkAccess(request: AccessCheckRequest): Promise<ApiResponse<AccessCheckResponse>>;
}

declare class AuditService {
    private readonly client;
    constructor(client: HostMKClient);
    list(query?: AuditQuery): Promise<PaginatedResponse<AuditEvent>>;
    get(id: string): Promise<ApiResponse<AuditEvent>>;
    create(event: Omit<AuditEvent, "id" | "createdAt">): Promise<ApiResponse<AuditEvent>>;
    delete(id: string): Promise<ApiResponse<boolean>>;
}

declare class FraudService {
    private readonly client;
    constructor(client: HostMKClient);
    list(query?: FraudQuery): Promise<PaginatedResponse<FraudAlert>>;
    get(id: string): Promise<ApiResponse<FraudAlert>>;
    create(alert: Omit<FraudAlert, "id" | "createdAt">): Promise<ApiResponse<FraudAlert>>;
    delete(id: string): Promise<ApiResponse<boolean>>;
}

interface CreateTenantRequest {
    name: string;
    slug: string;
}
interface UpdateTenantRequest {
    name?: string;
    slug?: string;
    active?: boolean;
}
interface TenantQuery {
    page?: number;
    limit?: number;
    search?: string;
    active?: boolean;
}
declare class TenancyAPI {
    private readonly client;
    constructor(client: HostMKClient);
    list(query?: TenantQuery): Promise<PaginatedResponse<Tenant>>;
    get(tenantId: string): Promise<ApiResponse<Tenant>>;
    create(payload: CreateTenantRequest): Promise<ApiResponse<Tenant>>;
    update(tenantId: string, payload: UpdateTenantRequest): Promise<ApiResponse<Tenant>>;
    delete(tenantId: string): Promise<ApiResponse<boolean>>;
    activate(tenantId: string): Promise<ApiResponse<Tenant>>;
    deactivate(tenantId: string): Promise<ApiResponse<Tenant>>;
}

export { type AccessCheckRequest, type AccessCheckResponse, type ApiResponse, type AuditEvent, type AuditQuery, AuditService, AuthAPI, type CreatePermissionRequest, type CreatePolicyRequest, type CreateRoleRequest, type CreateTenantRequest, type FraudAlert, type FraudQuery, FraudService, HostMKClient, type HostMKClientOptions, HostMKError, type LoginRequest, type LoginResponse, type PaginatedResponse, type PaginationMeta, type Permission, PermissionsAPI, type Policy, type RefreshTokenRequest, type RequestOptions, type Role, TenancyAPI, type Tenant, type TenantQuery, type UpdateTenantRequest, type ValidateTokenResponse };
