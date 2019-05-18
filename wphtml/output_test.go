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
		{
			Title: "two",
			GUID:  "b",
			Categories: []wputil.Category{
				{Name: "Moo"},
				{Name: "Bar"},
			},
		},
		{
			Title: "four",
			GUID:  "d",
			Categories: []wputil.Category{
				{Name: "Foo"},
				{Name: "Baz"},
			},
		},
		{
			Title: "three",
			GUID:  "c",
			Categories: []wputil.Category{
				{Name: "Foo"},
				{Name: "Bar"},
				{Name: "Baz"},
			},
		},
		{
			Title: "five",
			GUID:  "e",
			Categories: []wputil.Category{
				{Name: "Bar"},
				{Name: "Miz"},
			},
		},
	}
	t := feed.Tags{
		{
			Text:     "bar",
			Name:     "Bar",
			Priority: uint(1),
			Limit:    uint(0),
		},
		{
			Text:     "foo",
			Name:     "Foo",
			Priority: uint(0),
			Limit:    uint(0),
		},
	}
	return i, t
}

func TestPri(t *testing.T) {
	item, tags := data()

	t.Run("bar", func(t *testing.T) {
		n := priority(item[0], tags)
		if n != 1 {
			t.Error(n)
		}
	})
	t.Run("foo", func(t *testing.T) {
		n := priority(item[1], tags)
		if n != 0 {
			t.Error(n)
		}
	})

	t.Run("copy", func(t *testing.T) {
		cats := make(wpfeed.Tags, len(tags))
		copy(cats, tags)
		sort.Sort(cats)
		if tags[0].Priority != 1 {
			t.Error(cats, tags)
		}
	})
}

func TestList(t *testing.T) {
	items, tags := data()

	t.Run("len tags", func(t *testing.T) {
		o := makeTaggedList(items, tags)
		if len(o) != 2 {
			t.Error(len(o))
		}
	})

	t.Run("foo feed", func(t *testing.T) {
		o := makeTaggedList(items, tags)
		n := 0
		for _, v := range o {
			n += v.Len()
		}
		if n != 5 {
			t.Error("total: ", n)
		}
	})

	t.Run("foo feed", func(t *testing.T) {
		o := makeTaggedList(items, tags)
		if o["Foo"].Len() != 3 {
			t.Errorf("tag: %s / len %d", "Foo", o["Foo"].Len())
		}
	})

	t.Run("bar feed", func(t *testing.T) {
		o := makeTaggedList(items, tags)
		if o["Bar"].Len() != 2 {
			t.Errorf("tag: %s / len %d", "Bar", o["Bar"].Len())
		}
	})
}