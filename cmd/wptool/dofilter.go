package main

import (
	"bytes"
	"io"
	"log"

	"repo.local/wputil"
	"repo.local/wputil/wpimage"
)

func doFilter(list wpimage.ImageList, conf wputil.Config, out io.Writer) ([]byte, error) {
	log.SetOutput(out)
	log.SetPrefix("[  filter] ")

	log.Printf("> parsing image URLs [%s]", conf.Paths("images-json"))

	n := 0
	for k := range list {
		n += list[k].ParseImageURL(conf.SiteURL, conf.ImageDir, conf.Image404)
	}
	log.Printf("%d vaild URLs, %d errors", n, len(list)-n)

	buf := bytes.Buffer{}
	if err := list.Marshal(&buf); err != nil {
		return nil, err
	}

	log.Printf("> [%s/%d] %s", size(buf.Len()), len(list), conf.Paths("images-json"))
	return buf.Bytes(), nil
}
