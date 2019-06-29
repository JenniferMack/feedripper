package feedripper

import (
	"sort"
	"testing"
)

func Test(t *testing.T) {
	tg := tags{
		{Name: "tag3", Priority: 2},
		{Name: "tag1", Priority: 0},
		{Name: "tag2", Priority: 1},
	}

	t.Run("sort", func(t *testing.T) {
		sort.Sort(tg)
		if tg[0].Name != "tag1" {
			t.Error(tg)
		}
	})

	t.Run("reverse", func(t *testing.T) {
		sort.Sort(sort.Reverse(tg))
		if tg[0].Name != "tag3" {
			t.Error(tg)
		}
	})

	t.Run("range", func(t *testing.T) {
		tg[1].Priority = 3
		if !tg.priOutOfRange() {
			t.Error(tg)
		}
	})
}
