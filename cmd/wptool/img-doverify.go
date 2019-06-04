package main

import (
	"bytes"
	"io"
	"log"

	"repo.local/wputil"
	"repo.local/wputil/wpimage"
)

func doVerify(list wpimage.ImageList, conf wputil.Config, wr io.Writer) ([]byte, error) {
	log.SetOutput(wr)
	log.SetPrefix("[  verify] ")

	log.Printf("> verifying image URLs [%s]", conf.Paths("image-json"))

	ch := make(chan wpimage.ImageData, 10)
	go list.CheckStatus(ch, *flagImageVerbose, conf.Paths("image-404"))

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

	log.Printf("%d images checked, %d found, %d not found, %d on disk, %d errors",
		len(out), out.ValidNum(), len(out)-out.ValidNum(), out.SavedNum(), n)

	buf := bytes.Buffer{}
	if err := out.Marshal(&buf); err != nil {
		return nil, err
	}

	log.Printf("> [%s/%d/%d] %s", size(buf.Len()), len(out), out.ValidNum(), conf.Paths("image-json"))
	return buf.Bytes(), nil
}
