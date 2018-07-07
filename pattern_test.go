package horloge

import (
	"reflect"
	"testing"
	"time"
)

func TestPatternTime(t *testing.T) {
	p := Pattern{}

	if p.Time().IsZero() {
		t.Errorf("expected p.Time() not to return a zeroed time")
	}
}

func TestPatternOnWeekday(t *testing.T) {
	p := NewPattern("weekly").On(time.Tuesday)

	expected := []time.Weekday{time.Tuesday}
	actual := p.Days

	if !reflect.DeepEqual(expected, actual) {
		t.Errorf("expected pattern.Days to be %v, but got %v", expected, actual)
	}
}

func TestPatternOnMonths(t *testing.T) {
	p := NewPattern("monthly").On(time.January, time.October)

	expected := []time.Month{time.January, time.October}
	actual := p.Months

	if !reflect.DeepEqual(expected, actual) {
		t.Errorf("expected pattern.Days to be %v, but got %v", expected, actual)
	}
}

func TestPatternAt(t *testing.T) {
	p := NewPattern("every").At(14, 5, 5)

	expected := p.Time()
	actual := p.Months

	if !reflect.DeepEqual(expected, actual) {
		t.Errorf("expected pattern.Days to be %v, but got %v", expected, actual)
	}
}
