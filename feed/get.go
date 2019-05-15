package feed

import (
	"bytes"
	"fmt"
	"io"
	"log"

	"repo.local/wputil"
)

// Get is a convieniece function that saves feeds according to the
// provided config file.
func Get(conf io.Reader) error {
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

			n, s, err := writeJSON(feed, v.JSONDir, f.Name)
			if err != nil {
				return fmt.Errorf("json write: %s", err)
			}
			log.SetFlags(0)
			log.Printf("Feed: [%s/%d] %s -> %s/%s.json", size(n), s, f.URL, v.JSONDir, f.Name)
		}
	}
	return nil
}

func size(b int) string {
	const unit = 1024
	if b < unit {
		return fmt.Sprintf("%dB", b)
	}

	div, exp := int64(unit), 0
	for n := b / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}

	return fmt.Sprintf("%.1f%c", float64(b)/float64(div), "KMG"[exp])
}
