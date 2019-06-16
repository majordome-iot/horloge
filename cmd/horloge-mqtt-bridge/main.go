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
	app.Name = "horloge-mqtt-bridge"
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
			Name:  "mqtt-addr",
			Usage: "Address of the mqtt broker",
			Value: "tcp://localhost:1883",
		},
	}

	app.Action = func(c *cli.Context) error {
		mqttAddr := c.String("mqtt-addr")
		redisAddr, redisPasswd, redisDb := c.String("redis-addr"), c.String("redis-passwd"), c.Int("redis-db")

		mqttClient := horloge.NewMQTTClient(mqttAddr)
		redisClient := horloge.NewRedisClient(redisAddr, redisPasswd, redisDb)
		signalChan := make(chan os.Signal, 1)

		redisClient.AddPublishHandler(func(msg *redis.Message) {
			mqttClient.Publish(horloge.MQTT_CHANNEL, msg.Payload)
		})

		go redisClient.Wait()

		signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)
		<-signalChan

		return nil
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
