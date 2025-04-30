package main

import (
	"fmt"
	"github.com/compose-spec/compose-go/v2/types"
	"path/filepath"
)

type DeployBuildArgs struct {
	BaseArgs
}

func (args *DeployBuildArgs) Run() {
	cfg, configFile, err := loadConfig(args.Env)
	checkErrorPanic(err, fmt.Sprintf("❌ Error loading configuration file: %s", configFile))
	setApiEnvVars(args.Env, cfg, args.Verbose)

	project, err := composeLoadDeployProject(append(cfg.Project.Files, filepath.Join(args.Env, deployFile)))
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

func buildCommandFromService(svc types.ServiceConfig) []string {
	cmd := []string{"docker", "build", "--network", "host"}

	// If an image name is specified, use it to tag the built image.
	if svc.Image != "" {
		cmd = append(cmd, "-t", svc.Image)
	}

	// If the service has build configuration, process it.
	if svc.Build != nil {
		context := svc.Build.Context
		if context == "" {
			context = "."
		}

		// Specify the Dockerfile if provided.
		if svc.Build.Dockerfile != "" {
			cmd = append(cmd, "-f", filepath.Join(context, svc.Build.Dockerfile))
		}

		// Append build arguments if any.f
		if svc.Build.Args != nil {
			for key, valPtr := range svc.Build.Args {
				if valPtr != nil {
					cmd = append(cmd, "--build-arg", fmt.Sprintf("%s=%s", key, *valPtr))
				} else {
					cmd = append(cmd, "--build-arg", fmt.Sprintf("%s=", key))
				}
			}
		}

		// Use the provided build context. If empty, default to current directory.
		cmd = append(cmd, svc.Build.Context)
	} else {
		// If no build configuration is provided, default to current directory.
		cmd = append(cmd, ".")
	}

	return cmd
}
