package wphtml

import (
	"testing"
)

func TestHeader(t *testing.T) {
	h := makeHeader("World's News & Events")
	if string(h) != `<h1 class="section-header">World&rsquo;s News & Events</h1>`+"\n" {
		t.Error(string(h))
	}
}
