# Horloge

## Rationale


### Getting started

Launch Horloge

```
$ docker pull majordome/horloge
$ docker run -p 127.0.0.1:3000:3000/tcp majordome/horloge
```

Launch a Horloge task manager

```go
package main

import (
  "fmt"

  "github.com/majordome/horloge"
)


func main() {
	runner := horloge.NewRunner()

	task := horloge.NewTask("foobar", "foo", "bar")
	pattern := horloge.Pattern{Occurence: "daily", Hour: 1, Minute: 05, Second: 0}

	runner.Register(task, pattern)

	go func() {
		task := horloge.NewTask("foobar", "foo", "bar")
		pattern := horloge.Pattern{Occurence: "daily", Hour: 1, Minute: 06, Second: 0}
		runner.Register(task, pattern)
	}()

	runner.AddHandler("foobar", func(name string, args []string, t time.Time) {
		fmt.Printf("Running \"%s\" with args %+v at %s\n", name, args, t.String())
	})

	select {}
}

```

### What it is

Horloge is a cron like task runner

### What it isn't

Horloge is not an asyncronous task runner like Celery or Faktory.