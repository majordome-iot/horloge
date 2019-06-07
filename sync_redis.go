package horloge

import (
	"encoding/json"

	"github.com/go-redis/redis"
)

const HORLOGE_KEY = "horloge_jobs"

type SyncRedis struct {
	Addr     string
	Password string
	DB       int
}

func (s *SyncRedis) Read() []Job {
	var jobs []Job
	client := redis.NewClient(&redis.Options{
		Addr:     s.Addr,
		Password: s.Password,
		DB:       s.DB,
	})

	value := client.Get(HORLOGE_KEY)

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
	client := redis.NewClient(&redis.Options{
		Addr:     s.Addr,
		Password: s.Password,
		DB:       s.DB,
	})

	slice := make([]*Job, 0)

	for _, job := range jobs {
		slice = append(slice, job)
	}

	data, err := json.Marshal(slice)

	if err != nil {
		return err
	}

	err = client.Set(HORLOGE_KEY, data, 0).Err()
	if err != nil {
		return err
	}

	return nil
}
