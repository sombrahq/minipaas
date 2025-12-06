package main

import (
	"fmt"
	"path/filepath"
)

type DeployBuildArgs struct {
	BaseArgs
}

func (args *DeployBuildArgs) Run() {
	cfg, configFile, err := loadConfig(args.Env)
	checkErrorPanic(err, fmt.Sprintf("❌ Error loading configuration file: %s", configFile))
	setApiEnvVars(args.Env, cfg, args.Verbose)

	project, err := composeLoadDeployProject(append(cfg.Project.Files, filepath.Join(args.Env, appsFile)))
	checkErrorPanic(err, fmt.Sprintf("❌ Fail to load project files: %s", cfg.Project.Files))

	for name, svc := range project.Services {
		if svc.Build == nil {
			continue
		}
		buildArgs := buildCommandFromService(svc)
		err = runCommand(buildArgs, args.Verbose)
		if err != nil {
			fmt.Printf("❌ %s: %s\n", name, err.Error())
		} else {
			fmt.Printf("✅ %s: %s\n", name, cfg.Deploy.Version)
		}
	}
}
