package wpimage

import (
	"fmt"
	"io"

	"golang.org/x/net/html"
)

func ParseHTML(in io.Reader) (ImageList, error) {
	doc, err := html.Parse(in)
	if err != nil {
		return nil, fmt.Errorf("parse html: %s", err)
	}

	list := []ImageData{}
	fn := func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "img" {
			for _, v := range n.Attr {
				if v.Key == "src" {
					list = append(list, ImageData{Rawpath: v.Val})
				}
			}
		}
	}

	parse(doc, fn)
	return list, nil
}

func parse(n *html.Node, f func(n *html.Node)) {
	f(n)
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		parse(c, f)
	}
}
