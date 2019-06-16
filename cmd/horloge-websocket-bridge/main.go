package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/codegangsta/cli"
	"github.com/go-redis/redis"
	"github.com/majordome-iot/horloge"
)

func main() {
	app := cli.NewApp()
	app.Name = "horloge-websocket-bridge"
	app.Version = horloge.Version
	app.Authors = []cli.Author{
		cli.Author{
			Name:  "Samori Gorse",
			Email: "samorigorse@gmail.com",
		},
	}
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "redis-addr",
			Usage: "Address of the redis server",
			Value: "localhost:6379",
		},
		cli.StringFlag{
			Name:  "redis-passwd",
			Usage: "Password of the redis server",
			Value: "",
		},
		cli.IntFlag{
			Name:  "redis-db",
			Usage: "Which database to use",
			Value: 0,
		},
		cli.StringFlag{
			Name:  "websocket-addr",
			Usage: "Address to bind the websocket server to",
			Value: ":5000",
		},
	}

	app.Action = func(c *cli.Context) error {
		httpServer := horloge.NewWebsocketServer()
		redisClient := horloge.NewRedisClient(c.String("redis-addr"), c.String("redis-passwd"), c.Int("redis-db"))
		signalChan := make(chan os.Signal, 1)

		redisClient.AddPublishHandler(func(msg *redis.Message) {
			httpServer.Publish(msg.Payload)
		})

		go redisClient.Wait()
		httpServer.Run(c.String("websocket-addr"))

		signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)
		<-signalChan

		return nil
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}

}
