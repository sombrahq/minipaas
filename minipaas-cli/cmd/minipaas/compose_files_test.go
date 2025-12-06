package main

import (
	"os"
	"path/filepath"
	"testing"
)

func TestComposeLoadProject(t *testing.T) {
	dir := t.TempDir()
	f := filepath.Join(dir, "compose.yml")
	content := "version: '3.9'\nservices:\n  api:\n    image: busybox\n"
	if err := os.WriteFile(f, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}

	p, err := composeLoadProject([]string{f})
	if err != nil {
		t.Fatalf("composeLoadProject err: %v", err)
	}
	if _, ok := p.Services["api"]; !ok {
		t.Fatalf("service missing in project")
	}
}

func TestComposeLoadDeployProject_WithEnv(t *testing.T) {
	dir := t.TempDir()
	f := filepath.Join(dir, "compose.yml")
	// include environment variable that should be provided via OS env
	os.Setenv("MINIPAAS_DEPLOY_VERSION", "1.0.0")
	defer os.Unsetenv("MINIPAAS_DEPLOY_VERSION")
	content := "version: '3.9'\nservices:\n  api:\n    image: repo/app:${MINIPAAS_DEPLOY_VERSION}\n"
	if err := os.WriteFile(f, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}

	p, err := composeLoadDeployProject([]string{f})
	if err != nil {
		t.Fatalf("composeLoadDeployProject err: %v", err)
	}
	s := p.Services["api"]
	if s.Image == "" {
		t.Fatalf("expected image to be set from env, got empty")
	}
}
