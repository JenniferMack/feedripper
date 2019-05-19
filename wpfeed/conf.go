package wpfeed

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

type (
	// Config holds the information for saving WordPress feeds.
	Config struct {
		Name       string    `json:"name"`
		Number     string    `json:"number"`
		Deadline   time.Time `json:"deadline"` //RFC3339 = "2006-01-02T15:04:05Z07:00"
		Days       int       `json:"days"`
		JSONDir    string    `json:"json_dir"`
		RSSDir     string    `json:"rss_dir"`
		Language   string    `json:"language"`
		SiteURL    string    `json:"site_url"`
		MainTagNum uint      `json:"main_tag_num"`
		Tags       Tags      `json:"tags"`
		Exclude    []string  `json:"exclude"`
		Feeds      []Feed    `json:"feeds"`
	}

	Tags []Tag

	// Tag holds the tag name and priority
	Tag struct {
		Name     string `json:"name"`
		Text     string `json:"text"`
		Priority uint   `json:"priority"`
		Limit    uint   `json:"limit"`
	}

	// Feed holds the infomation for a particular feed.
	Feed struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	}
)

// sort interface for Tags
func (t Tags) Len() int           { return len(t) }
func (t Tags) Swap(i, j int)      { t[i], t[j] = t[j], t[i] }
func (t Tags) Less(i, j int) bool { return t[i].Priority < t[j].Priority }

func (t Tags) PriOutOfRange() bool {
	idx := uint(0)
	cmp := uint(len(t) - 1)

	for _, v := range t {
		if v.Priority > idx {
			idx = v.Priority
		}
	}
	return idx > cmp
}

func (f Feed) FetchURL() ([]byte, error) {
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

// ReadConfig returns a slice of `Config`.
func ReadConfig(in io.Reader) ([]Config, error) {
	c := []Config{}
	err := json.NewDecoder(in).Decode(&c)
	return c, err
}
