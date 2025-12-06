package main

// Strict, struct-only representation of the subset of Caddy JSON we use.
// Unknown fields are intentionally omitted and can be added later.

type CaddyConfig struct {
	Apps CaddyApps `json:"apps,omitempty"`
}

type CaddyApps struct {
	HTTP HTTPApp `json:"http,omitempty"`
}

// HTTPApp models the http/https app with servers.
type HTTPApp struct {
	Servers HTTPServers `json:"servers,omitempty"`
}

// HTTPServers models the set of servers we manage.
// Server name is fixed to "minipaas" in our workflow.
type HTTPServers struct {
	Minipaas Server `json:"minipaas,omitempty"`
}

// Server represents an HTTP server: listen addresses and routes.
type Server struct {
	Listen         []string       `json:"listen,omitempty"`
	Routes         []Route        `json:"routes,omitempty"`
	AutomaticHTTPS AutomaticHTTPS `json:"automatic_https,omitempty"`
}

// Route represents a route with matchers and handlers.
type Route struct {
	Match    []Match   `json:"match,omitempty"`
	Handle   []Handler `json:"handle,omitempty"`
	Terminal bool      `json:"terminal,omitempty"`
}

// Match matches host and path.
type Match struct {
	Host []string `json:"host,omitempty"`
	Path []string `json:"path,omitempty"`
}

// Handler supports reverse_proxy and subroute via fields we use.
type Handler struct {
	Type      string     `json:"handler"`
	Upstreams []Upstream `json:"upstreams,omitempty"`
	Routes    []Route    `json:"routes,omitempty"`
}

// Upstream for reverse_proxy
type Upstream struct {
	Dial string `json:"dial,omitempty"`
}

// AutomaticHTTPS controls Caddy's automatic HTTPS behavior per server.
type AutomaticHTTPS struct {
	Disable          bool `json:"disable,omitempty"`
	DisableRedirects bool `json:"disable_redirects,omitempty"`
}
