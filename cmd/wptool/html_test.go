package main

import (
	"bytes"
	"io/ioutil"
	"strings"
	"testing"
)

func TestTagOut(t *testing.T) {
	t.Skip("need update to config")

	b, err := ioutil.ReadFile("fixtures/config.json")
	if err != nil {
		t.Fatal(err)
	}
	cf := bytes.NewBuffer(b)
	wr := bytes.Buffer{}

	err = outputHTMLByTags(cf, nil, &wr)
	if err != nil {
		t.Error(err)
	}
	// wr.WriteTo(os.Stdout)
}

func TestRegexYT(t *testing.T) {
	re := regexDefault()
	for k := range re {
		re[k].Compile()
	}
	h := `<p><iframe width="618" height="464" src="https://www.youtube.com/embed/UvT0kzSpWfI?feature=oembed" frameborder="0" allow="accelerometer; autoplay; encrypted-media; gyroscope; picture-in-picture" allowfullscreen></iframe></p>`

	for _, v := range re {
		h = v.ReplaceAll(h)
	}
	if !strings.Contains(h, "UvT0kzSpWfI/default.jpg") {
		t.Error(h)
	}
}
