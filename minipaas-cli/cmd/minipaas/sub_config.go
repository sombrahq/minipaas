package main

import (
	"errors"
	"log"
)

type ConfigSubcommand struct {
	ConfigCreate *ConfigCreateArgs `arg:"subcommand:create"`
}

func (args *ConfigSubcommand) Run() {
	switch {
	case args.ConfigCreate != nil:
		args.ConfigCreate.Run()

	default:
		log.Fatal(errors.New("command not supported"))
	}

}
