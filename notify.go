package main

import (
	"log"
	"time"

	"github.com/philippgille/gokv"
)

var AnnounceNewEventsLock bool

func AnnounceNewEvents(store gokv.Store) {
	if AnnounceNewEventsLock {
		return
	}

	AnnounceNewEventsLock = true

	for _, cv := range channels {
		for ek, ev := range cv.Events {
			kvItem := new(string)
			found, err := store.Get(ev.Url, kvItem)

			// If there is desync or any error, remove this event from RAM
			if !found || err != nil {
				delete(cv.Events, ek)
				continue
			}

			// If this is already announced, skip this item
			if *kvItem == "true" {
				continue
			}

			// Announce this
			msg := `🆕 配信予定\n` + cv.Name + ": " + ev.Title + "<br />🔗 " + ev.Url + "<br />⏰ " + FormatDateTime(ev.StartsAt) + " 開始"
			PostToMastodon(msg)

			// Mark this as already announced
			err = store.Set(ev.Url, "true")
			if err != nil {
				log.Println(err)
			}
		}
	}

	AnnounceNewEventsLock = false
}

func AnnounceStarts(store gokv.Store) {
	for _, cv := range channels {
		for ek, ev := range cv.Events {
			// Check approx. start time in order to reduce request frequency
			if !(time.Now().After(ev.StartsAt)) {
				continue
			}
			// Check date & time this starts
			d := GetVideoDetails(ev.Url)
			if !IsStarted(d) {
				continue
			}

			// Announce this
			msg := `⏺️ 配信開始\n` + cv.Name + ": " + ev.Title + "\n🔗 " + ev.Url
			PostToMastodon(msg)

			// Delete this event from queue (on RAM & KV)
			delete(cv.Events, ek)
			store.Delete(ev.Url)
		}
	}
}
