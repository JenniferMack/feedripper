package wputil

import (
	"encoding/xml"
	"strings"
	"time"
)

type (
	rss struct {
		XMLName xml.Name `xml:"rss"`
		Channel channel  `xml:"channel"`
	}

	channel struct {
		XMLName xml.Name `xml:"channel"`
		Items   []Item   `xml:"item"`
	}

	Item struct {
		XMLName    xml.Name   `xml:"item" json:"-"`
		Title      string     `xml:"title" json:"title"`
		Link       string     `xml:"link" json:"link"`
		PubDate    xmlTime    `xml:"pubDate" json:"pub_date"`
		Categories []category `xml:"category" json:"categories"`
		GUID       string     `xml:"guid" json:"guid"`
		Body       body       `xml:"http://purl.org/rss/1.0/modules/content/ encoded" json:"body"`
	}

	category struct {
		XMLName xml.Name `xml:"category" json:"-"`
		Name    string   `xml:",cdata" json:"name"`
	}

	body struct {
		XMLName xml.Name `xml:"http://purl.org/rss/1.0/modules/content/ encoded" json:"-"`
		Text    string   `xml:",cdata" json:"text"`
	}

	xmlTime struct {
		time.Time
	}
)

func (i Item) hasTag(t string) bool {
	for _, v := range i.Categories {
		if strings.EqualFold(v.Name, t) {
			return true
		}
	}
	return false
}

func (i Item) hasTagList(t []string) bool {
	for _, tag := range t {
		if i.hasTag(tag) {
			return true
		}
	}
	return false
}

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
