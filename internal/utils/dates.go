package utils

import "time"

func GetWEDate(date time.Time) time.Time {
	for date.Weekday() != 6 {
		date = date.Add(time.Hour * 24)
	}
	return date
}

func GetPastWEDate() string {
	now := time.Now()

	now = now.Add(time.Hour * -24 * 5)

	for now.Weekday() != 6 {
		now = now.Add(time.Hour * -24)
	}

	return now.Format("2006-01-02")
}
