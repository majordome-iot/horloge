package horloge

import (
	"encoding/json"
	"io/ioutil"
	"os"
)

type SyncDisk struct {
	Path string
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
