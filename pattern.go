package horloge

import (
	"fmt"
	"reflect"
	"time"
)

type Pattern struct {
	Second    int    `json:"second" form:"second" query:"second"`
	Minute    int    `json:"minute" form:"minute" query:"minute"`
	Hour      int    `json:"hour" form:"hour" query:"hour"`
	Day       int    `json:"day" form:"day" query:"day"`
	Month     int    `json:"month" form:"month" query:"month"`
	Year      int    `json:"year" form:"year" query:"year"`
	Occurence string `json:"occurence" form:"occurence" query:"occurence"`
	Days      []time.Weekday
	Months    []time.Month
	Now       time.Time
}

func NewPattern(occurence string) *Pattern {
	return &Pattern{
		Occurence: occurence,
	}
}

func (p *Pattern) On(repeaters ...interface{}) *Pattern {
	var days []time.Weekday
	var months []time.Month

	for _, r := range repeaters {
		if p.Occurence == "weekly" {
			if reflect.ValueOf(r).String() != "time.Weekday" {
				days = append(days, r.(time.Weekday))
			}
		}

		if p.Occurence == "monthly" {
			if reflect.ValueOf(r).String() != "time.Month" {
				months = append(months, r.(time.Month))
			}
		}
	}

	if len(days) > 0 {
		p.Days = days
	}

	if len(months) > 0 {
		p.Months = months
	}

	return p
}

func (p *Pattern) At(hour, minute, second int) *Pattern {
	p.Hour = hour
	p.Minute = minute
	p.Second = second
	return p
}

func (p *Pattern) Time() time.Time {
	if p.Now.IsZero() {
		p.Now = time.Now()
	}

	return p.Now
}

func (p *Pattern) IsZero() bool {
	fmt.Println(p)

	return p.Second == 0 &&
		p.Minute == 0 &&
		p.Hour == 0 &&
		p.Day == 0 &&
		p.Month == 0 &&
		p.Year == 0 &&
		len(p.Months) == 0 &&
		len(p.Days) == 0 &&
		p.Occurence == ""
}

func (p *Pattern) alignClock(t time.Time) time.Time {
	h := time.Hour * time.Duration(p.Hour)
	m := time.Minute * time.Duration(p.Minute)
	s := time.Second * time.Duration(p.Second)

	return t.Add(h).Add(m).Add(s)
}

func (p *Pattern) alignDateTime(t time.Time) time.Time {
	return p.alignClock(t).AddDate(p.Year, p.Month, p.Day)
}
