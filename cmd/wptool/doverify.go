package main

import (
	"bytes"
	"log"

	"repo.local/wputil/wpimage"
)

func doVerify(list wpimage.ImageList, path string) ([]byte, error) {
	log.Printf("> verifying image URLs [%s]", path)

	ch := make(chan wpimage.ImageData, 10)
	go list.CheckStatus(ch, *flagImageVerbose)

	n := 0
	o := []wpimage.ImageData{}
	for v := range ch {
		if v.Err != "" {
			log.Printf("[error] %s", v.Err)
			n++
		}
		o = append(o, v)
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
