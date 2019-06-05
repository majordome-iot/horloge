package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/codegangsta/cli"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/majordome-iot/horloge"
	melody "gopkg.in/olahol/melody.v1"
)

type Event struct {
	Name string
	Args []string
	Time time.Time
}

func server(addr string) {
	e := echo.New()
	m := melody.New()

	// runner := horloge.NewRunner(horloge.SyncRedis{})
	runner := horloge.NewRunner()

	runner.AddHandler("*", func(arguments ...interface{}) {
		name, args, t := horloge.JobArgs(arguments)
		event := Event{
			Name: name,
			Args: args,
			Time: t,
		}
		b, err := json.Marshal(event)
		if err != nil {
			m.Broadcast([]byte("Fail"))
		} else {
			m.Broadcast(b)
		}
	})

	e.HideBanner = true
	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Routes
	e.GET("/ping", horloge.HTTPHandlerPing())
	e.GET("/health_check", horloge.HTTPHandlerHealthCheck())
	e.GET("/version", horloge.HTTPHandlerVersion())
	e.POST("/job", horloge.HTTPHandlerRegisterJob(runner))
	e.GET("/ws", func(c echo.Context) error {
		m.HandleRequest(c.Response().Writer, c.Request())
		return nil
	})

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

func main() {
	app := cli.NewApp()
	app.Name = "horloge"
	app.Version = horloge.Version
	app.Authors = []cli.Author{
		cli.Author{
			Name:  "Samori Gorse",
			Email: "samorigorse@gail.com",
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
	}

	app.Action = func(c *cli.Context) error {
		host := c.String("bind")
		port := c.Int("port")

		server(fmt.Sprintf("%s:%d", host, port))

		return nil
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
