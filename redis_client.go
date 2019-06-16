package horloge

import (
	"github.com/go-redis/redis"
)

const MQTT_CHANNEL = "horloge/job"

type RedisClient struct {
	pubsub   *redis.PubSub
	handlers []PublishHandler
}

type PublishHandler func(*redis.Message)

func NewRedisClient(addr string, passwd string, db int) *RedisClient {
	c := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: passwd,
		DB:       db,
	})

	var handlers []PublishHandler
	pubsub := c.Subscribe(PUBSUB_CHANNEL)

	_, err := pubsub.Receive()
	if err != nil {
		panic(err)
	}

	return &RedisClient{
		pubsub:   pubsub,
		handlers: handlers,
	}
}

func (c *RedisClient) Wait() {
	channel := c.pubsub.Channel()
	for msg := range channel {
		for _, handler := range c.handlers {
			handler(msg)
		}
	}
}

func (c *RedisClient) AddPublishHandler(handler PublishHandler) {
	c.handlers = append(c.handlers, handler)
}
