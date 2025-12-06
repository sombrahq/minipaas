package main

import (
	"errors"
	"fmt"
	"strings"
)

// indirection for testability
var _runCommand = runCommand
var _runCommandOutput = runCommandOutput

func dockerContainerExec(containerID string, args []string, verbose bool) error {
	allArgs := append([]string{"docker", "exec", "-i", containerID}, args...)
	return _runCommand(allArgs, verbose)
}

func dockerContainerExecOutput(containerID string, args []string, verbose bool) (string, error) {
	allArgs := append([]string{"docker", "exec", "-i", containerID}, args...)
	return _runCommandOutput(allArgs, verbose)
}

func getContainerID(serviceName string) (string, error) {
	cmd := []string{
		"docker", "ps",
		"--filter", "label=com.docker.swarm.service.name=" + serviceName,
		"--format", "{{.ID}}",
	}
	output, err := _runCommandOutput(cmd, false)
	if err != nil {
		return "", fmt.Errorf("docker ps command failed: %v", err)
	}
	ids := strings.Fields(output)
	if len(ids) == 0 {
		return "", errors.New("no running container found for service " + serviceName)
	}
	return ids[0], nil
}
