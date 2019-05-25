package main

import (
	"bytes"
	"log"
	"sync"

	"repo.local/wputil/wpimage"
)

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
