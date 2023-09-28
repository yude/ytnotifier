package main

import (
	"fmt"
	"strings"
	"time"
)

func FormatDateTime(t time.Time) string {
	return t.Format("2006/1/2 15:04 MST")
}

func IsMatchesIncludeFilter(title string) bool {
	if len(cfg.IncludeFilter) == 0 {
		return true
	}

	for _, v := range cfg.IncludeFilter {
		if strings.Contains(title, v) {
			return true
		}
	}

	return false
}

func IsMatchesExcludeFilter(title string) bool {
	if len(cfg.ExcludeFilter) == 0 {
		return true
	}

	for _, v := range cfg.ExcludeFilter {
		if strings.Contains(title, v) {
			return true
		}
	}

	return false
}

func PrettyPrintEvent(e Event) {
	fmt.Printf("Title: %v \nURL: %v \nStarts at: %v", e.Title, e.Url, FormatDateTime(e.StartsAt))
}

func RemoveBaseUrlFromYouTubeLink(l string) string {
	return strings.Replace(l, "https://www.youtube.com/watch?v=", "", -1)
}
