package feedpub

import (
	"bytes"
	"io/ioutil"
	"log"
	"path"
	"sort"
	"strconv"
	"strings"
	"text/template"

	"github.com/microcosm-cc/bluemonday"
	"github.com/russross/blackfriday"
	"golang.org/x/net/html"
)

const html_tmpl = `
{{- range $cont := . -}}
<h1>{{$cont.Header}}</h1>
{{range $post := $cont.PostList}}
<h2>{{$post.Title}}</h2>
<h3>{{$post.Description}}</h3>
<!-- pubDate: {{$post.PubDate}} -->

<div class="body-text">
{{$post.Body}}
</div>
{{$cont.Sep}}
{{end}}
{{end}}
`

type postData struct {
	Header   string
	Sep      string
	PostList items
}

func ExportHTML(conf Config, lg *log.Logger) error {
	tmpl := template.Must(template.New("html").Parse(html_tmpl))
	itms, n := readItems(conf.Names("json"))
	lg.Printf("Processing %d posts", n)

	p := makeTaggedList(itms, conf)
	buf := bytes.Buffer{}
	err := tmpl.Execute(&buf, p)
	if err != nil {
		return err
	}

	if err := ioutil.WriteFile(conf.Names("html"), buf.Bytes(), 0644); err != nil {
		return err
	}

	lg.Printf("[%s] => %s", sizeOf(buf.Len()), conf.Names("html"))
	return nil
}

func makeTaggedList(i items, conf Config) []postData {
	tagged := sortPosts(i, conf.Tags)

	posts := []postData{}
	for _, v := range conf.Tags {
		sort.Sort(sort.Reverse(tagged[v.Name]))
		p := postData{
			Header:   v.Name,
			Sep:      conf.Separator,
			PostList: tagged[v.Name],
		}
		posts = append(posts, p)
	}
	return posts
}

func sortPosts(i items, t tags) map[string]items {
	// sort by priority
	byPri := make(tags, len(t))
	copy(byPri, t)
	sort.Sort(byPri)

	pl := make([]items, len(t))
	for _, post := range i {
		post = sanitize(post)
		for k, tag := range byPri {
			if tag.Limit > 0 {
				if len(pl[k]) >= int(tag.Limit) {
					continue // next tag
				}
			}
			if post.hasTag(tag.Text) {
				pl[k] = append(pl[k], post)
				break // next post
			}
		}
	}

	list := make(map[string]items, len(t))
	for k, v := range byPri {
		if len(pl[k]) == 0 {
			continue
		}
		list[v.Name] = pl[k]
	}
	return list
}

func sanitize(i item) item {
	str, _ := Parse(strings.NewReader(i.Body),
		func(n *html.Node) {
			if n.Type == html.ElementNode && n.Data == "img" {
				for _, v := range n.Attr {
					if v.Key == "src" && strings.Contains(v.Val, "s.w.org") {
						x := path.Base(v.Val)
						x = strings.TrimSuffix(x, path.Ext(x))
						x = strings.ReplaceAll(x, "-", "")
						emo, _ := strconv.ParseInt(x, 16, 64)
						n.Type = html.TextNode
						n.Data = string(emo)
					}
				}
			}
		},
		ConvertToLink("iframe", "Video Link"),
		RewrapImg("a", "figure"),
		WrapElem("img", "figure"),
		AddCaption("youtube.com", "Video Link"),
		ReplaceElem("h1", "h3"),
		ReplaceElem("h2", "h3"),
		ReplaceElem("h4", "h3"),
		func(n *html.Node) {
			if n.Type == html.TextNode {
				n.Data = smartenString(n.Data)
			}
		},
		func(n *html.Node) {
			for _, v := range i.Images {
				ReplaceAttr("img", "src", v.RawPath, v.LocalPath)(n)
			}
		},
	)

	i.Body = clean(str)
	i.Title = smartenString(i.Title)
	i.Description = smartenString(i.Description)
	return i
}

func smartenString(s string) string {
	re := strings.NewReplacer(
		`“`, `"`,
		`”`, `"`,
		`‘`, `'`,
		`’`, `'`,
		`–`, `—`,
		` – `, `—`,
		` `, ` `,
	)
	s = re.Replace(s)
	sp := blackfriday.NewSmartypantsRenderer(
		blackfriday.Smartypants |
			blackfriday.SmartypantsDashes |
			blackfriday.SmartypantsLatexDashes,
	)

	out := bytes.Buffer{}
	sp.Process(&out, []byte(s))

	return html.UnescapeString(out.String())
}

func clean(h string) string {
	bm := bluemonday.NewPolicy()
	bm.AllowAttrs("href", "title").OnElements("a")
	bm.AllowAttrs("src").OnElements("img")
	bm.AllowElements("ul", "ol", "li", "br", "p", "em",
		"h1", "h2", "h3", "h4", "h5", "h6",
		"blockquote", "strong", "figure", "figcaption")

	return bm.Sanitize(h)
}
