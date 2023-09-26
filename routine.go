package main

import (
	"fmt"
	"log"
	"time"

	"github.com/philippgille/gokv"
)

var CheckNewEventLock bool

func CheckNewEvent(store gokv.Store) {
	if CheckNewEventLock {
		return
	}
	CheckNewEventLock = true

	for _, cv := range channels {
		feed := ParseRss(cv.Url)

		log.Println("Checking " + feed.Author.Name)
		for _, iv := range feed.Items {
			// Already added to event list (on RAM)
			_, exists := cv.Events[iv.Link]
			if exists {
				continue
			}
			// Already added to event list (on KV)
			kvItem := new(string)
			kvFound, err := store.Get(iv.Link, kvItem)
			if err != nil {
				log.Println(err)
			}
			if kvFound {
				continue
			}

			// Does not match to include filter words
			if !(IsMatchesIncludeFilter(iv.Title)) {
				continue
			}

			// Matches to exclude filter words
			if IsMatchesExcludeFilter(iv.Title) {
				continue
			}

			d := GetVideoDetails(iv.Link)

			// Not a live streaming
			if !(IsLiveStream(d)) {
				continue
			}

			// Not a upcoming live
			if !(IsUpcoming(d)) {
				continue
			}

			// Add this event to event list (on RAM)
			new := Event{
				Title:    iv.Title,
				Url:      iv.Link,
				StartsAt: GetScheduledTime(d),
			}
			cv.Events[iv.Link] = new

			// Add this event to event list (on KV)
			err = store.Set(iv.Link, "false")
			if err != nil {
				log.Println(err)
			}
		}

		for _, e := range cv.Events {
			fmt.Println("[New event]")
			PrettyPrintEvent(e)
			fmt.Print("\n")
		}

		cv.LastChecked = time.Now()
	}

	CheckNewEventLock = false
}
