package main

import (
	"bytes"
	"fmt"
	"io"

	"golang.org/x/net/html"
)

func printHeadings(in io.Reader) []byte {
	var cnt int
	out := bytes.Buffer{}

	h := func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "h1" {
			fmt.Fprintf(&out, n.FirstChild.Data)
			cnt = 1
			return
		}
		if n.Type == html.ElementNode && n.Data == "h2" {
			f := n.FirstChild.NextSibling
			fmt.Fprintf(&out, "%02d  - %s\n", cnt, f.FirstChild.Data)
			cnt++
			return
		}
	}

	doc, err := html.Parse(in)
	if err != nil {
		return nil
	}

	parseHTML(doc, h)
	return out.Bytes()
}

func parseHTML(n *html.Node, before func(n *html.Node)) {
	if before != nil {
		before(n)
	}

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		parseHTML(c, before)
	}
}
