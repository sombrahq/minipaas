package main

import (
	"fmt"
)

type CodeWorkerArgs struct {
	BaseArgs
	Services []string `arg:"positional,required" help:"Services to configure as a worker."`
}

func (args *CodeWorkerArgs) Run() {
	deployProject, composeFile, err := loadProject(args.Env)
	checkErrorPanic(err, fmt.Sprintf("❌ Failed to load build file: %s", composeFile))

	for _, service := range args.Services {
		err = addComposeWorkerDeploy(deployProject, service)
		checkErrorPanic(err, fmt.Sprintf("❌ Failed to update deployment file: %s", composeFile))
	}

	composeFile, err = saveProject(args.Env, deployProject)
	checkErrorPanic(err, fmt.Sprintf("❌ Failed to write file: %s", composeFile))
	fmt.Println("✅ ", composeFile)

}
