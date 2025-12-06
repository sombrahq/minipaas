package main

import (
	"encoding/json"
	"fmt"
	"net"
	"net/url"
	"os"
	"path/filepath"
	"strings"
)

const caddyFile = "caddy.json"

func caddyLoadServers(env string) (string, []byte, error) {
	serverFile := filepath.Join(env, caddyFile)
	payloadBytes, err := os.ReadFile(serverFile)
	return serverFile, payloadBytes, err
}

func parsePublicURL(input string) (*url.URL, error) {
	if !strings.Contains(input, "://") {
		input = "https://" + input
	}
	parse, err := url.Parse(input)
	if err != nil {
		return nil, err
	}
	scheme := parse.Scheme
	if parse.Port() == "" {
		if scheme == "http" {
			parse.Host = parse.Host + ":80"
		} else {
			parse.Host = parse.Host + ":443"
		}
	}
	return parse, err
}

// normalizeCaddyPath ensures a path is in the format Caddy expects for
// matchers, defaulting to "/*" and appending a trailing wildcard when needed.
func normalizeCaddyPath(p string) string {
	if p == "" || p == "/" {
		return "/*"
	}
	if !strings.HasPrefix(p, "/") {
		p = "/" + p
	}
	if !strings.HasSuffix(p, "*") {
		if strings.HasSuffix(p, "/") {
			p = p + "*"
		} else {
			p = p + "/*"
		}
	}
	return p
}

// splitTarget parses service[:port] and defaults the port to 80.
func splitTarget(target string) (service, port string) {
	service = target
	port = "80"
	if idx := strings.LastIndex(target, ":"); idx != -1 {
		service = target[:idx]
		port = target[idx+1:]
	}
	return
}

// ensureServerAndListen ensures apps.http.servers.minipaas exists in the root
// config and that it listens on the publicPort. It also toggles automatic HTTPS
// depending on scheme.
func ensureServerAndListen(root map[string]interface{}, scheme, publicPort string) map[string]interface{} {
	ensureMap := func(parent map[string]interface{}, key string) map[string]interface{} {
		if parent == nil {
			return nil
		}
		if v, ok := parent[key]; ok {
			if m, ok := v.(map[string]interface{}); ok {
				return m
			}
		}
		m := map[string]interface{}{}
		parent[key] = m
		return m
	}

	apps := ensureMap(root, "apps")
	httpApp := ensureMap(apps, "http")
	servers := ensureMap(httpApp, "servers")
	server, _ := servers["minipaas"].(map[string]interface{})
	if server == nil {
		server = map[string]interface{}{}
		servers["minipaas"] = server
	}

	// Ensure listen includes the desired port
	listenAddr := ":" + publicPort
	var listenSlice []interface{}
	if v, ok := server["listen"].([]interface{}); ok {
		listenSlice = v
	}
	hasListen := false
	for _, l := range listenSlice {
		if s, ok := l.(string); ok && s == listenAddr {
			hasListen = true
			break
		}
	}
	if !hasListen {
		listenSlice = append(listenSlice, listenAddr)
		server["listen"] = listenSlice
	}

	// If scheme is http, disable automatic HTTPS and redirects for this server
	if scheme == "http" {
		aut, _ := server["automatic_https"].(map[string]interface{})
		if aut == nil {
			aut = map[string]interface{}{}
		}
		aut["disable"] = true
		aut["disable_redirects"] = true
		server["automatic_https"] = aut
	}
	return server
}

// buildRoute constructs a Caddy route map with host+path match and a
// reverse_proxy handler to the given upstream dial.
func buildRoute(host, path, upstreamDial string) map[string]interface{} {
	return map[string]interface{}{
		"match": []interface{}{
			map[string]interface{}{
				"host": []interface{}{host},
				"path": []interface{}{path},
			},
		},
		"handle": []interface{}{
			map[string]interface{}{
				"handler": "reverse_proxy",
				"upstreams": []interface{}{
					map[string]interface{}{
						"dial": upstreamDial,
					},
				},
			},
		},
		"terminal": true,
	}
}

// routeMatch extracts the first host/path matcher from a generic route map.
func routeMatch(route map[string]interface{}) (host, path string) {
	v, ok := route["match"].([]interface{})
	if !ok || len(v) == 0 {
		return "", ""
	}
	m, ok := v[0].(map[string]interface{})
	if !ok {
		return "", ""
	}
	if hs, ok := m["host"].([]interface{}); ok && len(hs) > 0 {
		if s, ok := hs[0].(string); ok {
			host = s
		}
	}
	if ps, ok := m["path"].([]interface{}); ok && len(ps) > 0 {
		if s, ok := ps[0].(string); ok {
			path = s
		}
	}
	return
}

// replaceOrAppendRoute replaces an existing route with the same host+path or
// appends a new one if not found.
func replaceOrAppendRoute(server map[string]interface{}, newRoute map[string]interface{}) {
	var routes []interface{}
	if v, ok := server["routes"].([]interface{}); ok {
		routes = v
	}
	replaced := false
	nh, np := routeMatch(newRoute)
	for i := range routes {
		if r, ok := routes[i].(map[string]interface{}); ok {
			rh, rp := routeMatch(r)
			if rh == nh && rp == np {
				routes[i] = newRoute
				replaced = true
				break
			}
		}
	}
	if !replaced {
		routes = append(routes, newRoute)
	}
	server["routes"] = routes
}

// caddyUpdateConfigAddRoute reads env/caddy.json as a full Caddy config (wrapping
// server-only JSON if needed), ensures apps.http.servers.minipaas exists,
// adds or replaces the route matching the given domain+path, and writes back.
func caddyUpdateConfigAddRoute(env, url, target string) (string, error) {
	fn := filepath.Join(env, caddyFile)

	publicURL, err := parsePublicURL(url)
	if err != nil {
		return fn, err
	}

	scheme := publicURL.Scheme
	if scheme == "" {
		scheme = "https"
	}
	publicPort := publicURL.Port()
	domain := publicURL.Hostname()
	normPath := normalizeCaddyPath(publicURL.Path)

	// Parse target service[:port]
	service, svcPort := splitTarget(target)
	upstreamDial := fmt.Sprintf("minipaas_%s:%s", service, svcPort)

	// Load the file; handle both full config and server-only JSON
	data, err := os.ReadFile(fn)
	if err != nil {
		return fn, err
	}

	// Decode into a generic map to preserve unknown fields verbatim
	var root map[string]interface{}
	if err := json.Unmarshal(data, &root); err != nil {
		return fn, err
	}

	server := ensureServerAndListen(root, scheme, publicPort)

	hostOnly := domain
	if h, _, err := net.SplitHostPort(domain); err == nil {
		hostOnly = h
	}
	newRoute := buildRoute(hostOnly, normPath, upstreamDial)

	replaceOrAppendRoute(server, newRoute)

	// Write back preserving unrelated fields
	out, err := json.MarshalIndent(root, "", "  ")
	if err != nil {
		return fn, err
	}
	if err := os.WriteFile(fn, out, 0644); err != nil {
		return fn, err
	}
	return fn, nil
}
