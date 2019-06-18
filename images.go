package feedpub

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"sort"
	"strings"
	"sync"

	"golang.org/x/net/html"
)

func FetchImages(conf Config, loud bool, lg *log.Logger) error {
	lg.SetPrefix("[images  ] ")
	type comm struct {
		itm    item
		err    error
		imgNum int
		imgTot int
	}
	commCh := make(chan comm)
	token := make(chan struct{}, 5)

	wg := sync.WaitGroup{}
	itms, _ := readItems(conf.Names("json"))

	for _, v := range itms {
		wg.Add(1)

		go func(i comm) {
			defer func() { commCh <- i; <-token }()
			defer wg.Done()
			token <- struct{}{}

			for _, v := range i.itm.Images {
				i.imgTot++
				if isOnDisk(v.LocalPath) == true {
					if loud {
						lg.Printf("[%6s] %.80s", "skip", v.LocalPath)
					}
					continue
				}

				ib, err := FetchItem(v.URL, "image")
				if err != nil {
					i.err = err
					return
				}

				jb, err := MakeJPEG(ib, conf.ImageQual, conf.ImageWidth)
				if err != nil {
					i.err = err
					return
				}

				err = ioutil.WriteFile(v.LocalPath, jb, 0644)
				if err != nil {
					i.err = err
					return
				}

				if loud {
					lg.Printf("[% 6s] %.80s", sizeOf(len(jb)), v.LocalPath)
				}
				i.imgNum++
			}
		}(comm{itm: v})
	}
	go func() { wg.Wait(); close(commCh) }()

	it := items{}
	imgCnt, imgTot, errCnt := 0, 0, 0
	for v := range commCh {
		it = append(it, v.itm)
		imgCnt += v.imgNum
		imgTot += v.imgTot
		if v.err != nil {
			lg.Printf("[error] %s", v.err)
			errCnt++
		}
	}

	if len(it) != len(itms) {
		return fmt.Errorf("item count mismatch %d/%d", len(it), len(itms))
	}

	sort.Sort(it)
	lg.Printf("[%03d/%03d] images downloaded, %d errors", imgCnt, imgTot, errCnt)
	_, err := writeJSON(it, conf.Names("json"), true)
	if err != nil {
		return fmt.Errorf("json write: %s", err)
	}
	return nil
}

func ExtractImages(conf Config, pp bool, lg *log.Logger, fn ...func(*html.Node)) error {
	lg.SetPrefix("[images  ] ")
	itms, _ := readItems(conf.Names("json"))
	cnt := 0

	for k, v := range itms {
		u := []string{}
		fn = append(fn, ExtractAttr("img", "src", &u))
		str, err := Parse(strings.NewReader(v.Body), fn...)
		if err != nil {
			return fmt.Errorf("html parse: %s", err)
		}

		itms[k].Body = str
		it := []feedimage{}

		for _, i := range u {
			fp, err := parseRawPath(conf, i)
			if err != nil {
				return err
			}
			lp := makeLocPath(conf.Names("dir-images"), fp)

			it = append(it, feedimage{
				URL:       fp,
				LocalPath: lp,
				RawPath:   i,
			})
			cnt++
		}
		itms[k].Images = it
	}

	n, err := writeJSON(itms, conf.Names("json"), pp)
	if err != nil {
		return fmt.Errorf("json write: %s", err)
	}

	lg.Printf("[%03d/%s] images => %s", cnt, sizeOf(n), conf.Names("json"))
	return nil
}

func parseRawPath(conf Config, u string) (string, error) {
	data, err := url.Parse(u)
	if err != nil {
		return "", err
	}

	if data.Host == "" {
		data.Host = conf.SiteURL
	}

	if data.Scheme == "" {
		data.Scheme = "http"
	}

	if conf.UseTLS && data.Scheme == "http" {
		data.Scheme = "https"
	}

	data.RawQuery = ""
	if data.Host == "www.youtube.com" {
		data.Host = "img.youtube.com"
		p := path.Base(data.Path)
		data.Path = "/vi/" + p + "/default.jpg"
	}
	return data.String(), nil
}

func makeLocPath(d, p string) string {
	pth := path.Base(p)
	if strings.Contains(p, "img.youtube.com") {
		pth = path.Base(path.Dir(p)) + "-" + pth
	}
	ext := path.Ext(pth)
	pth = strings.TrimSuffix(pth, ext) + ".jpg"
	return filepath.Join(d, pth)
}

func isOnDisk(p string) bool {
	if _, err := os.Stat(p); !os.IsNotExist(err) {
		return true
	}
	return false
}
