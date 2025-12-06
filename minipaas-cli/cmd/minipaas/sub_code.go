package main

import (
	"errors"
	"log"
)

type CodeSubcommand struct {
	CodeInit   *CodeInitArgs   `arg:"subcommand:init"`
	CodeRoute  *CodeRouteArgs  `arg:"subcommand:route"`
	CodeJob    *CodeJobArgs    `arg:"subcommand:job"`
	CodeWorker *CodeWorkerArgs `arg:"subcommand:worker"`
	CodeCron   *CodeCronArgs   `arg:"subcommand:cron"`
}

func (args *CodeSubcommand) Run() {
	switch {
	case args.CodeInit != nil:
		args.CodeInit.Run()
	case args.CodeRoute != nil:
		args.CodeRoute.Run()
	case args.CodeJob != nil:
		args.CodeJob.Run()
	case args.CodeWorker != nil:
		args.CodeWorker.Run()
	case args.CodeCron != nil:
		args.CodeCron.Run()

	default:
		log.Fatal(errors.New("command not supported"))
	}

}
