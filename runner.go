package horloge

import (
	"fmt"
	"io/ioutil"
	"os"
	"sync"
	"time"

	"github.com/sirupsen/logrus"
)

const (
	JobExistsError string = "Cannot add task \"%s\", another task with the same name exists"
)

type Callback func(...interface{})

type Runner struct {
	jobs     map[string]*Job
	handlers map[string][]Callback
	log      *logrus.Logger
}

func JobArgs(a []interface{}) (string, []string, time.Time) {
	return a[0].(string), a[1].([]string), a[2].(time.Time)
}

func NewRunner() *Runner {
	log := logrus.New()
	if os.Getenv("HORLOGE_TEST") != "" {
		log.SetOutput(ioutil.Discard)
	}

	r := &Runner{
		jobs:     make(map[string]*Job),
		handlers: make(map[string][]Callback),
		log:      log,
	}

	return r
}

// ToJSON Returns a list of known jobs
func (r *Runner) ToJSON() ([]*Job, error) {
	s := make([]*Job, 0)

	for _, job := range r.jobs {
		s = append(s, job)
	}
	return s, nil
}

// AddJob Registers a job.
func (r *Runner) AddJob(job Job) ([]time.Time, error) {
	var nexts []time.Time

	if r.HasJob(&job) {
		return nexts, fmt.Errorf(JobExistsError, job.Name)
	}

	nexts = job.Calendar()

	for i, n := range nexts {
		go r.Schedule(job, n, i)
	}

	r.jobs[job.Name] = &job
	r.Emit("job:add", job)

	return nexts, nil
}

// AddJobs Registers a list of jobs
func (r *Runner) AddJobs(jobs []Job) {
	for _, job := range jobs {
		r.AddJob(job)
	}
}

// GetJob Returns a job with name `name`
func (r *Runner) GetJob(name string) *Job {
	job, ok := r.jobs[name]
	if ok {
		return job
	}
	return nil
}

// Emit Calls every handlers attached to `name` event.
func (r *Runner) Emit(name string, args ...interface{}) {
	handlers, ok := r.handlers[name]

	if ok {
		for _, handler := range handlers {
			handler(args...)
		}
	}
}

// Subscribe Registers an handler as a callback to `name` event.
func (r *Runner) Subscribe(name string, f Callback) {
	_, ok := r.handlers[name]

	if !ok {
		r.handlers[name] = make([]Callback, 0)
	}

	r.handlers[name] = append(r.handlers[name], f)
}

// AddHandler Registers a callback to job event.
func (r *Runner) AddHandler(name string, f Callback) {
	r.Subscribe("job:"+name, f)
}

// Schedule Schedules a job in the future
func (r *Runner) Schedule(j Job, n time.Time, index int) {
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
		r.Emit("job:"+j.Name, j.Name, j.Args, t)
		r.Emit("job:*", j.Name, j.Args, t)
		next := j.Calendar()[index]
		if r.HasJob(&j) {
			go r.Schedule(j, next, index)

		}
		wg.Done()
	}
	wg.Wait()
}

// HasJob Returns true if a Job with the same job name exists
func (r *Runner) HasJob(job *Job) bool {
	_, ok := r.jobs[job.Name]

	return ok
}

// RemoveJob Removes and cancels job
func (r *Runner) RemoveJob(job *Job) {
	contextLogger := r.log.WithFields(logrus.Fields{
		"action": "remove",
		"name":   job.Name,
	})
	if r.HasJob(job) {
		contextLogger.Info(fmt.Sprintf("Canceling job \"%s\"\n", job.Name))
		job.Cancel()
		delete(r.jobs, job.Name)
		r.Emit("job:remove", job)
	}
}

func writeSync(r *Runner, s Sync) func(...interface{}) {
	return func(args ...interface{}) {
		newJobs, _ := r.ToJSON()
		s.Write(newJobs)
	}
}

// Sync Registers a datasource to sync to.
func (r *Runner) Sync(s Sync) {
	jobs := s.Read()
	r.AddJobs(jobs)

	r.Subscribe("job:add", writeSync(r, s))
	r.Subscribe("job:remove", writeSync(r, s))
}
