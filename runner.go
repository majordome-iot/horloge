package horloge

import (
	"fmt"
	"sync"
	"time"
)

type TaskHandler func(name string, args []string, t time.Time)

type Runner struct {
	jobs     []*Job
	handlers map[string]TaskHandler
}

func NewRunner() Runner {
	r := Runner{
		handlers: make(map[string]TaskHandler),
	}

	return r
}

func (r *Runner) AddTask(task *Task, pattern Pattern) *Job {
	j := NewJob(task, pattern)

	nexts := j.Calendar()
	// TODO: Do this automatically
	j.tickers = make([]*time.Ticker, len(nexts))

	for i, n := range nexts {
		go r.Schedule(j, n, i)
	}

	r.jobs = append(r.jobs, j)

	return j
}

func (r *Runner) AddHandler(name string, f func(name string, args []string, t time.Time)) {
	r.handlers[name] = f
}

func (r *Runner) Execute(j *Job, t time.Time) {
	name := j.task.Name
	args := j.task.Args
	handler, ok := r.handlers[name]

	if ok {
		handler(name, args, t)
	} else {
		fmt.Printf("[WARN] No handlers for task %s\n", name)
	}
}

func (r *Runner) Schedule(j *Job, n time.Time, index int) {
	fmt.Printf("[INFO] Scheduling task \"%s\" at %s\n", j.task.Name, n.String())
	wg := &sync.WaitGroup{}
	wg.Add(1)

	duration := time.Until(n)
	ticker := time.NewTicker(duration)
	j.tickers[index] = ticker

	select {
	case t := <-ticker.C:
		r.Execute(j, t)
		next := j.Calendar()[index]
		if len(r.jobs) > 0 {
			go r.Schedule(j, next, index)

		}

		wg.Done()
	}

	wg.Wait()
}

func (r *Runner) Remove(job *Job) {
	for i, j := range r.jobs {
		if job.task.Name == j.task.Name {
			j.Cancel()
			r.jobs = removeElement(r.jobs, i)
		}
	}
}

func removeElement(a []*Job, i int) []*Job {
	a[i] = a[len(a)-1]
	a[len(a)-1] = nil
	a = a[:len(a)-1]

	return a
}
