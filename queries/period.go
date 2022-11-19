package queries

import (
	"github.com/ungame/timetrack/timeext"
	"time"
)

type Period interface {
	Range(location *time.Location) (time.Time, time.Time)
}

type PeriodType string

const (
	Today     PeriodType = "today"
	Yesterday PeriodType = "yesterday"
	Weekly    PeriodType = "weekly"
	Monthly   PeriodType = "monthly"
)

func (p PeriodType) String() string {
	return string(p)
}

func (p PeriodType) IsValid() bool {
	switch p {
	case Today:
	case Yesterday:
	case Weekly:
	case Monthly:
	default:
		return false
	}
	return true
}

func (p PeriodType) Range(location *time.Location) (time.Time, time.Time) {
	var (
		now   = time.Now()
		start time.Time
		end   time.Time
	)
	switch p {
	case Today:
		start = timeext.GetStartOfDay(location)
		end = now

	case Yesterday:
		yesterday := now.AddDate(0, 0, -1)
		start = timeext.GetStartOfDayFrom(yesterday, location)
		end = timeext.GetEndOfDayFrom(yesterday, location)

	case Weekly:
		sevenDaysAgo := now.AddDate(0, 0, -7)
		start = timeext.GetStartOfDayFrom(sevenDaysAgo, location)
		end = now

	case Monthly:
		thirtyDaysAgo := now.AddDate(0, 0, -30)
		start = timeext.GetStartOfDayFrom(thirtyDaysAgo, location)
		end = now
	}

	return start, end
}
