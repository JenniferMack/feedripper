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
	Config struct {
		Name       string    `json:"name"`
		Number     string    `json:"number"`
		Deadline   time.Time `json:"deadline"` //RFC3339 = "2006-01-02T15:04:05Z07:00"
		Days       uint      `json:"days"`
		JSONDir    string    `json:"json_dir"`
		RSSDir     string    `json:"rss_dir"`
		Language   string    `json:"language"`
		SiteURL    string    `json:"site_url"`
		MainTagNum uint      `json:"main_tag_num"`
		Tags       []tag     `json:"tags"`
		Exclude    []string  `json:"exclude"`
		Feeds      []Feed    `json:"feeds"`
	}

	Feed struct {
		Name string `json:"name"`
		URL  string `json:"url"`
		Type string `json:"type"`
	}

	tag struct {
		Priority uint   `json:"priority"`
		Name     string `json:"name"`
	}
)

func (f Feed) fetch() ([]byte, error) {
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

func readConfig(in io.Reader) ([]Config, error) {
	c := []Config{}
	err := json.NewDecoder(in).Decode(&c)
	return c, err
}
