package horloge

import (
	"time"
)

type Job struct {
	Name    string
	Args    []string
	Pattern Pattern
	tickers []*time.Ticker
}

func NewJob(name string, pattern Pattern, args ...string) *Job {
	return &Job{
		Name:    name,
		Args:    args,
		Pattern: pattern,
	}
}

func (j *Job) Repeat() []time.Time {
	var time []time.Time
	p := j.Pattern
	occurence := p.Occurence

	if occurence == "every" {
		time = j.every(p)
	} else if occurence == "yearly" {
		time = j.yearly(p)
	} else if occurence == "monthly" {
		time = j.monthly(p)
	} else if occurence == "weekly" {
		time = j.weekly(p)
	} else if occurence == "daily" {
		time = j.daily(p)
	}

	return time
}

func (j *Job) Cancel() {
	for _, t := range j.tickers {
		t.Stop()
	}
}

func (j *Job) Calendar() []time.Time {
	nexts := j.Repeat()

	j.tickers = make([]*time.Ticker, len(nexts))

	return nexts
}

func (j *Job) every(p Pattern) []time.Time {

	return []time.Time{
		alignDateTime(p.Time(), p),
	}
}

func (j *Job) daily(p Pattern) []time.Time {
	now := p.Time()
	midnight := Bod(now)
	tomorrowMidnight := Bod(tomorrow(now))

	next := alignClock(midnight, p)

	// Execution time has already passed for today
	// We scheduled it to the next day
	if next.Before(now) {
		next = alignClock(tomorrowMidnight, p)
	}

	return []time.Time{
		next,
	}
}

func (j *Job) weekly(p Pattern) []time.Time {
	now := p.Time()
	midnight := Bod(now)
	next := alignClock(midnight, p)

	nexts := make([]time.Time, len(p.Days))

	for i, wd := range p.Days {
		daysDiff := DiffToWeekday(midnight, wd)
		// Occurence happens later today
		if daysDiff == 7 && next.After(now) {
			daysDiff = 0
		}
		nextExecutionDate := midnight.AddDate(0, 0, daysDiff)
		nexts[i] = alignClock(nextExecutionDate, p)
	}

	return nexts
}

func (j *Job) monthly(p Pattern) []time.Time {
	now := p.Time()
	firstDayOfTheMonth := Bom(now)
	next := alignDateTime(firstDayOfTheMonth, p)
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
		nexts[i] = alignClock(nextExecutionDate, p)
	}

	return nexts
}

func (j *Job) yearly(p Pattern) []time.Time {
	now := p.Time()
	firstDayOfTheYear := Boy(now)
	next := alignDateTime(firstDayOfTheYear, p)

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
		alignClock(nextExecutionDate, p),
	}
}

func tomorrow(t time.Time) time.Time {
	return t.AddDate(0, 0, 1)
}

func alignClock(t time.Time, p Pattern) time.Time {
	h := time.Hour * time.Duration(p.Hour)
	m := time.Minute * time.Duration(p.Minute)
	s := time.Second * time.Duration(p.Second)

	return t.Add(h).Add(m).Add(s)
}

func alignDateTime(t time.Time, p Pattern) time.Time {
	return alignClock(t, p).AddDate(p.Year, p.Month, p.Day)
}
