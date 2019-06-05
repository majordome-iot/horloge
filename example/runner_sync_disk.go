package main

import (
	"fmt"

	"github.com/majordome-iot/horloge"
)

func main() {
	runner := horloge.NewRunner()

	runner.AddHandler("wake up", func(arguments ...interface{}) {
		name, args, time := horloge.JobArgs(arguments)
		fmt.Printf("[INFO] Running \"%s\" with args %+v at %s\n", name, args, time.String())
	})

	runner.Sync(&horloge.SyncDisk{
		Path: "./test.json",
	})

	select {}
}
