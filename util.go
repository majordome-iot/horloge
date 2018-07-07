package horloge

import "time"

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
