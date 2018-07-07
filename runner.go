package horloge

import (
	"errors"
	"fmt"
	"sync"
	"time"
)

type TaskHandler func(name string, args []string, t time.Time)

type Runner struct {
	jobs     map[string]*Job
	handlers map[string]TaskHandler
}

func NewRunner() Runner {
	r := Runner{
		jobs:     make(map[string]*Job),
		handlers: make(map[string]TaskHandler),
	}

	return r
}

func (r *Runner) AddJob(job *Job) error {
	nexts := job.Calendar()

	if r.HasJob(job) {
		return errors.New(
			fmt.Sprintf("[ERROR] Cannot add task \"%s\", another task with the same name exists", job.Name))
	}

	for i, n := range nexts {
		go r.Schedule(job, n, i)
	}

	r.jobs[job.Name] = job

	return nil
}

func (r *Runner) AddHandler(name string, f func(name string, args []string, t time.Time)) {
	r.handlers[name] = f
}

func (r *Runner) Execute(j *Job, t time.Time) {
	name := j.Name
	args := j.Args
	handler, ok := r.handlers[name]

	if ok {
		handler(name, args, t)
	} else {
		fmt.Printf("[WARN] No handlers for task %s\n", name)
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
		fmt.Printf("[INFO] Canceling job \"%s\"", job.Name)
		job.Cancel()
		delete(r.jobs, job.Name)
	}
}
