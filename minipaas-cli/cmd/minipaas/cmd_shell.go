package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

type ShellArgs struct {
	BaseArgs
}

func (args *ShellArgs) Run() {
	configFile := filepath.Join(args.Env, "minipaas.yaml")
	cfg, configFile, err := loadConfig(args.Env)
	checkErrorPanic(err, fmt.Sprintf("❌ Error loading configuration file: %s", configFile))
	setApiEnvVars(args.Env, cfg, args.Verbose)

	shell := os.Getenv("SHELL")
	if shell == "" {
		shell = "/bin/sh"
	}
	fmt.Printf("🔹 Launching shell: %s\n", shell)

	cmd := exec.Command(shell)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Run()
}
