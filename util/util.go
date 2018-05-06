package util

import "time"

func ParseLocation(locationStr string) (loc *time.Location, err error) {
	if locationStr == "" {
		return time.UTC, nil
	}
	return time.LoadLocation(locationStr)
}
