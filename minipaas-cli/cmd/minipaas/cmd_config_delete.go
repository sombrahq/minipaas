package main

import (
	"fmt"
)

type ConfigDeleteArgs struct {
	BaseArgs
	Name string `arg:"positional,required" help:"The name of the Docker config to delete."`
}

func (args *ConfigDeleteArgs) Run() {
	cfg, configFile, err := loadConfig(args.Env)
	checkErrorPanic(err, fmt.Sprintf("❌ Error loading configuration file: %s", configFile))
	setApiEnvVars(args.Env, cfg, args.Verbose)

	err = configInternalDelete(args.Name, args.Verbose)
	checkErrorPanic(err, fmt.Sprintf("❌ Error deleting config: %s", args.Name))
	fmt.Printf("✅ Config deleted: %s\n", args.Name)

	composeFile, err := updateDeployFileRemoveConfig(args.Env, args.Name)
	checkErrorPanic(err, fmt.Sprintf("❌ Failed to update deploy file: %s", composeFile))
	fmt.Printf("✅ Updated deploy file: %s\n", composeFile)
}

func configInternalDelete(name string, verbose bool) error {
	if name == "minipaas__config" {
		return nil
	}
	return runCommand([]string{"docker", "config", "rm", name}, verbose)
}
