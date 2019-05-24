package wpimage

import (
	"strings"
	"testing"
)

func TestList(t *testing.T) {
	h := `<p>This is foo and <a href="foo.jpg"><img src="foo.jpg"></a> and more.</p>
	<p>Another image <a href="bar-page.html"><img src="bar.png"></a>.`
	x, e := ParseHTML(strings.NewReader(h))
	if e != nil {
		t.Error(e)
	}
	if len(x) != 2 {
		t.Error(x)
	}
}
