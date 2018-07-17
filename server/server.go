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

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/majordome/horloge"
	melody "gopkg.in/olahol/melody.v1"
)

type Event struct {
	Name string
	Args []string
	Time time.Time
}

func main() {
	e := echo.New()
	m := melody.New()
	runner := horloge.NewRunner()

	runner.AddHandler("all", func(name string, args []string, t time.Time) {
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

	httpAddr := ""
	httPort := 8080
	addr := fmt.Sprintf("%s:%d", httpAddr, httPort)

	go func() {
		e.Logger.Infof("HTTP Server Listening to %s\n", addr)
		e.Logger.Fatal(e.Start(addr))
	}()

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)
	<-signalChan

	log.Println("Shutdown signal received, exiting...")
	e.Shutdown(context.Background())
}
