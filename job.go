package horloge

import (
	"time"
)

type Job struct {
	task    *Task
	pattern Pattern
	tickers []*time.Ticker
}

func NewJob(task *Task, pattern Pattern) *Job {
	return &Job{
		task:    task,
		pattern: pattern,
	}
}

func (j *Job) Cancel() {
	for _, t := range j.tickers {
		t.Stop()
	}
}

func (j *Job) Calendar() []time.Time {
	task := j.task
	pattern := j.pattern
	nexts := task.Repeat(pattern)

	j.tickers = make([]*time.Ticker, len(nexts))

	return nexts
}
