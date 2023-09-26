package main

import (
	"errors"
	"io"
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
