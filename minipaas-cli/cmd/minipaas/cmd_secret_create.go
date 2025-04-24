package main

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"path/filepath"
)

type SecretCreateArgs struct {
	BaseArgs
	Name string   `arg:"--name" help:"Name of the Docker secret to create. If not provided, a unique name is generated."`
	File string   `arg:"positional" help:"Path to file to use for secret content. If omitted, reads from STDIN."`
	For  []string `arg:"--for,separate" help:"Containers that use the secret"`
}

func (args *SecretCreateArgs) Run() {
	cfg, configFile, err := loadConfig(args.Env)
	checkErrorPanic(err, fmt.Sprintf("❌ Error loading configuration file: %s", configFile))
	setApiEnvVars(args.Env, cfg, args.Verbose)

	var content []byte
	var baseName string
	if args.File != "" {
		baseName = filepath.Base(args.File)
		content, err = os.ReadFile(args.File)
		checkErrorPanic(err, fmt.Sprintf("❌ Failed reading file: %s", args.File))
	} else {
		if args.Name == "" {
			checkErrorPanic(errors.New("when no file is provided, --name is mandatory"), "❌ Failed to get name")
		}
		content, err = io.ReadAll(os.Stdin)
		checkErrorPanic(err, "❌ Failed reading from STDIN")
		baseName = args.Name
	}

	if len(content) == 0 {
		log.Printf("⚠️ Input is empty, skipping.")
		return
	}

	secretCreateAndStore(args.Env, baseName, content, args.For, args.Verbose)
}

func secretCreateAndStore(env, baseName string, content []byte, service []string, verbose bool) {
	secretName, err := secretCreate(baseName, content, verbose)
	checkErrorPanic(err, fmt.Sprintf("❌ Error creating secret for input: %s", baseName))
	fmt.Printf("✅ Secret created: %s\n", secretName)

	deployProject, composeFile, err := loadProject(env)
	checkErrorPanic(err, fmt.Sprintf("❌ Failed to load deploy file: %s", composeFile))

	err = addComposeSecret(deployProject, secretName, baseName, service)
	checkErrorPanic(err, fmt.Sprintf("❌ Failed to update compose: %s", composeFile))

	composeFile, err = saveProject(env, deployProject)
	checkErrorPanic(err, fmt.Sprintf("❌ Failed to update compose file: %s", composeFile))
	fmt.Printf("✅ Updated compose file with secret: %s\n", composeFile)
}

func secretExists(name string) bool {
	cmd := exec.Command("docker", "secret", "inspect", name)
	return cmd.Run() == nil
}

func secretCreate(baseName string, secret []byte, verbose bool) (string, error) {
	hash := sha256.Sum256(secret)
	hashPrefix := hex.EncodeToString(hash[:])[:8]

	secretName := fmt.Sprintf("%s.%s", baseName, hashPrefix)

	return secretName, createLiteralSecret(secretName, secret, verbose)
}

func createLiteralSecret(secretName string, secret []byte, verbose bool) error {
	if secretExists(secretName) {
		return nil
	}

	return runCommandWithInput([]string{"docker", "secret", "create", secretName, "-"}, secret, verbose)
}
