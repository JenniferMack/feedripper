package wppub

import (
	"encoding/json"
	"encoding/xml"
	"io"
)

func ReadWPXML(in io.Reader) ([]Item, error) {
	r := RSS{}
	err := xml.NewDecoder(in).Decode(&r)
	return r.Channel.Items, err
}

func ReadWPJSON(in io.Reader) ([]Item, error) {
	i := []Item{}
	err := json.NewDecoder(in).Decode(&i)
	return i, err
}

func WriteWPJSON(i []Item, out io.Writer) (int, error) {
	items, err := json.MarshalIndent(i, "", " ")
	if err != nil {
		return 0, err
	}
	return out.Write(items)
}
