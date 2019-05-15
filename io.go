package wputil

import (
	"encoding/xml"
	"io"
)

// ReadWPXML reads WordPress RSS feed XML from a io.Reader and returns a populated Feed.
// Duplicates are removed and the internal list is sorted newest first.
func ReadWPXML(in io.Reader) (Feed, error) {
	r := rss{}
	f := Feed{}

	err := xml.NewDecoder(in).Decode(&r)
	if err != nil {
		return f, err
	}

	f.Merge(r.Channel.Items)
	return f, nil
}

// ReadWPJSON reads JSON from an io.Reader and returns a populated Feed.
// Duplicates are removed and the internal list is sorted newest first.
func ReadWPJSON(in io.Reader) (Feed, error) {
	f := Feed{}
	_, err := io.Copy(&f, in)
	if err != nil {
		return f, err
	}
	return f, nil
}
