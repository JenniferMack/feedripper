package main

import (
	"bytes"
	"fmt"
	"io"
	"log"

	"repo.local/wputil"
	"repo.local/wputil/feed"
)

func getFeeds(conf io.Reader) error {
	c, err := feed.ReadConfig(conf)
	if err != nil {
		return err
	}

	for _, v := range c {
		for _, f := range v.Feeds {
			b, err := f.FetchURL()
			if err != nil {
				return fmt.Errorf("fetching feed: %s", err)
			}

			err = feed.WriteRawXML(b, v.RSSDir, f.Name)
			if err != nil {
				return fmt.Errorf("xml write: %s", err)
			}

			fd, err := wputil.ReadWPXML(bytes.NewReader(b))
			if err != nil {
				return fmt.Errorf("json load: %s", err)
			}

			n, s, err := feed.WriteJSON(fd, v.JSONDir, f.Name)
			if err != nil {
				return fmt.Errorf("json write: %s", err)
			}
			log.SetFlags(0)
			log.SetPrefix("[fetching] ")
			log.Printf("[%s/%d] %s -> %s/%s.json", size(n), s, f.URL, v.JSONDir, f.Name)
		}
	}
	return nil
}
