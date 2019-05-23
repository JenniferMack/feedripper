package wpimage

import (
	"testing"
)

func TestILParse(t *testing.T) {
	il := ImageData{}
	il.Parse("http://photobin.com/img/colür/foo.png?width=500&dpi=300")
	if il.Path != "https://photobin.com/img/col%C3%BCr/foo.png" {
		t.Error(il)
	}
}
