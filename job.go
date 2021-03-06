package horloge

import (
	"time"
)

type Job struct {
	Name        string
	Description string
	Args        []string
	Pattern     Pattern
	tickers     []*time.Ticker
}

// NewJob Creates a new Job that will occur based on the provided pattern.
func NewJob(name string, pattern Pattern, args []string) *Job {
	return &Job{
		Name:    name,
		Args:    args,
		Pattern: pattern,
	}
}

// Bind Attachs arguments to a Job.
func (j *Job) Bind(args []string) {
	j.Args = args
}

// Repeat Returns next occurences.
func (j *Job) Repeat() []time.Time {
	var time []time.Time
	p := j.Pattern

	switch occurence := p.Occurence; occurence {
	case "every":
		time = j.every(p)
	case "yearly":
		time = j.yearly(p)
	case "monthly":
		time = j.monthly(p)
	case "weekly":
		time = j.weekly(p)
	case "daily":
		time = j.daily(p)
	}

	return time
}

// Cancel Cancels a job
func (j *Job) Cancel() {
	for _, t := range j.tickers {
		t.Stop()
	}
}

// Calendar Creates a calendar of every
func (j *Job) Calendar() []time.Time {
	nexts := j.Repeat()

	j.tickers = make([]*time.Ticker, len(nexts))

	return nexts
}

func (j *Job) every(p Pattern) []time.Time {

	return []time.Time{
		p.alignDateTime(p.Time()),
	}
}

func (j *Job) daily(p Pattern) []time.Time {
	now := p.Time()
	midnight := Bod(now)
	tomorrowMidnight := Bod(tomorrow(now))

	next := p.alignClock(midnight)

	// Execution time has already passed for today
	// We scheduled it to the next day
	if next.Before(now) {
		next = p.alignClock(tomorrowMidnight)
	}

	return []time.Time{
		next,
	}
}

func (j *Job) weekly(p Pattern) []time.Time {
	now := p.Time()
	midnight := Bod(now)
	next := p.alignClock(midnight)

	nexts := make([]time.Time, len(p.Days))

	for i, wd := range p.Days {
		daysDiff := DiffToWeekday(midnight, wd)
		// Occurence happens later today
		if daysDiff == 7 && next.After(now) {
			daysDiff = 0
		}
		nextExecutionDate := midnight.AddDate(0, 0, daysDiff)
		nexts[i] = p.alignClock(nextExecutionDate)
	}

	return nexts
}

func (j *Job) monthly(p Pattern) []time.Time {
	now := p.Time()
	firstDayOfTheMonth := Bom(now)
	next := p.alignDateTime(firstDayOfTheMonth)
	nexts := make([]time.Time, len(p.Months))

	days := p.Day
	if days > 0 {
		days = days - 1
	}

	for i, m := range p.Months {
		monthsDiff := DiffToMonth(firstDayOfTheMonth, m)

		// Occurence happens later this month
		if monthsDiff == 12 && next.After(now) {
			monthsDiff = 0
		}

		nextExecutionDate := firstDayOfTheMonth.AddDate(0, monthsDiff, days)
		nexts[i] = p.alignClock(nextExecutionDate)
	}

	return nexts
}

func (j *Job) yearly(p Pattern) []time.Time {
	now := p.Time()
	firstDayOfTheYear := Boy(now)
	next := p.alignDateTime(firstDayOfTheYear)

	days := p.Day
	if days > 0 {
		days = days - 1
	}

	months := p.Month
	if months > 0 {
		months = months - 1
	}

	year := 1
	if next.After(now) {
		year = 0
	}

	nextExecutionDate := firstDayOfTheYear.AddDate(year, months, days)
	return []time.Time{
		p.alignClock(nextExecutionDate),
	}
}

func tomorrow(t time.Time) time.Time {
	return t.AddDate(0, 0, 1)
}
