package feed

import (
	"bytes"
	"fmt"
	"io"
	"log"

	"repo.local/wputil"
)

func Run(conf io.Reader) error {
	c, err := ReadConfig(conf)
	if err != nil {
		return err
	}

	for _, v := range c {
		for _, f := range v.Feeds {
			b, err := f.fetch()
			if err != nil {
				return fmt.Errorf("fetching feed: %s", err)
			}

			err = writeRawXML(b, v.RSSDir, f.Name)
			if err != nil {
				return fmt.Errorf("xml write: %s", err)
			}

			feed, err := wputil.ReadWPXML(bytes.NewReader(b))
			if err != nil {
				return fmt.Errorf("json load: %s", err)
			}

			n, err := writeJSON(feed, v.JSONDir, f.Name)
			if err != nil {
				return fmt.Errorf("json write: %s", err)
			}
			log.Printf("[%d/%d] %s", n, len(feed.List()), f.URL)
		}
	}
	return nil
}
