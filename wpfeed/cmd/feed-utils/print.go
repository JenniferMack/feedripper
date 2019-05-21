package main

import (
	"bytes"
	"fmt"
	"os"

	"golang.org/x/net/html"
)

func printHeadings(feed string) ([]byte, error) {
	in, err := os.Open(feed)
	if err != nil {
		return nil, err
	}
	defer in.Close()

	var cnt int
	out := bytes.Buffer{}

	h := func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "h1" {
			fmt.Fprintf(&out, "-- %s --\n", n.FirstChild.Data)
			cnt = 1
			return
		}
		if n.Type == html.ElementNode && n.Data == "h2" {
			f := n.FirstChild.NextSibling
			fmt.Fprintf(&out, "%02d - %s\n", cnt, f.FirstChild.Data)
			cnt++
			return
		}
	}

	doc, err := html.Parse(in)
	if err != nil {
		return nil, err
	}
	in.Close()

	parseHTML(doc, h)
	return out.Bytes(), nil
}

func parseHTML(n *html.Node, before func(n *html.Node)) {
	if before != nil {
		before(n)
	}

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		parseHTML(c, before)
	}
}
