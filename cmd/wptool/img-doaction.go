package main

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"

	"repo.local/wputil"
	"repo.local/wputil/wpimage"
)

func images(c io.Reader, a []string) error {
	log.SetFlags(0)
	log.SetPrefix("[  images] ")

	conf, err := wputil.NewConfigList(c)
	if err != nil {
		return fmt.Errorf("reading config: %s", err)
	}

	for _, confItem := range conf {
		for _, action := range a {
			err := doAction(action, confItem)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func doAction(a string, c wputil.Config) error {
	// paths, err := c.Paths(os.Getwd())
	// if err != nil {
	// 	return err
	// }
	list, err := readExistingFile(c.Paths("image-json"))
	if err != nil {
		return err
	}
	// load early incase of changes
	list.SetDefaults(c.ImageQual, c.ImageWidth, c.UseTLS)

	var out []byte
	var e error
	outfile := c.Paths("image-json")
	wr := os.Stderr

	switch a {
	case "status":
		doStatus(list, c, wr)
		return nil

	case "parse":
		out, e = doParse(list, c, wr)

	case "filter":
		out, e = doFilter(list, c, wr)

	case "verify":
		out, e = doVerify(list, c, wr)

	case "fetch":
		out, e = doFetch(list, c, wr)

	case "update":
		out, e = doUpdate(list, c, wr)
		outfile = c.Paths("image-html")
	}

	if e != nil {
		return err
	}
	return ioutil.WriteFile(outfile, out, 0644)
}

func readExistingFile(p string) (wpimage.ImageList, error) {
	b, err := ioutil.ReadFile(p)
	if err != nil {
		b = []byte("[]")
	}

	list := wpimage.ImageList{}
	err = list.Unmarshal(bytes.NewBuffer(b))
	if err != nil {
		return nil, err
	}
	return list, nil
}