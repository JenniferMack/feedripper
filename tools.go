package feedpub

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

func BuildFeeds(conf Config, pp bool, lg *log.Logger) error {
	lg.SetPrefix("[building] ")
	errCh := make(chan error)
	wg := sync.WaitGroup{}

	for _, fd := range conf.Feeds {
		wg.Add(1)
		go func(ff feed) {
			defer wg.Done()

			path := filepath.Join(conf.Names("dir-rss"), ff.Name)
			gl, err := filepath.Glob(path + `_*.xml`)
			if err != nil {
				errCh <- err
				return
			}

			itms := items{}
			lg.Printf("[%03d] archives <= %s", len(gl), ff.Name)
			for _, v := range gl {
				fi, err := os.Open(v)
				if err != nil {
					errCh <- err
					return
				}
				defer fi.Close()

				rss := rss{}
				err = xml.NewDecoder(fi).Decode(&rss)
				if err != nil {
					errCh <- err
					return
				}

				itms.add(rss.Channel.Items...)
			}

			if len(itms) == 0 {
				errCh <- fmt.Errorf("no items found")
				return
			}

			sort.Sort(itms)
			out := conf.feedPath(ff.Name, "json")
			n, err := writeJSON(itms, out, pp)
			if err != nil {
				errCh <- err
				return
			}

			lg.Printf("[%03d/%s] items => %s", len(itms), sizeOf(n), out)
		}(fd)
	}
	go func() { wg.Wait(); close(errCh) }()

	errcnt := 0
	for v := range errCh {
		lg.SetPrefix("[   error] ")
		lg.Printf("%s", v)
		errcnt++
	}

	if errcnt > 0 {
		return fmt.Errorf("%d errors during recovery", errcnt)
	}
	return nil
}
