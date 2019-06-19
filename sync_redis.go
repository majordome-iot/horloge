package horloge

import (
	"encoding/json"
	"time"

	"github.com/go-redis/redis"
)

const HORLOGE_KEY = "horloge_jobs"

// SyncRedis
type SyncRedis struct {
	Client   *redis.Client
	runner   *Runner
	Addr     string
	Password string
	DB       int
}

// Event Used to serialize events published in Redis
type Event struct {
	Name string    `json:"name"`
	Args []string  `json:"args"`
	Time time.Time `json:"time"`
}

func NewSyncRedis(runner *Runner, addr string, password string, db int) *SyncRedis {
	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
	})

	runner.AddHandler("*", func(arguments ...interface{}) {
		name, args, time := JobArgs(arguments)
		data, err := json.Marshal(Event{
			name,
			args,
			time,
		})

		if err != nil {
			panic(err)
		}

		err = client.Publish(PUBSUB_CHANNEL, string(data)).Err()

		if err != nil {
			panic(err)
		}
	})

	return &SyncRedis{
		Client: client,
		runner: runner,
	}
}

func (s *SyncRedis) Read() []Job {
	var jobs []Job

	value := s.Client.Get(HORLOGE_KEY)
	err := value.Err()

	if err != redis.Nil && err != nil {
		panic(err)
	}

	data := []byte(value.Val())

	if len(data) > 0 {
		err := json.Unmarshal(data, &jobs)
		check(err)
	}

	return jobs
}

func (s *SyncRedis) Write(jobs []*Job) error {
	var slice []*Job

	for _, job := range jobs {
		slice = append(slice, job)
	}

	data, err := json.Marshal(slice)

	if err != nil {
		return err
	}

	err = s.Client.Set(HORLOGE_KEY, data, 0).Err()
	if err != nil {
		return err
	}

	return nil
}
