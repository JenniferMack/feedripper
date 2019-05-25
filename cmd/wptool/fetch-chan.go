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

func doAction(a string, c wpfeed.Config) error {
	paths, err := c.Paths(os.Getwd())
	if err != nil {
		return err
	}
	list, err := readExistingFile(paths["images"])
	if err != nil {
		return err
	}
	// load early incase of changes
	list.SetDefaults(c.ImageQual, c.ImageWidth, c.UseTLS)

	var out []byte
	var e error
	wr := os.Stderr
	switch a {
	case "parse":
		out, e = doParse(list, paths["images"], paths["html"])

	case "filter":
		out, e = doFilter(list, paths["images"], c.SiteURL)

	case "verify":
		out, e = doVerify(list, paths["images"])

	case "fetch":
		out, e = doFetch(list, paths["images"], paths["imageDir"], wr)
	}

	if e != nil {
		return err
	}
	return ioutil.WriteFile(paths["images"], out, 0644)
}

func doFetch(list wpimage.ImageList, path, dir string, wr io.Writer) ([]byte, error) {
	log.Printf("> fetching images [%s]", path)
	log.SetOutput(wr)

	type carrier struct {
		item  wpimage.ImageData
		image []byte
		err   error
	}
	wg := sync.WaitGroup{}
	ch := make(chan carrier, 10)

	for _, v := range list {
		wg.Add(1)
		go func(i wpimage.ImageData) {
			defer wg.Done()

			b, err := i.FetchImage(dir)
			out := carrier{item: i, image: b, err: err}
			ch <- out
		}(v)
	}
	go func() {
		wg.Wait()
		close(ch)
	}()

	o := []wpimage.ImageData{}
	var got, errs int
	for v := range ch {
		v.item.Saved = true

		if v.err != nil {
			log.Printf("[error] fetching: %s", v.err)
			v.item.Saved = false
			errs++
		}
		if v.image != nil {
			got++
		}

		err := saveImage(v.image, v.item.LocalPath, v.item.ImgWidth, v.item.ImgQual)
		if err != nil {
			log.Printf("[error] saving: %s", err)
			v.item.Saved = false
			errs++
		}

		o = append(o, v.item)
	}
	out := wpimage.ImageList(o)

	log.Printf("%d/%d downloaded, %d skipped, %d errors", got, len(out), len(out)-got, errs)

	buf := bytes.Buffer{}
	if err := out.Marshal(&buf); err != nil {
		return nil, err
	}

	log.Printf("> [%s/%d/%d] %s", size(buf.Len()), len(out), out.SavedNum(), path)
	return buf.Bytes(), nil
}

func doVerify(list wpimage.ImageList, path string) ([]byte, error) {
	log.Printf("> verifying image URLs [%s]", path)

	type carrier struct {
		err  error
		item wpimage.ImageData
	}
	wg := sync.WaitGroup{}
	ch := make(chan carrier, 10)

	for _, v := range list {
		wg.Add(1)
		go func(i wpimage.ImageData) {
			defer wg.Done()

			_, err := i.CheckImageStatus()
			re := carrier{item: i, err: err}
			ch <- re
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
			n++
		}
		o = append(o, v.item)
	}
	out := wpimage.ImageList(o)

	log.Printf("%d images checked, %d found, %d not found, %d errors",
		len(out), out.ValidNum(), len(out)-out.ValidNum(), n)

	buf := bytes.Buffer{}
	if err := out.Marshal(&buf); err != nil {
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
	if err := merged.Marshal(&buf); err != nil {
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
		n += list[k].ParseImageURL(u)
	}
	log.Printf("%d vaild URLs, %d errors", n, len(list)-n)

	buf := bytes.Buffer{}
	if err := list.Marshal(&buf); err != nil {
		return nil, err
	}

	log.Printf("> [%s/%d] %s", size(buf.Len()), len(list), path)
	return buf.Bytes(), nil
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

func saveImage(in []byte, p string, w uint, q int) error {
	if in == nil {
		return nil
	}
	j, err := wpimage.MakeJPEG(in, q, w)
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(p, j, 0644)
	if err != nil {
		return fmt.Errorf("disk %s: %s", p, err)
	}
	if *flagImageVerbose {
		log.Printf("[%s|%s]", size(len(j)), p)
	}
	return nil
}
