package main

import (
	"bytes"
	"io"
	"log"

	"repo.local/wputil"
	"repo.local/wputil/wpimage"
)

func doStatus(list wpimage.ImageList, name string, out io.Writer) {
	log.SetOutput(out)
	log.SetPrefix("[  status] ")

	log.Printf("> Status report for %s", name)
	log.Printf("%d image paths", len(list))
	log.Printf("%d valid images", list.ValidNum())
	log.Printf("%d saved images", list.SavedNum())

	buf := bytes.Buffer{}
	n := 0
	log.SetOutput(&buf)
	for k, v := range list {
		if v.Err != "" {
			log.Printf("%d - %s", k, wputil.Trim(80, v.Err))
			n++
		}
	}
	log.SetOutput(out)
	if n > 0 {
		log.Printf("%d errors recorded:", n)
		buf.WriteTo(out)
	} else {
		log.Print("no errors found")
	}

	n = 0
	for _, v := range list {
		if v.Resp >= 400 {
			n++
		}
	}
	if n > 0 {
		log.Printf("%d non-200 http responses found", n)
	} else {
		log.Print("no non-200 http responses found")
	}
}
