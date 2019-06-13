package horloge

import (
	"encoding/json"

	"github.com/go-redis/redis"
)

const HORLOGE_KEY = "horloge_jobs"

type SyncRedis struct {
	Client   *redis.Client
	runner   *Runner
	Addr     string
	Password string
	DB       int
}

func NewSyncRedis(runner *Runner, addr string, password string, db int) *SyncRedis {
	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
	})

	return &SyncRedis{
		Client: client,
		runner: runner,
	}
}

func (s *SyncRedis) Read() []Job {
	var jobs []Job

	value := s.Client.Get(HORLOGE_KEY)

	if value.Err() != nil {
		panic(value.Err())
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
