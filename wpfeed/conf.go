package wpfeed

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"path/filepath"
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
		WorkDir    string    `json:"work_dir"`
		JSONDir    string    `json:"json_dir"`
		RSSDir     string    `json:"rss_dir"`
		ImageDir   string    `json:"image_dir"`
		UseTLS     bool      `json:"use_tls"`
		ImageQual  int       `json:"image_qual"`
		ImageWidth uint      `json:"image_width"`
		Language   string    `json:"language"`
		SiteURL    string    `json:"site_url"`
		Separator  string    `json:"separator"`
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

func (c Config) Paths(d string, e error) (map[string]string, error) {
	if e != nil {
		return nil, fmt.Errorf("unable to resolve working directory %s", e)
	}

	if c.isWorkDir(d) {
		c.WorkDir = "."
	}

	name := fmt.Sprintf("%s-%s", c.Name, c.Number)
	paths := make(map[string]string)
	paths["json"] = filepath.Join(c.WorkDir, name+".json")
	paths["images"] = filepath.Join(c.WorkDir, name+"-images.json")
	paths["html"] = filepath.Join(c.WorkDir, name+".html")
	paths["html-img"] = filepath.Join(c.WorkDir, name+"-images.html")
	paths["jsonDir"] = filepath.Join(c.WorkDir, c.JSONDir)
	paths["rssDir"] = filepath.Join(c.WorkDir, c.RSSDir)
	paths["imageDir"] = filepath.Join(c.WorkDir, c.ImageDir)
	return paths, nil
}

func (c Config) isWorkDir(d string) bool {
	return filepath.Base(d) == filepath.Base(c.WorkDir)
}

func (c Config) IsWorkDir(d string, e error) bool {
	if e != nil {
		log.Printf("unable to resolve working directory %s", e)
		return false
	}
	return filepath.Base(d) == filepath.Base(c.WorkDir)
}

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
	if err != nil {
		return nil, err
	}

	for _, v := range c {
		if v.Tags.PriOutOfRange() {
			return nil, fmt.Errorf("[%s] priority out of range", v.Name)
		}
	}
	return c, err
}
