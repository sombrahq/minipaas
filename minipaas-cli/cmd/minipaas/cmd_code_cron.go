package main

import (
	"fmt"
)

type CodeCronArgs struct {
	BaseArgs
	Services []string `arg:"positional,required" help:"Services to configure as a job."`
	Cron     string   `arg:"--cron" help:"Cron schedule." default:"* * * * *"`
}

func (args *CodeCronArgs) Run() {
	deployProject, composeFile, err := loadProject(args.Env)
	checkErrorPanic(err, fmt.Sprintf("❌ Failed to load compose file: %s", composeFile))

	for _, service := range args.Services {
		composeEnsureDeploy(deployProject, service)
		err = addComposeCronDeploy(deployProject, service, args.Cron)
		checkErrorPanic(err, fmt.Sprintf("❌ Failed to update compose file: %s", composeFile))
	}

	composeFile, err = saveProject(args.Env, deployProject)
	checkErrorPanic(err, fmt.Sprintf("❌ Failed to write compose file: %s", composeFile))
	fmt.Println("✅ ", composeFile)

}
