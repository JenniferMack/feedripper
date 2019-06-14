package feedpub

import (
	"fmt"
	"time"
)

type (
	Config struct {
		Name       string    `json:"name"`
		Number     string    `json:"number"`
		Deadline   time.Time `json:"deadline"` //RFC3339 = "2006-01-02T15:04:05Z07:00"
		Days       int       `json:"days"`
		SeqName    string    `json:"seq_name"`
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
		Tags       tags      `json:"tags"`
		Exclude    []string  `json:"exclude"`
		Feeds      []feed    `json:"feeds"`
	}
)

func (c Config) Names(path string) string {
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
	case "json", "xml", "html", "img.json":
		return fmt.Sprintf("%s.%s", name, path)
	case "image-404":
		return c.Image404
	case "dir-images":
		return c.ImageDir
	case "dir-rss":
		return c.RSSDir
	case "dir-json":
		return c.JSONDir
	}
	return ""
}

func (c Config) DateRange() string {
	str := c.Deadline
	end := c.Deadline.AddDate(0, 0, c.Days)
	if c.Days < 0 {
		str, end = end, str
	}

	strFmt := "02"
	if str.Month() < end.Month() {
		strFmt += " Jan"
	}

	if str.Year() < end.Year() {
		strFmt = "02 Jan 2006"
	}
	return fmt.Sprintf("%s–%s", str.Format(strFmt), end.Format("02 Jan 2006"))
}
