package main

import (
	"fmt"
	"log"
	"strings"
)

type DeployRoutingArgs struct {
	BaseArgs
}

func (args *DeployRoutingArgs) Run() {
	cfg, configFile, err := loadConfig(args.Env)
	checkErrorPanic(err, fmt.Sprintf("❌ Fail to load configuration file: %s", configFile))
	setApiEnvVars(args.Env, cfg, args.Verbose)

	containerID, err := getContainerID(CaddyContainerName)
	checkErrorPanic(err, fmt.Sprintf("❌ Fail to obtain container ID for `%s`", CaddyContainerName))

	getAppsCmd := []string{
		"/usr/bin/wget",
		"-O", "-", "-q",
		"http://127.0.0.1:2019/config/apps",
	}
	appsOutput, err := dockerContainerExecOutput(containerID, getAppsCmd, args.Verbose)
	checkErrorPanic(err, "❌ Fail to get Caddy apps config")
	trimmed := strings.TrimSpace(appsOutput)
	if trimmed == "" || trimmed == "null" || trimmed == "{}" {
		emptyAppJSON := `{"http":{"servers":{}}}`
		postAppsCmd := []string{
			"/usr/bin/wget",
			"-O", "-", "-q",
			"--header=Content-Type: application/json",
			"--post-data=" + emptyAppJSON,
			"http://127.0.0.1:2019/config/apps",
		}
		err = dockerContainerExec(containerID, postAppsCmd, args.Verbose)
		if err != nil {
			log.Printf("⚠️ Warning: Fail to set default apps")
		}
	}
	serverFile, payload, err := caddyLoadServers(args.Env)
	checkErrorPanic(err, fmt.Sprintf("❌ Fail to load server config: %s", serverFile))

	url := "http://127.0.0.1:2019/config/apps/http/servers/minipaas"

	cmdArgs := []string{
		"/usr/bin/wget",
		"-O", "-", "-q",
		"--header=Content-Type: application/json",
		"--post-data=" + string(payload),
		url,
	}

	err = dockerContainerExec(containerID, cmdArgs, args.Verbose)
	checkErrorPanic(err, "❌ Fail to update MiniPaaS server")
	fmt.Printf("✅ Routing successful\n")
}
