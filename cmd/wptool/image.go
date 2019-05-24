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

func makeImageList(c io.Reader) error {
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

		path := filepath.Join(v.WorkDir, base+".html")
		b, err := ioutil.ReadFile(path)
		if err != nil {
			return err
		}

		img, err := wpimage.ParseHTML(bytes.NewReader(b))
		if err != nil {
			return err
		}
		log.Printf("%d images loaded from %s", len(img), path)

		path = filepath.Join(v.WorkDir, base+"-image.json")
		old := wpimage.ImageList{}
		b, err = ioutil.ReadFile(path)
		if err != nil {
			log.Printf("%s does not exist, skipping", path)
			b = []byte("[]")
		}

		err = old.Unmarshal(bytes.NewBuffer(b))
		if err != nil {
			return err
		}

		log.Printf("%d images loaded from %s", len(old), path)
		img.Merge(old)
		log.Printf("%d unique images to check", len(img))

		nimg, nerr := 0, 0
		for k := range img {
			e := img[k].ParseImageURL(img[k].Rawpath)
			if e != nil {
				log.Printf("[parse error] %s", e)
			}
			i, e := img[k].CheckImageStatus()
			nimg += i
			if e != nil {
				nerr += 1
				log.Printf("[URL error] %s", e)
			}
		}

		buf := bytes.Buffer{}
		err = img.Marshal(&buf)
		if err != nil {
			return err
		}

		log.Printf("> [%s %s] %d checked, %d skipped, %d errors", size(buf.Len()),
			path, nimg+nerr, len(img)-nimg, nerr)

		err = ioutil.WriteFile(path, buf.Bytes(), 0644)
		if err != nil {
			return err
		}
	}
	return nil
}
