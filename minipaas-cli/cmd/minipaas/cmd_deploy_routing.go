package main

import (
	"fmt"
)

type DeployRoutingArgs struct {
	BaseArgs
}

func (args *DeployRoutingArgs) Run() {
	cfg, configFile, err := loadConfig(args.Env)
	checkErrorPanic(err, fmt.Sprintf("❌ Fail to load configuration file: %s", configFile))
	setApiEnvVars(args.Env, cfg, args.Verbose)

	serverFile, payload, err := caddyLoadServers(args.Env)
	checkErrorPanic(err, fmt.Sprintf("❌ Fail to load server JSON: %s", serverFile))

	cmdArgs := []string{
		"/usr/bin/wget",
		"-O", "-", "-q",
		"--header=Content-Type: application/json",
		"--post-data=" + string(payload),
		"http://127.0.0.1:2019/load",
	}

	containerID, err := getContainerID(CaddyContainerName)
	checkErrorPanic(err, fmt.Sprintf("❌ Fail to obtain container ID for `%s`", CaddyContainerName))

	err = dockerContainerExec(containerID, cmdArgs, args.Verbose)
	checkErrorPanic(err, "❌ Fail to update server in Caddy")
	fmt.Printf("✅ Routing updated\n")
}
