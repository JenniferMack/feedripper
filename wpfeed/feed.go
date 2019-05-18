// Package wpfeed provides tools for getting and saving WordPress feeds.
package wpfeed

import (
	"io/ioutil"
	"path/filepath"
	"strconv"
	"time"
)

func WriteRawXML(b []byte, dir, name string) error {
	stamp := "_" + strconv.FormatInt(time.Now().Unix(), 10)
	path := filepath.Join(dir, name+stamp+".xml")
	return ioutil.WriteFile(path, b, 0644)
}
