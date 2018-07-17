package horloge

import (
	"errors"
	"fmt"
	"sync"
	"time"
)

const (
	ALL string = "all"
)

type JobCallback func(name string, args []string, t time.Time)

type Runner struct {
	jobs     map[string]*Job
	handlers map[string]JobCallback
}

func NewRunner() Runner {
	r := Runner{
		jobs:     make(map[string]*Job),
		handlers: make(map[string]JobCallback),
	}

	return r
}

// Multiple return values
func (r *Runner) AddJob(job *Job) ([]time.Time, error) {
	var nexts []time.Time

	if r.HasJob(job) {
		return nexts, errors.New(
			fmt.Sprintf("[ERROR] Cannot add task \"%s\", another task with the same name exists", job.Name))
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

	handler, ok = r.handlers[ALL]

	if ok {
		handler(name, args, t)
	}
}

func (r *Runner) Schedule(j *Job, n time.Time, index int) {
	fmt.Printf("[INFO] Scheduling task \"%s\" at %s\n", j.Name, n.String())
	wg := &sync.WaitGroup{}
	wg.Add(1)

	duration := time.Until(n)
	ticker := time.NewTicker(duration)
	j.tickers[index] = ticker

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
	if r.HasJob(job) {
		fmt.Printf("[INFO] Canceling job \"%s\"\n", job.Name)
		job.Cancel()
		delete(r.jobs, job.Name)
	}
}
