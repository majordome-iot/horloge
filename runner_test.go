package horloge

import (
	"testing"
	"time"
)

func TestAddJob(t *testing.T) {
	runner := NewRunner()

	pattern := Pattern{Occurence: "every", Second: 1}
	job := NewJob("foobar", pattern)

	err := runner.AddJob(job)
	if err != nil {
		t.Errorf("expected runner not to return an error")
	}

	err = runner.AddJob(job)
	if err == nil {
		t.Errorf("expected runner to return an error")
	}
}

func TestRemoveJob(t *testing.T) {
	runner := NewRunner()

	pattern := Pattern{}
	job := NewJob("foobar", pattern)

	err := runner.AddJob(job)
	if err != nil {
		t.Errorf("expected runner not to return an error")
	}

	runner.RemoveJob(job)

	err = runner.AddJob(job)
	if err != nil {
		t.Errorf("expected runner to return an error")
	}
}

func TestExecuteJob(t *testing.T) {
	called := false
	now := time.Now()

	runner := NewRunner()
	pattern := Pattern{}
	job := NewJob("foobar", pattern)

	runner.AddJob(job)
	runner.AddHandler("foobar", func(name string, args []string, t time.Time) {
		called = true
	})

	runner.Execute(job, now)

	if !called {
		t.Errorf("expected runner to call handler")
	}
}

func TestExecuteJobWithNoHandlers(t *testing.T) {
	called := false
	now := time.Now()

	runner := NewRunner()
	job := NewJob("foobar", Pattern{})

	runner.AddJob(job)
	runner.Execute(job, now)

	if called {
		t.Errorf("expected runner not to call handler")
	}
}
