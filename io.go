package wputil

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

// NewConfigList returns the config info.
func NewConfigList(in io.Reader) (Configs, error) {
	c := Configs{}
	err := json.NewDecoder(in).Decode(c)
	if err != nil {
		return nil, err
	}

	for _, v := range c {
		if v.Tags.priOutOfRange() {
			return nil, fmt.Errorf("[%s] priority out of range", v.Name)
		}
	}
	return c, nil
}

// ReadWPXML reads WordPress RSS feed XML from a io.Reader and returns a populated Feed.
// Duplicates are removed and the internal list is sorted newest first.
func ReadWPXML(in io.Reader) (Feed, error) {
	r := rss{}
	f := Feed{}

	err := xml.NewDecoder(in).Decode(&r)
	if err != nil {
		return f, err
	}

	f.Merge(r.Channel.Items)
	return f, nil
}

func WriteRawXML(b []byte, dir, name string) error {
	stamp := "_" + strconv.FormatInt(time.Now().Unix(), 10)
	path := filepath.Join(dir, name+stamp+".xml")
	return ioutil.WriteFile(path, b, 0644)
}

// ReadWPJSON reads JSON from an io.Reader and returns a populated Feed.
// Duplicates are removed and the internal list is sorted newest first.
func ReadWPJSON(in io.Reader) (Feed, error) {
	f := Feed{}
	_, err := io.Copy(&f, in)
	if err != nil {
		return f, err
	}
	return f, nil
}

// FetchFeed returns an XML feed.
func FetchFeed(f Feed) ([]byte, error) {
	if status, ok := statusOK(f.URL); !ok {
		return nil, fmt.Errorf("[%s] %s", status, f.URL)
	}

	resp, err := http.Get(f.URL)
	if err != nil {
		return nil, fmt.Errorf("getting feed: %s", err)
	}
	defer resp.Body.Close()

	b := bytes.Buffer{}

	_, err = b.ReadFrom(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("reading feed: %s", err)
	}
	return b.Bytes(), nil
}

func statusOK(u string) (string, bool) {
	resp, err := http.Head(u)
	if err != nil {
		return err.Error(), false
	}
	if resp.StatusCode != 200 {
		return resp.Status, false
	}
	if !strings.Contains(resp.Header.Get("content-type"), "xml") {
		return resp.Header.Get("content-type"), false
	}
	return resp.Status, true
}
