package models

import "time"

func CurMonth() string {
	return time.Now().Format("2006-01")
}
func CurDay() string {
	return time.Now().Format(time.DateOnly)
}
func Month(handledAt uint64) string {
	return time.Unix(int64(handledAt), 0).Format("2006-01")
}
func Day(handledAt uint64) string {
	return time.Unix(int64(handledAt), 0).Format(time.DateOnly)
}
func Hour(handledAt uint64) string {
	return time.Unix(int64(handledAt), 0).Format("2006-01-02-15")
}
