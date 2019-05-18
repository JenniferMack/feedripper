package wphtml

import (
	"bytes"
	"fmt"
	"regexp"
	"sort"

	"gopkg.in/russross/blackfriday.v2"
	"repo.local/wputil"
	wpfeed "repo.local/wputil/feed"
)

type RegexList struct {
	Pattern string
	Replace string
	Re      *regexp.Regexp
}

func (r *RegexList) Compile() {
	r.Re = regexp.MustCompile(r.Pattern)
}

func TaggedOutput(feed wputil.Feed, tags []wpfeed.Tag, sep string, reg []RegexList) ([]byte, error) {
	f := makeTaggedList(feed.List(), tags)

	html := bytes.Buffer{}
	for _, t := range tags {
		for _, i := range f[t.Name].List() {
			html.Write(makeHeader(t.Name))
			html.Write(makePost(i, reg))
			fmt.Fprintf(&html, "\n%s\n\n", sep)
		}
	}
	return html.Bytes(), nil
}

func makeTaggedList(items []wputil.Item, tags wpfeed.Tags) map[string]wputil.Feed {
	tmp := make([]wputil.Feed, len(tags))
	for _, p := range items {
		t1, t2 := priority(p, tags)
		if tags[t1].Limit > 0 {
			if tmp[t1].Len() < int(tags[t1].Limit) {
				tmp[t1].AppendItem(p)
				continue
			} else {
				tmp[t2].AppendItem(p)
				continue
			}
		}
		tmp[t1].AppendItem(p)
	}

	f := make(map[string]wputil.Feed)
	for k, v := range tmp {
		f[tags[k].Name] = v
	}
	return f
}

func priority(i wputil.Item, t wpfeed.Tags) (uint, uint) {
	st := make(wpfeed.Tags, len(t))
	copy(st, t)
	sort.Sort(st)

	var first, second uint
	finder := struct {
		found bool
		pos   int
	}{}
	for k, v := range st {
		if i.HasTag(v.Text) {
			// map sorted priority to original priority
			first = t[k].Priority
			finder.found = true
			finder.pos = k
			break
		}
	}
	for k, v := range st {
		if finder.found && k == finder.pos {
			continue
		}
		if i.HasTag(v.Text) {
			second = t[k].Priority
			break
		}
	}
	return first, second
}

func makeHeader(h string) []byte {
	s := fmt.Sprintf(`<h1 class="section-header">%s</h1>`+"\n", smartenString(h))
	return []byte(s)
}

func makePost(i wputil.Item, re []RegexList) []byte {
	h := i.Body.Text
	for _, r := range re {
		h = r.Re.ReplaceAllString(h, r.Replace)
	}
	s := fmt.Sprintf(`
<h2 class="item-title">
  <a href="%s">%s</a>
</h2>

<div class="body-text">
%s
</div>
`, i.Link, i.Title, makeHTML(h))
	return []byte(s)
}

func smartenString(s string) string {
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
