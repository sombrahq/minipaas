package main

import "testing"

func TestReadEmbeddedSwarm_Success(t *testing.T) {
	files := []string{"common.yml", "registry.yml", "caddy.yml", "caddy.json", "postgres.yaml"}
	for _, f := range files {
		if _, err := readEmbeddedSwarm(f); err != nil {
			t.Fatalf("readEmbeddedSwarm(%s) error: %v", f, err)
		}
	}
}

func TestReadEmbeddedSwarm_NotFound(t *testing.T) {
	if _, err := readEmbeddedSwarm("nope.yaml"); err == nil {
		t.Fatalf("expected error for missing embedded file")
	}
}
