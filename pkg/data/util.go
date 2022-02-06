package data

import (
	"time"
)

func parseDateString(dateStr string) (time.Time, error) {

	date, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		return time.Time{}, err
	}
	return date, nil

}
