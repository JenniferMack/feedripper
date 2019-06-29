package feedripper

import (
	"testing"
)

func TestLP(t *testing.T) {
	p := `https://foo.bar.org/path/file-d%E2%80%99%C3%A9an.png`
	g := makeLocPath("images", p)
	w := "images/file-d’éan.jpg"
	if g != w {
		t.Error(g)
	}
}

func TestLPYT(t *testing.T) {
	p := `https://img.youtube.com/vi/sDsDsDsD/default.jpg`
	pp := makeLocPath("images", p)
	if pp != `images/sDsDsDsD-default.jpg` {
		t.Error(pp)
	}
}

func TestOnDisk(t *testing.T) {
	p := "hosts"
	if !isOnDisk("/etc/" + p) {
		t.Error(p)
	}
	if isOnDisk("/cte/" + p) {
		t.Error(p)
	}
}

func TestParseRP(t *testing.T) {
	c := Config{
		SiteURL: "foo.com",
		UseTLS:  true,
	}
	t.Run("yt", func(t *testing.T) {
		p := `https://www.youtube.com/embed/qwfpgjluyar?feature=oembed`
		pp, _ := parseRawPath(c, p)
		if pp != `https://img.youtube.com/vi/qwfpgjluyar/default.jpg` {
			t.Error(pp)
		}
	})

	t.Run("url", func(t *testing.T) {
		p := `http://photobin.com/img/colür/foo.png?width=500&dpi=300`
		pp, _ := parseRawPath(c, p)
		if pp != `https://photobin.com/img/col%C3%BCr/foo.png` {
			t.Error(pp)
		}
	})
}
