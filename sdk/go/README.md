# HostMK Security Framework Go SDK

Typed Go client for the HTTP endpoints registered by this repository. The SDK module is separate from the server module and targets the Go version declared in `go.mod`.

## Install

```sh
go get github.com/azezaz/HostMK-Security-Framework/sdk/go
```

## Create a client

```go
ctx := context.Background()
sdk, err := hostmk.NewClient("https://security.example.com", client.WithAPIKey(apiKey))
if err != nil {
	return err
}

health, err := sdk.Monitoring.Health(ctx)
if err != nil {
	return err
}
fmt.Println(health.Status)
```

The base URL is the server origin or an optional reverse-proxy prefix. Resource paths include `/api/v1`; a base URL already ending in `/api/v1` is normalized without duplicating that segment. Use HTTPS outside local development.

## Authentication

```go
sdk, err := hostmk.NewClient(apiOrigin)
if err != nil {
	return err
}

tokens, err := sdk.Auth.Login(ctx, auth.LoginRequest{
	TenantSlug: tenantSlug,
	Email:      email,
	Password:   password,
	MFACode:    mfaCode,
})
if err != nil {
	return err
}

// Login stores the returned access token in the shared client.
tokens, err = sdk.Auth.Refresh(ctx, tokens.RefreshToken)
```

`client.WithAPIKey` sends `X-API-Key`; `client.WithBearerToken` sends `Authorization: Bearer`. If both are configured, the API key takes precedence. Successful Login/Refresh switches the client to the returned bearer token. The server has no logout endpoint; `sdk.Auth.ClearToken()` clears SDK credentials locally.

Bootstrap is available only while the backend has no users:

```go
result, err := sdk.Auth.Bootstrap(ctx, auth.BootstrapRequest{
	TenantName: tenantName,
	TenantSlug: tenantSlug,
	Name:       name,
	Email:      email,
	Password:   password,
})
```

## Services

- `Auth`: bootstrap, login, refresh, local token clear.
- `Users`: list/create/update and activate, suspend, verify-email, enable-MFA, delete actions.
- `Tenants`: list/create/update/activate/suspend.
- `Sessions`: list/revoke.
- `Roles`: list/create/update/delete.
- `Permissions`: list/create/delete/check.
- `APIKeys`: list/create/rotate/revoke. The plaintext key is returned only by create/rotate; store it securely.
- `Secrets`: list/create/rotate/history/revoke. Secret values are sent only to create/rotate and are not returned by backend persistence responses; history is masked.
- `Audit`: list (the backend accepts `limit`, maximum 500) and create.
- `Fraud`: list/create/update/delete rules, analyze attributes, list violations. The backend provides no violation mutation endpoint.
- `Settings`: list/save.
- `Monitoring`: health, ready, version.

All network methods accept `context.Context`. List methods return `client.ListResult[T]` with the backend count and items. Non-2xx responses return `*client.Error` with status code, message and raw body:

```go
var apiErr *client.Error
if errors.As(err, &apiErr) {
	log.Printf("API status=%d: %s", apiErr.StatusCode, apiErr.Message)
}
```

The SDK intentionally does not add server operations that are absent from the registered API (for example logout, role assignment APIs beyond the user action contract, audit export, or fraud-violation updates).

## Test

```sh
cd sdk/go
go test ./...
```
