package wputil

import (
	"time"
)

type (
	// Config holds the information for saving WordPress feeds.
	config struct {
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
		Feeds      []feed    `json:"feeds"`
	}

	// Tag holds the tag name and priority
	tag struct {
		Name     string `json:"name"`
		Text     string `json:"text"`
		Priority uint   `json:"priority"`
		Limit    uint   `json:"limit"`
	}

	//
	Configs []config

	//
	Tags []tag
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
