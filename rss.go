package main

import (
	"log"

	"github.com/mmcdole/gofeed"
)

func ParseRss(url string) *gofeed.Feed {
	fp := gofeed.NewParser()
	feed, err := fp.ParseURL(url)
	if err != nil {
		log.Println(err)
	}

	return feed
}
