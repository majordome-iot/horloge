package horloge

import (
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/sirupsen/logrus"
)

const (
	All            string = "all"
	JobExistsError string = "Cannot add task \"%s\", another task with the same name exists"
)

type JobCallback func(name string, args []string, t time.Time)

type Runner struct {
	jobs     map[string]*Job
	handlers map[string]JobCallback
	log      *logrus.Logger
}

func NewRunner() *Runner {
	log := logrus.New()

	r := &Runner{
		jobs:     make(map[string]*Job),
		handlers: make(map[string]JobCallback),
		log:      log,
	}

	return r
}

func (r *Runner) AddJob(job *Job) ([]time.Time, error) {
	var nexts []time.Time

	if r.HasJob(job) {
		return nexts, fmt.Errorf(JobExistsError, job.Name)
	}

	nexts = job.Calendar()

	for i, n := range nexts {
		go r.Schedule(job, n, i)
	}

	r.jobs[job.Name] = job

	return nexts, nil
}

func (r *Runner) AddHandler(name string, f JobCallback) {
	r.handlers[name] = f
}

func (r *Runner) Execute(j *Job, t time.Time) {
	name := j.Name
	args := j.Args
	handler, ok := r.handlers[name]

	if ok {
		handler(name, args, t)
	}

	handler, ok = r.handlers[All]

	if ok {
		handler(name, args, t)
	}
}

func (r *Runner) Schedule(j *Job, n time.Time, index int) {
	contextLogger := r.log.WithFields(logrus.Fields{
		"at":     n,
		"name":   j.Name,
		"action": "schedule",
	})

	wg := &sync.WaitGroup{}
	wg.Add(1)

	duration := time.Until(n)
	ticker := time.NewTicker(duration)
	j.tickers[index] = ticker

	contextLogger.Info(fmt.Sprintf("Scheduling task \"%s\" at %s\n", j.Name, n.String()))

	select {
	case t := <-ticker.C:
		r.Execute(j, t)
		next := j.Calendar()[index]
		if r.HasJob(j) {
			go r.Schedule(j, next, index)

		}
		wg.Done()
	}
	wg.Wait()
}

func (r *Runner) HasJob(job *Job) bool {
	_, ok := r.jobs[job.Name]

	return ok
}

func (r *Runner) RemoveJob(job *Job) {
	contextLogger := r.log.WithFields(logrus.Fields{
		"action": "remove",
		"name":   job.Name,
	})
	if r.HasJob(job) {
		contextLogger.Info(fmt.Sprintf("Canceling job \"%s\"\n", job.Name))
		job.Cancel()
		delete(r.jobs, job.Name)
	}
}
