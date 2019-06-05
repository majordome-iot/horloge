package horloge

import (
	"testing"
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
