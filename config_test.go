package feedpub

import (
	"testing"
	"time"
)

func TestFeedPath(t *testing.T) {
	c := Config{JSONDir: "json"}
	p := c.feedPath("foo", "json")
	if p != "json/foo.json" {
		t.Error(p)
	}
}

func TestConfig(t *testing.T) {
	c, e := ReadConfig("fixtures/config.json")
	if e != nil {
		t.Fatal(e)
	}

	t.Run("name", func(t *testing.T) {
		if c.Names("name") != "tester-en_1" {
			t.Error(c.Names("name"))
		}
	})

	t.Run("image.json", func(t *testing.T) {
		if c.Names("img.json") != "tester-en_1.img.json" {
			t.Error(c.Names("img.json"))
		}
	})

	t.Run("deadline", func(t *testing.T) {
		tm, _ := time.Parse(time.RFC3339, "2019-12-31T18:00:00-05:00")
		if !c.Deadline.Equal(tm) {
			t.Errorf("got %#v, want %#v", c.Deadline, tm)
		}
	})

	t.Run("tags", func(t *testing.T) {
		if len(c.Tags) != 3 {
			t.Error(c.Tags)
		}
	})

	t.Run("feeds", func(t *testing.T) {
		if len(c.Feeds) != 2 {
			t.Error(c.Feeds)
		}
	})

	t.Run("range 1", func(t *testing.T) {
		r := c.DateRange()
		// warning: en dash
		if r != "24–31 Dec 2019" {
			t.Error(r)
		}
	})

	t.Run("range 2", func(t *testing.T) {
		c.Deadline = c.Deadline.AddDate(0, 0, 1)
		r := c.DateRange()
		// warning: en dash
		if r != "25 Dec 2019–01 Jan 2020" {
			t.Error(r)
		}
	})
}
