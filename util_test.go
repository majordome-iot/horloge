package horloge

import (
	"testing"
	"time"
)

func TestDiffToUpcomingDay(t *testing.T) {
	then := time.Date(2015, time.October, 21, 0, 0, 0, 0, time.UTC) // This is a Wednesday

	expected := 2
	actual := DiffToWeekday(then, time.Friday)

	if expected != actual {
		t.Errorf("expected %d to be %d", expected, actual)
	}
}

func TestDiffToPastDay(t *testing.T) {
	then := time.Date(2015, time.October, 21, 0, 0, 0, 0, time.UTC) // This is a Wednesday

	expected := 5
	actual := DiffToWeekday(then, time.Monday)

	if expected != actual {
		t.Errorf("expected %d to be %d", expected, actual)
	}
}

func TestDiffToUpcomingMonth(t *testing.T) {
	then := time.Date(2015, time.October, 21, 0, 0, 0, 0, time.UTC)

	expected := 2
	actual := DiffToMonth(then, time.December)

	if expected != actual {
		t.Errorf("expected %d to be %d", expected, actual)
	}
}

func TestDiffToPastMonth(t *testing.T) {
	then := time.Date(2015, time.October, 21, 0, 0, 0, 0, time.UTC)

	expected := 3
	actual := DiffToMonth(then, time.January)

	if expected != actual {
		t.Errorf("expected %d to be %d", expected, actual)
	}
}
