package feedpub

import (
	"strings"
	"testing"
)

func TestDropElem(t *testing.T) {
	h := `<p>Text <img src="https://s.w.org/images/core/emoji/2.4/72x72/1f914.png" class="wp-smiley"/><img src="https://s.w.org/images/core/emoji/2.4/72x72/1f915.png" class="wp-smiley"/> more text.</p>`
	r, e := Parse(strings.NewReader(h), DropElemIf("img", "src", "s.w.org"))
	if e != nil {
		t.Fatal(e)
	}

	if r != `<p>Text  more text.</p>` {
		t.Errorf("%s", r)
	}
}

func TestNoClean(t *testing.T) {
	h := `<html><head><title>foo</title></head><body><p class="style">three</p></body></html>`
	r, e := Parse(strings.NewReader(h), ReplaceElem("p", "span"))
	if e != nil {
		t.Fatal(e)
	}

	if r != `<html><head><title>foo</title></head><body><span class="style">three</span></body></html>` {
		t.Errorf("%s", r)
	}
}

func TestClean(t *testing.T) {
	h := `<p class="style">three</p>`
	r, e := Parse(strings.NewReader(h), ReplaceElem("p", "span"))
	if e != nil {
		t.Fatal(e)
	}

	if r != `<span class="style">three</span>` {
		t.Errorf("%s", r)
	}
}

func TestReplaceAttr(t *testing.T) {
	h := `<a href="foo.org">bar</a><a href="foo.com">bar</a>`
	r, e := Parse(strings.NewReader(h), ReplaceAttr("a", "href", "foo.com", "bar.com"))
	if e != nil {
		t.Fatal(e)
	}

	if !strings.Contains(r, `<a href="foo.org">bar</a><a href="bar.com">bar</a>`) {
		t.Errorf("%s", r)
	}
}

func TestReplace(t *testing.T) {
	h := `one two <p class="style">three</p> four`
	r, e := Parse(strings.NewReader(h), ReplaceElem("p", "span"))
	if e != nil {
		t.Fatal(e)
	}

	if !strings.Contains(r, `<span class="style">three</span>`) {
		t.Errorf("%s", r)
	}
}

func TestWrap(t *testing.T) {
	h := `This is <strong>bold</strong> and this is <em>italic</em>, both need spans.`
	r, e := Parse(strings.NewReader(h), WrapElem("em", "span"), WrapElem("strong", "span"))
	if e != nil {
		t.Fatal(e)
	}

	if !strings.Contains(r, `This is <span><strong>bold</strong></span> and this is <span><em>italic</em></span>, both need spans.`) {
		t.Errorf("%s", r)
	}
}

func TestWrapImg(t *testing.T) {
	h := `<a href="/img/foo.jpg"><img src="/img/foo.jpg"/></a>`
	r, e := Parse(strings.NewReader(h), WrapElem("img", "figure"))
	if e != nil {
		t.Fatal(e)
	}

	if !strings.Contains(r, `<a href="/img/foo.jpg"><figure><img src="/img/foo.jpg"/></figure></a>`) {
		t.Errorf("%s", r)
	}
}

func TestUnwrap(t *testing.T) {
	h := `<p><a href="/img/foo.jpg"><img src="/img/foo.jpg"/></a></p>`
	r, e := Parse(strings.NewReader(h), UnwrapElem("img", "a"))
	if e != nil {
		t.Fatal(e)
	}

	if !strings.Contains(r, `<p><img src="/img/foo.jpg"/></p>`) {
		t.Errorf("%s", r)
	}
}

func TestFull(t *testing.T) {
	h := ` <p><iframe src="https://www.youtube.com/embed/BbuhJCIP1xI?feature=oembed"></iframe></p>`
	r, e := Parse(strings.NewReader(h),
		ReplaceElem("iframe", "img"),
		WrapElem("img", "figure"),
		AddCaption("youtube.com", "Video Link"),
	)
	if e != nil {
		t.Fatal(e)
	}

	if !strings.Contains(r, `<figcaption>Video Link</figcaption></figure>`) {
		t.Errorf("%s", r)
	}
}

func TestConvert(t *testing.T) {
	h := `<p><iframe src="https://www.youtube.com/embed/BbUhJVIP1xI?feature=oembed"></iframe></p>\n<p><iframe src="https://www.bootube.com/embed/BbuhJCIP1xI?feature=oembed"></iframe></p>`
	r, e := Parse(strings.NewReader(h),
		ConvertElemIf("iframe", "img", "src", "youtube.com"),
		ConvertToLink("iframe", "Video Link"),
	)
	if e != nil {
		t.Fatal(e)
	}

	if r != `<p><img src="https://www.youtube.com/embed/BbUhJVIP1xI?feature=oembed"/></p>\n<p><a href="https://www.bootube.com/embed/BbuhJCIP1xI?feature=oembed">Video Link</a></p>` {
		t.Errorf("%s", r)
	}
}
