package feedpub

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
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

func FetchFeeds(conf Config, lg *log.Logger) error {
	lg.SetPrefix("[fetching] ")
	wg := sync.WaitGroup{}
	errCh := make(chan error)
	clock := time.Now()

	for _, v := range conf.Feeds {
		wg.Add(1)
		go func(v feed) {
			defer wg.Done()
			lg.Printf("feed: %s", v.URL)
			b, err := FetchItem(v.URL, "xml")
			if err != nil {
				errCh <- fmt.Errorf("feed %s: %s", v.Name, err)
				return
			}

			// xml
			// name := fmt.Sprintf("%s-%s_%d", v.Name, conf.Names("name"), time.Now().Unix())
			tm := strconv.FormatInt(time.Now().Unix(), 10)
			loc := conf.feedPath(v.Name, tm, "xml")
			err = ioutil.WriteFile(loc, b, 0644)
			if err != nil {
				errCh <- fmt.Errorf("write xml: %s", err)
				return
			}
			lg.Printf("save: %s", loc)
		}(v)
	}
	go func() {
		wg.Wait()
		close(errCh)
	}()

	errcnt := 0
	for e := range errCh {
		errcnt++
		lg.SetPrefix("[   error] ")
		lg.Printf("%.85s", e)
	}

	lg.Printf("[%03d] feeds fetched in %s, %d errors", len(conf.Feeds),
		time.Since(clock).Round(time.Millisecond), errcnt)

	if errcnt > 0 {
		return fmt.Errorf("%d errors, check the log", errcnt)
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

func readItems(p string) (items, int) {
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
		return nil, fmt.Errorf("http head: %s", err)
	}
	defer resp.Body.Close()

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

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("body read: %s", err)
	}
	return b, nil
}
