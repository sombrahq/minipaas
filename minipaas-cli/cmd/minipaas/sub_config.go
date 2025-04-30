package main

import (
	"errors"
	"log"
)

type ConfigSubcommand struct {
	ConfigCreate *ConfigCreateArgs `arg:"subcommand:create"`
	ConfigDelete *ConfigDeleteArgs `arg:"subcommand:delete"`
	ConfigPrune  *ConfigPruneArgs  `arg:"subcommand:prune"`
}

func (args *ConfigSubcommand) Run() {
	switch {
	case args.ConfigCreate != nil:
		args.ConfigCreate.Run()
	case args.ConfigDelete != nil:
		args.ConfigDelete.Run()
	case args.ConfigPrune != nil:
		args.ConfigPrune.Run()

	default:
		log.Fatal(errors.New("command not supported"))
	}

}
