package main

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"os/exec"
)

// Make external calls overridable for tests
var dockerConfigInspect = func(name string) error {
	return exec.Command("docker", "config", "inspect", name).Run()
}

var dockerConfigCreate = func(name string, content []byte, verbose bool) error {
	return runCommandWithInput([]string{"docker", "config", "create", name, "-"}, content, verbose)
}

func configExists(name string) bool {
	return dockerConfigInspect(name) == nil
}

func configCreate(baseName string, config []byte, verbose bool) (string, error) {
	hash := sha256.Sum256(config)
	hashPrefix := hex.EncodeToString(hash[:])[:8]

	configName := fmt.Sprintf("%s.%s", baseName, hashPrefix)

	return configName, configCreateLiteral(configName, config, verbose)
}

func configCreateLiteral(configName string, config []byte, verbose bool) error {
	if configExists(configName) {
		return nil
	}
	return dockerConfigCreate(configName, config, verbose)
}
