package feedripper

import (
	"sort"
	"testing"
	"time"
)

func TestAdd(t *testing.T) {
	tm, _ := time.Parse(time.RFC3339, "2019-01-01T12:00:00Z")

	list := items{
		{GUID: "5", PubDate: xmlTime{tm.Add(5 * time.Hour)}},
		{GUID: "4", PubDate: xmlTime{tm.Add(4 * time.Hour)}},
		{GUID: "2", PubDate: xmlTime{tm.Add(2 * time.Hour)}},
		{GUID: "2", PubDate: xmlTime{tm.Add(2*time.Hour + 30*time.Minute)}},
		{GUID: "1", PubDate: xmlTime{tm.Add(1 * time.Hour)}},
		{GUID: "3", PubDate: xmlTime{tm.Add(3 * time.Hour)}},
	}
	list.add(list...)
	if len(list) != 5 {
		t.Error(len(list))
	}

	sort.Sort(list)
	if list[1].GUID != "2" {
		t.Errorf("%+v", list[1].GUID)
	}

	if !list[1].PubDate.Equal(tm.Add((120 + 30) * time.Minute)) {
		t.Errorf("%+v", list[1])
	}
}
