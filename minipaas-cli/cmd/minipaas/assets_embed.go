package main

import (
	"embed"
	"fmt"
)

//go:embed embed/minipaas-swarm/*
var embeddedSwarmFS embed.FS

func readEmbeddedSwarm(name string) ([]byte, error) {
	path := "embed/minipaas-swarm/" + name
	b, err := embeddedSwarmFS.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("embedded file %s not found: %w", name, err)
	}
	return b, nil
}
