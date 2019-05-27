package main

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/url"

	"golang.org/x/net/html"
	"repo.local/wputil/wpfeed"
	"repo.local/wputil/wpimage"
)

func doUpdate(list wpimage.ImageList, paths wpfeed.Paths, c wpfeed.Config, out io.Writer) ([]byte, error) {
	htm, err := ioutil.ReadFile(paths["html"])
	if err != nil {
		return nil, err
	}

	log.SetOutput(out)
	log.Printf("> updating %s", paths["html-img"])

	doc, err := html.Parse(bytes.NewBuffer(htm))
	if err != nil {
		return nil, fmt.Errorf("html parse: %s", err)
	}

	// replace links with local links
	var imgcnt, chgcnt int
	ln := func(n *html.Node) {
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

	im := func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "a" {
			if n.FirstChild != nil {
				if n.FirstChild.Type == html.ElementNode && n.FirstChild.Data == "img" {
					n.Data = "figure"
					n.Attr = []html.Attribute{}
				}
			}
			for k, v := range n.Attr {
				if v.Key == "href" {
					u, _ := url.Parse(v.Val)
					if u.Host == "" {
						u.Host = c.SiteURL
						if c.UseTLS {
							u.Scheme = "https"
						} else {
							u.Scheme = "http"
						}
						n.Attr[k].Val = u.String()
					}
				}
			}
		}
	}

	parseHTML(doc, ln)
	parseHTML(doc, im)
	buf := bytes.Buffer{}
	err = html.Render(&buf, doc)
	if err != nil {
		return nil, fmt.Errorf("render: %s", err)
	}
	log.Printf("%d images found, %d URLs modified", imgcnt, chgcnt)
	log.Printf("> [%s/%d] %s", size(buf.Len()), imgcnt, paths["html-img"])

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
