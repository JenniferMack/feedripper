package feedpub

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"sort"
	"strings"
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

	for _, v := range conf.Feeds {
		l.Printf("url: %s", v.URL)
		b, err := FetchItem(v.URL, "xml")
		if err != nil {
			return fmt.Errorf("feed %s: %s", v.Name, err)
		}

		// xml
		name := fmt.Sprintf("%s_%d", v.Name, time.Now().Unix())
		loc := conf.feedPath(name, "xml")
		err = ioutil.WriteFile(loc, b, 0644)
		if err != nil {
			return fmt.Errorf("write xml: %s", err)
		}
		l.Printf("wrote: %s", loc)

		// json
		x := rss{}
		err = xml.Unmarshal(b, &x)
		if err != nil {
			return fmt.Errorf("decode xml: %s", err)
		}

		// merge json
		loc = conf.feedPath(v.Name, "json")
		oi := oldItems(loc)
		l.Printf("found %d items in %s", len(oi), loc)
		oi.add(x.Channel.Items)

		b, err = json.Marshal(oi)
		if err != nil {
			return fmt.Errorf("encode json: %s", err)
		}

		err = ioutil.WriteFile(loc, b, 0644)
		if err != nil {
			return fmt.Errorf("write json: %s", err)
		}
		l.Printf("[%d] wrote: %s", len(oi), loc)
	}
	return nil
}

func mergeFeeds(conf Config, lg *log.Logger) items {
	feed := items{}
	n := 0
	for _, v := range conf.Feeds {
		path := conf.feedPath(v.Name, "json")
		oi := oldItems(path)
		n += len(oi)
		lg.Printf("[%03d/%03d] total / items from %s", n, len(oi), path)
		feed.add(oi)
	}
	return feed
}

func WriteItemList(conf Config, lg *log.Logger) error {
	lg.SetPrefix("[merging ] ")
	list := mergeFeeds(conf, lg)

	n := len(list)
	list.trim(conf)
	lg.Printf("[%03d] items outside of date range", n-len(list))
	sort.Sort(list)

	name := conf.Names("json")
	f, err := os.Create(name)
	if err != nil {
		return fmt.Errorf("merge: %s", err)
	}

	enc := json.NewEncoder(f)
	err = enc.Encode(list)
	if err != nil {
		return fmt.Errorf("encode json: %s", err)
	}

	lg.Printf("[%03d] items in %s", len(list), name)
	err = f.Close()
	if err != nil {
		return fmt.Errorf("write json: %s", err)
	}
	return nil
}

func oldItems(p string) items {
	b, err := ioutil.ReadFile(p)
	if err != nil {
		return items{}
	}

	it := items{}
	err = json.Unmarshal(b, &it)
	if err != nil {
		return items{}
	}
	return it
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
