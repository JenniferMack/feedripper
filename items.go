package feedpub

import (
	"encoding/xml"
	"time"
)

type (
	items []item

	item struct {
		XMLName    xml.Name   `xml:"item" json:"-"`
		Title      string     `xml:"title" json:"title"`
		Link       string     `xml:"link" json:"link"`
		PubDate    xmlTime    `xml:"pubDate" json:"pub_date"`
		Categories []category `xml:"category" json:"categories"`
		GUID       string     `xml:"guid" json:"guid"`
		Body       body       `xml:"http://purl.org/rss/1.0/modules/content/ encoded" json:"body"`
	}

	body struct {
		XMLName xml.Name `xml:"http://purl.org/rss/1.0/modules/content/ encoded" json:"-"`
		Text    string   `xml:",cdata" json:"text"`
	}

	category struct {
		XMLName xml.Name `xml:"category" json:"-"`
		Name    string   `xml:",cdata" json:"name"`
	}

	xmlTime struct {
		time.Time
	}
)

// sort items by date
func (i items) Len() int           { return len(i) }
func (i items) Swap(j, k int)      { i[j], i[k] = i[k], i[j] }
func (i items) Less(j, k int) bool { return i[j].PubDate.Before(i[k].PubDate.Time) }

func (i *items) trim(conf Config) {
	it := items{}
	for _, v := range *i {
		if conf.inRange(v) {
			it = append(it, v)
		}
	}
	*i = it
}

func (i *items) add(list items) {
	it := append(*i, list...)
	dup := make(map[string]item)

	for _, v := range it {
		d, ok := dup[v.GUID]
		if !ok {
			dup[v.GUID] = v
			continue
		}

		if v.PubDate.After(d.PubDate.Time) {
			dup[v.GUID] = v
		}
	}

	ii := items{}
	for _, v := range dup {
		ii = append(ii, v)
	}
	*i = ii
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
