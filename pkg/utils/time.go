package utils

import "time"

func IsValidTimeFormat(timeString string) bool {
	_, err := time.Parse("2006-01-02", timeString)
	return err == nil
}
