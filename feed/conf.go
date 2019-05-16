package feed

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
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
		Tags       []Tag     `json:"tags"`
		Exclude    []string  `json:"exclude"`
		Feeds      []Feed    `json:"feeds"`
	}

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

func (f Feed) FetchURL() ([]byte, error) {
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

// ReadConfig returns a slice of `Config`.
func ReadConfig(in io.Reader) ([]Config, error) {
	c := []Config{}
	err := json.NewDecoder(in).Decode(&c)
	return c, err
}
