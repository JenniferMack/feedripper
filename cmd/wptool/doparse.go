package main

import (
	"bytes"
	"log"
	"os"

	"repo.local/wputil/wpimage"
)

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
