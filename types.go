package feedpub

import (
	"encoding/xml"
	"time"
)

type (
	body struct {
		XMLName xml.Name `xml:"http://purl.org/rss/1.0/modules/content/ encoded" json:"-"`
		Text    string   `xml:",cdata" json:"text"`
	}

	category struct {
		XMLName xml.Name `xml:"category" json:"-"`
		Name    string   `xml:",cdata" json:"name"`
	}

	channel struct {
		XMLName xml.Name `xml:"channel"`
		Items   []item   `xml:"item"`
	}

	feed struct {
		Name  string `json:"name"`
		URL   string `json:"url"`
		items items
		json  []byte
		index int
	}

	item struct {
		XMLName    xml.Name   `xml:"item" json:"-"`
		Title      string     `xml:"title" json:"title"`
		Link       string     `xml:"link" json:"link"`
		PubDate    xmlTime    `xml:"pubDate" json:"pub_date"`
		Categories []category `xml:"category" json:"categories"`
		GUID       string     `xml:"guid" json:"guid"`
		Body       body       `xml:"http://purl.org/rss/1.0/modules/content/ encoded" json:"body"`
	}

	rss struct {
		XMLName xml.Name `xml:"rss"`
		Channel channel  `xml:"channel"`
	}

	xmlTime struct {
		time.Time
	}

	items []item
)
