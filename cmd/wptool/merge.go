package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"path/filepath"

	"repo.local/wputil"
	"repo.local/wputil/wpfeed"
)

func mergeFeeds(conf io.Reader, pretty bool) error {
	log.SetFlags(0)
	log.SetPrefix("[ merging] ")

	c, err := wpfeed.ReadConfig(conf)
	if err != nil {
		return fmt.Errorf("reading config: %s", err)
	}

	for _, v := range c {
		log.Printf("> Merging %s #%s...", v.Name, v.Number)
		// holds all feeds
		feeds := wputil.Feed{}
		for _, f := range v.Feeds {
			path := filepath.Join(v.JSONDir, f.Name+".json")
			d, err := ioutil.ReadFile(path)
			if err != nil {
				return fmt.Errorf("reading %s: %s", path, err)
			}
			_, err = feeds.Write(d)
			if err != nil {
				return fmt.Errorf("loading json: %s", err)
			}
			log.Printf("%d posts loaded, dupes removed [%s]", feeds.Len(), path)
		}

		// check deadline
		feeds.Deadline(v.Deadline, v.Days)
		if feeds.Len() == 0 {
			return fmt.Errorf("no posts found within deadline: %s, %d", v.Deadline, v.Days)
		} else {
			log.Printf("%d posts within deadline range", feeds.Len())
		}

		// Include
		tags := []string{}
		for _, t := range v.Tags {
			tags = append(tags, t.Text)
		}
		feeds = feeds.Include(tags)
		log.Printf("%d posts with included tags", feeds.Len())

		// Exclude
		feeds = feeds.Exclude(v.Exclude)
		log.Printf("%d posts after excluding tags", feeds.Len())

		//output
		b, err := formatFeed(feeds, pretty)
		if err != nil {
			return fmt.Errorf("json format: %s", err)
		}

		path := v.Name + ".json"
		err = ioutil.WriteFile(path, b, 0644)
		if err != nil {
			return fmt.Errorf("json write: %s", err)
		}
		log.Printf("> [%s/%d] %s", size(len(b)), feeds.Len(), path)
	}
	return nil
}
