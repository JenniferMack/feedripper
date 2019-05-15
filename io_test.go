package wputil

import (
	"bytes"
	"io/ioutil"
	"testing"
)

func TestXMLRead(t *testing.T) {
	b, err := ioutil.ReadFile("fixtures/test.xml")
	if err != nil {
		t.Fatal(err)
	}
	f, err := ReadWPXML(bytes.NewReader(b))
	if err != nil {
		t.Error(err)
	}
	if len(f.List()) != 25 {
		t.Error(len(f.List()))
	}
}

func TestJSONRead(t *testing.T) {
	b, err := ioutil.ReadFile("fixtures/test.json")
	if err != nil {
		t.Fatal(err)
	}
	f, err := ReadWPJSON(bytes.NewReader(b))
	if err != nil {
		t.Error(err)
	}
	if len(f.List()) != 25 {
		t.Error(len(f.List()))
	}
}
