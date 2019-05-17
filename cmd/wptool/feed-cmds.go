package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"sync"
	"time"

	"repo.local/wputil"
	"repo.local/wputil/feed"
)

type comm struct {
	err error
	msg string
}

func getFeeds(conf io.Reader, indent bool) error {
	c, err := feed.ReadConfig(conf)
	if err != nil {
		return fmt.Errorf("reading config: %s", err)
	}

	log.SetFlags(0)
	log.SetPrefix("[fetching] ")
	clock := time.Now()

	for _, v := range c {
		log.Printf("> starting %s, number %s...", v.Name, v.Number)
		wg := sync.WaitGroup{}

		commChan := make(chan comm)

		for _, f := range v.Feeds {
			wg.Add(1)
			go func(fd feed.Feed) {
				defer wg.Done()
				commChan <- fetch(fd, v.RSSDir, v.JSONDir)
			}(f)
		}
		go func() {
			wg.Wait()
			close(commChan)
		}()

		errflag := 0
		for v := range commChan {
			if v.err != nil {
				log.SetPrefix("[   error] ")
				log.Print(v.err)
				errflag++
			}
			if v.msg != "" {
				log.SetPrefix("[fetching] ")
				log.Print(v.msg)
			}
		}

		log.SetPrefix("[fetching] ")
		log.Printf("> ...processed %d feeds in %s", len(v.Feeds), time.Since(clock))
		if errflag > 0 {
			return fmt.Errorf("%d error(s) occured, check the log", errflag)
		}
	}
	return nil
}

func fetch(f feed.Feed, xDir, jDir string) comm {
	b, err := f.FetchURL()
	if err != nil {
		return comm{err: fmt.Errorf("%s", err)}
	}

	err = feed.WriteRawXML(b, xDir, f.Name)
	if err != nil {
		return comm{err: fmt.Errorf("xml write: %s", err)}
	}

	fd, err := wputil.ReadWPXML(bytes.NewReader(b))
	if err != nil {
		return comm{err: fmt.Errorf("json load: %s", err)}
	}

	path := filepath.Join(jDir, f.Name+".json")
	n, l, err := mergeAndSave(fd, path)
	if err != nil {
		return comm{err: fmt.Errorf("merge and save: %s", err)}
	}

	return comm{msg: fmt.Sprintf("[%s/%d] %s -> %s", size(n), l, f.URL, path)}
}

func mergeFeeds(conf io.Reader, indent bool) error {
	c, err := feed.ReadConfig(conf)
	if err != nil {
		return fmt.Errorf("reading config: %s", err)
	}

	log.SetFlags(0)
	log.SetPrefix("[ merging] ")

	for _, v := range c {
		log.Printf("starting %s, number %s...", v.Name, v.Number)
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
			log.Printf("%d posts loaded from %s", feeds.Len(), path)
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
		fi, err := os.Create(v.Name + ".json")
		if err != nil {
			return fmt.Errorf("opening file: %s", err)
		}
		// n, err := io.Copy(fi, &feeds)
		// n, err := fi.Write([]byte(feeds.String()))
		enc := json.NewEncoder(fi)
		if indent {
			enc.SetIndent("", "\t")
			enc.SetEscapeHTML(false)
		}
		err = enc.Encode(feeds.List())
		if err != nil {
			return fmt.Errorf("encoding: %s", err)
		}

		stat, err := fi.Stat()
		if err != nil {
			return fmt.Errorf("stat: %s", err)
		}
		n := stat.Size()

		err = fi.Close()
		if err != nil {
			return fmt.Errorf("closing: %s", err)
		}
		log.Printf("[%s/%d] %s", size(int(n)), feeds.Len(), fi.Name())
	}
	return nil
}
