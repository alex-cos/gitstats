package main

import (
	"strings"
	"time"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

func truncateToDay(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
}

func Capitalize(s string) string {
	caser := cases.Title(language.AmericanEnglish)
	return caser.String(strings.ToLower(strings.Trim(strings.ReplaceAll(s, "_", " "), " ")))
}
