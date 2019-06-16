package feedpub

import (
	"encoding/xml"
)

type (
	channel struct {
		XMLName xml.Name `xml:"channel"`
		Items   []item   `xml:"item"`
	}

	rss struct {
		XMLName xml.Name `xml:"rss"`
		Channel channel  `xml:"channel"`
	}
)
