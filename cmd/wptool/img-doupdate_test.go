package main

import (
	"bytes"
	"strings"
	"testing"

	"repo.local/wputil/wpimage"
)

func TestDoUpdate(t *testing.T) {
	h := `<p>This is text with a <a href="/images/foo.png"><img src="/images/foo.png"/>link<a> and more text.
<img src="/images/index?123"> some text <img src="/images/index?456"> and done.</p>`
	l := wpimage.ImageList{
		{
			Rawpath:   "/images/index?789",
			LocalPath: "img/404.jpg",
		},
		{
			Rawpath:   "/images/foo.png",
			LocalPath: "img/foo.jpg",
		},
	}
	out := bytes.Buffer{}
	b, err := doUpdate(l, []byte(h), "test file", &out)
	if err != nil {
		t.Fatal(err)
	}
	s := string(b)
	if !strings.Contains(s, `src="img/foo.jpg"`) {
		t.Error(s)
	}
	if strings.Contains(s, `src="/images/foo.png"`) {
		t.Error(s)
	}
	if strings.Contains(s, `?123`) {
		t.Error(s)
	}
	if strings.Contains(s, `?456`) {
		t.Error(s)
	}
}
