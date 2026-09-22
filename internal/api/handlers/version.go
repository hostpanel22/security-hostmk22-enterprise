package handlers

import (
	"net/http"
	"runtime"
	"time"

	"github.com/azezaz/HostMK-Security-Framework/internal/api/response"
)

const (
	ApplicationName    = "HostMK Security Framework"
	ApplicationVersion = "0.1.0"
	ApplicationLicense = "MIT"
	ApplicationAuthor  = "HostMK"
)

type VersionInfo struct {
	Name        string   `json:"name"`
	Version     string   `json:"version"`
	Environment string   `json:"environment"`
	BuildTime   string   `json:"build_time"`
	GoVersion   string   `json:"go_version"`
	OS          string   `json:"os"`
	Arch        string   `json:"arch"`
	License     string   `json:"license"`
	Author      string   `json:"author"`
	Modules     []string `json:"modules"`
}

func Version(
	w http.ResponseWriter,
	r *http.Request,
) {
	if r.Method != http.MethodGet {
		response.Error(
			w,
			http.StatusMethodNotAllowed,
			"method not allowed",
		)
		return
	}

	response.OK(
		w,
		"version information",
		buildVersionInfo(),
	)
}

func buildVersionInfo() VersionInfo {
	return VersionInfo{
		Name:        ApplicationName,
		Version:     ApplicationVersion,
		Environment: "production",
		BuildTime:   time.Now().UTC().Format(time.RFC3339),
		GoVersion:   runtime.Version(),
		OS:          runtime.GOOS,
		Arch:        runtime.GOARCH,
		License:     ApplicationLicense,
		Author:      ApplicationAuthor,
		Modules: []string{
			"audit",
			"fraud",
			"rbac",
			"permissions",
			"secrets",
			"tenancy",
		},
	}
}
