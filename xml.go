package wppub

import (
	"encoding/xml"
	"time"
)

type (
	RSS struct {
		XMLName xml.Name `xml:"rss"`
		Channel Channel  `xml:"channel"`
	}

	Channel struct {
		XMLName xml.Name `xml:"channel"`
		Items   []Item   `xml:"item"`
	}

	Item struct {
		XMLName    xml.Name   `xml:"item" json:"-"`
		Title      string     `xml:"title" json:"title"`
		Link       string     `xml:"link" json:"link"`
		PubDate    xmlTime    `xml:"pubDate" json:"pub_date"`
		Categories []Category `xml:"category" json:"categories"`
		GUID       string     `xml:"guid" json:"guid"`
		Body       Body       `xml:"http://purl.org/rss/1.0/modules/content/ encoded" json:"body"`
	}

	Category struct {
		XMLName xml.Name `xml:"category" json:"-"`
		Name    string   `xml:",cdata"`
	}

	Body struct {
		XMLName xml.Name `xml:"http://purl.org/rss/1.0/modules/content/ encoded" json:"-"`
		Text    string   `xml:",cdata"`
	}

	xmlTime struct {
		time.Time
	}
)

func (t *xmlTime) UnmarshalXML(d *xml.Decoder, s xml.StartElement) error {
	var v string
	d.DecodeElement(&v, &s)
	parse, err := time.Parse(time.RFC1123Z, v)
	if err != nil {
		return err
	}
	*t = xmlTime{parse.UTC()}
	return nil
}

func (x *xmlTime) Set(t time.Time) {
	*x = xmlTime{t}
}
