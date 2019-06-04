package wpimage

import (
	"bytes"
	"strings"
	"testing"
)

func TestILParse(t *testing.T) {
	il := ImageData{Rawpath: "http://photobin.com/img/colür/foo.png?width=500&dpi=300", UseTLS: true}
	il.ParseImageURL("foo", "bar", "baz")
	if il.Path != "https://photobin.com/img/col%C3%BCr/foo.png" {
		t.Error(il.Path)
	}
}

func TestYuT(t *testing.T) {
	// <img src="https://img.youtube.com/vi/sDOo5nDJwgA/default.jpg"/>
	d := ImageData{
		Rawpath:   "https://img.youtube.com/vi/sDOo5nDJwgA/default.jpg",
		Path:      "",
		LocalPath: "",
	}
	d.ParseImageURL("foo", "images", "images/404.jpg")
	if d.LocalPath != "images/sDOo5nDJwgA-default.jpg" {
		t.Error(d.LocalPath)
	}
	if d.Path != "https://img.youtube.com/vi/sDOo5nDJwgA/default.jpg" {
		t.Error(d.Path)
	}
}

func TestUnmar(t *testing.T) {
	d := `[
  {
    "Rawpath": "foo",
    "Path": "",
    "Host": "",
    "Valid": false,
    "Resp": 0,
    "Err": null
  },
  {
    "Rawpath": "bar",
    "Path": "",
    "Host": "",
    "Valid": false,
    "Resp": 0,
    "Err": null
  },
  {
    "Rawpath": "baz",
    "Path": "",
    "Host": "",
    "Valid": false,
    "Resp": 0,
    "Err": null
  }
]`

	i := ImageList{}
	dd := strings.NewReader(d)
	err := i.Unmarshal(dd)
	if err != nil {
		t.Error(err)
	}

	t.Run("merge", func(t *testing.T) {
		a := ImageList{}
		dd.Reset(d)
		a.Unmarshal(dd)

		b := ImageList{}
		dd.Reset(d)
		b.Unmarshal(dd)
		for k := range a {
			a[k].Rawpath = string(k + 63)
		}
		cnt1 := b.Merge(a)
		cnt2 := b.Merge(a)
		if len(cnt1) != 6 || len(cnt1) != len(cnt2) {
			t.Error(cnt1, len(a))
			t.Error(cnt2, len(b))
		}
	})

	t.Run("load json", func(t *testing.T) {
		if len(i) != 3 {
			t.Error(i)
		}
	})

	t.Run("dump json", func(t *testing.T) {
		b := bytes.Buffer{}
		err = i.Marshal(&b)
		if err != nil {
			t.Error(err)
		}
		if b.Len() < 100 {
			t.Error(b.Len())
		}
	})
}