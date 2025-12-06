package main

import (
	"encoding/json"
	"net/url"
	"os"
	"path/filepath"
	"testing"
)

func readCaddyConfig(t *testing.T, fn string) CaddyConfig {
	t.Helper()
	data, err := os.ReadFile(fn)
	if err != nil {
		t.Fatalf("read file: %v", err)
	}
	var cfg CaddyConfig
	if err := json.Unmarshal(data, &cfg); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	return cfg
}

func TestCaddyUpdateConfigAddRoute_HTTPS(t *testing.T) {
	dir := t.TempDir()
	fn := filepath.Join(dir, "caddy.json")
	if err := os.WriteFile(fn, []byte(`{}`), 0644); err != nil {
		t.Fatalf("write: %v", err)
	}

	out, err := caddyUpdateConfigAddRoute(dir, "example.com/app", "api:8080")
	if err != nil {
		t.Fatalf("caddyUpdateConfigAddRoute error: %v", err)
	}
	if out != fn {
		t.Fatalf("returned filename mismatch: got %q want %q", out, fn)
	}

	cfg := readCaddyConfig(t, fn)
	srv := cfg.Apps.HTTP.Servers.Minipaas
	found := false
	for _, l := range srv.Listen {
		if l == ":443" {
			found = true
			break
		}
	}
	if !found {
		t.Fatalf("expected listen :443, got %#v", srv.Listen)
	}
	if srv.AutomaticHTTPS.Disable || srv.AutomaticHTTPS.DisableRedirects {
		t.Fatalf("automatic https should be enabled for https")
	}
	if len(srv.Routes) == 0 {
		t.Fatalf("expected at least one route")
	}
	r := srv.Routes[len(srv.Routes)-1]
	if len(r.Match) == 0 || len(r.Match[0].Host) == 0 || r.Match[0].Host[0] != "example.com" {
		t.Fatalf("host match mismatch: %#v", r.Match)
	}
	if len(r.Match[0].Path) == 0 || r.Match[0].Path[0] != "/app/*" {
		t.Fatalf("path match mismatch: %#v", r.Match)
	}
	if len(r.Handle) == 0 || len(r.Handle[0].Upstreams) == 0 || r.Handle[0].Upstreams[0].Dial != "minipaas_api:8080" {
		t.Fatalf("upstream mismatch: %#v", r.Handle)
	}
}

func TestCaddyUpdateConfigAddRoute_HTTP_CustomPort(t *testing.T) {
	dir := t.TempDir()
	fn := filepath.Join(dir, "caddy.json")
	if err := os.WriteFile(fn, []byte(`{}`), 0644); err != nil {
		t.Fatalf("write: %v", err)
	}

	if _, err := caddyUpdateConfigAddRoute(dir, "http://example.com:8081/root", "web:80"); err != nil {
		t.Fatalf("caddyUpdateConfigAddRoute error: %v", err)
	}
	cfg := readCaddyConfig(t, fn)
	srv := cfg.Apps.HTTP.Servers.Minipaas
	found := false
	for _, l := range srv.Listen {
		if l == ":8081" {
			found = true
			break
		}
	}
	if !found {
		t.Fatalf(":8081 missing in listens: %#v", srv.Listen)
	}
	if !srv.AutomaticHTTPS.Disable || !srv.AutomaticHTTPS.DisableRedirects {
		t.Fatalf("automatic https should be disabled for http")
	}
	if len(srv.Routes) == 0 {
		t.Fatalf("expected at least one route")
	}
	r := srv.Routes[len(srv.Routes)-1]
	if r.Match[0].Host[0] != "example.com" {
		t.Fatalf("host mismatch: %#v", r.Match)
	}
	if r.Match[0].Path[0] != "/root/*" {
		t.Fatalf("path mismatch: %#v", r.Match)
	}
}

func TestCaddyUpdateConfigAddRoute_ReplaceAndRootPath(t *testing.T) {
	dir := t.TempDir()
	fn := filepath.Join(dir, "caddy.json")
	if err := os.WriteFile(fn, []byte(`{}`), 0644); err != nil {
		t.Fatal(err)
	}

	if _, err := caddyUpdateConfigAddRoute(dir, "https://example.org/", "web:80"); err != nil {
		t.Fatalf("add route: %v", err)
	}
	var cfg CaddyConfig
	data, _ := os.ReadFile(fn)
	_ = json.Unmarshal(data, &cfg)
	if len(cfg.Apps.HTTP.Servers.Minipaas.Routes) != 1 {
		t.Fatalf("expected 1 route, got %d", len(cfg.Apps.HTTP.Servers.Minipaas.Routes))
	}
	r := cfg.Apps.HTTP.Servers.Minipaas.Routes[0]
	if r.Match[0].Path[0] != "/*" {
		t.Fatalf("root path normalization failed: %#v", r.Match[0].Path)
	}

	if _, err := caddyUpdateConfigAddRoute(dir, "https://example.org/", "web:81"); err != nil {
		t.Fatalf("replace route: %v", err)
	}
	data, _ = os.ReadFile(fn)
	_ = json.Unmarshal(data, &cfg)
	if len(cfg.Apps.HTTP.Servers.Minipaas.Routes) != 1 {
		t.Fatalf("route should be replaced, not duplicated: %d", len(cfg.Apps.HTTP.Servers.Minipaas.Routes))
	}
	r = cfg.Apps.HTTP.Servers.Minipaas.Routes[0]
	if r.Handle[0].Upstreams[0].Dial != "minipaas_web:81" {
		t.Fatalf("upstream not updated: %#v", r.Handle[0].Upstreams)
	}
}

func TestCaddyUpdateConfigAddRoute_AddSecondRouteDifferentPath(t *testing.T) {
	dir := t.TempDir()
	fn := filepath.Join(dir, "caddy.json")
	if err := os.WriteFile(fn, []byte(`{}`), 0644); err != nil {
		t.Fatal(err)
	}

	if _, err := caddyUpdateConfigAddRoute(dir, "https://example.org/a", "web:80"); err != nil {
		t.Fatalf("first route: %v", err)
	}
	if _, err := caddyUpdateConfigAddRoute(dir, "https://example.org/b", "api:80"); err != nil {
		t.Fatalf("second route: %v", err)
	}
	cfg := readCaddyConfig(t, fn)
	rts := cfg.Apps.HTTP.Servers.Minipaas.Routes
	if len(rts) != 2 {
		t.Fatalf("expected 2 routes, got %d", len(rts))
	}
	// both should listen on :443 only once
	ls := cfg.Apps.HTTP.Servers.Minipaas.Listen
	count := 0
	for _, l := range ls {
		if l == ":443" {
			count++
		}
	}
	if count != 1 {
		t.Fatalf(":443 should appear once, got %d", count)
	}
}

func TestCaddyLoadServers(t *testing.T) {
	dir := t.TempDir()
	fn := filepath.Join(dir, "caddy.json")
	payload := []byte(`{"apps": {}}`)
	if err := os.WriteFile(fn, payload, 0644); err != nil {
		t.Fatal(err)
	}
	gotFile, gotPayload, err := caddyLoadServers(dir)
	if err != nil {
		t.Fatalf("caddyLoadServers error: %v", err)
	}
	if gotFile != fn {
		t.Fatalf("file mismatch: %q vs %q", gotFile, fn)
	}
	if string(gotPayload) != string(payload) {
		t.Fatalf("payload mismatch: %q vs %q", string(gotPayload), string(payload))
	}
}

func TestParsePublicURL(t *testing.T) {
	tests := []struct {
		in   string
		host string
		port string
		path string
	}{
		{"example.com", "example.com", "443", ""},
		{"http://example.com", "example.com", "80", ""},
		{"https://example.com", "example.com", "443", ""},
		{"https://example.com/app", "example.com", "443", "/app"},
		{"http://example.com:8080/app", "example.com", "8080", "/app"},
	}
	for _, tt := range tests {
		u, err := parsePublicURL(tt.in)
		if err != nil {
			t.Fatalf("parsePublicURL(%q) error: %v", tt.in, err)
		}
		if u.Hostname() != tt.host {
			t.Fatalf("host mismatch for %q: got %q want %q", tt.in, u.Hostname(), tt.host)
		}
		if u.Port() != tt.port {
			t.Fatalf("port mismatch for %q: got %q want %q", tt.in, u.Port(), tt.port)
		}
		if u.Path != tt.path {
			t.Fatalf("path mismatch for %q: got %q want %q", tt.in, u.Path, tt.path)
		}
		if _, err := url.Parse(u.String()); err != nil {
			t.Fatalf("stdlib parse failed for %q: %v", tt.in, err)
		}
	}
}

func TestNormalizeCaddyPath(t *testing.T) {
	cases := map[string]string{
		"":       "/*",
		"/":      "/*",
		"api":    "/api/*",
		"/api":   "/api/*",
		"/api/":  "/api/*",
		"/api/*": "/api/*",
	}
	for in, want := range cases {
		if got := normalizeCaddyPath(in); got != want {
			t.Fatalf("normalizeCaddyPath(%q)=%q want %q", in, got, want)
		}
	}
}

func TestSplitTarget(t *testing.T) {
	s, p := splitTarget("web")
	if s != "web" || p != "80" {
		t.Fatalf("default port failed: %s %s", s, p)
	}
	s, p = splitTarget("web:8080")
	if s != "web" || p != "8080" {
		t.Fatalf("explicit port failed: %s %s", s, p)
	}
}

func TestEnsureServerAndListen(t *testing.T) {
	root := map[string]interface{}{}
	// https adds :443 and leaves automatic_https enabled
	srv := ensureServerAndListen(root, "https", "443")
	listens, _ := srv["listen"].([]interface{})
	found := false
	for _, l := range listens {
		if l == ":443" {
			found = true
		}
	}
	if !found {
		t.Fatalf(":443 not present: %#v", listens)
	}
	if _, ok := srv["automatic_https"].(map[string]interface{}); ok {
		// should not set disable fields for https
		aut := srv["automatic_https"].(map[string]interface{})
		if aut["disable"] == true || aut["disable_redirects"] == true {
			t.Fatalf("automatic_https should not be disabled for https")
		}
	}
	// calling again shouldn't duplicate listens
	_ = ensureServerAndListen(root, "https", "443")
	listens, _ = srv["listen"].([]interface{})
	count := 0
	for _, l := range listens {
		if l == ":443" {
			count++
		}
	}
	if count != 1 {
		t.Fatalf("listen duplicated: %#v", listens)
	}

	// http enables disable flags and adds :8081
	srv = ensureServerAndListen(root, "http", "8081")
	aut, _ := srv["automatic_https"].(map[string]interface{})
	if aut["disable"] != true || aut["disable_redirects"] != true {
		t.Fatalf("automatic_https flags missing for http: %#v", aut)
	}
}

func TestBuildRouteAndRouteMatch(t *testing.T) {
	r := buildRoute("ex.com", "/a/*", "minipaas_web:80")
	h, p := routeMatch(r)
	if h != "ex.com" || p != "/a/*" {
		t.Fatalf("match mismatch: %s %s", h, p)
	}
	handles, _ := r["handle"].([]interface{})
	if len(handles) == 0 {
		t.Fatalf("handle missing")
	}
}

func TestReplaceOrAppendRoute(t *testing.T) {
	srv := map[string]interface{}{}
	r1 := buildRoute("ex.com", "/a/*", "minipaas_web:80")
	replaceOrAppendRoute(srv, r1)
	routes, _ := srv["routes"].([]interface{})
	if len(routes) != 1 {
		t.Fatalf("expected 1 route, got %d", len(routes))
	}

	// replace same host/path
	r1b := buildRoute("ex.com", "/a/*", "minipaas_web:81")
	replaceOrAppendRoute(srv, r1b)
	routes, _ = srv["routes"].([]interface{})
	if len(routes) != 1 {
		t.Fatalf("should still be 1 route after replace")
	}
	// append different path
	r2 := buildRoute("ex.com", "/b/*", "minipaas_api:80")
	replaceOrAppendRoute(srv, r2)
	routes, _ = srv["routes"].([]interface{})
	if len(routes) != 2 {
		t.Fatalf("expected 2 routes, got %d", len(routes))
	}
}
