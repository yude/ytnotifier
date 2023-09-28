package main

import (
	"errors"
	"io"
	"log"
	"net/http"
	"net/url"
	"time"
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

func InitMastodon() {
	CheckMastodonCredentials()
	SetupMastodonProfile()
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

func SetupMastodonProfile() {
	val := url.Values{}
	val.Set("bot", "true")
	val.Set("source[language]", "ja")

	res, err := http.PostForm(cfg.Mastodon.Domain+"/api/v1/accounts/update_credentials?access_token="+cfg.Mastodon.Token, val)
	if err != nil {
		log.Println("Failed to set proper profile to my Mastodon account. Ignoring.")
	}

	defer res.Body.Close()
}

func SetLastUpdatedToMastodon(t time.Time) {
	val := url.Values{}
	val.Set("fields_attributes[3][🔃 最終更新]", FormatDateTime(t))

	res, err := http.PostForm(cfg.Mastodon.Domain+"/api/v1/accounts/update_credentials?access_token="+cfg.Mastodon.Token, val)
	if err != nil {
		log.Println("Failed to update last update timestamp. Ignoring.")
	}

	defer res.Body.Close()
}
