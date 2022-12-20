package utils

import "time"

func WithinTimeSpan(start, end, check time.Time) bool {
	return check.After(start) && check.Before(end)
}

func CheckTimeIsAfter(lastUpdate time.Time, delay time.Duration) (timeNow time.Time, valid bool) {
	loc, _ := time.LoadLocation("UTC")

	lastUpdate = lastUpdate.In(loc)
	timeNow = time.Now().In(loc)

	canUpdateAfter := lastUpdate.Add(delay)

	valid = !WithinTimeSpan(
		lastUpdate,
		canUpdateAfter,
		timeNow,
	)
	return
}
