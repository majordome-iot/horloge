package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/codegangsta/cli"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/majordome-iot/horloge"
)

func server(addr string, runner *horloge.Runner) {
	e := echo.New()

	e.HideBanner = true
	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Routes
	e.GET("/ping", horloge.HTTPHandlerPing())
	e.GET("/health_check", horloge.HTTPHandlerHealthCheck())
	e.GET("/version", horloge.HTTPHandlerVersion())
	e.POST("/jobs", horloge.HTTPHandlerRegisterJob(runner))
	e.GET("/jobs", horloge.HTTPHandlerListJobs(runner))
	e.GET("/jobs/:id", horloge.HTTPHandlerJobDetail(runner))
	e.DELETE("/jobs/:id", horloge.HTTPHandlerDeleteJob(runner))

	go func() {
		fmt.Printf("🕒 Horloge v%s\n", horloge.Version)
		fmt.Printf("Http server powered by Echo v%s\n", echo.Version)
		fmt.Printf("Websocket server powered by Melody\n")
		e.Logger.Fatal(e.Start(addr))
	}()

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)
	<-signalChan

	fmt.Println("Shutdown signal received, exiting...")
	e.Shutdown(context.Background())
}

func sync(c *cli.Context, runner *horloge.Runner) horloge.Sync {
	switch s := c.String("sync"); s {
	case "redis":
		addr := c.String("redis-addr")
		db := c.Int("redis-db")
		passwd := c.String("redis-passwd")
		fmt.Printf("Syncing with redis %s with db %d \n", addr, db)

		return horloge.NewSyncRedis(runner, addr, passwd, db)
	case "file":
		path := c.String("file-path")
		fmt.Printf("Syncing with file: %s\n", path)

		return horloge.NewSyncDisk(runner, path)
	default:
		fmt.Println("No sync")
		return horloge.NewSyncNone()
	}
}

func bind(c *cli.Context) string {
	return fmt.Sprintf("%s:%d", c.String("bind"), c.Int("port"))
}

func main() {
	app := cli.NewApp()
	app.Name = "horloge"
	app.Version = horloge.Version
	app.Authors = []cli.Author{
		cli.Author{
			Name:  "Samori Gorse",
			Email: "samorigorse@gmail.com",
		},
	}
	app.Flags = []cli.Flag{
		cli.IntFlag{
			Name:  "port, p",
			Usage: "Port to listen to",
			Value: 6432,
		},
		cli.StringFlag{
			Name:  "bind, b",
			Usage: "Address to bind to",
			Value: "127.0.0.1",
		},
		cli.StringFlag{
			Name:  "sync",
			Usage: "Sync method to use (redis, file)",
			Value: "none",
		},
		cli.StringFlag{
			Name:  "file-path, f",
			Usage: "Output file path (used with `file` sync)",
			Value: "none",
		},

		cli.StringFlag{
			Name:  "redis-addr",
			Usage: "Address of the redis server (used with `redis` sync)",
			Value: "localhost:6379",
		},
		cli.StringFlag{
			Name:  "redis-passwd",
			Usage: "Password of the redis server (used with `redis` sync)",
			Value: "",
		},
		cli.IntFlag{
			Name:  "redis-db",
			Usage: "Which database to use (used with `redis` sync)",
			Value: 0,
		},
	}

	app.Action = func(c *cli.Context) error {
		bindTo := bind(c)

		runner := horloge.NewRunner()
		sync := sync(c, runner)
		runner.Sync(sync)

		server(bindTo, runner)

		return nil
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
