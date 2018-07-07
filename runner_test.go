package horloge

import (
	"reflect"
	"testing"
	"time"
)

func TestAddJob(t *testing.T) {
	runner := NewRunner()

	pattern := Pattern{Occurence: "every", Second: 1}
	job := NewJob("foobar", pattern)

	err := runner.AddJob(job)
	if err != nil {
		t.Errorf("expected runner not to return an error when adding a job")
	}

	err = runner.AddJob(job)
	if err == nil {
		t.Errorf("expected runner to return an error when adding a job with the same name")
	}
}

func TestHasJob(t *testing.T) {
	runner := NewRunner()
	job := NewJob("foobar", Pattern{})
	runner.AddJob(job)
	actual := runner.HasJob(job)

	if !actual {
		t.Errorf("expected HasJob to return true, but got %v", actual)
	}
}

func TestRemoveJob(t *testing.T) {
	runner := NewRunner()
	job := NewJob("foobar", Pattern{})

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
	var actual []string
	expected := []string{"foo", "bar"}

	now := time.Now()
	runner := NewRunner()
	job := NewJob("foobar", Pattern{}, expected)

	runner.AddJob(job)
	runner.AddHandler("foobar", func(name string, args []string, t time.Time) {
		actual = args
	})
	runner.Execute(job, now)

	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("expected runner to be called with args %v, but got %v", expected, actual)
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
