package main

import (
	"fmt"
	"time"

	"github.com/hsh/horloge"
)

func main() {
	runner := horloge.NewRunner()

	pattern := horloge.Pattern{Occurence: "every", Second: 5}
	job := horloge.NewJob("foobar", pattern, "foo", "bar")

	runner.AddJob(job)

	go func() {
		time.Sleep(8 * time.Second)
		fmt.Println("Cancel bro")
		runner.RemoveJob(job)
	}()

	// job := horloge.CreateJob("foo", "@every 1 second")
	// job.Bind([string]{"foo", "bar"})
	// runner.AddJob(job)

	runner.AddHandler("foobar", func(name string, args []string, t time.Time) {
		fmt.Printf("[INFO] Running \"%s\" with args %+v at %s\n", name, args, t.String())
	})

	select {}
}
