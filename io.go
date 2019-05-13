package wppub

import (
	"encoding/json"
	"encoding/xml"
	"io"
)

func ReadWPXML(in io.Reader) ([]WPItem, error) {
	r := RSS{}
	err := xml.NewDecoder(in).Decode(&r)
	return r.Channel.Items, err
}

func ReadWPJSON(in io.Reader) ([]WPItem, error) {
	i := []WPItem{}
	err := json.NewDecoder(in).Decode(&i)
	return i, err
}

func WriteWPJSON(i []WPItem, out io.Writer) (int, error) {
	items, err := json.MarshalIndent(i, "", " ")
	if err != nil {
		return 0, err
	}
	return out.Write(items)
}
