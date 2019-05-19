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

	_, err := runner.AddJob(job)
	if err != nil {
		t.Errorf("expected runner not to return an error when adding a job")
	}

	_, err = runner.AddJob(job)
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

func TestToJSON(t *testing.T) {
	runner := NewRunner()
	job := NewJob("foobar", Pattern{})
	runner.AddJob(job)

	json, _ := runner.ToJSON()
	newRunner, err := NewFromJson(json)
	_, ok := newRunner.jobs[job.Name]

	if err != nil {
		t.Errorf("Error while loading from JSON: %v", err)
	}

	if ok != true {
		t.Errorf("Job %s not found in runner", job.Name)
	}
}

func TestRemoveJob(t *testing.T) {
	runner := NewRunner()
	job := NewJob("foobar", Pattern{})

	_, err := runner.AddJob(job)
	if err != nil {
		t.Errorf("expected runner not to return an error")
	}

	runner.RemoveJob(job)

	_, err = runner.AddJob(job)
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
		t.Errorf("expected handler to be called with args %v, but got %v", expected, actual)
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
