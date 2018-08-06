package main

import (
	"fmt"
	"time"

	"github.com/majordome/horloge"
)

func main() {
	runner := horloge.NewRunner()
	pattern := horloge.NewPattern("daily").At(9, 30, 0)
	job := horloge.NewJob("wake up", *pattern)

	runner.AddJob(job)

	runner.AddHandler("wake up", func(name string, args []string, t time.Time) {
		fmt.Printf("[INFO] Running \"%s\" with args %+v at %s\n", name, args, t.String())
	})

	select {}
}
