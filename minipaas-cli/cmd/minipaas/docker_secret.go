package main

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"os/exec"
)

// Make external calls overridable for tests
var dockerSecretInspect = func(name string) error {
	return exec.Command("docker", "secret", "inspect", name).Run()
}

var dockerSecretCreate = func(name string, content []byte, verbose bool) error {
	return runCommandWithInput([]string{"docker", "secret", "create", name, "-"}, content, verbose)
}

func secretExists(name string) bool {
	return dockerSecretInspect(name) == nil
}

func secretCreate(baseName string, secret []byte, verbose bool) (string, error) {
	hash := sha256.Sum256(secret)
	hashPrefix := hex.EncodeToString(hash[:])[:8]

	secretName := fmt.Sprintf("%s.%s", baseName, hashPrefix)

	return secretName, secretCreateLiteral(secretName, secret, verbose)
}

func secretCreateLiteral(secretName string, secret []byte, verbose bool) error {
	if secretExists(secretName) {
		return nil
	}
	return dockerSecretCreate(secretName, secret, verbose)
}
