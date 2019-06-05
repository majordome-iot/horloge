package main

import (
	"fmt"

	"github.com/majordome-iot/horloge"
)

func main() {
	runner := horloge.NewRunner()
	pattern := horloge.NewPattern("daily").At(2, 17, 30)
	job := horloge.NewJob("wake up", *pattern, []string{"super", "coo"})

	runner.AddJob(job)
	runner.AddHandler("wake up", func(arguments ...interface{}) {
		name, args, time := horloge.JobArgs(arguments)
		fmt.Printf("[INFO] Running \"%s\" with args %+v at %s\n", name, args, time.String())
	})

	select {}
}
