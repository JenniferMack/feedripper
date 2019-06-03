package main

import (
	"bytes"
	"io"
	"io/ioutil"
	"log"

	"repo.local/wputil"
	"repo.local/wputil/wpimage"
)

func doParse(list wpimage.ImageList, conf wputil.Config, out io.Writer) ([]byte, error) {
	log.SetOutput(out)
	log.SetPrefix("[   parse] ")

	b, err := ioutil.ReadFile(conf.Paths("html"))
	if err != nil {
		return nil, err
	}
	imgs, err := wpimage.ParseHTML(bytes.NewReader(b))
	if err != nil {
		return nil, err
	}

	merged := imgs.Merge(list)
	log.Printf("> parsing HTML [%s]", conf.Paths("html"))
	log.Printf("%-3d images found in %s", len(imgs), conf.Paths("html"))
	log.Printf("%-3d images found in %s", len(list), conf.Paths("images-json"))

	buf := bytes.Buffer{}
	if err := merged.Marshal(&buf); err != nil {
		return nil, err
	}

	log.Printf("%-3d unique images recorded", len(merged))
	log.Printf("> [%s/%d] %s", wputil.FileSize(buf.Len()), len(merged), conf.Paths("images-json"))
	return buf.Bytes(), nil
}
