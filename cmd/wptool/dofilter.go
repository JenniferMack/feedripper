package main

import (
	"bytes"
	"log"

	"repo.local/wputil/wpimage"
)

func doFilter(list wpimage.ImageList, path, u string) ([]byte, error) {
	log.Printf("> parsing image URLs [%s]", path)

	n := 0
	for k := range list {
		n += list[k].ParseImageURL(u)
	}
	log.Printf("%d vaild URLs, %d errors", n, len(list)-n)

	buf := bytes.Buffer{}
	if err := list.Marshal(&buf); err != nil {
		return nil, err
	}

	log.Printf("> [%s/%d] %s", size(buf.Len()), len(list), path)
	return buf.Bytes(), nil
}
