package main

import (
	"errors"
	"io"
	"log"
	"net/http"
	"net/url"
)

func PostToMastodon(toot string) error {
	val := url.Values{}
	val.Set("status", toot)
	val.Set("visibility", "unlisted")
	val.Set("language", "ja")
	val.Set("sensitive", "true")

	res, err := http.PostForm(cfg.Mastodon.Domain+"/api/v1/statuses?access_token="+cfg.Mastodon.Token, val)
	if err != nil {
		return err
	}

	defer res.Body.Close()

	if res.StatusCode != http.StatusCreated {
		_, err := io.ReadAll(res.Body)
		if err != nil {
			return errors.New("failed to retrieve error")
		}
		return nil
	}

	return nil
}

func CheckMastodonCredentials() {
	if cfg.Mastodon.Domain == "" {
		log.Fatal("The domain of Mastodon instance is not specified.")
	}

	req, _ := http.NewRequest(
		"GET",
		cfg.Mastodon.Domain+"/api/v1/accounts/verify_credentials",
		nil,
	)
	req.Header.Set("Authorization", "Bearer "+cfg.Mastodon.Token)
	client := new(http.Client)

	res, err := client.Do(req)

	if err != nil || res.StatusCode == 401 {
		if err != nil {
			log.Println(err)
		}
		log.Fatal("Failed to log in to Mastodon instance.")
	}
}
