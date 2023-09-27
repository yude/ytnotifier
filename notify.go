package main

import (
	"fmt"
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
			var msg string
			if ev.StartsAt == time.Unix(0, 0) {
				msg = fmt.Sprintf("🆕 配信予定\n%v: %v\n🔗 %v\n⏰ 開始時刻未定", cv.Name, ev.Title, ev.Url)
			} else {
				msg = fmt.Sprintf("🆕 配信予定\n%v: %v\n🔗 %v\n⏰ %v 開始", cv.Name, ev.Title, ev.Url, FormatDateTime(ev.StartsAt))
			}

			err = PostToMastodon(msg)
			if err != nil {
				log.Println(err)
				continue
			}

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
			d, err := GetVideoDetails(ev.Url)
			if err != nil {
				continue
			}
			if !IsStarted(d) {
				continue
			}

			// Announce this
			msg := fmt.Sprintf("⏺️ 配信開始\n%v: %v\n🔗 %v", cv.Name, ev.Title, ev.Url)
			PostToMastodon(msg)

			// Delete this event from queue (on RAM & KV)
			delete(cv.Events, ek)
			store.Delete(ev.Url)
		}
	}
}
