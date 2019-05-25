package main

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"sync"

	"repo.local/wputil/wpfeed"
	"repo.local/wputil/wpimage"
)

func images(c io.Reader, a []string) error {
	log.SetFlags(0)
	log.SetPrefix("[  images] ")

	conf, err := wpfeed.ReadConfig(c)
	if err != nil {
		return fmt.Errorf("reading config: %s", err)
	}

	for _, confItem := range conf {
		for _, action := range a {
			err := doAction(action, confItem)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func readExistingFile(p string) (wpimage.ImageList, error) {
	b, err := ioutil.ReadFile(p)
	if err != nil {
		b = []byte("[]")
	}

	list := wpimage.ImageList{}
	err = list.Unmarshal(bytes.NewBuffer(b))
	if err != nil {
		return nil, err
	}
	return list, nil
}

func doAction(a string, c wpfeed.Config) error {
	paths, err := c.Paths(os.Getwd())
	if err != nil {
		return err
	}
	list, err := readExistingFile(paths["images"])
	if err != nil {
		return err
	}

	switch a {
	case "parse":
		l, err := doParse(list, paths["images"], paths["html"])
		if err != nil {
			return err
		}
		return ioutil.WriteFile(paths["images"], l, 0644)

	case "filter":
		l, err := doFilter(list, paths["images"], c.SiteURL)
		if err != nil {
			return err
		}
		return ioutil.WriteFile(paths["images"], l, 0644)

	case "verify":
		l, err := doVerify(list, paths["images"])
		if err != nil {
			return err
		}
		return ioutil.WriteFile(paths["images"], l, 0644)

	case "fetch":
		l, err := doFetch(list, paths["images"], paths["imageDir"])
		if err != nil {
			return err
		}
		return ioutil.WriteFile(paths["images"], l, 0644)
	}
	return nil
}

func doFetch(list wpimage.ImageList, path, dir string) ([]byte, error) {
	log.Printf("> fetching images [%s]", path)

	type retChan struct {
		item  wpimage.ImageData
		image []byte
		err   error
	}
	wg := sync.WaitGroup{}
	ch := make(chan retChan)
	token := make(chan struct{}, 10)

	for _, v := range list {
		wg.Add(1)
		go func(i wpimage.ImageData) {
			defer wg.Done()
			token <- struct{}{}

			b, err := i.FetchImage(dir)
			out := retChan{item: i, image: b, err: err}
			ch <- out
			<-token
		}(v)
	}
	go func() {
		wg.Wait()
		close(ch)
	}()

	o := []wpimage.ImageData{}
	for v := range ch {
		o = append(o, v.item)
		if v.err != nil {
			log.Printf("[error] fetching: %s", v.err)
		}

		if v.image == nil {
			continue
		}
		j, err := wpimage.MakeJPEG(v.image, 40, 250)
		if err != nil {
			log.Printf("[error] jpeg %s: %s", v.item.LocalPath, v.err)
			continue
		}
		err = ioutil.WriteFile(v.item.LocalPath, j, 0644)
		log.Printf("[%s|%s] saved", size(len(j)), v.item.LocalPath)
		if err != nil {
			log.Printf("[error] saving %s: %s", v.item.LocalPath, v.err)
		}
	}
	out := wpimage.ImageList(o)

	log.Printf("%d images saved, %d skipped",
		out.SavedNum(), len(out)-out.SavedNum())

	buf := bytes.Buffer{}
	err := out.Marshal(&buf)
	if err != nil {
		return nil, err
	}

	log.Printf("> [%s/%d/%d] %s", size(buf.Len()), len(out), out.SavedNum(), path)
	return buf.Bytes(), nil

	return nil, nil
}

func doVerify(list wpimage.ImageList, path string) ([]byte, error) {
	log.Printf("> verifying image URLs [%s]", path)

	type retChan struct {
		num  int
		err  error
		item wpimage.ImageData
	}
	wg := sync.WaitGroup{}
	ch := make(chan retChan)
	token := make(chan struct{}, 10)

	for _, v := range list {
		wg.Add(1)
		go func(i wpimage.ImageData) {
			defer wg.Done()
			token <- struct{}{}

			n, err := i.CheckImageStatus()
			re := retChan{num: n, err: err, item: i}
			ch <- re
			<-token
		}(v)
	}
	go func() {
		wg.Wait()
		close(ch)
	}()

	n := 0
	o := []wpimage.ImageData{}
	for v := range ch {
		if v.err != nil {
			log.Printf("[error] %s", v.err)
		}
		n += v.num
		o = append(o, v.item)
	}
	out := wpimage.ImageList(o)

	log.Printf("%d images checked, %d found, %d not found, %d errors",
		len(out), out.ValidNum(), len(out)-out.ValidNum(), n)

	buf := bytes.Buffer{}
	err := out.Marshal(&buf)
	if err != nil {
		return nil, err
	}

	log.Printf("> [%s/%d/%d] %s", size(buf.Len()), len(out), out.ValidNum(), path)
	return buf.Bytes(), nil
}

func doParse(list wpimage.ImageList, imgPath, htmlPath string) ([]byte, error) {
	h, err := os.Open(htmlPath)
	if err != nil {
		return nil, err
	}
	defer h.Close()

	imgs, err := wpimage.ParseHTML(h)
	if err != nil {
		return nil, err
	}
	h.Close()

	merged := imgs.Merge(list)
	log.Printf("> parsing HTML [%s]", htmlPath)
	log.Printf("%-3d images found in %s", len(imgs), htmlPath)
	log.Printf("%-3d images found in %s", len(list), imgPath)

	buf := bytes.Buffer{}
	err = merged.Marshal(&buf)
	if err != nil {
		return nil, err
	}

	log.Printf("%-3d unique images recorded", len(merged))
	log.Printf("> [%s/%d] %s", size(buf.Len()), len(merged), imgPath)
	return buf.Bytes(), nil
}

func doFilter(list wpimage.ImageList, path, u string) ([]byte, error) {
	log.Printf("> parsing image URLs [%s]", path)

	n := 0
	for k := range list {
		n += list[k].ParseImageURL(u, false)
	}
	log.Printf("%d vaild URLs, %d errors", n, len(list)-n)

	buf := bytes.Buffer{}
	err := list.Marshal(&buf)
	if err != nil {
		return nil, err
	}
	log.Printf("> [%s/%d] %s", size(buf.Len()), len(list), path)
	return buf.Bytes(), nil
}
