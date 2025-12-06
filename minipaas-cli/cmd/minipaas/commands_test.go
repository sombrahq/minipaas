package main

import (
	"os"
	"path/filepath"
	"testing"
)

func makeScript(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	p := filepath.Join(dir, "script.sh")
	if err := os.WriteFile(p, []byte("#!/bin/sh\n"+content+"\n"), 0755); err != nil {
		t.Fatalf("write script: %v", err)
	}
	return p
}

func TestRunCommand_SuccessAndFailure(t *testing.T) {
	ok := makeScript(t, "exit 0")
	if err := runCommand([]string{ok}, false); err != nil {
		t.Fatalf("expected success, got: %v", err)
	}
	fail := makeScript(t, "exit 1")
	if err := runCommand([]string{fail}, false); err == nil {
		t.Fatalf("expected failure, got nil")
	}
}

func TestRunCommandWithInput(t *testing.T) {
	// script exits 0 only if stdin equals "hello"
	sc := makeScript(t, `
in=$(cat)
if [ "$in" = "hello" ]; then
  exit 0
else
  exit 1
fi`)
	if err := runCommandWithInput([]string{sc}, []byte("hello"), false); err != nil {
		t.Fatalf("expected success, got: %v", err)
	}
	if err := runCommandWithInput([]string{sc}, []byte("nope"), false); err == nil {
		t.Fatalf("expected failure for wrong input")
	}
}

func TestRunCommandOutput(t *testing.T) {
	sc := makeScript(t, "printf 'world'")
	out, err := runCommandOutput([]string{sc}, false)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out != "world" {
		t.Fatalf("output mismatch: %q", out)
	}
}
