package feedpub

import "testing"

func TestLP(t *testing.T) {
	p := `https://foo.bar.org/path/file.png`
	p = makeLocPath(p)
	if p != "file.jpg" {
		t.Error(p)
	}
}

func TestOnDisk(t *testing.T) {
	p := "hosts"
	if !isOnDisk("/etc/", p) {
		t.Error(p)
	}
	if isOnDisk("/cte/", p) {
		t.Error(p)
	}
}
