package wputil

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestFetch(t *testing.T) {
	s := httptest.NewServer(http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("content-type", "application/json")
			fmt.Fprint(w, `{"foo":"bar"}`)
		}))

	b, e := FetchItem(s.URL, "json")
	if e != nil {
		t.Error(e)
	}

	var i interface{}
	e = json.Unmarshal(b, &i)
	if e != nil {
		t.Error(e)
	}
}

func Test404(t *testing.T) {
	s := httptest.NewServer(http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			http.NotFound(w, r)
		}))

	_, e := FetchItem(s.URL, "json")
	if e.Error() != "status: 404" {
		t.Error(e)
	}
}

func TestCT(t *testing.T) {
	s := httptest.NewServer(http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprint(w, `{"foo":"bar"}`)
		}))

	_, e := FetchItem(s.URL, "json")
	if e.Error() != "content-type: text/plain; charset=utf-8" {
		t.Error(e)
	}
}
