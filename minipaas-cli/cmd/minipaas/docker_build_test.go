package main

import (
	"path/filepath"
	"testing"

	"github.com/compose-spec/compose-go/v2/types"
)

func strptr(s string) *string { return &s }

func TestBuildCommandFromService(t *testing.T) {
	svc := types.ServiceConfig{Image: "repo/app:old"}
	dockerfile := "Dockerfile.prod"
	svc.Build = &types.BuildConfig{
		Context:    "./app",
		Dockerfile: dockerfile,
		Args: map[string]*string{
			"FOO":   strptr("bar"),
			"EMPTY": nil,
		},
	}

	cmd := buildCommandFromService(svc)
	wantCtx := svc.Build.Context
	wantDf := filepath.Join(svc.Build.Context, svc.Build.Dockerfile)

	has := func(x string) bool {
		for _, a := range cmd {
			if a == x {
				return true
			}
		}
		return false
	}
	if !has("docker") || !has("build") {
		t.Fatalf("missing docker build: %#v", cmd)
	}
	if !has("-t") || !has("repo/app:old") {
		t.Fatalf("missing tag: %#v", cmd)
	}
	if !has("-f") || !has(wantDf) {
		t.Fatalf("missing dockerfile: %#v", cmd)
	}
	if !has(wantCtx) {
		t.Fatalf("missing context: %#v", cmd)
	}
	if !has("--build-arg") || !has("FOO=bar") || !has("EMPTY=") {
		t.Fatalf("missing build args: %#v", cmd)
	}
}
