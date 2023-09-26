package main

import (
	"log"

	"github.com/philippgille/gokv/file"
	"github.com/robfig/cron/v3"
)

func main() {
	options := file.DefaultOptions

	kvs, err := file.NewStore(options)
	if err != nil {
		log.Fatal(err)
	}
	defer kvs.Close()

	CheckNewEventLock = false
	AnnounceNewEventsLock = false
	InitConfig()

	c := cron.New()
	c.AddFunc("@every 5s", func() {
		CheckNewEvent(kvs)
		AnnounceNewEvents(kvs)
		AnnounceStarts(kvs)
	})

	log.Println("Started monitoring configured channels.")
	c.Start()

	select {}
}
