package main

import (
	"strings"
	"time"
)

func truncateToDay(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
}

func Capitalize(s string) string {
	return strings.Title(strings.ToLower(strings.Trim(strings.ReplaceAll(s, "_", " "), " ")))
}
