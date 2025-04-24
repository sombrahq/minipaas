package main

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"os/exec"
)

func runCommand(cmd []string, verbose bool) error {
	if verbose {
		fmt.Printf("🔹 Running: %v\n", cmd)
	}
	ctx := context.Background()
	process := exec.CommandContext(ctx, cmd[0], cmd[1:]...)
	if verbose {
		process.Stdout = os.Stdout
		process.Stderr = os.Stderr
	}
	return process.Run()
}

func runCommandWithInput(cmd []string, input []byte, verbose bool) error {
	if verbose {
		fmt.Printf("🔹 Running: %v\n", cmd)
	}
	ctx := context.Background()
	process := exec.CommandContext(ctx, cmd[0], cmd[1:]...)
	process.Stdin = bytes.NewReader(input)
	if verbose {
		process.Stdout = os.Stdout
		process.Stderr = os.Stderr
	}
	return process.Run()
}

func runCommandOutput(cmd []string, verbose bool) (string, error) {
	if verbose {
		fmt.Printf("🔹 Running: %v\n", cmd)
	}
	ctx := context.Background()
	process := exec.CommandContext(ctx, cmd[0], cmd[1:]...)
	output, err := process.CombinedOutput()
	return string(output), err
}
