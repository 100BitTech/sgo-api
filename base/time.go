package base

import (
	"time"
)

const (
	DateLayout     = time.DateOnly
	DateTimeLayout = time.DateTime
)

func ToStartOfMinute(t time.Time) time.Time {
	if t.IsZero() {
		return t
	}

	return time.Date(t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), 0, 0, t.Location())
}

func ToStartOfHour(t time.Time) time.Time {
	if t.IsZero() {
		return t
	}

	return time.Date(t.Year(), t.Month(), t.Day(), t.Hour(), 0, 0, 0, t.Location())
}

func ToStartOfDay(t time.Time) time.Time {
	if t.IsZero() {
		return t
	}

	return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
}

// 一周开始（周一）
func ToStartOfWeek(t time.Time) time.Time {
	if t.IsZero() {
		return t
	}

	return time.Date(t.Year(), t.Month(), t.Day()-int(t.Weekday())+1, 0, 0, 0, 0, t.Location())
}

func ToStartOfMonth(t time.Time) time.Time {
	if t.IsZero() {
		return t
	}

	return time.Date(t.Year(), t.Month(), 1, 0, 0, 0, 0, t.Location())
}

func ToStartOfYear(t time.Time) time.Time {
	if t.IsZero() {
		return t
	}

	return time.Date(t.Year(), 1, 1, 0, 0, 0, 0, t.Location())
}
