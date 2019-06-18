package feedpub

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"
	"path/filepath"
	"sort"
	"strings"
	"sync"
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

			for k, v := range i.itm.Images {
				i.imgTot++
				if v.OnDisk == true {
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
				i.itm.Images[k].OnDisk = true
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

func ExtractImages(conf Config, pp bool, lg *log.Logger) error {
	lg.SetPrefix("[images  ] ")
	itms, _ := readItems(conf.Names("json"))
	cnt, ondi := 0, 0

	for k, v := range itms {
		u := []string{}
		str, err := Parse(strings.NewReader(v.Body),
			ConvertElemIf("iframe", "img", "src", "youtube.com"),
			ExtractAttr("img", "src", &u),
		)
		if err != nil {
			return fmt.Errorf("html parse: %s", err)
		}

		itms[k].Body = str
		it := []feedimage{}

		for _, i := range u {
			if strings.Contains(i, "?") {
				continue
			}
			lp := makeLocPath(conf.Names("dir-images"), i)
			od := isOnDisk(lp)
			if od {
				ondi++
			}

			it = append(it, feedimage{
				URL:       i,
				LocalPath: lp,
				OnDisk:    od,
			})
			cnt++
		}
		itms[k].Images = it
	}

	lg.Printf("[%03d/%03d] images => %s", ondi, cnt, conf.Names("json"))

	_, err := writeJSON(itms, conf.Names("json"), pp)
	if err != nil {
		return fmt.Errorf("json write: %s", err)
	}
	return nil
}

func makeLocPath(d, p string) string {
	pth := path.Base(p)
	ext := path.Ext(pth)
	pth = strings.TrimSuffix(pth, ext)
	pth = filepath.Join(d, pth+".jpg")
	return pth
}

func isOnDisk(p string) bool {
	if _, err := os.Stat(p); !os.IsNotExist(err) {
		return true
	}
	return false
}
