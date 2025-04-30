package main

import (
	"errors"
	"github.com/alexflint/go-arg"
	"log"
)

/***********
COMMANDS
************/

var args struct {
	CertsSubcommand  *CertsSubcommand  `arg:"subcommand:certs"`
	CodeSubcommand   *CodeSubcommand   `arg:"subcommand:code"`
	SecretSubcommand *SecretSubcommand `arg:"subcommand:secret"`
	ConfigSubcommand *ConfigSubcommand `arg:"subcommand:config"`
	DeploySubcommand *DeploySubcommand `arg:"subcommand:deploy"`

	Shell *ShellArgs `arg:"subcommand:shell"`
}

/***********
CONFIG
************/

func main() {
	arg.MustParse(&args)

	switch {
	case args.CertsSubcommand != nil:
		args.CertsSubcommand.Run()

	case args.CodeSubcommand != nil:
		args.CodeSubcommand.Run()

	case args.SecretSubcommand != nil:
		args.SecretSubcommand.Run()

	case args.ConfigSubcommand != nil:
		args.ConfigSubcommand.Run()

	case args.DeploySubcommand != nil:
		args.DeploySubcommand.Run()

	case args.Shell != nil:
		args.Shell.Run()
	default:
		log.Fatal(errors.New("command not supported"))
	}
}
