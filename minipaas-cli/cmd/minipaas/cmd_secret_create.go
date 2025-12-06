package main

import (
	"errors"
	"fmt"
	"io"
	"log"
	"os"
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

	env := args.Env
	service := args.For
	verbose := args.Verbose

	secretName, err := secretCreate(baseName, content, verbose)
	checkErrorPanic(err, fmt.Sprintf("❌ Error creating secret for input: %s", baseName))
	fmt.Printf("✅ Secret created: %s\n", secretName)

	// Load project config to discover all compose files
	cfg, _, err = loadConfig(env)
	checkErrorPanic(err, "❌ Failed to load MiniPaaS configuration")

	// Determine compose files to scan and patch
	orderedFiles := composeFilesForEnv(env, cfg)

	// Group services by owning compose file
	svcPerFile, missing := groupServicesByComposeFile(orderedFiles, service)
	if len(missing) > 0 {
		checkErrorPanic(fmt.Errorf("services not found: %v", missing), "❌ Failed to find services in compose files")
	}

	// Patch each compose file that owns at least one target service
	for file, svcs := range svcPerFile {
		project, _, lerr := loadComposeFile(file)
		checkErrorPanic(lerr, fmt.Sprintf("❌ Failed to load compose file: %s", file))

		lerr = addComposeSecret(project, secretName, baseName, svcs)
		checkErrorPanic(lerr, fmt.Sprintf("❌ Failed to update compose: %s", file))

		_, lerr = saveComposeFile(file, project)
		checkErrorPanic(lerr, fmt.Sprintf("❌ Failed to update compose file: %s", file))
		fmt.Printf("✅ Updated compose file with secret: %s\n", file)
	}
}
