# Security Package

## Overview

The `core/security` package provides comprehensive security middleware and utilities for protecting your Go applications at the host level. It includes protection against common web vulnerabilities, security headers, rate limiting, brute force attacks, and more.

## Features

### Security Middleware

| Middleware | Description | Default Status |
|------------|-------------|----------------|
| **CORS** | Cross-Origin Resource Sharing protection | Enabled |
| **CSRF** | Cross-Site Request Forgery protection | Enabled |
| **Security Headers** | Adds security headers (CSP, HSTS, etc.) | Enabled |
| **Rate Limiting** | Prevents abuse and DoS attacks | Enabled |
| **IP Filtering** | Whitelist/blacklist IP addresses | Disabled |
| **Request Size Limit** | Limits request body and file upload sizes | Enabled |
| **Brute Force Protection** | Prevents brute force attacks | Enabled |
| **Input Validation** | Validates and sanitizes user input | Enabled |

### Security Scanning

- **SAST** (Static Application Security Testing) - Scans source code for vulnerabilities
- **Dependency Scanning** - Checks for vulnerable dependencies
- **Secret Scanning** - Detects accidentally committed secrets
- **Vulnerability Database** - Checks against known CVEs

## Quick Start

### Basic Usage

```go
package main

import (
    "net/http"
    "github.com/hostpanel22/security-hostmk22-enterprise/core/security"
)

func main() {
    // Create a simple handler
    handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        w.Write([]byte("Hello, World!"))
    })

    // Wrap with security middleware
    secureHandler := security.BuildSecureMiddleware(security.DefaultSecurityConfig())(handler)

    // Start server
    http.ListenAndServe(":8080", secureHandler)
}
```

### Using Predefined Security Levels

```go
// Minimal security (basic protection)
secureHandler := security.WrapWithSecurity(handler, security.SecurityLevelMinimal)

// Standard security (recommended for most applications)
secureHandler := security.WrapWithSecurity(handler, security.SecurityLevelStandard)

// Production security (maximum protection)
secureHandler := security.WrapWithSecurity(handler, security.SecurityLevelProduction)

// Paranoid security (maximum protection with strict defaults)
secureHandler := security.WrapWithSecurity(handler, security.SecurityLevelParanoid)
```

### Custom Configuration

```go
config := security.SecurityConfig{
    CORS: security.CORSConfig{
        Enabled:          true,
        AllowedOrigins:   []string{"https://myapp.com", "https://api.myapp.com"},
        AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
        AllowedHeaders:   []string{"Authorization", "Content-Type"},
        AllowCredentials: true,
        MaxAge:           86400,
    },
    CSRF: security.CSRFConfig{
        Enabled:        true,
        TokenLength:    32,
        TokenExpiry:    24 * time.Hour,
        CookieSecure:   true,
        CookieHttpOnly: true,
    },
    RateLimit: security.RateLimitConfig{
        Enabled:          true,
        RequestsPerMinute: 100,
        BurstSize:        20,
        TrustProxy:       true,
    },
    Headers: security.HeadersConfig{
        Enabled:             true,
        ContentSecurityPolicy: security.StrictCSP(),
        HSTS:                 "max-age=63072000; includeSubDomains; preload",
        XFrameOptions:        "DENY",
        XContentTypeOptions:  "nosniff",
        ReferrerPolicy:       "strict-origin-when-cross-origin",
    },
}

secureHandler := security.BuildSecureMiddleware(config)(handler)
```

## Middleware Details

### CORS Middleware

```go
// Simple CORS with default settings
corsMiddleware := security.CORS(security.DefaultCORSConfig())

// Strict CORS for production
corsMiddleware := security.CORS(security.CORSConfig{
    Enabled:          true,
    AllowedOrigins:   []string{"https://myapp.com"},
    AllowedMethods:   []string{"GET", "POST"},
    AllowCredentials: true,
})

// Permissive CORS for development (use with caution!)
corsMiddleware := security.PermissiveCORS()
```

### CSRF Middleware

```go
// Create CSRF token store
store := security.NewMemoryCSRFTokenStore(32, 24*time.Hour)

// Create CSRF middleware
csrfMiddleware := security.CSRF(security.DefaultCSRFConfig(), store)

// Generate CSRF token for a user
token, err := store.Generate(userID)

// Validate CSRF token
isValid := store.Validate(userID, token)
```

### Security Headers Middleware

```go
// Default security headers
headersMiddleware := security.SecurityHeaders(security.DefaultHeadersConfig())

// Custom CSP
config := security.DefaultHeadersConfig()
config.ContentSecurityPolicy = security.BuildCSP(map[string][]string{
    "default-src": {"'self'"},
    "script-src": {"'self'", "https://cdn.jsdelivr.net"},
    "style-src":  {"'self'", "'unsafe-inline'"},
    "img-src":    {"'self'", "data:", "https:"},
})
headersMiddleware := security.SecurityHeaders(config)

// Predefined CSP policies
config.ContentSecurityPolicy = security.StrictCSP()      // Maximum security
config.ContentSecurityPolicy = security.APICSP()         // For REST APIs
config.ContentSecurityPolicy = security.WebAppCSP()      // For web applications
```

### Rate Limiting Middleware

```go
// Simple rate limiting
rateLimitMiddleware := security.RateLimit(security.RateLimitConfig{
    Enabled:          true,
    RequestsPerMinute: 60,
    BurstSize:        10,
    TrustProxy:       true,
})

// With IP whitelist/blacklist
rateLimitMiddleware := security.RateLimit(security.RateLimitConfig{
    Enabled:    true,
    Whitelist:  []string{"127.0.0.1", "192.168.1.1"},
    Blacklist:  []string{"10.0.0.1"},
})
```

### IP Filtering Middleware

```go
// Whitelist only specific IPs
ipFilter := security.IPFilter(security.IPFilterConfig{
    Enabled:       true,
    Whitelist:     []string{"192.168.1.0/24", "10.0.0.1"},
    DefaultAction: "deny",
})

// Blacklist specific IPs
ipFilter := security.IPFilter(security.IPFilterConfig{
    Enabled:       true,
    Blacklist:     []string{"10.0.0.1", "192.168.1.100"},
    DefaultAction: "allow",
})

// Check if IP is in range
isInRange := security.IPInRange("192.168.1.5", "192.168.1.0/24")

// Check if IP is private
isPrivate := security.IsPrivateIP("192.168.1.1")
isPublic := security.IsPublicIP("8.8.8.8")
```

### Brute Force Protection

```go
// Create brute force store
store := security.NewMemoryBruteForceStore(5, 15*time.Minute, 5*time.Minute)

// Create brute force middleware
bfMiddleware := security.BruteForceProtection(security.BruteForceConfig{
    Enabled:         true,
    MaxAttempts:     5,
    LockoutDuration: 15 * time.Minute,
    ProtectedPaths:  []string{"/login", "/api/auth"},
}, store)

// Use in handler
func loginHandler(w http.ResponseWriter, r *http.Request) {
    handler := security.NewBruteForceHandler(store, config)
    
    // Check and record attempt
    if !handler.CheckAndRecord(w, r, isSuccess) {
        return // Already handled
    }
    
    // ... authentication logic ...
}
```

### Input Validation

```go
// Create input validator
validator := security.NewInputValidator(security.DefaultInputValidationConfig())

// Validate string
safeInput, err := validator.ValidateString(userInput)

// Validate query parameters
if err := validator.ValidateQueryParams(r); err != nil {
    // Handle error
}

// Validate form values
if err := validator.ValidateFormValues(r); err != nil {
    // Handle error
}

// Sanitize all input
security.SanitizeAll(r)

// Simple sanitization
safe := security.SanitizeInput(userInput)
```

### File Upload Security

```go
// Create file upload security checker
fileSecurity := security.NewFileUploadSecurity(security.DefaultRequestSizeConfig())

// Check file size
if err := fileSecurity.CheckFileSize(fileSize); err != nil {
    // File too large
}

// Check MIME type
if err := fileSecurity.CheckMIMEType(mimeType); err != nil {
    // Invalid MIME type
}

// Check file extension
if err := fileSecurity.CheckFileExtension(filename); err != nil {
    // Dangerous extension
}

// Check file content
if err := fileSecurity.CheckFileContent(content); err != nil {
    // Suspicious content
}
```

## Security Scanning

### Running Scans

```go
// Create scanner with default configuration
scanner := security.NewSecurityScanner(security.DefaultScannerConfig())

// Run comprehensive scan
report, err := scanner.Scan(context.Background())

// Check for vulnerabilities
vulns, err := security.CheckForVulnerabilities()

// Get security status
status := security.GetSecurityStatus()
```

### HTTP Handlers

```go
// Create scan handler
scanHandler := security.NewScanHandler(scanner)

// Add routes
mux.HandleFunc("/api/security/scan", scanHandler.RunScan)
mux.HandleFunc("/api/security/reports", scanHandler.GetReport)
mux.HandleFunc("/api/security/vulnerabilities", scanHandler.CheckVulnerabilities)

// Create security check handler
checkHandler := security.NewSecurityCheckHandler()
mux.HandleFunc("/api/security/status", checkHandler.Status)
mux.HandleFunc("/api/security/health", checkHandler.HealthCheck)
```

## Security Levels

The package provides predefined security levels for different environments:

### Minimal Security
- CORS: Disabled
- CSRF: Disabled
- Rate Limiting: Disabled
- IP Filtering: Disabled
- Brute Force: Disabled
- Security Headers: Basic

### Standard Security (Default)
- CORS: Enabled with safe defaults
- CSRF: Enabled
- Rate Limiting: Enabled (60 req/min)
- IP Filtering: Disabled
- Brute Force: Enabled (5 attempts, 15 min lockout)
- Security Headers: Comprehensive

### Production Security
- CORS: Enabled with specific origins
- CSRF: Enabled with secure cookies
- Rate Limiting: Enabled (100 req/min)
- IP Filtering: Disabled (handle at network level)
- Brute Force: Enabled
- Security Headers: Maximum protection

### Paranoid Security
- CORS: Disabled
- CSRF: Enabled
- Rate Limiting: Very strict (10 req/min)
- IP Filtering: Whitelist only
- Brute Force: Very strict (3 attempts, 1 hour lockout)
- Security Headers: Maximum protection

## Best Practices

### 1. Always Use HTTPS
```go
// In production, always use HTTPS
config.Headers.HSTS = "max-age=63072000; includeSubDomains; preload"
```

### 2. Validate All Input
```go
// Always validate and sanitize user input
safeInput, err := validator.ValidateString(userInput)
if err != nil {
    // Handle error
}
```

### 3. Use Secure Cookies
```go
config.CSRF.CookieSecure = true
config.CSRF.CookieHttpOnly = true
config.CSRF.CookieSameSite = "Lax"
```

### 4. Rate Limit All Endpoints
```go
// Apply rate limiting to all endpoints
config.RateLimit.Enabled = true
config.RateLimit.RequestsPerMinute = 100
```

### 5. Protect Against Brute Force
```go
// Protect authentication endpoints
config.BruteForce.Enabled = true
config.BruteForce.MaxAttempts = 5
config.BruteForce.LockoutDuration = 15 * time.Minute
```

### 6. Use Content Security Policy
```go
// Use a strict CSP
config.Headers.ContentSecurityPolicy = security.StrictCSP()
```

### 7. Regular Security Scans
```go
// Run regular security scans
scanner := security.NewSecurityScanner(config)
go func() {
    ticker := time.NewTicker(24 * time.Hour)
    for range ticker.C {
        report, _ := scanner.Scan(context.Background())
        // Log or alert on findings
    }
}()
```

## Configuration Reference

### SecurityConfig

```go
type SecurityConfig struct {
    CORS        CORSConfig
    CSRF        CSRFConfig
    RateLimit   RateLimitConfig
    IPFilter    IPFilterConfig
    RequestSize RequestSizeConfig
    Headers     HeadersConfig
    BruteForce  BruteForceConfig
    Session     SessionConfig
}
```

### CORSConfig

```go
type CORSConfig struct {
    Enabled           bool
    AllowedOrigins    []string
    AllowedMethods    []string
    AllowedHeaders    []string
    ExposedHeaders    []string
    AllowCredentials  bool
    MaxAge            time.Duration
    AllowPrivateNetwork bool
}
```

### CSRFConfig

```go
type CSRFConfig struct {
    Enabled         bool
    TokenLength     int
    TokenExpiry     time.Duration
    CookieName      string
    CookieSecure    bool
    CookieHttpOnly  bool
    CookieSameSite  string
    HeaderName      string
    SafeMethods     []string
}
```

### RateLimitConfig

```go
type RateLimitConfig struct {
    Enabled          bool
    RequestsPerMinute int
    BurstSize        int
    WindowDuration   time.Duration
    TrustProxy       bool
    Whitelist        []string
    Blacklist        []string
}
```

### IPFilterConfig

```go
type IPFilterConfig struct {
    Enabled       bool
    Whitelist     []string
    Blacklist     []string
    TrustProxy    bool
    DefaultAction string // "allow" or "deny"
}
```

### RequestSizeConfig

```go
type RequestSizeConfig struct {
    MaxBodySize         int64
    MaxHeaderSize       int64
    MaxMultipartMemory  int64
    MaxFileSize         int64
    AllowedMIMETypes    []string
    BlockedMIMETypes    []string
}
```

### HeadersConfig

```go
type HeadersConfig struct {
    Enabled                     bool
    ContentSecurityPolicy       string
    ContentSecurityPolicyReportOnly string
    HSTS                        string
    HSTSMaxAge                  int
    HSTSIncludeSubdomains       bool
    HSTSPreload                 bool
    XContentTypeOptions         string
    XFrameOptions               string
    XSSProtection               string
    ReferrerPolicy              string
    PermissionsPolicy           string
    ExpectCT                    string
    Server                      string
}
```

### BruteForceConfig

```go
type BruteForceConfig struct {
    Enabled          bool
    MaxAttempts      int
    LockoutDuration  time.Duration
    WindowDuration   time.Duration
    LockoutKeyPrefix string
    ProtectedPaths   []string
    TrustProxy       bool
}
```

## Error Handling

All security middleware returns appropriate HTTP status codes:

- `400 Bad Request` - Invalid input, missing CSRF token
- `403 Forbidden` - CSRF token invalid, IP not allowed, rate limited
- `429 Too Many Requests` - Rate limited, brute force lockout
- `413 Request Entity Too Large` - Request too large

## Performance Considerations

- **Memory Usage**: In-memory stores (CSRF, brute force) consume memory. For high-traffic applications, consider using Redis.
- **Rate Limiting**: The default in-memory rate limiter is suitable for small applications. For production, consider a distributed rate limiter.
- **Scanning**: Security scans can be resource-intensive. Run them during low-traffic periods.

## Testing

```go
// Test security middleware
func TestSecurityMiddleware(t *testing.T) {
    handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        w.Write([]byte("OK"))
    })
    
    secureHandler := security.BuildSecureMiddleware(security.DefaultSecurityConfig())(handler)
    
    // Test with valid request
    req := httptest.NewRequest("GET", "/test", nil)
    w := httptest.NewRecorder()
    secureHandler.ServeHTTP(w, req)
    
    if w.Code != http.StatusOK {
        t.Errorf("Expected status 200, got %d", w.Code)
    }
}
```

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Run tests: `go test ./core/security/...`
5. Submit a pull request

## License

This package is licensed under the Apache License 2.0. See the [LICENSE](../LICENSE) file for details.

## Security

If you discover a security vulnerability in this package, please report it following the security policy in [SECURITY.md](../../SECURITY.md).

---

**Note**: This package is designed for production use but should be thoroughly tested in your environment before deployment. Always follow the principle of defense in depth and implement security at multiple levels (network, application, database).
