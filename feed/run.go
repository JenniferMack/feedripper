package feed

import (
	"io"
	"io/ioutil"
	"log"
	"path/filepath"
	"strconv"
	"time"
)

func Run(conf io.Reader) error {
	c, err := readConfig(conf)
	if err != nil {
		return err
	}
	// get feed data
	for _, v := range c {
		for _, f := range v.Feeds {
			b, err := f.fetch()
			if err != nil {
				log.Printf("fetching feed: %s", err)
				continue
			}
			stamp := "_" + strconv.FormatInt(time.Now().Unix(), 10)
			err = ioutil.WriteFile(filepath.Join(v.RSSDir, f.Name+stamp+f.Type), b, 0644)
			if err != nil {
				log.Printf("writing feed: %s", err)
				continue
			}
			err = mergeJSON(b, v.JSONDir, f.Name)
			if err != nil {
				log.Printf("merging: %s", err)
				continue
			}
		}
	}
	return nil
}
