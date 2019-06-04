package wphtml

import (
	"bytes"
	"fmt"
	"regexp"
	"sort"
	"time"

	"repo.local/wputil"
)

type Post struct {
	Title string
	Date  string
	Link  string
	Body  string
}

type RegexList struct {
	Pattern string
	Replace string
	re      *regexp.Regexp
}

func (r *RegexList) Compile() {
	r.re = regexp.MustCompile(r.Pattern)
}

func (r RegexList) ReplaceAll(s string) string {
	return r.re.ReplaceAllString(s, r.Replace)
}

func TaggedOutput(feed wputil.Feed, tags []wputil.Tag, sep string, reg []RegexList) ([]byte, error) {
	f := makeTaggedList(feed.List(), tags)
	htm := bytes.Buffer{}

	for _, t := range tags {
		if f[t.Name].Len() == 0 {
			continue
		}

		posts := formatPosts(f[t.Name].List(), reg)
		fmt.Fprintf(&htm, `<h1 class="section-header">%s</h1>`, t.Name)

		for _, p := range posts {
			fmt.Fprintf(&htm, `

<h2 class="item-title">
  <a href="%s">
  %s
  </a>
</h2>

<div class="body-text">
<!-- pubDate: %s -->
%s</div>

%s
`, p.Link, p.Title, p.Date, p.Body, sep)
		}
	}
	return htm.Bytes(), nil
}

func formatPosts(items []wputil.Item, re []RegexList) []Post {
	list := []Post{}
	for _, v := range items {
		post := Post{
			Title: smartenString(v.Title),
			Link:  v.Link,
			Date:  v.PubDate.Format(time.RFC3339),
			Body:  cleanHTML(v.Body.Text, re),
		}
		list = append(list, post)
	}
	return list
}

func makeTaggedList(items []wputil.Item, tags wputil.Tags) map[string]wputil.Feed {
	// priority sorted copy of tags
	byPri := make(wputil.Tags, len(tags))
	copy(byPri, tags)
	sort.Sort(byPri)

	list := make(map[string][]wputil.Item, len(tags))
	for _, post := range items {
		for _, tag := range byPri {
			if tag.Limit > 0 {
				if len(list[tag.Name]) >= int(tag.Limit) {
					continue // next tag
				}
			}
			if post.HasTag(tag.Text) {
				list[tag.Name] = append(list[tag.Name], post)
				break // next post
			}
		}
	}

	out := make(map[string]wputil.Feed)
	for k, v := range list {
		t := wputil.Feed{}
		t.Merge(v)
		out[k] = t
	}
	return out
}
