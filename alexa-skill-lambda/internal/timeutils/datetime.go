package timeutils

import "time"

func Today(location *time.Location) time.Time {
	if location == nil {
		location = time.UTC
	}
	now := time.Now().In(location)
	year, month, day := now.Date()
	today := time.Date(year, month, day, 0, 0, 0, 0, location)
	return today
}

func Tomorrow(location *time.Location) time.Time {
	return Today(location).Add(24 * time.Hour)
}

func LocalDate(t time.Time, location *time.Location) time.Time {
	if location == nil {
		location = time.UTC
	}
	startDateYear, startDateMonth, startDateDay := t.In(location).Date()
	localDate := time.Date(startDateYear, startDateMonth, startDateDay, 0, 0, 0, 0, location)
	return localDate
}
