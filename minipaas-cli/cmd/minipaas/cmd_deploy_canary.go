package main

import (
	"fmt"
	"github.com/compose-spec/compose-go/v2/types"
	"path/filepath"
)

type DeployCanaryArgs struct {
	BaseArgs
	Services []string `arg:"positional,required" help:"Services to release with canary."`
	Replicas int      `arg:"--replicas" help:"Replicas to add for canary release." default:"1"`
}

func (args *DeployCanaryArgs) Run() {
	cfg, configFile, err := loadConfig(args.Env)
	checkErrorPanic(err, fmt.Sprintf("❌ Fail to load configuration file: %s", configFile))
	setApiEnvVars(args.Env, cfg, args.Verbose)

	composeFiles := append(cfg.Project.Files, filepath.Join(args.Env, appsFile))
	deployment, err := composeLoadDeployProject(composeFiles)
	checkErrorPanic(err, fmt.Sprintf("❌ Fail to load project files: %s", composeFiles))

	var minipaasName string
	var srv types.ServiceConfig

	for _, service := range args.Services {
		srv, err = deployment.GetService(service)
		if err != nil {
			fmt.Println(fmt.Sprintf("❌ Fail to get service: %s", service))
			continue
		}
		minipaasName = fmt.Sprintf("minipaas_%s", service)

		err = runCommand([]string{"docker", "service", "update", "--replicas", fmt.Sprintf("%d", args.Replicas), "--image", srv.Image, minipaasName}, args.Verbose)
		if err != nil {
			fmt.Println(fmt.Sprintf("❌ Fail to do canary release for service: %s", service))
			continue
		}

		fmt.Printf("✅ %s: %s\n", service, srv.Image)
	}

}
