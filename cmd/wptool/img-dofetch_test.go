package main

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"repo.local/wputil/wpfeed"
	"repo.local/wputil/wpimage"
)

func TestDoFetch404(t *testing.T) {
	buf := bytes.Buffer{}
	s := httptest.NewServer(http.NotFoundHandler())
	list := wpimage.ImageList{
		{
			Path:  s.URL + "/img/test.png",
			Valid: true,
		},
	}

	fi := make(wpfeed.Paths)
	fi["one"] = "filepath.json"
	fi["two"] = "dirpath"

	b, e := doFetch(list, fi, &buf)
	if b == nil {
		t.Error(b)
	}
	if e != nil {
		t.Error(e)
	}
	if !bytes.Contains(buf.Bytes(), []byte("[error] fetching: 404: test.png")) {
		t.Error(buf.String())
	}
}

func TestDoFetchImg(t *testing.T) {
	buf := bytes.Buffer{}
	s := httptest.NewServer(http.HandlerFunc(tester(t)))
	list := wpimage.ImageList{
		{
			Path:  s.URL + "/img/test.png",
			Valid: true,
		},
		{
			Path:  s.URL + "/img/foo.png",
			Valid: true,
			Saved: true,
		},
	}

	fi := make(wpfeed.Paths)
	fi["one"] = "filepath.json"
	fi["two"] = "dirpath"

	b, e := doFetch(list, fi, &buf)
	if b == nil {
		t.Error(b)
	}
	if e != nil {
		t.Error(e)
	}
	t.Run("gif err", func(t *testing.T) {
		if !bytes.Contains(buf.Bytes(), []byte("unknown block type: 0x00")) {
			t.Error(buf.String())
		}
	})
	t.Run("errors", func(t *testing.T) {
		if !bytes.Contains(buf.Bytes(), []byte("1 error,")) {
			t.Error(buf.String())
		}
	})
	t.Run("skipped", func(t *testing.T) {
		if !bytes.Contains(buf.Bytes(), []byte("1/2 downloaded, 1 error, 1 prev. saved")) {
			t.Error(buf.String())
		}
	})
}

func tester(t *testing.T) func(http.ResponseWriter, *http.Request) {
	t.Helper()
	b, e := ioutil.ReadFile("fixtures/bad.gif")
	if e != nil {
		t.Fatal(e)
	}
	return func(w http.ResponseWriter, r *http.Request) {
		w.Write(b)
	}
}
