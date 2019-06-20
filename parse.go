package wputil

import (
	"io"
	"strings"

	"golang.org/x/net/html"
)

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

func walkSiblingsIf(n *html.Node, attr, cond string) []*html.Node {
	list := []*html.Node{}
	for s := n.NextSibling; s != nil; s = s.NextSibling {
		if s.Type == n.Type && s.Data == n.Data {
			for _, v := range s.Attr {
				if v.Key == attr && strings.Contains(v.Val, cond) {
					list = append(list, s)
				}
			}
		}
	}
	return list
}

func DropElem(elem string) func(*html.Node) {
	return func(n *html.Node) {
		if isNode(n, html.ElementNode, elem) {
			n.Parent.RemoveChild(n)
		}
	}
}

func DropElemIf(elem, attr, cond string) func(*html.Node) {
	return func(n *html.Node) {
		if isNode(n, html.ElementNode, elem) {
			for _, v := range n.Attr {
				if v.Key == attr && strings.Contains(v.Val, cond) {
					list := []*html.Node{n}
					list = append(list, walkSiblingsIf(n, attr, cond)...)
					for _, v := range list {
						v.Parent.RemoveChild(v)
					}
				}
				break
			}
		}
	}
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
		if isNode(n, html.ElementNode, inner) && !isNode(n.Parent, html.ElementNode, outer) {
			nu := newNode(html.ElementNode, outer)
			n.Parent.InsertBefore(nu, n)
			n.Parent.RemoveChild(n)
			nu.AppendChild(n)
			if nu.NextSibling != nil {
				WrapElem(inner, outer)(nu.NextSibling)
			}
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
				if n.NextSibling != nil {
					UnwrapElem(inner, outer)(n.NextSibling)
				}
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
						ca := newNode(html.ElementNode, "a")
						ca.Attr = []html.Attribute{
							{Key: "href", Val: v.Val},
						}
						ca.AppendChild(newNode(html.TextNode, caption))
						nu.FirstChild = ca
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
			for _, v := range n.Attr {
				if v.Key == "src" {
					s := v.Val
					n.Attr = []html.Attribute{
						{Key: "href", Val: s},
					}
					break
				}
			}
			for c := n.FirstChild; c != nil; c = c.NextSibling {
				n.RemoveChild(c)
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

func RewrapImg(from, to string) func(*html.Node) {
	return func(n *html.Node) {
		if isNode(n, html.ElementNode, from) {
			if n.FirstChild != nil && isNode(n.FirstChild, html.ElementNode, "img") {
				n.Data = to
				n.Attr = nil
			}
		}
	}
}
