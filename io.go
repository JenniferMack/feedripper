package feedpub

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"sync"
	"time"
)

func ReadConfig(file string) (*Config, error) {
	b, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}

	c := Config{}
	err = json.Unmarshal(b, &c)
	if err != nil {
		return nil, err
	}

	if c.Tags.priOutOfRange() {
		return nil, fmt.Errorf("tag priority out of range")
	}
	return &c, nil
}

func FetchFeeds(conf Config, l *log.Logger) error {
	l.SetPrefix("[fetching] ")
	wg := sync.WaitGroup{}
	errCh := make(chan error)
	clock := time.Now()

	for _, v := range conf.Feeds {
		wg.Add(1)
		go func(v feed) {
			defer wg.Done()
			l.Printf("feed: %s", v.URL)
			b, err := FetchItem(v.URL, "xml")
			if err != nil {
				errCh <- fmt.Errorf("feed %s: %s", v.Name, err)
				return
			}

			// xml
			name := fmt.Sprintf("%s_%d", v.Name, time.Now().Unix())
			loc := conf.feedPath(name, "xml")
			err = ioutil.WriteFile(loc, b, 0644)
			if err != nil {
				errCh <- fmt.Errorf("write xml: %s", err)
				return
			}
			l.Printf("save: %s", loc)

			// json
			x := rss{}
			err = xml.Unmarshal(b, &x)
			if err != nil {
				errCh <- fmt.Errorf("decode xml: %s", err)
				return
			}

			// merge json
			loc = conf.feedPath(v.Name, "json")
			oi, n := oldItems(loc)
			l.Printf("read: [%03d/%s] items <= %s", len(oi), sizeOf(n), loc)
			oi.add(x.Channel.Items...)

			n, err = writeJSON(oi, loc, false)
			if err != nil {
				errCh <- fmt.Errorf("write json: %s", err)
				return
			}
			l.Printf("save: [%03d/%s] items => %s", len(oi), sizeOf(n), loc)
		}(v)
	}
	go func() {
		wg.Wait()
		close(errCh)
	}()

	errcnt := 0
	for e := range errCh {
		errcnt++
		l.Printf("[error] %.80s", e)
	}

	l.Printf("[%03d] feeds fetched in %s, %d errors", len(conf.Feeds),
		time.Since(clock).Round(time.Millisecond), errcnt)

	if errcnt > 0 {
		plural := "s"
		if errcnt == 1 {
			plural = ""
		}
		return fmt.Errorf("%d error%s, check the log", errcnt, plural)
	}
	return nil
}

func sizeOf(b int) string {
	const unit = 1000
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

func mergeFeeds(conf Config, lg *log.Logger) items {
	feed := items{}
	n := 0
	for _, v := range conf.Feeds {
		path := conf.feedPath(v.Name, "json")
		oi, _ := oldItems(path)
		n += len(oi)
		lg.Printf("[%03d/%03d] total / items from %s", n, len(oi), path)
		feed.add(oi...)
	}
	return feed
}

func writeJSON(obj interface{}, path string, pretty bool) (int, error) {
	buf := bytes.Buffer{}
	enc := json.NewEncoder(&buf)
	if pretty {
		enc.SetIndent("", "  ")
		enc.SetEscapeHTML(false)
	}

	err := enc.Encode(obj)
	if err != nil {
		return 0, err
	}

	err = ioutil.WriteFile(path, buf.Bytes(), 0644)
	if err != nil {
		return 0, err
	}
	return buf.Len(), nil
}

func oldItems(p string) (items, int) {
	b, err := ioutil.ReadFile(p)
	if err != nil {
		return items{}, 0
	}

	it := items{}
	err = json.Unmarshal(b, &it)
	if err != nil {
		return items{}, 0
	}
	return it, len(b)
}

func FetchItem(url, typ string) ([]byte, error) {
	resp, err := http.Head(url)
	if err != nil {
		resp.Body.Close()
		return nil, fmt.Errorf("http head: %s", err)
	}

	if resp.StatusCode >= 400 {
		resp.Body.Close()
		return nil, fmt.Errorf("status: %d", resp.StatusCode)
	}

	ht := resp.Header.Get("Content-Type")
	if !strings.Contains(ht, typ) {
		resp.Body.Close()
		return nil, fmt.Errorf("content-type: %s", ht)
	}

	resp, err = http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("http get: %s", err)
	}
	defer resp.Body.Close()

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("body read: %s", err)
	}
	return b, nil
}
