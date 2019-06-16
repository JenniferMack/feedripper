package feedpub

import (
	"fmt"
	"io"
	"log"
	"strings"

	"golang.org/x/net/html"
)

func ExtractImages(conf Config, pp bool, lg *log.Logger) error {
	itms, _ := oldItems(conf.Names("json"))
	cnt := 0

	for k, v := range itms {
		u := []string{}
		str, err := Parse(strings.NewReader(v.Body),
			ConvertElemIf("iframe", "img", "src", "youtube.com"),
			ExtractAttr("img", "src", &u),
		)
		if err != nil {
			return fmt.Errorf("html parse: %s", err)
		}

		itms[k].Body = str

		for _, i := range u {
			if strings.Contains(i, "?") {
				continue
			}
			lp := makeLocPath(i)
			itms[k].Images = append(itms[k].Images, image{
				URL:       i,
				LocalPath: lp,
				OnDisk:    isOnDisk(lp),
			})
			cnt++
		}
	}

	lg.Printf("[%03d] images => %s", cnt, conf.Names("json"))

	_, err := writeJSON(itms, conf.Names("json"), pp)
	if err != nil {
		return fmt.Errorf("json write: %s", err)
	}
	return nil
}

func makeLocPath(p string) string {
	return p
}

func isOnDisk(p string) bool {
	return false
}

func Parse(h io.Reader, opts ...func(*html.Node)) (string, error) {
	doc, err := html.Parse(h)
	if err != nil {
		return "", err
	}

	for _, opt := range opts {
		parser(doc, opt)
	}

	ret := strings.Builder{}
	html.Render(&ret, doc)

	htm := strings.TrimPrefix(ret.String(), "<html><head></head><body>")
	if htm != ret.String() {
		htm = strings.TrimSuffix(htm, "</body></html>")
	}
	return htm, nil
}

func parser(n *html.Node, fn func(*html.Node)) {
	fn(n)
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		parser(c, fn)
	}
}

func isNode(n *html.Node, t html.NodeType, d string) bool {
	return n.Type == t && n.Data == d
}

func newNode(t html.NodeType, d string) *html.Node {
	n := &html.Node{
		Type: t,
		Data: d,
	}
	return n
}

func ExtractAttr(elem, attr string, arr *[]string) func(*html.Node) {
	return func(n *html.Node) {
		if isNode(n, html.ElementNode, elem) {
			for _, v := range n.Attr {
				if v.Key == attr {
					*arr = append(*arr, v.Val)
				}
			}
		}
	}
}

func ReplaceAttr(elem, attr, old, repl string) func(*html.Node) {
	return func(n *html.Node) {
		if isNode(n, html.ElementNode, elem) {
			for k, v := range n.Attr {
				if v.Key == attr && v.Val == old {
					n.Attr[k].Val = repl
				}
			}
		}
	}
}

func ReplaceElem(from, to string) func(*html.Node) {
	return func(n *html.Node) {
		if isNode(n, html.ElementNode, from) {
			n.Data = to
		}
	}
}

func WrapElem(inner, outer string) func(*html.Node) {
	return func(n *html.Node) {
		if isNode(n, html.ElementNode, inner) {
			nu := newNode(html.ElementNode, outer)
			n.Parent.InsertBefore(nu, n)
			n.Parent.RemoveChild(n)
			nu.AppendChild(n)
		}
	}
}

func UnwrapElem(inner, outer string) func(*html.Node) {
	return func(n *html.Node) {
		if isNode(n, html.ElementNode, outer) {
			if n.FirstChild != nil && isNode(n.FirstChild, html.ElementNode, inner) {
				fc := n.FirstChild
				n.RemoveChild(fc)
				n.Parent.InsertBefore(fc, n)
				n.Parent.RemoveChild(n)
			}
		}
	}
}

func AddCaption(find, caption string) func(*html.Node) {
	return func(n *html.Node) {
		if isNode(n, html.ElementNode, "figure") {
			if n.FirstChild != nil && isNode(n.FirstChild, html.ElementNode, "img") {

				for _, v := range n.FirstChild.Attr {
					if v.Key == "src" && strings.Contains(v.Val, find) {

						nu := newNode(html.ElementNode, "figcaption")
						nu.FirstChild = newNode(html.TextNode, caption)
						n.AppendChild(nu)
						break
					}
				}
			}
		}
	}
}

func ConvertToLink(from, link string) func(*html.Node) {
	return func(n *html.Node) {
		if isNode(n, html.ElementNode, from) {
			n.Data = "a"
			for k, v := range n.Attr {
				if v.Key == "src" {
					n.Attr[k].Key = "href"
					break
				}
			}
			nu := newNode(html.TextNode, link)
			n.AppendChild(nu)
		}
	}
}

func ConvertElemIf(from, to, attr, cond string) func(*html.Node) {
	return func(n *html.Node) {
		if isNode(n, html.ElementNode, from) {
			for _, v := range n.Attr {
				if v.Key == attr && strings.Contains(v.Val, cond) {
					n.Data = to
				}
			}
		}
	}
}
