package main

import (
	"fmt"
	"net/url"
	"strings"
)

func parsePublicURL(input string) (*url.URL, error) {
	if !strings.Contains(input, "://") {
		input = "https://" + input
	}
	return url.Parse(input)
}

type CodeExposeArgs struct {
	BaseArgs
	Target string `arg:"positional,required" help:"Which service to expose. It can also contain the port. Default port to 80."`
	URL    string `arg:"positional,required" help:"Public URL that will be used to expose the service."`
}

func (args *CodeExposeArgs) Run() {
	deployProject, composeFile, err := loadProject(args.Env)
	checkErrorPanic(err, fmt.Sprintf("❌ Fail to load build file: %s", composeFile))

	parts, err := parsePublicURL(args.URL)
	checkErrorPanic(err, fmt.Sprintf("❌ Fail to parse URL: %s", args.URL))

	server, err := caddyCreateRouteGeneric(args.Target, parts.Host, parts.Path)
	checkErrorPanic(err, fmt.Sprintf("❌ Fail to create route for: %s", args.Target))

	serverFile, err := caddyStoreRouteToServer(args.Env, server)
	checkErrorPanic(err, fmt.Sprintf("❌ Fail to write file: %s", serverFile))
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
