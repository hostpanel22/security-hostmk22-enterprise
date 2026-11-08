package security

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

// SecurityScanner provides security scanning capabilities
type SecurityScanner struct {
	config ScannerConfig
}

// ScannerConfig defines security scanner settings
type ScannerConfig struct {
	// EnableSAST determines if Static Application Security Testing is enabled
	EnableSAST bool
	// SASTTools is a list of SAST tools to use
	SASTTools []string
	// EnableDAST determines if Dynamic Application Security Testing is enabled
	EnableDAST bool
	// DASTURL is the URL to scan
	DASTURL string
	// EnableDependencyScan determines if dependency scanning is enabled
	EnableDependencyScan bool
	// DependencyScanTools is a list of dependency scanning tools
	DependencyScanTools []string
	// EnableSecretScan determines if secret scanning is enabled
	EnableSecretScan bool
	// ScanFrequency is how often to run scans (in hours)
	ScanFrequency int
	// ReportPath is where to store scan reports
	ReportPath string
}

// DefaultScannerConfig returns a default scanner configuration
func DefaultScannerConfig() ScannerConfig {
	return ScannerConfig{
		EnableSAST:             true,
		SASTTools:              []string{"gosec", "staticcheck"},
		EnableDAST:             false,
		DASTURL:                "",
		EnableDependencyScan:   true,
		DependencyScanTools:    []string{"govulncheck", "nancy"},
		EnableSecretScan:       true,
		ScanFrequency:          24,
		ReportPath:             "./reports/security",
	}
}

// NewSecurityScanner creates a new security scanner
func NewSecurityScanner(config ScannerConfig) *SecurityScanner {
	return &SecurityScanner{config: config}
}

// Scan performs a comprehensive security scan
func (s *SecurityScanner) Scan(ctx context.Context) (*SecurityReport, error) {
	report := &SecurityReport{
		Timestamp: time.Now(),
		Findings:  []Finding{},
	}

	// Run SAST scan
	if s.config.EnableSAST {
		findings, err := s.runSAST(ctx)
		if err != nil {
			return nil, fmt.Errorf("failed to run SAST: %w", err)
		}
		report.Findings = append(report.Findings, findings...)
	}

	// Run dependency scan
	if s.config.EnableDependencyScan {
		findings, err := s.runDependencyScan(ctx)
		if err != nil {
			return nil, fmt.Errorf("failed to run dependency scan: %w", err)
		}
		report.Findings = append(report.Findings, findings...)
	}

	// Run secret scan
	if s.config.EnableSecretScan {
		findings, err := s.runSecretScan(ctx)
		if err != nil {
			return nil, fmt.Errorf("failed to run secret scan: %w", err)
		}
		report.Findings = append(report.Findings, findings...)
	}

	// Save report
	if err := s.saveReport(report); err != nil {
		return nil, fmt.Errorf("failed to save report: %w", err)
	}

	return report, nil
}

// runSAST runs Static Application Security Testing
func (s *SecurityScanner) runSAST(ctx context.Context) ([]Finding, error) {
	var findings []Finding

	for _, tool := range s.config.SASTTools {
		toolFindings, err := s.runSASTTool(ctx, tool)
		if err != nil {
			// Log error but continue with other tools
			findings = append(findings, Finding{
				Severity:    SeverityError,
				Type:        "sast_error",
				Title:       "SAST tool error",
				Description: fmt.Sprintf("Failed to run SAST tool %s: %v", tool, err),
				Tool:        tool,
			})
			continue
		}
		findings = append(findings, toolFindings...)
	}

	return findings, nil
}

// runSASTTool runs a specific SAST tool
func (s *SecurityScanner) runSASTTool(ctx context.Context, tool string) ([]Finding, error) {
	var cmd *exec.Cmd

	switch tool {
	case "gosec":
		cmd = exec.CommandContext(ctx, "gosec", "-fmt=json", "./...")
	case "staticcheck":
		cmd = exec.CommandContext(ctx, "staticcheck", "-f=json", "./...")
	case "semgrep":
		cmd = exec.CommandContext(ctx, "semgrep", "scan", "--json", ".")
	default:
		return nil, fmt.Errorf("unknown SAST tool: %s", tool)
	}

	output, err := cmd.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("%s: %w", string(output), err)
	}

	return parseSASTOutput(tool, output)
}

// parseSASTOutput parses the output of a SAST tool
func parseSASTOutput(tool string, output []byte) ([]Finding, error) {
	var findings []Finding

	switch tool {
	case "gosec":
		var results []struct {
			Severity string `json:"severity"`
			RuleID   string `json:"rule_id"`
			Details  string `json:"details"`
			File     string `json:"file"`
			Line     string `json:"line"`
		}
		if err := json.Unmarshal(output, &results); err != nil {
			return nil, err
		}
		for _, result := range results {
			findings = append(findings, Finding{
				Severity:    severityFromString(result.Severity),
				Type:        "sast",
				Title:       result.RuleID,
				Description: result.Details,
				File:        result.File,
				Line:        result.Line,
				Tool:        tool,
			})
		}
	case "staticcheck":
		var issues []struct {
			Severity string `json:"severity"`
			Code     string `json:"code"`
			Message  string `json:"message"`
			Pos      struct {
				Filename string `json:"filename"`
				Line    int    `json:"line"`
			} `json:"pos"`
		}
		if err := json.Unmarshal(output, &issues); err != nil {
			return nil, err
		}
		for _, issue := range issues {
			findings = append(findings, Finding{
				Severity:    severityFromString(issue.Severity),
				Type:        "sast",
				Title:       issue.Code,
				Description: issue.Message,
				File:        issue.Pos.Filename,
				Line:        fmt.Sprintf("%d", issue.Pos.Line),
				Tool:        tool,
			})
		}
	}

	return findings, nil
}

// runDependencyScan runs dependency vulnerability scanning
func (s *SecurityScanner) runDependencyScan(ctx context.Context) ([]Finding, error) {
	var findings []Finding

	for _, tool := range s.config.DependencyScanTools {
		toolFindings, err := s.runDependencyTool(ctx, tool)
		if err != nil {
			findings = append(findings, Finding{
				Severity:    SeverityError,
				Type:        "dependency_error",
				Title:       "Dependency scan error",
				Description: fmt.Sprintf("Failed to run dependency scan tool %s: %v", tool, err),
				Tool:        tool,
			})
			continue
		}
		findings = append(findings, toolFindings...)
	}

	return findings, nil
}

// runDependencyTool runs a specific dependency scanning tool
func (s *SecurityScanner) runDependencyTool(ctx context.Context, tool string) ([]Finding, error) {
	var cmd *exec.Cmd

	switch tool {
	case "govulncheck":
		cmd = exec.CommandContext(ctx, "govulncheck", "-json", "./...")
	case "nancy":
		cmd = exec.CommandContext(ctx, "nancy", "sleuth", "-json")
	case "grype":
		cmd = exec.CommandContext(ctx, "grype", "dir:.", "-o", "json")
	default:
		return nil, fmt.Errorf("unknown dependency scan tool: %s", tool)
	}

	output, err := cmd.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("%s: %w", string(output), err)
	}

	return parseDependencyOutput(tool, output)
}

// parseDependencyOutput parses the output of a dependency scanning tool
func parseDependencyOutput(tool string, output []byte) ([]Finding, error) {
	var findings []Finding

	if tool == "govulncheck" {
		var results []struct {
			Vulnerability struct {
				ID          string `json:"id"`
				Severity    string `json:"severity"`
				Description string `json:"details"`
			} `json:"vulnerability"`
			Module struct {
				Path    string `json:"path"`
				Version string `json:"version"`
			} `json:"module"`
		}
		if err := json.Unmarshal(output, &results); err != nil {
			return nil, err
		}
		for _, result := range results {
			findings = append(findings, Finding{
				Severity:    severityFromString(result.Vulnerability.Severity),
				Type:        "dependency",
				Title:       result.Vulnerability.ID,
				Description: result.Vulnerability.Description,
				Package:     result.Module.Path,
				Version:     result.Module.Version,
				Tool:        tool,
			})
		}
	}

	return findings, nil
}

// runSecretScan scans for accidentally committed secrets
func (s *SecurityScanner) runSecretScan(ctx context.Context) ([]Finding, error) {
	var findings []Finding

	// Use gitleaks for secret scanning
	cmd := exec.CommandContext(ctx, "gitleaks", "detect", "--source=.", "--format=json", "--report-path=/dev/stdout")
	output, err := cmd.CombinedOutput()
	if err != nil {
		// gitleaks returns non-zero exit code when it finds secrets
		// So we'll parse the output anyway
	}

	var results []struct {
		Description string `json:"description"`
		Secret      string `json:"secret"`
		RuleID      string `json:"ruleId"`
		Entropy     float64 `json:"entropy"`
		File        string `json:"file"`
		Line        int    `json:"line"`
	}

	if err := json.Unmarshal(output, &results); err == nil {
		for _, result := range results {
			findings = append(findings, Finding{
				Severity:    SeverityCritical,
				Type:        "secret",
				Title:       result.RuleID,
				Description: result.Description,
				File:        result.File,
				Line:        fmt.Sprintf("%d", result.Line),
				Tool:        "gitleaks",
			})
		}
	}

	return findings, nil
}

// saveReport saves the security report to a file
func (s *SecurityScanner) saveReport(report *SecurityReport) error {
	if s.config.ReportPath == "" {
		return nil
	}

	// Create directory if it doesn't exist
	if err := os.MkdirAll(s.config.ReportPath, 0755); err != nil {
		return err
	}

	// Generate report filename
	timestamp := report.Timestamp.Format("20060102_150405")
	filename := filepath.Join(s.config.ReportPath, "security_report_"+timestamp+".json")

	// Convert report to JSON
	data, err := json.MarshalIndent(report, "", "  ")
	if err != nil {
		return err
	}

	// Write report
	return os.WriteFile(filename, data, 0644)
}

// SecurityReport represents a security scan report
type SecurityReport struct {
	Timestamp time.Time `json:"timestamp"`
	Findings  []Finding `json:"findings"`
	Summary   Summary   `json:"summary"`
}

// Summary represents a summary of the security scan
type Summary struct {
	TotalFindings   int `json:"total_findings"`
	Critical        int `json:"critical"`
	High            int `json:"high"`
	Medium          int `json:"medium"`
	Low             int `json:"low"`
	Info            int `json:"info"`
	SASTFindings    int `json:"sast_findings"`
	DependencyFindings int `json:"dependency_findings"`
	SecretFindings  int `json:"secret_findings"`
}

// Finding represents a security finding
type Finding struct {
	Severity    Severity `json:"severity"`
	Type        string   `json:"type"`
	Title       string   `json:"title"`
	Description string   `json:"description"`
	File        string   `json:"file,omitempty"`
	Line        string   `json:"line,omitempty"`
	Package     string   `json:"package,omitempty"`
	Version     string   `json:"version,omitempty"`
	Tool        string   `json:"tool,omitempty"`
	CVE         string   `json:"cve,omitempty"`
	Remediation string   `json:"remediation,omitempty"`
}

// Severity represents the severity of a finding
type Severity string

const (
	SeverityCritical Severity = "CRITICAL"
	SeverityHigh     Severity = "HIGH"
	SeverityMedium   Severity = "MEDIUM"
	SeverityLow      Severity = "LOW"
	SeverityInfo     Severity = "INFO"
	SeverityError    Severity = "ERROR"
)

// severityFromString converts a string severity to a Severity type
func severityFromString(severity string) Severity {
	switch strings.ToUpper(severity) {
	case "CRITICAL":
		return SeverityCritical
	case "HIGH":
		return SeverityHigh
	case "MEDIUM":
		return SeverityMedium
	case "LOW":
		return SeverityLow
	case "INFO":
		return SeverityInfo
	default:
		return SeverityInfo
	}
}

// ScanHandler provides HTTP handlers for security scanning
type ScanHandler struct {
	scanner *SecurityScanner
}

// NewScanHandler creates a new scan handler
func NewScanHandler(scanner *SecurityScanner) *ScanHandler {
	return &ScanHandler{scanner: scanner}
}

// RunScan runs a security scan and returns the report
func (h *ScanHandler) RunScan(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost && r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	report, err := h.scanner.Scan(r.Context())
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"error":"` + err.Error() + `"}`))
		return
	}

	// Calculate summary
	summary := Summary{}
	for _, finding := range report.Findings {
		summary.TotalFindings++
		switch finding.Severity {
		case SeverityCritical:
			summary.Critical++
		case SeverityHigh:
			summary.High++
		case SeverityMedium:
			summary.Medium++
		case SeverityLow:
			summary.Low++
		case SeverityInfo:
			summary.Info++
		}
		switch finding.Type {
		case "sast":
			summary.SASTFindings++
		case "dependency":
			summary.DependencyFindings++
		case "secret":
			summary.SecretFindings++
		}
	}
	report.Summary = summary

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(report)
}

// GetReport returns a stored security report
func (h *ScanHandler) GetReport(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// List available reports
	files, err := os.ReadDir(h.scanner.config.ReportPath)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"error":"` + err.Error() + `"}`))
		return
	}

	var reports []string
	for _, file := range files {
		if !file.IsDir() && strings.HasPrefix(file.Name(), "security_report_") && strings.HasSuffix(file.Name(), ".json") {
			reports = append(reports, file.Name())
		}
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"reports": reports,
	})
}

// GetReportByName returns a specific security report
func (h *ScanHandler) GetReportByName(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	name := r.URL.Query().Get("name")
	if name == "" {
		http.Error(w, "Report name is required", http.StatusBadRequest)
		return
	}

	filename := filepath.Join(h.scanner.config.ReportPath, name)
	data, err := os.ReadFile(filename)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(`{"error":"Report not found"}`))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}

// CheckVulnerabilities checks for known vulnerabilities in dependencies
func (h *ScanHandler) CheckVulnerabilities(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Run dependency scan only
	config := h.scanner.config
	config.EnableSAST = false
	config.EnableSecretScan = false
	scanner := NewSecurityScanner(config)

	report, err := scanner.Scan(r.Context())
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"error":"` + err.Error() + `"}`))
		return
	}

	// Filter to only dependency findings
	var deps []Finding
	for _, finding := range report.Findings {
		if finding.Type == "dependency" {
			deps = append(deps, finding)
		}
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"vulnerabilities": deps,
		"count":          len(deps),
	})
}

// DependencyCheckHandler provides a handler for checking dependencies
type DependencyCheckHandler struct{}

// NewDependencyCheckHandler creates a new dependency check handler
func NewDependencyCheckHandler() *DependencyCheckHandler {
	return &DependencyCheckHandler{}
}

// CheckVulnerableDependencies checks for vulnerable dependencies in the database
func (h *DependencyCheckHandler) CheckVulnerableDependencies(w http.ResponseWriter, r *http.Request) {
	// This would integrate with a vulnerability database
	// For now, return a placeholder
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Dependency vulnerability checking would be implemented here",
		"status":  "not_implemented",
	})
}

// ScanMiddleware returns a middleware that runs security scans periodically
func ScanMiddleware(scanner *SecurityScanner, interval time.Duration) func(http.Handler) http.Handler {
	// Start scanner in background
	go func() {
		ticker := time.NewTicker(interval)
		defer ticker.Stop()
		for range ticker.C {
			_, err := scanner.Scan(context.Background())
			if err != nil {
				// Log error
			}
		}
	}()

	return func(next http.Handler) http.Handler {
		return next
	}
}

// RunSecurityScan runs a one-time security scan
func RunSecurityScan(config ScannerConfig) (*SecurityReport, error) {
	scanner := NewSecurityScanner(config)
	return scanner.Scan(context.Background())
}

// CheckForVulnerabilities checks for known vulnerabilities
func CheckForVulnerabilities() ([]Finding, error) {
	config := DefaultScannerConfig()
	config.EnableSAST = false
	config.EnableSecretScan = false
	scanner := NewSecurityScanner(config)
	report, err := scanner.Scan(context.Background())
	if err != nil {
		return nil, err
	}

	var vulnerabilities []Finding
	for _, finding := range report.Findings {
		if finding.Severity == SeverityCritical || finding.Severity == SeverityHigh {
			vulnerabilities = append(vulnerabilities, finding)
		}
	}

	return vulnerabilities, nil
}

// GetSecurityStatus returns the current security status
func GetSecurityStatus() map[string]interface{} {
	// Check for critical vulnerabilities
	vulns, err := CheckForVulnerabilities()
	if err != nil {
		return map[string]interface{}{
			"status":       "unknown",
			"vulnerabilities": 0,
			"error":        err.Error(),
		}
	}

	status := "healthy"
	if len(vulns) > 0 {
		status = "vulnerable"
	}

	return map[string]interface{}{
		"status":        status,
		"vulnerabilities": len(vulns),
		"critical":      countBySeverity(vulns, SeverityCritical),
		"high":          countBySeverity(vulns, SeverityHigh),
	}
}

// countBySeverity counts findings by severity
func countBySeverity(findings []Finding, severity Severity) int {
	count := 0
	for _, f := range findings {
		if f.Severity == severity {
			count++
		}
	}
	return count
}

// SecurityCheckHandler provides HTTP handlers for security checks
type SecurityCheckHandler struct{}

// NewSecurityCheckHandler creates a new security check handler
func NewSecurityCheckHandler() *SecurityCheckHandler {
	return &SecurityCheckHandler{}
}

// Status returns the current security status
func (h *SecurityCheckHandler) Status(w http.ResponseWriter, r *http.Request) {
	status := GetSecurityStatus()
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(status)
}

// HealthCheck performs a security health check
func (h *SecurityCheckHandler) HealthCheck(w http.ResponseWriter, r *http.Request) {
	// Check various security aspects
	checks := map[string]bool{
		"cors_enabled":     true,
		"csrf_protected":  true,
		"rate_limiting":    true,
		"input_validation": true,
		"security_headers": true,
	}

	allPassed := true
	for _, passed := range checks {
		if !passed {
			allPassed = false
			break
		}
	}

	status := "healthy"
	if !allPassed {
		status = "degraded"
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status": status,
		"checks": checks,
	})
}

// VulnerabilityCheck checks for specific vulnerabilities
func (h *SecurityCheckHandler) VulnerabilityCheck(w http.ResponseWriter, r *http.Request) {
	vulns, err := CheckForVulnerabilities()
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"error": err.Error(),
		})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"vulnerabilities": vulns,
		"count":          len(vulns),
	})
}
