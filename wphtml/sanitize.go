package wphtml

import (
	"bytes"
	"html"
	"strings"

	"github.com/microcosm-cc/bluemonday"
	"gopkg.in/russross/blackfriday.v2"
)

func cleanHTML(h string, re []RegexList) string {
	d := html.UnescapeString(h)
	for _, r := range re {
		d = r.ReplaceAll(d)
	}

	clean := sanitize(d)
	return makeHTML(clean)
}

func smartenString(s string) string {
	re := strings.NewReplacer(
		`“`, `"`,
		`”`, `"`,
		`‘`, `'`,
		`’`, `'`,
	)
	s = re.Replace(s)
	sp := blackfriday.NewSmartypantsRenderer(
		blackfriday.Smartypants |
			blackfriday.SmartypantsDashes |
			blackfriday.SmartypantsLatexDashes,
	)

	out := bytes.Buffer{}
	sp.Process(&out, []byte(s))

	return out.String()
}

func makeHTML(s string) string {
	bf := blackfriday.Run([]byte(s), blackfriday.WithRenderer(
		blackfriday.NewHTMLRenderer(
			blackfriday.HTMLRendererParameters{
				Flags: blackfriday.UseXHTML |
					blackfriday.Smartypants |
					blackfriday.SmartypantsDashes |
					blackfriday.SmartypantsLatexDashes,
			},
		),
	))
	return string(bf)
}

func sanitize(h string) string {
	bm := bluemonday.NewPolicy()
	bm.AllowAttrs("href", "title").OnElements("a")
	bm.AllowAttrs("src").OnElements("img")
	bm.AllowElements("ul", "ol", "li", "br", "p", "em",
		"h1", "h2", "h3", "h4", "h5", "h6",
		"blockquote", "strong", "figure", "figcaption")

	return bm.Sanitize(h)
}
