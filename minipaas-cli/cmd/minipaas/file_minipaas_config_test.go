package main

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadAndSaveConfig(t *testing.T) {
	dir := t.TempDir()
	// Write a config and read it back
	cfg := Config{
		Project: ProjectConfig{Files: []string{"a.yml", "b.yml"}},
		Api:     ApiConfig{Local: true},
		Deploy:  DeployConfig{Version: "2.3.4"},
	}
	fn, err := saveConfig(dir, cfg)
	if err != nil {
		t.Fatalf("saveConfig: %v", err)
	}
	if filepath.Base(fn) != "minipaas.yaml" {
		t.Fatalf("unexpected filename: %s", fn)
	}
	read, gotFn, err := loadConfig(dir)
	if err != nil {
		t.Fatalf("loadConfig: %v", err)
	}
	if gotFn != fn {
		t.Fatalf("filename mismatch: %s vs %s", gotFn, fn)
	}
	if read.Deploy.Version != "2.3.4" || !read.Api.Local || len(read.Project.Files) != 2 {
		t.Fatalf("config roundtrip mismatch: %#v", read)
	}
}

func TestSetApiEnvVars_LocalAndTLS(t *testing.T) {
	os.Unsetenv("DOCKER_CERT_PATH")
	os.Unsetenv("DOCKER_HOST")
	os.Unsetenv("DOCKER_TLS_VERIFY")
	os.Unsetenv("MINIPAAS_DEPLOY_VERSION")

	cfg := Config{Deploy: DeployConfig{Version: "1.2.3"}, Api: ApiConfig{Local: true}}
	setApiEnvVars("/tmp/env", cfg, false)
	if v := os.Getenv("MINIPAAS_DEPLOY_VERSION"); v != "1.2.3" {
		t.Fatalf("MINIPAAS_DEPLOY_VERSION=%q", v)
	}
	if os.Getenv("DOCKER_TLS_VERIFY") != "" || os.Getenv("DOCKER_HOST") != "" || os.Getenv("DOCKER_CERT_PATH") != "" {
		t.Fatalf("docker envs should be empty in local mode")
	}

	os.Unsetenv("DOCKER_CERT_PATH")
	os.Unsetenv("DOCKER_HOST")
	os.Unsetenv("DOCKER_TLS_VERIFY")
	os.Unsetenv("MINIPAAS_DEPLOY_VERSION")

	cfg = Config{Deploy: DeployConfig{Version: "9.9.9"}, Api: ApiConfig{Local: false, Host: "tcp://example:2376", Certs: ".tls"}}
	envDir := t.TempDir()
	setApiEnvVars(envDir, cfg, false)
	if os.Getenv("MINIPAAS_DEPLOY_VERSION") != "9.9.9" {
		t.Fatalf("version not set")
	}
	if os.Getenv("DOCKER_TLS_VERIFY") != "1" {
		t.Fatalf("tls verify not set")
	}
	if os.Getenv("DOCKER_HOST") != "tcp://example:2376" {
		t.Fatalf("host not set")
	}
	wantCert := filepath.Join(envDir, ".tls")
	if os.Getenv("DOCKER_CERT_PATH") != wantCert {
		t.Fatalf("cert path mismatch: %q", os.Getenv("DOCKER_CERT_PATH"))
	}
}
