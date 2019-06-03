package wputil

import (
	"fmt"
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
		ImageDir   string    `json:"image_dir"`
		UseTLS     bool      `json:"use_tls"`
		ImageQual  int       `json:"image_qual"`
		ImageWidth uint      `json:"image_width"`
		Image404   string    `json:"image_404"`
		Language   string    `json:"language"`
		SiteURL    string    `json:"site_url"`
		Separator  string    `json:"separator"`
		Tags       Tags      `json:"tags"`
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

	//
	Configs []Config

	//
	Tags []Tag
)

// sort interface for Tags
func (t Tags) Len() int           { return len(t) }
func (t Tags) Swap(i, j int)      { t[i], t[j] = t[j], t[i] }
func (t Tags) Less(i, j int) bool { return t[i].Priority < t[j].Priority }

func (t Tags) priOutOfRange() bool {
	idx := uint(0)
	cmp := uint(len(t) - 1)

	for _, v := range t {
		if v.Priority > idx {
			idx = v.Priority
		}
	}
	return idx > cmp
}

func (c Config) FeedName(feed string) string {
	return fmt.Sprintf("%s-%s", feed, c.Paths("json"))
}

func (c Config) Paths(path string) string {
	name := c.Name
	if c.Language != "" {
		name += "-" + c.Language
	}
	if c.Number != "" {
		name += "_" + c.Number
	}

	switch path {
	case "name":
		return name
	case "json":
		return fmt.Sprintf("%s.%s", name, path)
	case "xml":
		return fmt.Sprintf("%s.%s", name, path)
	case "html":
		return fmt.Sprintf("%s.%s", name, path)
	case "image-json":
		return fmt.Sprintf("%s-%s.%s", name, "image", "json")
	case "image-html":
		return fmt.Sprintf("%s-%s.%s", name, "image", "html")
	case "image-dir":
		return c.ImageDir
	case "image-404":
		return c.Image404
	case "rss-dir":
		return c.RSSDir
	case "json-dir":
		return c.JSONDir
	}
	return ""
}
