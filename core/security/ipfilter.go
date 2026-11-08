package security

import (
	"context"
	"net"
	"net/http"
	"strings"
)

// IPFilter returns a middleware that filters requests based on IP address
func IPFilter(config IPFilterConfig) func(http.Handler) http.Handler {
	if !config.Enabled {
		return func(next http.Handler) http.Handler {
			return next
		}
	}

	// Parse IP ranges
	whitelist := parseIPRanges(config.Whitelist)
	blacklist := parseIPRanges(config.Blacklist)

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			clientIP := getClientIP(r, config.TrustProxy)

			// Check if IP is in blacklist
			if isIPInRanges(clientIP, blacklist) {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusForbidden)
				w.Write([]byte(`{"error":"IP address not allowed"}`))
				return
			}

			// If whitelist is configured, check if IP is in whitelist
			if len(whitelist) > 0 {
				if !isIPInRanges(clientIP, whitelist) {
					w.Header().Set("Content-Type", "application/json")
					w.WriteHeader(http.StatusForbidden)
					w.Write([]byte(`{"error":"IP address not allowed"}`))
					return
				}
			}

			// Add IP to request context for logging/auditing
			ctx := context.WithValue(r.Context(), ipContextKey, clientIP)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// ipRange represents an IP address or CIDR range
type ipRange struct {
	ip     net.IP
	mask   net.IPMask
	isCIDR bool
	isIPv6 bool
}

// parseIPRanges parses a list of IP addresses and CIDR ranges
func parseIPRanges(ips []string) []ipRange {
	var ranges []ipRange
	for _, ipStr := range ips {
		ipStr = strings.TrimSpace(ipStr)
		if ipStr == "" {
			continue
		}

		// Try to parse as CIDR
		_, ipNet, err := net.ParseCIDR(ipStr)
		if err == nil {
			ranges = append(ranges, ipRange{
				ip:     ipNet.IP,
				mask:   ipNet.Mask,
				isCIDR: true,
				isIPv6: ipNet.IP.To4() == nil,
			})
			continue
		}

		// Try to parse as IP
		ip := net.ParseIP(ipStr)
		if ip != nil {
			var mask net.IPMask
			if ip.To4() != nil {
				mask = net.CIDRMask(32, 32)
			} else {
				mask = net.CIDRMask(128, 128)
			}
			ranges = append(ranges, ipRange{
				ip:     ip,
				mask:   mask,
				isCIDR: false,
				isIPv6: ip.To4() == nil,
			})
		}
	}
	return ranges
}

// isIPInRanges checks if an IP address is in any of the ranges
func isIPInRanges(ipStr string, ranges []ipRange) bool {
	ip := net.ParseIP(ipStr)
	if ip == nil {
		return false
	}

	for _, r := range ranges {
		if r.isIPv6 != (ip.To4() == nil) {
			continue // IPv4 vs IPv6 mismatch
		}

		if r.isCIDR {
			if r.mask == nil {
				if ip.Equal(r.ip) {
					return true
				}
			} else {
				ipNet := net.IPNet{IP: r.ip, Mask: r.mask}
				if ipNet.Contains(ip) {
					return true
				}
			}
		} else {
			if ip.Equal(r.ip) {
				return true
			}
		}
	}
	return false
}

// getClientIP extracts the client IP address from the request
func getClientIP(r *http.Request, trustProxy bool) string {
	if trustProxy {
		// Check X-Forwarded-For
		if forwarded := r.Header.Get("X-Forwarded-For"); forwarded != "" {
			ips := strings.Split(forwarded, ",")
			if len(ips) > 0 {
				return strings.TrimSpace(ips[0])
			}
		}
		// Check X-Real-IP
		if realIP := r.Header.Get("X-Real-IP"); realIP != "" {
			return strings.TrimSpace(realIP)
		}
	}

	// Fall back to RemoteAddr
	ip, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return r.RemoteAddr
	}
	return ip
}

// contextKey is a custom type for context keys
type contextKey string

// ipContextKey is the key for storing IP in context
const ipContextKey contextKey = "ip"

// WithIP adds the client IP to the request context
func WithIP(ctx context.Context, ip string) context.Context {
	return context.WithValue(ctx, ipContextKey, ip)
}

// GetIP retrieves the client IP from the request context
func GetIP(ctx context.Context) (string, bool) {
	ip, ok := ctx.Value(ipContextKey).(string)
	return ip, ok
}

// IPInRange checks if an IP address is in a specific CIDR range
func IPInRange(ipStr, cidr string) bool {
	_, ipNet, err := net.ParseCIDR(cidr)
	if err != nil {
		return false
	}
	ip := net.ParseIP(ipStr)
	if ip == nil {
		return false
	}
	return ipNet.Contains(ip)
}

// IsPrivateIP checks if an IP address is in a private range
func IsPrivateIP(ipStr string) bool {
	ip := net.ParseIP(ipStr)
	if ip == nil {
		return false
	}

	// IPv4 private ranges
	if ip4 := ip.To4(); ip4 != nil {
		// 10.0.0.0/8
		if ip4[0] == 10 {
			return true
		}
		// 172.16.0.0/12
		if ip4[0] == 172 && (ip4[1]&0xF0) == 16 {
			return true
		}
		// 192.168.0.0/16
		if ip4[0] == 192 && ip4[1] == 168 {
			return true
		}
		// 169.254.0.0/16 (link-local)
		if ip4[0] == 169 && ip4[1] == 254 {
			return true
		}
		// 127.0.0.0/8 (localhost)
		if ip4[0] == 127 {
			return true
		}
	}

	// IPv6 private ranges
	if ip.To4() == nil {
		// ::1/128 (localhost)
		if ip.IsLoopback() {
			return true
		}
		// fc00::/7 (unique local)
		if ip[0] == 0xfc || ip[0] == 0xfd {
			return true
		}
		// fe80::/10 (link-local)
		if (ip[0] == 0xfe && (ip[1]&0xC0) == 0x80) {
			return true
		}
	}

	return false
}

// IsPublicIP checks if an IP address is public
func IsPublicIP(ipStr string) bool {
	return !IsPrivateIP(ipStr)
}

// SimpleIPFilter returns a simple IP filter middleware
func SimpleIPFilter(whitelist, blacklist []string) func(http.Handler) http.Handler {
	config := DefaultIPFilterConfig()
	config.Enabled = true
	config.Whitelist = whitelist
	config.Blacklist = blacklist
	return IPFilter(config)
}

// WhitelistOnly returns a middleware that only allows specified IPs
func WhitelistOnly(ips []string) func(http.Handler) http.Handler {
	config := IPFilterConfig{
		Enabled:       true,
		Whitelist:     ips,
		Blacklist:     []string{},
		TrustProxy:    true,
		DefaultAction: "deny",
	}
	return IPFilter(config)
}

// BlacklistOnly returns a middleware that blocks specified IPs
func BlacklistOnly(ips []string) func(http.Handler) http.Handler {
	config := IPFilterConfig{
		Enabled:       true,
		Whitelist:     []string{},
		Blacklist:     ips,
		TrustProxy:    true,
		DefaultAction: "allow",
	}
	return IPFilter(config)
}
