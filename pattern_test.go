package horloge

import "testing"

func TestPatternTime(t *testing.T) {
	p := Pattern{}

	if p.Time().IsZero() {
		t.Errorf("expected p.Time() not to return a zeroed time")
	}
}
