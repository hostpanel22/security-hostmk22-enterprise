package identity

import "context"

type contextKey string

const key contextKey = "identity"

type Identity struct {
	UserID   string
	TenantID string
	APIKeyID string
}

func WithIdentity(ctx context.Context, id Identity) context.Context {
	return context.WithValue(ctx, key, id)
}

func FromContext(ctx context.Context) (Identity, bool) {
	value, ok := ctx.Value(key).(Identity)
	return value, ok
}
