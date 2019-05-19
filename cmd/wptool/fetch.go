package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"path/filepath"
	"sync"
	"time"

	"repo.local/wputil"
	"repo.local/wputil/wpfeed"
)

type comm struct {
	err error
	msg string
}

func getFeeds(conf io.Reader, pretty bool) error {
	c, err := wpfeed.ReadConfig(conf)
	if err != nil {
		return fmt.Errorf("reading config: %s", err)
	}

	log.SetFlags(0)
	log.SetPrefix("[fetching] ")
	clock := time.Now()

	for _, v := range c {
		log.Printf("> Fetching %s #%s...", v.Name, v.Number)
		wg := sync.WaitGroup{}

		commChan := make(chan comm)

		for _, f := range v.Feeds {
			wg.Add(1)
			go func(fd wpfeed.Feed) {
				defer wg.Done()
				commChan <- fetch(fd, pretty, v.RSSDir, v.JSONDir)
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
		log.Printf("> Processed %d feeds in %s", len(v.Feeds), time.Since(clock))
		if errflag > 0 {
			return fmt.Errorf("%d error(s) occured, check the log", errflag)
		}
	}
	return nil
}

func fetch(f wpfeed.Feed, pretty bool, xDir, jDir string) comm {
	b, err := f.FetchURL()
	if err != nil {
		return comm{err: fmt.Errorf("%s", err)}
	}

	err = wpfeed.WriteRawXML(b, xDir, f.Name)
	if err != nil {
		return comm{err: fmt.Errorf("xml write: %s", err)}
	}

	fd, err := wputil.ReadWPXML(bytes.NewReader(b))
	if err != nil {
		return comm{err: fmt.Errorf("json load: %s", err)}
	}

	path := filepath.Join(jDir, f.Name+".json")
	n, l, err := mergeAndSave(fd, pretty, path)
	if err != nil {
		return comm{err: fmt.Errorf("merge and save: %s", err)}
	}

	return comm{msg: fmt.Sprintf("[%s/%d] %s -> %s", size(n), l, f.URL, path)}
}
