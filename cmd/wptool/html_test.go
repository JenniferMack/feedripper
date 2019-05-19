package main

import (
	"bytes"
	"io/ioutil"
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
