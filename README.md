# Horloge

## Rationale


### Getting started

Launch Horloge using docker

```
$ docker pull shinuza/horloge
$ docker run -p 127.0.0.1:6432:6432/tcp shinuza/horloge
```

Using the CLI

```
$ go get github.com/marjordome/horloge
$ horloge --bind 0.0.0.0 --port 1234
```

### CLI Usage

```
NAME:
   horloge - A new cli application

USAGE:
   main.exe [global options] command [command options] [arguments...]

VERSION:
   0.1.0

AUTHOR:
   Samori Gorse <samorigorse@gail.com>

COMMANDS:
     help, h  Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --port value, -p value  Port to listen to (default: 6432)
   --bind value, -b value  Address to bind to (default: "127.0.0.1")
   --help, -h              show help
   --version, -v           print the version
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
	job := horloge.NewJob("wake up", *pattern)

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