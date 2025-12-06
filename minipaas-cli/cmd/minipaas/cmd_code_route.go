package main

import (
	"fmt"
	"strings"
)

type CodeRouteArgs struct {
	BaseArgs
	URL    string `arg:"positional,required" help:"Public URL that will be used to expose the service."`
	Target string `arg:"positional,required" help:"Which service to expose. It can also contain the port. Default port to 80."`
}

func (args *CodeRouteArgs) Run() {
	deployProject, composeFile, err := loadProject(args.Env)
	checkErrorPanic(err, fmt.Sprintf("❌ Fail to load build file: %s", composeFile))

	serverFile, err := caddyUpdateConfigAddRoute(args.Env, args.URL, args.Target)
	checkErrorPanic(err, fmt.Sprintf("❌ Fail to update caddy config: %s", serverFile))
	fmt.Println("✅ ", serverFile)

	components := strings.Split(args.Target, ":")
	container := components[0]
	port := "80"
	if len(components) == 2 {
		port = components[1]
	}
	err = addComposeResilientDeploy(deployProject, container, port)
	checkErrorPanic(err, fmt.Sprintf("❌ Fail to update deployment file: %s", composeFile))

	composeFile, err = saveProject(args.Env, deployProject)
	checkErrorPanic(err, fmt.Sprintf("❌ Fail to write file: %s", composeFile))
	fmt.Println("✅ ", composeFile)
}
