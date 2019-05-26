package main

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"sync"

	"repo.local/wputil/wpimage"
)

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

	suffix := "s"
	if errs == 1 {
		suffix = ""
	}
	log.Printf("%d/%d downloaded, %d error%s, %d prev. saved", got, len(out), errs, suffix, len(out)-got)

	buf := bytes.Buffer{}
	if err := out.Marshal(&buf); err != nil {
		return nil, err
	}

	log.Printf("> [%s/%d/%d] %s", size(buf.Len()), len(out), out.SavedNum(), path)
	return buf.Bytes(), nil
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
