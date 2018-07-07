package horloge

import (
	"reflect"
	"testing"
	"time"
)

func TestBind(t *testing.T) {
	job := NewJob("Run outside", Pattern{})
	job.Bind([]string{"in", "the", "rain"})

	actual := job.Args
	expected := []string{"in", "the", "rain"}

	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("expected handler to be called with args %v, but got %v", expected, actual)
	}
}

func TestEvery(t *testing.T) {
	then := time.Date(2015, time.October, 21, 15, 0, 0, 0, time.UTC)
	p := Pattern{Occurence: "every", Hour: 1, Minute: 30, Now: then}
	job := NewJob("Bat eyes", p)

	result := job.Repeat()
	actual := result[0].String()
	expected := "2015-10-21 16:30:00 +0000 UTC"

	if expected != actual {
		t.Errorf("expected %s to be %s", expected, actual)
	}
}

func TestDailyPassedTime(t *testing.T) {
	then := time.Date(2015, time.October, 21, 15, 0, 0, 0, time.UTC)
	p := Pattern{Occurence: "daily", Hour: 9, Minute: 30, Second: 0, Now: then}
	job := NewJob("Breakfast", p)

	result := job.Repeat()
	actual := result[0].String()
	expected := "2015-10-22 09:30:00 +0000 UTC"

	if expected != actual {
		t.Errorf("expected %s to be %s", expected, actual)
	}
}

func TestDailyFutureTime(t *testing.T) {
	then := time.Date(2015, time.October, 21, 7, 30, 0, 0, time.UTC)
	p := Pattern{Occurence: "daily", Hour: 9, Minute: 30, Second: 0, Now: then}
	job := NewJob("Breakfast", p)

	result := job.Repeat()
	actual := result[0].String()
	expected := "2015-10-21 09:30:00 +0000 UTC"

	if expected != actual {
		t.Errorf("expected %s to be %s", expected, actual)
	}
}

func TestWeekly(t *testing.T) {
	then := time.Date(2015, time.October, 21, 15, 23, 0, 0, time.UTC) // This is a Wednesday
	days := []time.Weekday{time.Monday, time.Thursday}

	p := Pattern{Occurence: "weekly", Days: days, Hour: 20, Minute: 0, Second: 0, Now: then}
	job := NewJob("Take out the trash", p)

	actual := job.Repeat()
	expected := []string{
		"2015-10-26 20:00:00 +0000 UTC",
		"2015-10-22 20:00:00 +0000 UTC",
	}

	for i, e := range expected {
		if e != actual[i].String() {
			t.Errorf("expected %s to be %s", actual[i].String(), e)
		}
	}
}

func TestWeeklySameDay(t *testing.T) {
	then := time.Date(2015, time.October, 21, 8, 0, 0, 0, time.UTC) // This is a Wednesday
	days := []time.Weekday{time.Wednesday}

	p := Pattern{Occurence: "weekly", Days: days, Hour: 12, Minute: 0, Second: 0, Now: then}
	job := NewJob("Picnic with the park", p)

	result := job.Repeat()
	actual := result[0].String()
	expected := "2015-10-21 12:00:00 +0000 UTC"

	if expected != actual {
		t.Errorf("expected %s to be %s", expected, actual)
	}
}
func TestMonthly(t *testing.T) {
	then := time.Date(2015, time.October, 21, 15, 0, 0, 0, time.UTC)
	months := []time.Month{time.January, time.December}

	p := Pattern{Occurence: "monthly", Months: months, Day: 20, Now: then}
	job := NewJob("Fill out my W-2", p)

	result := job.Repeat()
	expected := []string{
		"2016-01-20 00:00:00 +0000 UTC",
		"2015-12-20 00:00:00 +0000 UTC",
	}

	for i, e := range expected {
		actual := result[i].String()
		if actual != e {
			t.Errorf("expected %s to be %s", actual, e)
		}
	}
}

func TestMonthlyWithDatetime(t *testing.T) {
	then := time.Date(2015, time.October, 21, 15, 0, 0, 0, time.UTC)
	months := []time.Month{time.January, time.December}

	p := Pattern{Occurence: "monthly", Months: months, Day: 20, Hour: 12, Minute: 30, Now: then}
	job := NewJob("Lunch with my step mom", p)

	result := job.Repeat()
	expected := []string{
		"2016-01-20 12:30:00 +0000 UTC",
		"2015-12-20 12:30:00 +0000 UTC",
	}

	for i, e := range expected {
		actual := result[i].String()
		if actual != e {
			t.Errorf("expected %s to be %s", actual, e)
		}
	}
}

func TestMonthlyWithoutDate(t *testing.T) {
	then := time.Date(2015, time.October, 21, 15, 0, 0, 0, time.UTC)
	months := []time.Month{time.February, time.September}

	p := Pattern{Occurence: "monthly", Months: months, Now: then}
	job := NewJob("Fill out my W-2", p)

	actual := job.Repeat()
	expected := []string{
		"2016-02-01 00:00:00 +0000 UTC",
		"2016-09-01 00:00:00 +0000 UTC",
	}

	for i, e := range expected {
		if e != actual[i].String() {
			t.Errorf("expected %s to be %s", actual[i].String(), e)
		}
	}
}

func TestMonthlySameMonth(t *testing.T) {
	then := time.Date(2015, time.October, 21, 15, 23, 0, 0, time.UTC)
	months := []time.Month{time.October}

	p := Pattern{Occurence: "monthly", Months: months, Day: 23, Now: then}
	job := NewJob("Fill out my W-2", p)

	result := job.Repeat()
	actual := result[0].String()
	expected := "2015-10-23 00:00:00 +0000 UTC"

	if expected != actual {
		t.Errorf("expected %s to be %s", expected, actual)
	}
}

func TestYearly(t *testing.T) {
	then := time.Date(2015, time.October, 21, 15, 0, 0, 0, time.UTC)

	p := Pattern{Occurence: "yearly", Month: 7, Day: 4, Hour: 22, Minute: 35, Now: then}
	job := NewJob("Mow the lawn", p)

	result := job.Repeat()
	actual := result[0].String()
	expected := "2016-07-04 22:35:00 +0000 UTC"

	if expected != actual {
		t.Errorf("expected %s to be %s", expected, actual)
	}
}

func TestYearlyWithoutDate(t *testing.T) {
	then := time.Date(2015, time.October, 21, 15, 0, 0, 0, time.UTC)

	p := Pattern{Occurence: "yearly", Now: then}
	job := NewJob("Mow the lawn", p)

	result := job.Repeat()
	actual := result[0].String()
	expected := "2016-01-01 00:00:00 +0000 UTC"

	if expected != actual {
		t.Errorf("expected %s to be %s", expected, actual)
	}
}

func TestYearlyUpcomingDate(t *testing.T) {
	then := time.Date(2015, time.October, 21, 15, 0, 0, 0, time.UTC)

	p := Pattern{Occurence: "yearly", Month: int(time.October), Day: 22, Now: then}
	job := NewJob("Mow the lawn", p)

	result := job.Repeat()
	actual := result[0].String()
	expected := "2015-10-22 00:00:00 +0000 UTC"

	if expected != actual {
		t.Errorf("expected %s to be %s", expected, actual)
	}
}
