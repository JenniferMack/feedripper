// Package feed provides tools for getting and saving WordPress feeds.
package feed

import (
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"repo.local/wputil"
)

func WriteRawXML(b []byte, dir, name string) error {
	stamp := "_" + strconv.FormatInt(time.Now().Unix(), 10)
	path := filepath.Join(dir, name+stamp+".xml")
	return ioutil.WriteFile(path, b, 0644)
}

func WriteJSON(feed wputil.Feed, dir, name string) (int, int, error) {
	path := filepath.Join(dir, name+".json")
	// Read saved JSON if any
	b := readJSONFile(path)
	_, err := feed.Write(b)
	if err != nil {
		return 0, 0, err
	}

	f, err := os.Create(path)
	if err != nil {
		return 0, 0, err
	}
	defer f.Close()

	n, err := io.Copy(f, &feed)
	if err != nil {
		return 0, 0, err
	}

	err = f.Close()
	if err != nil {
		return 0, 0, err
	}
	return int(n), feed.Len(), nil
}

func readJSONFile(p string) []byte {
	b, err := ioutil.ReadFile(p)
	if err != nil {
		return nil
	}
	return b
}
