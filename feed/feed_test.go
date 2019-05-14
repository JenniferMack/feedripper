package feed

import (
	"testing"
	"time"

	wppub "repo.local/wp-pub"
)

func TestMerge(t *testing.T) {
	l := []wppub.WPItem{
		{Title: "one", GUID: "foo"},
		{Title: "two", GUID: "foo"},
		{Title: "three", GUID: "bar"},
	}
	l[0].PubDate.Set(time.Now())
	l[1].PubDate.Set(time.Now().Add(10 * time.Minute))
	l[2].PubDate.Set(time.Now().Add(5 * time.Minute))

	m := mergeItems(nil, l)
	if len(m) != 2 || m[0].Title != "two" {
		t.Errorf("%#v", m)
	}
}

func TestDeadline(t *testing.T) {
	l := []wppub.WPItem{
		{Title: "one", GUID: "foo"},
		{Title: "two", GUID: "bar"},
		{Title: "three", GUID: "baz"},
	}
	l[0].PubDate.Set(time.Now().AddDate(0, 0, -2))
	l[1].PubDate.Set(time.Now().Add(-10 * time.Minute))
	l[2].PubDate.Set(time.Now().Add(-5 * time.Minute))

	d, err := dropExpired(l, time.Now(), -1)
	if err != nil {
		t.Fatal(err)
	}
	if len(d) != 2 {
		t.Errorf("%#v", d)
	}
}
