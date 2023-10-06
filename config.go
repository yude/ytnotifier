package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/BurntSushi/toml"
)

type (
	Config struct {
		Channels      []string
		ChannelsUrl   string   `toml:"channels_url"`
		IncludeFilter []string `toml:"include_filter"`
		ExcludeFilter []string `toml:"exclude_filter"`
		Mastodon      Mastodon
	}

	Mastodon struct {
		Domain string
		Token  string
	}

	Channel struct {
		Url         string
		Name        string
		LastChecked time.Time
		Events      map[string]Event
	}

	Event struct {
		Title    string
		Url      string
		StartsAt time.Time
	}
)

var cfg Config
var channels []Channel

func InitConfig() {
	if _, err := toml.DecodeFile("config.toml", &cfg); err != nil {
		log.Fatal(err)
	}

	if cfg.ChannelsUrl == "" {
		GetChannelsFromConfig()
	} else {
		GetChannelsFromUrl()
	}

	log.Println("Loaded " + fmt.Sprint(len(channels)) + " channels.")
}

func GetChannelsFromConfig() {
	for _, v := range cfg.Channels {
		f := ParseRss(v)
		if f == nil {
			log.Printf("RSS feed %v no longer exist, skipping.\n", v)
			continue
		}
		new := Channel{
			Url:         v,
			Name:        f.Author.Name,
			LastChecked: time.Now(),
			Events:      make(map[string]Event),
		}
		channels = append(channels, new)
	}
}

func GetChannelsFromUrl() {
	resp, err := http.Get(cfg.ChannelsUrl)
	if err != nil {
		log.Fatalf("Failed to get channel list source: %v\n", err)
	}
	defer resp.Body.Close()

	b, _ := io.ReadAll(resp.Body)
	c := strings.Split(string(b), "\n")

	for _, v := range c {
		if !(strings.HasPrefix(v, "http")) {
			continue
		}
		f := ParseRss(v)
		if f == nil {
			log.Printf("RSS feed %v no longer exist, skipping.\n", v)
			continue
		}
		new := Channel{
			Url:         v,
			Name:        f.Author.Name,
			LastChecked: time.Now(),
			Events:      make(map[string]Event),
		}
		channels = append(channels, new)
	}
}
