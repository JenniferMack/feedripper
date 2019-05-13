package feed

import (
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

func Run(conf string) error {
	// open and read conf file
	f, err := os.Open(conf)
	defer f.Close()
	if err != nil {
		return err
	}
	c, err := readConfig(f)
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
		}
	}
	// open old json
	// convert raw to json
	// merge old and new
	// write merged
	return nil
}
