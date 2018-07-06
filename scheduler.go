package horloge

import (
	"time"
)

type Pattern struct {
	Days      []time.Weekday
	Months    []time.Month
	Second    int
	Minute    int
	Hour      int
	Day       int
	Month     int
	Year      int
	Occurence string
	Now       time.Time
}

type Task struct {
	Name    string
	Args    []string
	pattern Pattern
}

func DiffToMonth(t time.Time, m time.Month) int {
	cm := t.Month()
	if cm < m {
		return int(m - cm)
	} else {
		return int(cm)*-1 + int(m) + 12
	}
}

func DiffToWeekday(t time.Time, wd time.Weekday) int {
	cwd := t.Weekday()
	if cwd < wd {
		return int(wd - cwd)
	} else {
		return int(cwd)*-1 + int(wd) + 7
	}
}

func Bod(t time.Time) time.Time {
	year, month, day := t.Date()
	return time.Date(year, month, day, 0, 0, 0, 0, t.Location())
}

func Bom(t time.Time) time.Time {
	year, month, _ := t.Date()
	return time.Date(year, month, 1, 0, 0, 0, 0, t.Location())
}

func Boy(t time.Time) time.Time {
	year, _, _ := t.Date()
	return time.Date(year, 1, 1, 0, 0, 0, 0, t.Location())
}

func (p *Pattern) Time() time.Time {
	if p.Now.IsZero() {
		p.Now = time.Now()
	}

	return p.Now
}

func NewTask(Name string, Args ...string) *Task {
	return &Task{
		Name: Name,
		Args: Args,
	}
}

func (t *Task) Repeat(p Pattern) []time.Time {
	var time []time.Time
	occurence := p.Occurence

	if occurence == "every" {
		time = t.every(p)
	} else if occurence == "yearly" {
		time = t.yearly(p)
	} else if occurence == "monthly" {
		time = t.monthly(p)
	} else if occurence == "weekly" {
		time = t.weekly(p)
	} else if occurence == "daily" {
		time = t.daily(p)
	}

	return time

}

func (t *Task) every(p Pattern) []time.Time {

	return []time.Time{
		alignDateTime(p.Time(), p),
	}
}

func (t *Task) daily(p Pattern) []time.Time {
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

func (t *Task) weekly(p Pattern) []time.Time {
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

func (t *Task) monthly(p Pattern) []time.Time {
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

func (t *Task) yearly(p Pattern) []time.Time {
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
