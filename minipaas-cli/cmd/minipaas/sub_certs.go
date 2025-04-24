package main

import (
	"errors"
	"log"
)

type CertsSubcommand struct {
	ServerCerts *CertsServerArgs `arg:"subcommand:server"`
	ClientCerts *CertsClientArgs `arg:"subcommand:client"`
}

func (args *CertsSubcommand) Run() {
	switch {
	case args.ServerCerts != nil:
		args.ServerCerts.Run()
	case args.ClientCerts != nil:
		args.ClientCerts.Run()

	default:
		log.Fatal(errors.New("command not supported"))
	}

}
