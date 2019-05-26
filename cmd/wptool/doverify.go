package main

import (
	"bytes"
	"log"

	"repo.local/wputil/wpfeed"
	"repo.local/wputil/wpimage"
)

func doVerify(list wpimage.ImageList, paths wpfeed.Paths) ([]byte, error) {
	log.Printf("> verifying image URLs [%s]", paths["images"])

	ch := make(chan wpimage.ImageData, 10)
	go list.CheckStatus(ch, *flagImageVerbose, paths["404-img"])

	n := 0
	o := []wpimage.ImageData{}
	for v := range ch {
		if v.Err != "" {
			log.Printf("[  error] %s", v.Err)
			n++
		}
		o = append(o, v)
	}
	out := wpimage.ImageList(o)

	log.Printf("%d images checked, %d found, %d not found, %d errors",
		len(out), out.SavedNum(), len(out)-out.SavedNum(), n)

	buf := bytes.Buffer{}
	if err := out.Marshal(&buf); err != nil {
		return nil, err
	}

	log.Printf("> [%s/%d/%d] %s", size(buf.Len()), len(out), out.ValidNum(), paths["images"])
	return buf.Bytes(), nil
}
