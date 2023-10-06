package main

import (
	"log"

	"github.com/mmcdole/gofeed"
)

func ParseRss(url string) *gofeed.Feed {
	fp := gofeed.NewParser()
	feed, err := fp.ParseURL(url)
	if err != nil {
		log.Printf("Failed to parse RSS feed for %v: %v\n", url, err)
		return nil
	}

	return feed
}
