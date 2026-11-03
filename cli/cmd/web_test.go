package cmd

import (
	"net/url"
	"testing"
)

func TestBrowserAddr(t *testing.T) {
	cases := map[string]string{
		"127.0.0.1:3005": "127.0.0.1:3005",
		"0.0.0.0:3005":   "127.0.0.1:3005",
		"[::]:3005":      "127.0.0.1:3005",
		"[::1]:3005":     "[::1]:3005",
	}

	for input, want := range cases {
		if got := browserAddr(input); got != want {
			t.Fatalf("browserAddr(%q) = %q, want %q", input, got, want)
		}
	}
}

func TestBootstrapURLCarriesGatewayAndToken(t *testing.T) {
	raw := bootstrapURL("127.0.0.1:3005", "http://localhost:9832", "tok-123")

	u, err := url.Parse(raw)
	if err != nil {
		t.Fatalf("failed to parse bootstrap URL: %v", err)
	}
	if u.Scheme != "http" || u.Host != "127.0.0.1:3005" {
		t.Fatalf("unexpected base URL: %s", raw)
	}

	query := u.Query()
	if got := query.Get("gateway"); got != "http://localhost:9832" {
		t.Fatalf("gateway = %q", got)
	}
	if got := query.Get("token"); got != "tok-123" {
		t.Fatalf("token = %q", got)
	}
}

func TestBootstrapURLWithoutToken(t *testing.T) {
	raw := bootstrapURL("0.0.0.0:3005", "http://localhost:9832", "")

	u, err := url.Parse(raw)
	if err != nil {
		t.Fatalf("failed to parse bootstrap URL: %v", err)
	}
	if u.Host != "127.0.0.1:3005" {
		t.Fatalf("expected wildcard host to become loopback, got %q", u.Host)
	}
	if _, ok := u.Query()["token"]; ok {
		t.Fatalf("expected no token parameter")
	}
}

func TestDefaultGatewayURL(t *testing.T) {
	t.Setenv("NAGARE_GATEWAY_URL", "")
	t.Setenv("NAGARE_GATEWAY_PORT", "")

	if got := defaultGatewayURL(); got != "http://localhost:9832" {
		t.Fatalf("defaultGatewayURL() = %q", got)
	}

	t.Setenv("NAGARE_GATEWAY_PORT", "9444")
	if got := defaultGatewayURL(); got != "http://localhost:9444" {
		t.Fatalf("defaultGatewayURL() = %q", got)
	}

	t.Setenv("NAGARE_GATEWAY_URL", "https://gateway.example.com")
	if got := defaultGatewayURL(); got != "https://gateway.example.com" {
		t.Fatalf("defaultGatewayURL() = %q", got)
	}
}
