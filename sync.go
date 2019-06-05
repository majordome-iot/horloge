package horloge

import (
	"encoding/json"
	"io/ioutil"
	"os"
)

type Sync interface {
	Write([]*Job) error
	Read() []Job
}

type SyncRedis struct {
	Addr     string
	Password string
	DB       int
}

type SyncDisk struct {
	Path string
}

type SyncMemory struct {
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func (s *SyncDisk) Read() []Job {
	var jobs []Job
	data, err := ioutil.ReadFile(s.Path)

	if os.IsNotExist(err) {
		data = []byte("[]")
		ioutil.WriteFile(s.Path, data, 0644)
	} else {
		check(err)
	}

	err = json.Unmarshal(data, &jobs)

	check(err)

	return jobs
}

func (s *SyncDisk) Write(jobs []*Job) error {
	slice := make([]*Job, 0)

	for _, job := range jobs {
		slice = append(slice, job)
	}

	data, err := json.Marshal(slice)

	if err != nil {
		return err
	}

	return ioutil.WriteFile(s.Path, data, 0644)
}

func (s *SyncRedis) Write(jobs []*Job) error {
	return nil
}
