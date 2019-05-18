package wphtml

import (
	"sort"
	"testing"

	"repo.local/wputil"
	"repo.local/wputil/feed"
	wpfeed "repo.local/wputil/feed"
)

func TestHeader(t *testing.T) {
	h := makeHeader("World's News & Events")
	if string(h) != `<h1 class="section-header">World&rsquo;s News & Events</h1>`+"\n" {
		t.Error(string(h))
	}
}

func data() ([]wputil.Item, feed.Tags) {
	i := []wputil.Item{
		{
			Title: "one",
			GUID:  "a",
			Categories: []wputil.Category{
				{Name: "Foo"},
				{Name: "Bar"},
			},
		},
	}
	t := feed.Tags{
		{
			Text:     "bar",
			Name:     "Bar",
			Priority: uint(2),
		},
		{
			Text:     "foo",
			Name:     "Foo",
			Priority: uint(1),
		},
	}
	return i, t
}

func TestPri(t *testing.T) {
	item, tags := data()

	t.Run("bar", func(t *testing.T) {
		n := priority(item[0], tags, "Bar")
		if !n {
			t.Error(n)
		}
	})
	t.Run("foo", func(t *testing.T) {
		n := priority(item[0], tags, "Foo")
		if !n {
			t.Error(n)
		}
	})

	t.Run("copy", func(t *testing.T) {
		cats := make(wpfeed.Tags, len(tags))
		copy(cats, tags)
		sort.Sort(cats)
		if tags[0].Priority != 2 {
			t.Error(tags)
		}
	})
}
