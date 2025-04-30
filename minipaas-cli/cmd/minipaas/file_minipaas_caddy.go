package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

const caddyFile = "caddy.minipaas.json"

func caddyLoadServers(env string) (string, []byte, error) {
	serverFile := filepath.Join(env, caddyFile)
	payloadBytes, err := os.ReadFile(serverFile)
	return serverFile, payloadBytes, err
}

func caddyCheckMatchersAreEqual(a, b interface{}) bool {
	aj, err := json.Marshal(a)
	if err != nil {
		return false
	}
	bj, err := json.Marshal(b)
	if err != nil {
		return false
	}
	return string(aj) == string(bj)
}

func caddyStoreRouteToServer(env string, newRoute map[string]interface{}) (string, error) {
	fn := filepath.Join(env, caddyFile)
	data, err := os.ReadFile(fn)
	if err != nil {
		return fn, err
	}

	var server map[string]interface{}
	if err = json.Unmarshal(data, &server); err != nil {
		return fn, err
	}

	var routes []interface{}
	if r, exists := server["routes"]; exists {
		arr, ok := r.([]interface{})
		if !ok {
			return fn, fmt.Errorf("the 'routes' key exists but is not an array")
		}
		routes = arr
	} else {
		routes = []interface{}{}
	}

	newMatcher, newHasMatcher := newRoute["match"]
	replaced := false
	if newHasMatcher {
		for i, existing := range routes {
			existingRoute, ok := existing.(map[string]interface{})
			if !ok {
				continue
			}
			if existingMatcher, ok := existingRoute["match"]; ok {
				if caddyCheckMatchersAreEqual(existingMatcher, newMatcher) {
					routes[i] = newRoute
					replaced = true
					break
				}
			}
		}
	}

	if !replaced {
		routes = append(routes, newRoute)
	}

	server["routes"] = routes
	updated, err := json.MarshalIndent(server, "", "  ")
	if err != nil {
		return fn, err
	}

	if err = os.WriteFile(fn, updated, 0644); err != nil {
		return fn, err
	}

	return fn, nil
}

func caddyCreateMiniPaasServerFile(env string) (string, error) {
	fn := filepath.Join(env, caddyFile)
	server := map[string]interface{}{
		"listen": []string{":80", ":443"},
		"routes": []map[string]interface{}{},
	}
	updated, err := json.MarshalIndent(server, "", "  ")
	if err != nil {
		return fn, err
	}

	return fn, os.WriteFile(fn, updated, 0644)
}

func caddyCreateRouteGeneric(target, domain, path string) (map[string]interface{}, error) {
	if path == "" {
		path = "/*"
	}
	matcherSet := map[string]interface{}{
		"host": []string{domain},
		"path": []string{path},
	}

	handlerData := map[string]interface{}{
		"handler": "reverse_proxy",
		"upstreams": []map[string]interface{}{
			{
				"dial": fmt.Sprintf("minipaas_%s", target),
			},
		},
	}

	route := map[string]interface{}{
		"match":  []interface{}{matcherSet},
		"handle": []interface{}{handlerData},
	}

	return route, nil
}
