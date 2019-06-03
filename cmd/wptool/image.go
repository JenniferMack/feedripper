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

func makeImageList(c io.Reader) error {
	log.SetFlags(0)
	log.SetPrefix("[  images] ")

	conf, err := wputil.NewConfigList(c)
	if err != nil {
		return fmt.Errorf("reading config: %s", err)
	}

	for _, v := range conf {
		b, err := ioutil.ReadFile(v.Paths("html"))
		if err != nil {
			return err
		}

		img, err := wpimage.ParseHTML(bytes.NewReader(b))
		if err != nil {
			return err
		}
		log.Printf("%d images loaded from %s", len(img), v.Paths("html"))

		old := wpimage.ImageList{}
		b, err = ioutil.ReadFile(v.Paths("image-json"))
		if err != nil {
			log.Printf("%s does not exist, skipping", v.Paths("image-json"))
			b = []byte("[]")
		}

		err = old.Unmarshal(bytes.NewBuffer(b))
		if err != nil {
			return err
		}

		log.Printf("%d images loaded from %s", len(old), v.Paths("image-json"))
		img.Merge(old)
		log.Printf("%d unique images to check", len(img))

		nimg, nerr := 0, 0
		// for k := range img {
		// 	e := img[k].ParseImageURL("")
		// 	if e != nil {
		// 		log.Printf("[parse error] %s", e)
		// 	}
		// 	i, e := img[k].CheckImageStatus()
		// 	nimg += i
		// 	if e != nil {
		// 		nerr += 1
		// 		log.Printf("[URL error] %s", e)
		// 	}
		// }

		buf := bytes.Buffer{}
		err = img.Marshal(&buf)
		if err != nil {
			return err
		}

		log.Printf("> [%s %s] %d checked, %d skipped, %d errors", size(buf.Len()),
			v.Paths("image-json"), nimg+nerr, len(img)-nimg, nerr)

		err = ioutil.WriteFile(v.Paths("image-json"), buf.Bytes(), 0644)
		if err != nil {
			return err
		}
	}
	return nil
}
