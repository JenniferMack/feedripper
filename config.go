package feedpub

import (
	"fmt"
	"path/filepath"
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

	feed struct {
		Name string `json:"name"`
		URL  string `json:"url"`
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
	case "dir-rss", "dir-xml":
		return c.RSSDir
	case "dir-json":
		return c.JSONDir
	}
	return ""
}

func (c Config) feedPath(name, suf, typ string) string {
	if suf != "" {
		suf = `_` + suf
	}
	n := fmt.Sprintf("%s-%s%s.%s", name, c.Names("name"), suf, typ)
	return filepath.Join(c.Names("dir-"+typ), n)
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

func (c Config) inRange(itm item) bool {
	str := c.Deadline
	end := c.Deadline.AddDate(0, 0, c.Days)
	if c.Days < 0 {
		str, end = end, str
	}

	if itm.PubDate.After(str) && itm.PubDate.Before(end) {
		return true
	}
	return false
}
