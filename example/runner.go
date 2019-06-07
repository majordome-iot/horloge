package main

import (
	"fmt"

	"github.com/majordome-iot/horloge"
)

func main() {
	runner := horloge.NewRunner()

	runner.Sync(&horloge.SyncRedis{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})

	pattern := horloge.NewPattern("daily").At(9, 30, 0)
	job := horloge.NewJob("wake up", *pattern, []string{"super", "cool"})
	runner.AddJob(job)
	runner.AddHandler("wake up", func(arguments ...interface{}) {
		name, args, time := horloge.JobArgs(arguments)
		fmt.Printf("[INFO] Running \"%s\" with args %+v at %s\n", name, args, time.String())
	})

	select {}
}
