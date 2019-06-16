package horloge

import (
	"log"
	"os"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

type MQTTClient struct {
	client mqtt.Client
}

func NewMQTTClient(addr string) *MQTTClient {
	mqtt.DEBUG = log.New(os.Stdout, "", 0)
	mqtt.ERROR = log.New(os.Stdout, "", 0)

	opts := mqtt.NewClientOptions().AddBroker(addr)
	opts.SetKeepAlive(2 * time.Second)
	opts.SetPingTimeout(1 * time.Second)

	c := mqtt.NewClient(opts)

	if token := c.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}

	return &MQTTClient{
		client: c,
	}
}

func (c *MQTTClient) Publish(channel string, payload string) {
	c.client.Publish(channel, 0, false, payload)
}
