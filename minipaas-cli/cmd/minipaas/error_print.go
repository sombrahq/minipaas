package main

import "log"

func checkErrorPanic(err error, msg string) {
	if err != nil {
		log.Println(msg)
		log.Fatalf("❌ Error: %v", err)
	}
}
