package main

import (
	"errors"
	"log"
)

type SecretSubcommand struct {
	SecretCreate *SecretCreateArgs `arg:"subcommand:create"`
	SecretDelete *SecretDeleteArgs `arg:"subcommand:delete"`
	SecretPrune  *SecretPruneArgs  `arg:"subcommand:prune"`
}

func (args *SecretSubcommand) Run() {
	switch {
	case args.SecretCreate != nil:
		args.SecretCreate.Run()
	case args.SecretDelete != nil:
		args.SecretDelete.Run()
	case args.SecretPrune != nil:
		args.SecretPrune.Run()

	default:
		log.Fatal(errors.New("command not supported"))
	}

}
