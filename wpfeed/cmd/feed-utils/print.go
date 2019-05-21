package main

import (
	"fmt"
	"os"

	"golang.org/x/net/html"
)

func printHeadings() {
	var cnt int
	h := func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "h1" {
			fmt.Println(n.FirstChild.Data)
			cnt = 1
			return
		}
		if n.Type == html.ElementNode && n.Data == "h2" {
			f := n.FirstChild.NextSibling
			fmt.Printf("%02d  - %s\n", cnt, f.FirstChild.Data)
			cnt++
			return
		}
	}

	doc, err := html.Parse(os.Stdin)
	if err != nil {
		return
	}
	parseHTML(doc, h, nil)
}

func parseHTML(n *html.Node, before, after func(n *html.Node)) {
	if before != nil {
		before(n)
	}

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		parseHTML(c, before, after)
	}

	if after != nil {
		after(n)
	}
}
