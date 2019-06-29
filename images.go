package feedripper

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"strings"
	"sync"

	"golang.org/x/net/html"
)

type comm struct {
	post item
	errs []error
}

func FetchImages(conf Config, loud bool, lg *log.Logger) error {
	lg.SetPrefix("[images  ] ")
	ch := make(chan comm)
	token := make(chan struct{}, 5)

	wg := sync.WaitGroup{}
	posts, _ := readItems(conf.Names("json"))
	lg.Printf("downloading <= %s", conf.Names("json"))

	for _, v := range posts {
		wg.Add(1)
		go func(p item) {
			defer wg.Done()
			defer func() { <-token }()
			token <- struct{}{}
			fetchPostImgs(p, ch, lg, conf, loud)
		}(v)
	}
	go func() { wg.Wait(); close(ch) }()

	imgCnt, imgTot, errCnt := 0, 0, 0
	postList := items{}
	for v := range ch {
		postList.add(v.post)
		imgTot += len(v.post.Images)
		imgCnt += len(v.post.Images) - len(v.errs)
		for _, err := range v.errs {
			if loud {
				lg.Printf("[error] %s", err)
			}
			errCnt++
		}
	}
	lg.Printf("[%03d/%03d] images downloaded, %d errors", imgCnt, imgTot, errCnt)
	if _, err := writeJSON(postList, conf.Names("json"), true); err != nil {
		return fmt.Errorf("update: %s", err)
	}
	return nil
}

func fetchPostImgs(i item, ch chan comm, lg *log.Logger, conf Config, loud bool) {
	c := comm{
		post: i,
	}
	defer func() { ch <- c }()

	for k, v := range i.Images {
		if isOnDisk(v.LocalPath) {
			if loud {
				lg.Printf("[%6s] %.79s", "skip", v.LocalPath)
			}
			continue
		}

		imgbyts, err := FetchItem(v.URL, "image")
		if err != nil {
			c.errs = append(c.errs, fmt.Errorf("%s: %.80s", err, v.URL))
			i.Images[k].LocalPath = conf.Names("image-404")
			continue
		}

		jpgbyts, err := MakeJPEG(imgbyts, conf.ImageQual, conf.ImageWidth)
		if err != nil {
			c.errs = append(c.errs, err)
			continue
		}

		err = ioutil.WriteFile(v.LocalPath, jpgbyts, 0644)
		if err != nil {
			c.errs = append(c.errs, err)
			continue
		}
		lg.Printf("[% 6s] %.80s", sizeOf(len(jpgbyts)), v.LocalPath)
	}
}

func ExtractImages(conf Config, pp bool, lg *log.Logger, fn ...func(*html.Node)) error {
	lg.SetPrefix("[images  ] ")
	itms, _ := readItems(conf.Names("json"))
	cnt := 0

	for k, v := range itms {
		u := []string{}
		str, err := Parse(strings.NewReader(v.Body),
			// fix relative links while here
			func(n *html.Node) {
				if n.Type == html.ElementNode && n.Data == "a" {
					for k, v := range n.Attr {
						if v.Key == "href" {
							u, err := url.Parse(v.Val)
							if err != nil {
								return
							}

							if u.Host == "" {
								u.Host = conf.SiteURL
							}
							if u.Scheme == "" {
								u.Scheme = "http"
							}
							if conf.UseTLS {
								u.Scheme = "https"
							}
							n.Attr[k].Val = u.String()
						}
					}
				}
			},
			ConvertElemIf("iframe", "img", "src", "youtube.com"),
			ExtractAttr("img", "src", &u),
		)
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
	u, _ := url.PathUnescape(p)
	pth := path.Base(u)

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
