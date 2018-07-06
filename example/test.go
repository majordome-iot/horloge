package main

import (
	"fmt"
	"time"

	"github.com/hsh/horloge"
)

func main() {
	runner := horloge.NewRunner()

	task := horloge.NewTask("foobar", "foo", "bar")
	pattern := horloge.Pattern{Occurence: "every", Second: 5}
	job := runner.AddTask(task, pattern)

	go func() {
		time.Sleep(8 * time.Second)
		runner.Remove(job)
	}()

	// job := horloge.CreateJob("foo", "@every 1 second")
	// job.Bind([string]{"foo", "bar"})
	// runner.AddJob(job)

	runner.AddHandler("foobar", func(name string, args []string, t time.Time) {
		fmt.Printf("[INFO] Running \"%s\" with args %+v at %s\n", name, args, t.String())
	})

	select {}
}
