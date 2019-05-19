package wpfeed

import "testing"

func TestOORange(t *testing.T) {
	t1 := Tags{
		{Priority: 1},
		{Priority: 2},
	}
	t2 := Tags{
		{Priority: 1},
		{Priority: 0},
	}
	t3 := Tags{
		{Priority: 2},
		{Priority: 1},
	}
	t4 := Tags{
		{Priority: 0},
		{Priority: 1},
		{Priority: 2},
	}
	if !t1.PriOutOfRange() {
		t.Error("out of range")
	}
	if t2.PriOutOfRange() {
		t.Error("out of range")
	}
	if !t3.PriOutOfRange() {
		t.Error("out of range")
	}
	if t4.PriOutOfRange() {
		t.Error("out of range")
	}
}