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
	if !t1.PriorityOutOfRange() {
		t.Error("out of range")
	}
	if t2.PriorityOutOfRange() {
		t.Error("out of range")
	}
	if !t3.PriorityOutOfRange() {
		t.Error("out of range")
	}
}
