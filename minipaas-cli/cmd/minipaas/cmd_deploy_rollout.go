package main

import (
	"fmt"
	"path/filepath"
)

type DeployRolloutArgs struct {
	BaseArgs
}

func (args *DeployRolloutArgs) Run() {
	cfg, configFile, err := loadConfig(args.Env)
	checkErrorPanic(err, fmt.Sprintf("❌ Error loading configuration file: %s", configFile))
	setApiEnvVars(args.Env, cfg, args.Verbose)

	composeFiles := cfg.Project.Files
	composeFiles = append(composeFiles, filepath.Join(args.Env, appsFile))

	var files []string
	for _, fn := range composeFiles {
		files = append(files, "-c", fn)
	}

	deployArgs := append([]string{"docker", "stack", "deploy"}, files...)
	deployArgs = append(deployArgs, "minipaas")
	err = runCommand(deployArgs, args.Verbose)
	checkErrorPanic(err, fmt.Sprintf("❌ Error deploying version %s", cfg.Deploy.Version))
	fmt.Printf("✅ Deployment successful: %s\n", cfg.Deploy.Version)
}
