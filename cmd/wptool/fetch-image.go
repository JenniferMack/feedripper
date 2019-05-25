package main

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"repo.local/wputil/wpfeed"
	"repo.local/wputil/wpimage"
)

func fetchImages(c io.Reader) error {
	log.SetFlags(0)
	log.SetPrefix("[  images] ")

	conf, err := wpfeed.ReadConfig(c)
	if err != nil {
		return fmt.Errorf("reading config: %s", err)
	}

	for _, v := range conf {
		base := v.Name + "-" + v.Number
		if v.IsWorkDir(os.Getwd()) {
			v.WorkDir = "."
		}

		path := filepath.Join(v.WorkDir, base+"-image.json")
		b, err := ioutil.ReadFile(path)
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

		log.Printf("%d images loaded from %s", le, path)
		log.Printf("%d/%d images already saved", sa, va)
		log.Printf("%d images to download", va-sa)

		if va-sa == 0 {
			log.Print("> nothing to do, exiting")
			return nil
		}

		imgDir := "images"
		dlpath := filepath.Join(v.WorkDir, imgDir)
		num, err := img.FetchImages(dlpath)
		if err != nil {
			return err
		}
		log.Printf("> %d images downloaded, %d errors", num, va-sa-num)

		buf := bytes.Buffer{}
		err = img.Marshal(&buf)
		if err != nil {
			return err
		}
		err = ioutil.WriteFile(path, buf.Bytes(), 0644)
		if err != nil {
			return err
		}
	}
	return nil
}
