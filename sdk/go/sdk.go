// Package hostmk composes typed services for the HostMK Security Framework API.
package hostmk

import (
	"github.com/azezaz/HostMK-Security-Framework/sdk/go/apikeys"
	"github.com/azezaz/HostMK-Security-Framework/sdk/go/audit"
	"github.com/azezaz/HostMK-Security-Framework/sdk/go/auth"
	"github.com/azezaz/HostMK-Security-Framework/sdk/go/client"
	"github.com/azezaz/HostMK-Security-Framework/sdk/go/fraud"
	"github.com/azezaz/HostMK-Security-Framework/sdk/go/monitoring"
	"github.com/azezaz/HostMK-Security-Framework/sdk/go/permissions"
	"github.com/azezaz/HostMK-Security-Framework/sdk/go/rbac"
	"github.com/azezaz/HostMK-Security-Framework/sdk/go/secrets"
	"github.com/azezaz/HostMK-Security-Framework/sdk/go/sessions"
	"github.com/azezaz/HostMK-Security-Framework/sdk/go/settings"
	"github.com/azezaz/HostMK-Security-Framework/sdk/go/tenancy"
	"github.com/azezaz/HostMK-Security-Framework/sdk/go/users"
)

type Client struct {
	HTTP        *client.Client
	Auth        *auth.Service
	Users       *users.Service
	Tenants     *tenancy.Service
	Sessions    *sessions.Service
	Roles       *rbac.Roles
	Permissions *permissions.Service
	APIKeys     *apikeys.Service
	Secrets     *secrets.Service
	Audit       *audit.Service
	Fraud       *fraud.Service
	Settings    *settings.Service
	Monitoring  *monitoring.Service
}

func NewClient(baseURL string, options ...client.Option) (*Client, error) {
	httpClient, err := client.New(baseURL, options...)
	if err != nil {
		return nil, err
	}
	return &Client{
		HTTP:        httpClient,
		Auth:        auth.New(httpClient),
		Users:       users.New(httpClient),
		Tenants:     tenancy.New(httpClient),
		Sessions:    sessions.New(httpClient),
		Roles:       rbac.NewRoles(httpClient),
		Permissions: permissions.New(httpClient),
		APIKeys:     apikeys.New(httpClient),
		Secrets:     secrets.New(httpClient),
		Audit:       audit.New(httpClient),
		Fraud:       fraud.New(httpClient),
		Settings:    settings.New(httpClient),
		Monitoring:  monitoring.New(httpClient),
	}, nil
}
