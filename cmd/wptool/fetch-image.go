package main

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"log"

	"repo.local/wputil"
	"repo.local/wputil/wpimage"
)

func fetchImages(c io.Reader) error {
	log.SetFlags(0)
	log.SetPrefix("[  images] ")

	conf, err := wputil.NewConfigList(c)
	if err != nil {
		return fmt.Errorf("reading config: %s", err)
	}

	for _, v := range conf {
		b, err := ioutil.ReadFile(v.Paths("image-json"))
		if err != nil {
			return err
		}

		img := wpimage.ImageList{}
		err = img.Unmarshal(bytes.NewReader(b))
		if err != nil {
			return err
		}
		le := len(img)
		sa := img.SavedNum()
		va := img.ValidNum()

		log.Printf("%d images loaded from %s", le, v.Paths("image-json"))
		log.Printf("%d/%d images already saved", sa, va)
		log.Printf("%d images to download", va-sa)

		if va-sa == 0 {
			log.Print("> nothing to do, exiting")
			return nil
		}

		num, err := img.FetchImages(v.Paths("image-dir"))
		if err != nil {
			return err
		}
		log.Printf("> %d images downloaded, %d errors", num, va-sa-num)

		buf := bytes.Buffer{}
		err = img.Marshal(&buf)
		if err != nil {
			return err
		}
		err = ioutil.WriteFile(v.Paths("image-json"), buf.Bytes(), 0644)
		if err != nil {
			return err
		}
	}
	return nil
}
