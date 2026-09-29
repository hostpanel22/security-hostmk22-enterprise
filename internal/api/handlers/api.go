package handlers

import (
	"encoding/json"
	"errors"
	"io"
	"log/slog"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/pquerna/otp/totp"

	"github.com/azezaz/HostMK-Security-Framework/internal/api/response"
	"github.com/azezaz/HostMK-Security-Framework/internal/auth"
	"github.com/azezaz/HostMK-Security-Framework/internal/config"
	"github.com/azezaz/HostMK-Security-Framework/internal/crypto"
	"github.com/azezaz/HostMK-Security-Framework/internal/engine"
	"github.com/azezaz/HostMK-Security-Framework/internal/identity"
	"github.com/azezaz/HostMK-Security-Framework/internal/models"
	"github.com/azezaz/HostMK-Security-Framework/internal/persistence"
)

type API struct {
	Config *config.Config
	Store  *persistence.Store
	Tokens auth.Tokens
	AEAD   *crypto.AEAD
	Log    *slog.Logger
}

func (a *API) tenantScope(r *http.Request) (identity.Identity, bool) {
	id, ok := identity.FromContext(r.Context())
	if !ok {
		return identity.Identity{}, false
	}
	if header := strings.TrimSpace(r.Header.Get("X-Tenant-ID")); header != "" && id.TenantID == "" {
		id.TenantID = header
	}
	return id, id.TenantID != ""
}

func decodeJSON(r *http.Request, dest any) error {
	return json.NewDecoder(io.LimitReader(r.Body, 1<<20)).Decode(dest)
}

func writeStoreError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, models.ErrNotFound):
		response.NotFound(w, err.Error())
	case errors.Is(err, models.ErrValidation), errors.Is(err, models.ErrInvalidInput):
		response.BadRequest(w, err.Error())
	default:
		response.InternalServerError(w, "persistence error")
	}
}

type bootstrapRequest struct {
	TenantName string `json:"tenant_name"`
	TenantSlug string `json:"tenant_slug"`
	Name       string `json:"name"`
	Email      string `json:"email"`
	Password   string `json:"password"`
}

func (a *API) AuthBootstrap(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		response.Error(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}
	count, err := a.Store.CountUsers(r.Context())
	if err != nil {
		writeStoreError(w, err)
		return
	}
	if count > 0 {
		response.Error(w, http.StatusConflict, "bootstrap already completed")
		return
	}

	var req bootstrapRequest
	if err := decodeJSON(r, &req); err != nil {
		response.BadRequest(w, "invalid json payload")
		return
	}

	hash, err := crypto.HashPassword(req.Password)
	if err != nil {
		response.InternalServerError(w, "password hashing failed")
		return
	}

	now := time.Now().UTC()
	tenant := &models.Tenant{
		Name:   req.TenantName,
		Slug:   req.TenantSlug,
		Status: models.TenantStatusActive,
	}
	createdTenant, err := a.Store.CreateTenant(r.Context(), tenant)
	if err != nil {
		writeStoreError(w, err)
		return
	}

	user := &models.User{
		TenantID:        createdTenant.ID,
		Email:           req.Email,
		PasswordHash:    hash,
		Name:            req.Name,
		Status:          models.UserStatusActive,
		EmailVerifiedAt: &now,
	}
	createdUser, err := a.Store.CreateUser(r.Context(), user)
	if err != nil {
		writeStoreError(w, err)
		return
	}
	createdTenant.OwnerID = createdUser.ID
	if _, err := a.Store.UpdateTenant(r.Context(), createdTenant); err != nil {
		writeStoreError(w, err)
		return
	}

	response.Created(w, "platform bootstrap completed", map[string]any{
		"tenant": createdTenant,
		"user":   createdUser,
	})
}

type loginRequest struct {
	TenantSlug string `json:"tenant_slug"`
	Email      string `json:"email"`
	Password   string `json:"password"`
	MFACode    string `json:"mfa_code"`
}

func (a *API) AuthLogin(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		response.Error(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}
	var req loginRequest
	if err := decodeJSON(r, &req); err != nil {
		response.BadRequest(w, "invalid json payload")
		return
	}
	tenant, err := a.Store.GetTenantBySlug(r.Context(), req.TenantSlug)
	if err != nil {
		response.Unauthorized(w, "invalid credentials")
		return
	}
	if !tenant.CanAccessPlatform() {
		response.Forbidden(w, "tenant is not active")
		return
	}
	user, err := a.Store.GetUserByEmail(r.Context(), tenant.ID, req.Email)
	if err != nil || !crypto.VerifyPassword(req.Password, user.PasswordHash) {
		response.Unauthorized(w, "invalid credentials")
		return
	}
	if !user.CanAuthenticate() {
		response.Forbidden(w, "user is not active")
		return
	}
	if user.MFAEnabled {
		plain, err := a.AEAD.Decrypt(user.MFASecretEnc)
		if err != nil || !totp.Validate(req.MFACode, string(plain)) {
			response.Unauthorized(w, "invalid mfa code")
			return
		}
	}

	tokens, err := a.issueTokens(r, user)
	if err != nil {
		response.InternalServerError(w, "token issuance failed")
		return
	}
	now := time.Now().UTC()
	user.LastLogin = &now
	_, _ = a.Store.UpdateUser(r.Context(), user)
	response.OK(w, "authenticated", tokens)
}

type refreshRequest struct {
	RefreshToken string `json:"refresh_token"`
}

func (a *API) AuthRefresh(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		response.Error(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}
	var req refreshRequest
	if err := decodeJSON(r, &req); err != nil {
		response.BadRequest(w, "invalid json payload")
		return
	}
	stored, err := a.Store.GetRefreshToken(r.Context(), persistence.HashToken(req.RefreshToken))
	if err != nil || stored.RevokedAt != nil || time.Now().UTC().After(stored.ExpiresAt) {
		response.Unauthorized(w, "invalid refresh token")
		return
	}
	user, err := a.Store.GetUser(r.Context(), stored.UserID)
	if err != nil {
		response.Unauthorized(w, "invalid refresh token")
		return
	}
	_ = a.Store.RevokeRefreshToken(r.Context(), stored.ID)
	tokens, err := a.issueTokens(r, user)
	if err != nil {
		response.InternalServerError(w, "token issuance failed")
		return
	}
	response.OK(w, "token refreshed", tokens)
}

func (a *API) issueTokens(r *http.Request, user *models.User) (map[string]any, error) {
	access, err := a.Tokens.SignAccess(user.ID, user.TenantID)
	if err != nil {
		return nil, err
	}
	plain, hash, err := persistence.RandomToken()
	if err != nil {
		return nil, err
	}
	session := &models.Session{
		TenantID:  user.TenantID,
		UserID:    user.ID,
		UserAgent: r.UserAgent(),
		IPAddress: r.RemoteAddr,
		ExpiresAt: time.Now().UTC().Add(a.Tokens.RefreshTTL),
	}
	if err := a.Store.CreateSession(r.Context(), session); err != nil {
		return nil, err
	}
	refresh := &models.RefreshToken{
		TenantID:  user.TenantID,
		UserID:    user.ID,
		SessionID: session.ID,
		TokenHash: hash,
		ExpiresAt: session.ExpiresAt,
	}
	if err := a.Store.CreateRefreshToken(r.Context(), refresh); err != nil {
		return nil, err
	}
	return map[string]any{
		"access_token":  access,
		"refresh_token": plain,
		"token_type":    "Bearer",
		"expires_in":    int(a.Tokens.AccessTTL.Seconds()),
	}, nil
}

func (a *API) Health(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		response.Error(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}
	dbStatus := "healthy"
	dbMessage := "postgres reachable"
	if err := a.Store.Ping(r.Context()); err != nil {
		dbStatus = "unhealthy"
		dbMessage = err.Error()
	}
	status := "healthy"
	code := http.StatusOK
	if dbStatus != "healthy" {
		status = "degraded"
		code = http.StatusServiceUnavailable
	}
	payload := map[string]any{
		"status":  status,
		"service": a.Config.App.Name,
		"version": a.Config.App.Version,
		"env":     a.Config.App.Environment,
		"components": []map[string]string{
			{"name": "postgres", "status": dbStatus, "message": dbMessage},
			{"name": "api", "status": "healthy", "message": "http server operational"},
		},
	}
	response.JSON(w, code, "health", payload)
}

func (a *API) Version(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		response.Error(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}
	response.OK(w, "version information", map[string]any{
		"name":    a.Config.App.Name,
		"version": a.Config.App.Version,
		"env":     a.Config.App.Environment,
	})
}

func (a *API) Users(w http.ResponseWriter, r *http.Request) {
	id, ok := a.tenantScope(r)
	if !ok {
		response.Unauthorized(w, "tenant required")
		return
	}
	switch r.Method {
	case http.MethodGet:
		users, err := a.Store.ListUsers(r.Context(), id.TenantID)
		if err != nil {
			writeStoreError(w, err)
			return
		}
		response.OK(w, "users retrieved", map[string]any{"count": len(users), "users": users})
	case http.MethodPost:
		var req struct {
			Email    string   `json:"email"`
			Name     string   `json:"name"`
			Password string   `json:"password"`
			Status   string   `json:"status"`
			RoleIDs  []string `json:"role_ids"`
		}
		if err := decodeJSON(r, &req); err != nil {
			response.BadRequest(w, "invalid json payload")
			return
		}
		hash, err := crypto.HashPassword(req.Password)
		if err != nil {
			response.InternalServerError(w, "password hashing failed")
			return
		}
		user := &models.User{
			TenantID:     id.TenantID,
			Email:        req.Email,
			Name:         req.Name,
			PasswordHash: hash,
			Status:       models.UserStatus(req.Status),
			RoleIDs:      req.RoleIDs,
		}
		created, err := a.Store.CreateUser(r.Context(), user)
		if err != nil {
			writeStoreError(w, err)
			return
		}
		a.audit(r, id, "users.create", "users", created.ID)
		response.Created(w, "user created", created)
	default:
		response.Error(w, http.StatusMethodNotAllowed, "method not allowed")
	}
}

func (a *API) UserAction(w http.ResponseWriter, r *http.Request) {
	id, ok := a.tenantScope(r)
	if !ok {
		response.Unauthorized(w, "tenant required")
		return
	}
	userID := strings.TrimSpace(r.URL.Query().Get("id"))
	if userID == "" {
		response.BadRequest(w, "id is required")
		return
	}
	user, err := a.Store.GetUser(r.Context(), userID)
	if err != nil {
		writeStoreError(w, err)
		return
	}
	if user.TenantID != id.TenantID {
		response.Forbidden(w, "tenant mismatch")
		return
	}

	switch r.URL.Query().Get("action") {
	case "suspend":
		user.Status = models.UserStatusSuspended
	case "activate":
		user.Status = models.UserStatusActive
	case "verify-email":
		now := time.Now().UTC()
		user.EmailVerifiedAt = &now
	case "delete":
		if err := a.Store.DeleteUser(r.Context(), user.ID); err != nil {
			writeStoreError(w, err)
			return
		}
		a.audit(r, id, "users.delete", "users", user.ID)
		response.OK(w, "user deleted", map[string]string{"id": user.ID})
		return
	case "enable-mfa":
		key, err := totp.Generate(totp.GenerateOpts{Issuer: a.Config.JWT.Issuer, AccountName: user.Email})
		if err != nil {
			response.InternalServerError(w, "mfa enrollment failed")
			return
		}
		enc, err := a.AEAD.Encrypt([]byte(key.Secret()))
		if err != nil {
			response.InternalServerError(w, "mfa encryption failed")
			return
		}
		user.MFAEnabled = true
		user.MFASecretEnc = enc
		if _, err := a.Store.UpdateUser(r.Context(), user); err != nil {
			writeStoreError(w, err)
			return
		}
		response.OK(w, "mfa enabled", map[string]string{"otpauth_url": key.URL(), "secret": key.Secret()})
		return
	default:
		var req struct {
			Name    string   `json:"name"`
			Status  string   `json:"status"`
			RoleIDs []string `json:"role_ids"`
		}
		if err := decodeJSON(r, &req); err != nil {
			response.BadRequest(w, "invalid json payload")
			return
		}
		if req.Name != "" {
			user.Name = req.Name
		}
		if req.Status != "" {
			user.Status = models.UserStatus(req.Status)
		}
		if req.RoleIDs != nil {
			user.RoleIDs = req.RoleIDs
		}
	}
	updated, err := a.Store.UpdateUser(r.Context(), user)
	if err != nil {
		writeStoreError(w, err)
		return
	}
	a.audit(r, id, "users.update", "users", user.ID)
	response.OK(w, "user updated", updated)
}

func (a *API) Tenants(w http.ResponseWriter, r *http.Request) {
	id, ok := identity.FromContext(r.Context())
	if !ok {
		response.Unauthorized(w, "authentication required")
		return
	}
	switch r.Method {
	case http.MethodGet:
		tenants, err := a.Store.ListTenants(r.Context())
		if err != nil {
			writeStoreError(w, err)
			return
		}
		response.OK(w, "tenants retrieved", map[string]any{"count": len(tenants), "tenants": tenants})
	case http.MethodPost:
		var req models.Tenant
		if err := decodeJSON(r, &req); err != nil {
			response.BadRequest(w, "invalid json payload")
			return
		}
		created, err := a.Store.CreateTenant(r.Context(), &req)
		if err != nil {
			writeStoreError(w, err)
			return
		}
		a.audit(r, id, "tenants.create", "tenants", created.ID)
		response.Created(w, "tenant created", created)
	default:
		response.Error(w, http.StatusMethodNotAllowed, "method not allowed")
	}
}

func (a *API) TenantAction(w http.ResponseWriter, r *http.Request) {
	id, ok := identity.FromContext(r.Context())
	if !ok {
		response.Unauthorized(w, "authentication required")
		return
	}
	tenantID := strings.TrimSpace(r.URL.Query().Get("id"))
	tenant, err := a.Store.GetTenant(r.Context(), tenantID)
	if err != nil {
		writeStoreError(w, err)
		return
	}
	switch r.URL.Query().Get("action") {
	case "suspend":
		tenant.Suspend()
	case "activate":
		tenant.Activate()
	default:
		var req models.Tenant
		if err := decodeJSON(r, &req); err != nil {
			response.BadRequest(w, "invalid json payload")
			return
		}
		req.ID = tenant.ID
		req.CreatedAt = tenant.CreatedAt
		tenant = &req
	}
	updated, err := a.Store.UpdateTenant(r.Context(), tenant)
	if err != nil {
		writeStoreError(w, err)
		return
	}
	a.audit(r, id, "tenants.update", "tenants", tenant.ID)
	response.OK(w, "tenant updated", updated)
}

func (a *API) Roles(w http.ResponseWriter, r *http.Request) {
	id, ok := a.tenantScope(r)
	if !ok {
		response.Unauthorized(w, "tenant required")
		return
	}
	switch r.Method {
	case http.MethodGet:
		roles, err := a.Store.ListRoles(r.Context(), id.TenantID)
		if err != nil {
			writeStoreError(w, err)
			return
		}
		response.OK(w, "roles retrieved", map[string]any{"count": len(roles), "roles": roles})
	case http.MethodPost:
		var role models.Role
		if err := decodeJSON(r, &role); err != nil {
			response.BadRequest(w, "invalid json payload")
			return
		}
		role.TenantID = id.TenantID
		created, err := a.Store.CreateRole(r.Context(), &role)
		if err != nil {
			writeStoreError(w, err)
			return
		}
		response.Created(w, "role created", created)
	case http.MethodPut:
		var role models.Role
		if err := decodeJSON(r, &role); err != nil {
			response.BadRequest(w, "invalid json payload")
			return
		}
		role.TenantID = id.TenantID
		updated, err := a.Store.UpdateRole(r.Context(), &role)
		if err != nil {
			writeStoreError(w, err)
			return
		}
		response.OK(w, "role updated", updated)
	case http.MethodDelete:
		if err := a.Store.DeleteRole(r.Context(), r.URL.Query().Get("id")); err != nil {
			writeStoreError(w, err)
			return
		}
		response.OK(w, "role deleted", map[string]string{"id": r.URL.Query().Get("id")})
	default:
		response.Error(w, http.StatusMethodNotAllowed, "method not allowed")
	}
}

func (a *API) Permissions(w http.ResponseWriter, r *http.Request) {
	id, ok := a.tenantScope(r)
	if !ok {
		response.Unauthorized(w, "tenant required")
		return
	}
	switch r.Method {
	case http.MethodGet:
		items, err := a.Store.ListPermissions(r.Context(), id.TenantID)
		if err != nil {
			writeStoreError(w, err)
			return
		}
		response.OK(w, "permissions retrieved", map[string]any{"count": len(items), "permissions": items})
	case http.MethodPost:
		var permission models.Permission
		if err := decodeJSON(r, &permission); err != nil {
			response.BadRequest(w, "invalid json payload")
			return
		}
		permission.TenantID = id.TenantID
		created, err := a.Store.CreatePermission(r.Context(), &permission)
		if err != nil {
			writeStoreError(w, err)
			return
		}
		response.Created(w, "permission created", created)
	case http.MethodDelete:
		if err := a.Store.DeletePermission(r.Context(), r.URL.Query().Get("id")); err != nil {
			writeStoreError(w, err)
			return
		}
		response.OK(w, "permission deleted", map[string]string{"id": r.URL.Query().Get("id")})
	default:
		response.Error(w, http.StatusMethodNotAllowed, "method not allowed")
	}
}

func (a *API) PermissionCheck(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		response.Error(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}
	var req struct {
		UserID     string `json:"user_id"`
		Permission string `json:"permission"`
	}
	if err := decodeJSON(r, &req); err != nil {
		response.BadRequest(w, "invalid json payload")
		return
	}
	allowed, err := a.Store.UserHasPermission(r.Context(), req.UserID, req.Permission)
	if err != nil {
		writeStoreError(w, err)
		return
	}
	response.OK(w, "permission check completed", map[string]any{
		"user_id":    req.UserID,
		"permission": req.Permission,
		"allowed":    allowed,
	})
}

func (a *API) APIKeys(w http.ResponseWriter, r *http.Request) {
	id, ok := a.tenantScope(r)
	if !ok {
		response.Unauthorized(w, "tenant required")
		return
	}
	switch r.Method {
	case http.MethodGet:
		keys, err := a.Store.ListAPIKeys(r.Context(), id.TenantID)
		if err != nil {
			writeStoreError(w, err)
			return
		}
		response.OK(w, "api keys retrieved", map[string]any{"count": len(keys), "api_keys": keys})
	case http.MethodPost:
		if r.URL.Query().Get("action") == "rotate" {
			rotated, err := a.Store.RotateAPIKey(r.Context(), r.URL.Query().Get("id"))
			if err != nil {
				writeStoreError(w, err)
				return
			}
			response.OK(w, "api key rotated", rotated)
			return
		}
		var key models.APIKey
		if err := decodeJSON(r, &key); err != nil {
			response.BadRequest(w, "invalid json payload")
			return
		}
		key.TenantID = id.TenantID
		created, err := a.Store.CreateAPIKey(r.Context(), &key)
		if err != nil {
			writeStoreError(w, err)
			return
		}
		response.Created(w, "api key created — store the plaintext now, it will not be shown again", created)
	case http.MethodDelete:
		if err := a.Store.RevokeAPIKey(r.Context(), r.URL.Query().Get("id")); err != nil {
			writeStoreError(w, err)
			return
		}
		response.OK(w, "api key revoked", map[string]string{"id": r.URL.Query().Get("id")})
	default:
		response.Error(w, http.StatusMethodNotAllowed, "method not allowed")
	}
}

func (a *API) Secrets(w http.ResponseWriter, r *http.Request) {
	id, ok := a.tenantScope(r)
	if !ok {
		response.Unauthorized(w, "tenant required")
		return
	}
	switch r.Method {
	case http.MethodGet:
		if secretID := r.URL.Query().Get("id"); secretID != "" && r.URL.Query().Get("history") == "true" {
			versions, err := a.Store.ListSecretVersions(r.Context(), secretID)
			if err != nil {
				writeStoreError(w, err)
				return
			}
			response.OK(w, "secret history retrieved", versions)
			return
		}
		secrets, err := a.Store.ListSecrets(r.Context(), id.TenantID)
		if err != nil {
			writeStoreError(w, err)
			return
		}
		response.OK(w, "secrets retrieved", map[string]any{"count": len(secrets), "secrets": secrets})
	case http.MethodPost:
		if r.URL.Query().Get("action") == "rotate" {
			rotated, err := a.Store.RotateSecret(r.Context(), a.AEAD, r.URL.Query().Get("id"), r.URL.Query().Get("value"), id.UserID)
			if err != nil {
				writeStoreError(w, err)
				return
			}
			response.OK(w, "secret rotated", rotated)
			return
		}
		var secret models.Secret
		if err := decodeJSON(r, &secret); err != nil {
			response.BadRequest(w, "invalid json payload")
			return
		}
		secret.TenantID = id.TenantID
		created, err := a.Store.CreateSecret(r.Context(), a.AEAD, &secret, id.UserID)
		if err != nil {
			writeStoreError(w, err)
			return
		}
		response.Created(w, "secret created", created)
	case http.MethodDelete:
		if err := a.Store.RevokeSecret(r.Context(), r.URL.Query().Get("id")); err != nil {
			writeStoreError(w, err)
			return
		}
		response.OK(w, "secret revoked", map[string]string{"id": r.URL.Query().Get("id")})
	default:
		response.Error(w, http.StatusMethodNotAllowed, "method not allowed")
	}
}

func (a *API) AuditEvents(w http.ResponseWriter, r *http.Request) {
	id, ok := a.tenantScope(r)
	if !ok {
		response.Unauthorized(w, "tenant required")
		return
	}
	switch r.Method {
	case http.MethodGet:
		limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
		events, err := a.Store.ListAuditEvents(r.Context(), id.TenantID, limit)
		if err != nil {
			writeStoreError(w, err)
			return
		}
		response.OK(w, "audit events retrieved", map[string]any{"count": len(events), "events": events})
	case http.MethodPost:
		var event models.AuditEvent
		if err := decodeJSON(r, &event); err != nil {
			response.BadRequest(w, "invalid json payload")
			return
		}
		event.TenantID = id.TenantID
		event.UserID = id.UserID
		event.IPAddress = r.RemoteAddr
		event.UserAgent = r.UserAgent()
		created, err := a.Store.CreateAuditEvent(r.Context(), &event)
		if err != nil {
			writeStoreError(w, err)
			return
		}
		response.Created(w, "audit event created", created)
	default:
		response.Error(w, http.StatusMethodNotAllowed, "method not allowed")
	}
}

func (a *API) FraudRules(w http.ResponseWriter, r *http.Request) {
	id, ok := a.tenantScope(r)
	if !ok {
		response.Unauthorized(w, "tenant required")
		return
	}
	switch r.Method {
	case http.MethodGet:
		rules, err := a.Store.ListFraudRules(r.Context(), id.TenantID)
		if err != nil {
			writeStoreError(w, err)
			return
		}
		response.OK(w, "fraud rules retrieved", map[string]any{"count": len(rules), "rules": rules})
	case http.MethodPost:
		var rule models.FraudRule
		if err := decodeJSON(r, &rule); err != nil {
			response.BadRequest(w, "invalid json payload")
			return
		}
		rule.TenantID = id.TenantID
		created, err := a.Store.CreateFraudRule(r.Context(), &rule)
		if err != nil {
			writeStoreError(w, err)
			return
		}
		response.Created(w, "fraud rule created", created)
	case http.MethodPut:
		var rule models.FraudRule
		if err := decodeJSON(r, &rule); err != nil {
			response.BadRequest(w, "invalid json payload")
			return
		}
		rule.TenantID = id.TenantID
		updated, err := a.Store.UpdateFraudRule(r.Context(), &rule)
		if err != nil {
			writeStoreError(w, err)
			return
		}
		response.OK(w, "fraud rule updated", updated)
	case http.MethodDelete:
		if err := a.Store.DeleteFraudRule(r.Context(), r.URL.Query().Get("id")); err != nil {
			writeStoreError(w, err)
			return
		}
		response.OK(w, "fraud rule deleted", map[string]string{"id": r.URL.Query().Get("id")})
	default:
		response.Error(w, http.StatusMethodNotAllowed, "method not allowed")
	}
}

func (a *API) FraudAnalyze(w http.ResponseWriter, r *http.Request) {
	id, ok := a.tenantScope(r)
	if !ok {
		response.Unauthorized(w, "tenant required")
		return
	}
	if r.Method != http.MethodPost {
		response.Error(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}
	var attributes map[string]any
	if err := decodeJSON(r, &attributes); err != nil {
		response.BadRequest(w, "invalid json payload")
		return
	}
	rules, err := a.Store.ListFraudRules(r.Context(), id.TenantID)
	if err != nil {
		writeStoreError(w, err)
		return
	}

	score := 0
	matched := []string{}
	blocked := false
	for _, rule := range rules {
		if !rule.Enabled {
			continue
		}
		ok, err := engine.EvaluateRule(attributes, rule.Field, rule.Operator, rule.Threshold)
		if err != nil || !ok {
			continue
		}
		score += rule.Score
		matched = append(matched, rule.Name)
		event := &models.FraudViolation{
			TenantID: id.TenantID,
			UserID:   id.UserID,
			RuleID:   rule.ID,
			Score:    float64(rule.Score),
			Reason:   rule.Name,
			Status:   rule.Action,
			Metadata: attributes,
		}
		if strings.EqualFold(rule.Action, "block") {
			blocked = true
			event.Level = models.FraudRiskCritical
		} else if rule.Score >= 40 {
			event.Level = models.FraudRiskHigh
		} else {
			event.Level = models.FraudRiskMedium
		}
		_, _ = a.Store.CreateFraudEvent(r.Context(), event)
	}

	level := "low"
	if blocked || score >= a.Store.IntSetting(r.Context(), id.TenantID, "engine.fraud.high_score", 70) {
		level = "high"
	} else if score >= a.Store.IntSetting(r.Context(), id.TenantID, "engine.fraud.medium_score", 40) {
		level = "medium"
	}

	response.OK(w, "fraud analysis completed", map[string]any{
		"risk_score":  score,
		"risk_level":  level,
		"blocked":     blocked,
		"matched":     matched,
		"analyzed_at": time.Now().UTC(),
	})
}

func (a *API) FraudViolations(w http.ResponseWriter, r *http.Request) {
	id, ok := a.tenantScope(r)
	if !ok {
		response.Unauthorized(w, "tenant required")
		return
	}
	if r.Method != http.MethodGet {
		response.Error(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}
	events, err := a.Store.ListFraudEvents(r.Context(), id.TenantID, 100)
	if err != nil {
		writeStoreError(w, err)
		return
	}
	response.OK(w, "fraud events retrieved", map[string]any{"count": len(events), "events": events})
}

func (a *API) Settings(w http.ResponseWriter, r *http.Request) {
	id, ok := a.tenantScope(r)
	if !ok {
		response.Unauthorized(w, "tenant required")
		return
	}
	switch r.Method {
	case http.MethodGet:
		items, err := a.Store.ListSettings(r.Context(), id.TenantID)
		if err != nil {
			writeStoreError(w, err)
			return
		}
		response.OK(w, "settings retrieved", map[string]any{"count": len(items), "settings": items})
	case http.MethodPut, http.MethodPost:
		var setting models.Setting
		if err := decodeJSON(r, &setting); err != nil {
			response.BadRequest(w, "invalid json payload")
			return
		}
		if setting.TenantID == "" {
			setting.TenantID = id.TenantID
		}
		saved, err := a.Store.UpsertSetting(r.Context(), &setting)
		if err != nil {
			writeStoreError(w, err)
			return
		}
		response.OK(w, "setting saved", saved)
	default:
		response.Error(w, http.StatusMethodNotAllowed, "method not allowed")
	}
}

func (a *API) audit(r *http.Request, id identity.Identity, action, resource, message string) {
	event := &models.AuditEvent{
		TenantID:  id.TenantID,
		UserID:    id.UserID,
		Action:    action,
		Resource:  resource,
		Message:   message,
		IPAddress: r.RemoteAddr,
		UserAgent: r.UserAgent(),
		Severity:  models.AuditSeverityInfo,
		Status:    models.AuditStatusSuccess,
	}
	_, _ = a.Store.CreateAuditEvent(r.Context(), event)
}
