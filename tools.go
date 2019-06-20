package wputil

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"sort"
	"sync"

	"golang.org/x/net/html"
)

func Titles(conf Config, out io.Writer) error {
	b, err := ioutil.ReadFile(conf.Names("json"))
	if err != nil {
		return err
	}

	feed := items{}
	err = json.Unmarshal(b, &feed)
	if err != nil {
		return err
	}

	for k, v := range feed {
		fmt.Fprintf(out, "%02d. [%s] %.59s\n", k+1, v.PubDate.Format("0102|15:04:05"), v.Title)
	}
	return nil
}

func Tags(conf Config, out io.Writer) error {
	f, err := os.Open(conf.Names("html"))
	if err != nil {
		return err
	}
	cnt := 0
	_, err = Parse(f,
		func(n *html.Node) {
			if n.Type == html.ElementNode && n.Data == "h1" {
				fmt.Fprintf(out, "--- %s ---\n", n.FirstChild.Data)
				cnt = 1
			}
			if n.Type == html.ElementNode && n.Data == "h2" {
				fmt.Fprintf(out, "%02d. %.75s\n", cnt, n.FirstChild.Data)
				cnt++
			}
		})
	return err
}

func BuildFeeds(conf Config, pp bool, lg *log.Logger) error {
	lg.SetPrefix("[building] ")
	type commCh struct {
		items items
		path  string
		err   error
	}
	ch := make(chan commCh)
	wg := sync.WaitGroup{}

	for _, fd := range conf.Feeds {
		wg.Add(1)
		go func(ff feed) {
			defer wg.Done()
			ii := commCh{}

			gl, err := filepath.Glob(conf.feedPath(ff.Name, `*`, "xml"))
			if err != nil {
				ii.err = err
				ch <- ii
				return
			}

			lg.Printf("[%03d] archives <= %s", len(gl), ff.Name)
			for _, v := range gl {
				fi, err := os.Open(v)
				if err != nil {
					ii.err = err
					ch <- ii
					return
				}
				defer fi.Close()

				rss := rss{}
				err = xml.NewDecoder(fi).Decode(&rss)
				if err != nil {
					ii.err = err
					ch <- ii
					return
				}

				ii.items.add(rss.Channel.Items...)
			}

			if len(ii.items) == 0 {
				ii.err = fmt.Errorf("no items found")
				ch <- ii
				return
			}

			sort.Sort(ii.items)
			ii.path = conf.feedPath(ff.Name, "", "json")
			ch <- ii
		}(fd)
	}
	go func() { wg.Wait(); close(ch) }()

	errcnt := 0
	for v := range ch {
		if v.err != nil {
			lg.Printf("[error] %s", v.err)
			errcnt++
			continue
		}
		n, err := writeJSON(v.items, v.path, pp)
		if err != nil {
			lg.Printf("[error] %s", v.err)
			errcnt++
			continue
		}
		lg.Printf("[%03d/%s] items => %s", len(v.items), sizeOf(n), v.path)
	}

	if errcnt > 0 {
		return fmt.Errorf("%d errors during building", errcnt)
	}
	return nil
}
