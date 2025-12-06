package main

import (
	"errors"
	"log"
)

type SecretSubcommand struct {
	SecretCreate *SecretCreateArgs `arg:"subcommand:create"`
}

func (args *SecretSubcommand) Run() {
	switch {
	case args.SecretCreate != nil:
		args.SecretCreate.Run()

	default:
		log.Fatal(errors.New("command not supported"))
	}

}
