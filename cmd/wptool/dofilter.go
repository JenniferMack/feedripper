package main

import (
	"bytes"
	"log"

	"repo.local/wputil/wpfeed"
	"repo.local/wputil/wpimage"
)

func doFilter(list wpimage.ImageList, host string, paths wpfeed.Paths) ([]byte, error) {
	log.Printf("> parsing image URLs [%s]", paths["images"])

	n := 0
	for k := range list {
		n += list[k].ParseImageURL(host, paths["imageDir"], paths["404-img"])
	}
	log.Printf("%d vaild URLs, %d errors", n, len(list)-n)

	buf := bytes.Buffer{}
	if err := list.Marshal(&buf); err != nil {
		return nil, err
	}

	log.Printf("> [%s/%d] %s", size(buf.Len()), len(list), paths["images"])
	return buf.Bytes(), nil
}
