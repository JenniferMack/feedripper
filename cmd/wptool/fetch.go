package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"path/filepath"
	"sync"
	"time"

	"gitlab.com/dsse/wputils/wpfeed"
	"repo.local/wputil"
)

type comm struct {
	err error
	msg string
}

func getFeeds(conf io.Reader, pretty bool) error {
	log.SetFlags(0)
	log.SetPrefix("[fetching] ")
	clock := time.Now()

	c, err := wputil.NewConfigList(conf)
	if err != nil {
		return fmt.Errorf("reading config: %s", err)
	}

	for _, v := range c {
		log.Printf("> %s, fetching %s #%s...", time.Now().Format("Jan 02, 15:04"), v.Name, v.Number)
		wg := sync.WaitGroup{}

		commChan := make(chan comm)

		for _, f := range v.Feeds {
			wg.Add(1)
			go func(fd wputil.Feed) {
				defer wg.Done()
				commChan <- fetch(fd, pretty, v)
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

func fetch(f wputil.Feed, pretty bool, c wputil.Config) comm {
	b, err := wputil.FetchFeed(f)
	if err != nil {
		return comm{err: fmt.Errorf("%s", err)}
	}

	err = wpfeed.WriteRawXML(b, c.RSSDir, c.Paths("name"))
	if err != nil {
		return comm{err: fmt.Errorf("xml write: %s", err)}
	}

	fd, err := wputil.ReadWPXML(bytes.NewReader(b))
	if err != nil {
		return comm{err: fmt.Errorf("json load: %s", err)}
	}

	path := filepath.Join(c.JSONDir, c.Paths("json"))
	n, l, err := mergeAndSave(fd, pretty, path)
	if err != nil {
		return comm{err: fmt.Errorf("merge and save: %s", err)}
	}

	return comm{msg: fmt.Sprintf("[%s/%d] %s -> %s", size(n), l, f.URL, path)}
}
