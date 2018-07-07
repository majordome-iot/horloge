package horloge

import "time"

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

func (p *Pattern) Time() time.Time {
	if p.Now.IsZero() {
		p.Now = time.Now()
	}

	return p.Now
}
