package main

import (
	"bytes"
	"io/ioutil"
	"testing"
)

func TestTagOut(t *testing.T) {
	b, err := ioutil.ReadFile("fixtures/ars.json")
	if err != nil {
		t.Fatal(err)
	}
	fd := bytes.NewBuffer(b)

	b, err = ioutil.ReadFile("fixtures/config.json")
	if err != nil {
		t.Fatal(err)
	}
	cf := bytes.NewBuffer(b)
	wr := bytes.Buffer{}

	err = outputHTMLByTags(fd, cf, nil, &wr)
	if err != nil {
		t.Error(err)
	}
	// wr.WriteTo(os.Stdout)
}
