package main

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"

	"repo.local/wputil/wpfeed"
	"repo.local/wputil/wpimage"
)

func images(c io.Reader, a []string) error {
	log.SetFlags(0)
	log.SetPrefix("[  images] ")

	conf, err := wpfeed.ReadConfig(c)
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

func doAction(a string, c wpfeed.Config) error {
	paths, err := c.Paths(os.Getwd())
	if err != nil {
		return err
	}
	list, err := readExistingFile(paths["images"])
	if err != nil {
		return err
	}
	// load early incase of changes
	list.SetDefaults(c.ImageQual, c.ImageWidth, c.UseTLS)

	var out []byte
	var e error
	outfile := paths["images"]
	wr := os.Stderr

	switch a {
	case "status":
		doStatus(list, paths["images"], wr)
		return nil

	case "parse":
		out, e = doParse(list, paths["images"], paths["html"])

	case "filter":
		out, e = doFilter(list, paths["images"], c.SiteURL)

	case "verify":
		out, e = doVerify(list, paths["images"])

	case "fetch":
		out, e = doFetch(list, paths["images"], paths["imageDir"], wr)

	case "update":
		htm, err := ioutil.ReadFile(paths["html"])
		if err != nil {
			return err
		}

		out, e = doUpdate(list, htm, paths["html"], wr)
		outfile = paths["html"]
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
