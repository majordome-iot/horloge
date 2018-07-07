package main

import (
	"fmt"
	"time"

	"github.com/hsh/horloge"
)

func main() {
	runner := horloge.NewRunner()
	pattern := horloge.NewPattern("daily").At(15, 30, 0)
	job := horloge.NewJob("foobar", *pattern)

	runner.AddJob(job)

	go func() {
		time.Sleep(8 * time.Second)
		fmt.Println("Cancel bro")
		runner.RemoveJob(job)
	}()

	runner.AddHandler("foobar", func(name string, args []string, t time.Time) {
		fmt.Printf("[INFO] Running \"%s\" with args %+v at %s\n", name, args, t.String())
	})

	select {}
}
