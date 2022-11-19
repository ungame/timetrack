package timeext

import "time"

const (
	DateOnlyFormat = "2006-01-02"
	DateTimeFormat = "2006-01-02 15:04:05"
)

func GetStartOfDay(location *time.Location) time.Time {
	return GetStartOfDayFrom(time.Now(), location)
}

func GetEndOfDay(location *time.Location) time.Time {
	return GetEndOfDayFrom(time.Now(), location)
}

func GetStartOfDayFrom(t time.Time, location *time.Location) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, location)
}

func GetEndOfDayFrom(t time.Time, location *time.Location) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), 23, 59, 59, 0, location)
}
