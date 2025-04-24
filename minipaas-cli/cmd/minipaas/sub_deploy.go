package main

import (
	"errors"
	"log"
)

type DeploySubcommand struct {
	DeployBuild   *DeployBuildArgs   `arg:"subcommand:build"`
	DeployRollout *DeployRolloutArgs `arg:"subcommand:rollout"`
	DeployCanary  *DeployCanaryArgs  `arg:"subcommand:canary"`
	DeployRouting *DeployRoutingArgs `arg:"subcommand:routing"`
}

func (args *DeploySubcommand) Run() {
	switch {
	case args.DeployRollout != nil:
		args.DeployRollout.Run()
	case args.DeployBuild != nil:
		args.DeployBuild.Run()
	case args.DeployRouting != nil:
		args.DeployRouting.Run()
	case args.DeployCanary != nil:
		args.DeployCanary.Run()

	default:
		log.Fatal(errors.New("command not supported"))
	}

}
