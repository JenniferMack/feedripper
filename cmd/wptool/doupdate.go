package main

import (
	"bytes"
	"io"
	"log"

	"golang.org/x/net/html"
	"repo.local/wputil/wpimage"
)

func doUpdate(list wpimage.ImageList, in []byte, n string, out io.Writer) ([]byte, error) {
	log.SetOutput(out)
	log.Printf("> updating %s", n)

	doc, err := html.Parse(bytes.NewBuffer(in))
	if err != nil {
		return nil, err
	}

	var imgcnt, chgcnt int
	fn := func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "img" {
			imgcnt++
			for k, v := range n.Attr {
				if v.Key == "src" {
					if i, ok := list.MatchRawPath(v.Val); ok {
						n.Attr[k].Val = i
						chgcnt++
					}
				}
			}
		}
	}

	parseHTML(doc, fn)
	buf := bytes.Buffer{}
	err = html.Render(&buf, doc)
	if err != nil {
		return nil, err
	}
	log.Printf("%d images found, %d URLs modified", imgcnt, chgcnt)
	log.Printf("> [%s/%d] %s", size(buf.Len()), imgcnt, n)

	return buf.Bytes(), nil
}

func parseHTML(n *html.Node, before func(n *html.Node)) {
	if before != nil {
		before(n)
	}

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		parseHTML(c, before)
	}
}
