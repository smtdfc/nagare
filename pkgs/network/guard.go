package network

import (
	"fmt"
	"net"
	"net/url"
	"strings"
)

// GuardWebAccess blocks the URL when it targets loopback, link-local,
// private, or CGNAT infrastructure. When allowPrivate is engaged explicitly
// by the operator or tool arguments, all of these targets are accepted —
// needed for legitimately private workloads like local dev servers.
//
// Rationale: callers execute next to the operator's real machine, where a
// request to 169.254.169.254 (cloud metadata) or 127.0.0.1 (local services)
// can bypass firewalls and reach credentials the application must never
// touch.
func GuardWebAccess(rawURL string, allowPrivate bool) (*url.URL, error) {
	rawURL = strings.TrimSpace(rawURL)
	if rawURL == "" {
		return nil, fmt.Errorf("URL cannot be empty")
	}

	parsed, err := url.Parse(rawURL)
	if err != nil {
		return nil, fmt.Errorf("invalid URL: %w", err)
	}
	if parsed.Scheme != "http" && parsed.Scheme != "https" {
		return nil, fmt.Errorf("only http and https schemes are allowed, got '%s'", parsed.Scheme)
	}

	host := parsed.Hostname()
	if host == "" {
		return nil, fmt.Errorf("URL is missing a host")
	}
	if allowPrivate {
		return parsed, nil
	}

	ips, err := net.LookupIP(host)
	if err != nil {
		return nil, fmt.Errorf("failed to resolve host '%s': %w", host, err)
	}

	for _, ip := range ips {
		if reason := webAccessViolation(ip); reason != "" {
			return nil, fmt.Errorf("host '%s' resolves to %s which is not allowed (%s); set allow_private=true to override", host, ip, reason)
		}
	}

	return parsed, nil
}

// webAccessViolation returns the reason an IP is denied, or an empty string
// when the target is acceptable.
func webAccessViolation(ip net.IP) string {
	switch {
	case ip.IsLoopback():
		return "loopback address"
	case ip.IsLinkLocalUnicast():
		return "link-local address (cloud metadata endpoints live here)"
	case ip.IsLinkLocalMulticast(), ip.IsMulticast(), ip.IsUnspecified():
		return "reserved address"
	case ip.IsPrivate():
		return "private network address"
	case isCGNAT(ip):
		return "CGNAT address"
	default:
		return ""
	}
}

// isCGNAT matches the 100.64.0.0/10 carrier-grade NAT range.
func isCGNAT(ip net.IP) bool {
	v4 := ip.To4()
	if v4 == nil {
		return false
	}
	return v4[0] == 100 && v4[1] >= 64 && v4[1] < 128
}
