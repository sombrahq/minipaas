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

type ConfigCreateArgs struct {
	BaseArgs
	Name string   `arg:"--name" help:"Name of the Docker config to create. If not provided, a unique name is generated."`
	File string   `arg:"positional" help:"Path to file to use for config content. If omitted, reads from STDIN."`
	For  []string `arg:"--for,separate" help:"Containers that use the config"`
}

func (args *ConfigCreateArgs) Run() {
	cfg, configFile, err := loadConfig(args.Env)
	checkErrorPanic(err, fmt.Sprintf("❌ Failed to load configuration file: %s", configFile))
	setApiEnvVars(args.Env, cfg, args.Verbose)

	var content []byte
	var baseName string
	if args.File != "" {
		baseName = filepath.Base(args.File)
		content, err = os.ReadFile(args.File)
		checkErrorPanic(err, fmt.Sprintf("❌ Failed to read file: %s", args.File))
	} else {
		if args.Name == "" {
			checkErrorPanic(errors.New("when no file is provided, --name is mandatory"), "❌ Failed to get name")
		}
		content, err = io.ReadAll(os.Stdin)
		checkErrorPanic(err, "❌ Failed to read from STDIN")
		baseName = args.Name
	}

	if len(content) == 0 {
		log.Printf("⚠️ Input is empty, skipping.")
		return
	}

	configCreateAndStore(args.Env, baseName, content, args.For, args.Verbose)
}

func configCreateAndStore(env, baseName string, content []byte, service []string, verbose bool) {
	configName, err := createConfig(baseName, content, verbose)
	checkErrorPanic(err, fmt.Sprintf("❌ Failed to create config for input: %s", baseName))
	fmt.Printf("✅ Config created: %s\n", configName)

	deployProject, composeFile, err := loadProject(env)
	checkErrorPanic(err, fmt.Sprintf("❌ Failed to load compose file: %s", composeFile))

	err = addComposeConfig(deployProject, configName, baseName, service)
	checkErrorPanic(err, fmt.Sprintf("❌ Failed to update compose: %s", composeFile))

	composeFile, err = saveProject(env, deployProject)
	checkErrorPanic(err, fmt.Sprintf("❌ Failed to update compose file: %s", composeFile))
	fmt.Printf("✅ Updated deploy file with config: %s\n", configName)
}

func configExists(name string) bool {
	cmd := exec.Command("docker", "config", "inspect", name)
	return cmd.Run() == nil
}

func createConfig(baseName string, config []byte, verbose bool) (string, error) {
	hash := sha256.Sum256(config)
	hashPrefix := hex.EncodeToString(hash[:])[:8]

	configName := fmt.Sprintf("%s.%s", baseName, hashPrefix)

	return configName, createLiteralConfig(configName, config, verbose)
}

func createLiteralConfig(configName string, config []byte, verbose bool) error {
	if configExists(configName) {
		return nil
	}

	return runCommandWithInput([]string{"docker", "config", "create", configName, "-"}, config, verbose)
}

func createLock(configName string, verbose bool) error {
	if configExists(configName) {
		return fmt.Errorf("❌ Failed to acquire %s lock", configName)
	}

	return runCommandWithInput([]string{"docker", "config", "create", configName, "-"}, []byte("lock"), verbose)
}
