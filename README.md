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
	"time"

	"github.com/majordome/horloge"
)

func main() {
	runner := horloge.NewRunner()
	pattern := horloge.NewPattern("daily").At(9, 30, 0)
	job := horloge.NewJob("wake up", pattern)

	runner.AddJob(job)

	runner.AddHandler("wake up", func(name string, args []string, t time.Time) {
		fmt.Printf("[INFO] Running \"%s\" with args %+v at %s\n", name, args, t.String())
	})

	select {}
}

```

### What it is

Horloge is a cron like task runner

### What it isn't

Horloge is not an asyncronous task runner like Celery or Faktory.