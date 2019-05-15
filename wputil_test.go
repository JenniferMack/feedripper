package wputil

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"
	"testing"
	"time"
)

const s = `[
{ "title": "one", "guid": "a" },
{ "title": "two", "guid": "b" },
{ "title": "three", "guid": "c" }
]`
const s2 = `[
{ "title": "one", "guid": "a", "pub_date": "2019-05-14T11:15:00Z" },
{ "title": "two", "guid": "d", "pub_date": "2019-05-15T11:59:59Z" },
{ "title": "three", "guid": "e", "pub_date": "2019-05-15T12:10:00Z" }
]`

func TestLen(t *testing.T) {
	f, err := ReadWPJSON(strings.NewReader(s2))
	if err != nil {
		t.Error(err)
	}
	if f.Len() != 3 {
		t.Error(f.Len())
	}
}

func TestDates(t *testing.T) {
	d, err := time.Parse(time.RFC3339, "2019-05-15T12:00:00Z")
	if err != nil {
		t.Error(err)
	}

	t.Run("range +", func(t *testing.T) {
		f, err := ReadWPJSON(strings.NewReader(s2))
		if err != nil {
			t.Error(err)
		}
		err = f.Deadline(d, 1)
		if err != nil {
			t.Error(err)
		}
		if f.List()[0].GUID != "e" {
			t.Error(f.String())
		}
	})

	t.Run("range -", func(t *testing.T) {
		f, err := ReadWPJSON(strings.NewReader(s2))
		if err != nil {
			t.Error(err)
		}
		err = f.Deadline(d, -1)
		if err != nil {
			t.Error(err)
		}
		if f.List()[0].GUID != "d" {
			t.Error(f.String())
		}
	})

	t.Run("range 0", func(t *testing.T) {
		f, err := ReadWPJSON(strings.NewReader(s2))
		if err != nil {
			t.Error(err)
		}
		err = f.Deadline(d, 0)
		if err == nil {
			t.Error(err)
		}
	})
}

func TestNil(t *testing.T) {
	f := Feed{}
	_, err := f.Write(nil)
	if err != nil {
		t.Error(err)
	}
}

func TestString(t *testing.T) {
	f := Feed{}
	f.Write([]byte(s))
	str := fmt.Sprint(f)
	if len(str) != 436 {
		t.Error(len(str))
	}
}

func TestInterface(t *testing.T) {
	t.Run("writer", func(t *testing.T) {
		f := Feed{}
		f.Write([]byte(s))
		if len(f.List()) != 3 {
			t.Error(len(f.List()))
		}
		// io.Copy(os.Stdout, &f)
	})

	t.Run("merge items", func(t *testing.T) {
		f := Feed{}
		f.Write([]byte(s))
		f.Merge([]item{{Title: "a string of words", GUID: "foo"}})
		if len(f.List()) != 4 {
			t.Error(len(f.List()))
		}
	})

	t.Run("many writer", func(t *testing.T) {
		f := Feed{}
		f.Write([]byte(s))
		f.Write([]byte(s2))
		if len(f.List()) != 5 {
			t.Error(len(f.List()))
		}
	})

	t.Run("reader", func(t *testing.T) {
		g := Feed{}
		amt := []int{1, 4, 12, 45, 205, 1500}
		for _, v := range amt {
			g.Reset()
			for i := 0; i < v; i++ {
				g.Merge([]item{{Title: "a string of words", GUID: strconv.Itoa(i)}})
			}
			b, err := ioutil.ReadAll(&g)
			if err != nil {
				t.Error(err)
			}
			s := []item{}
			json.Unmarshal(b, &s)
			if len(s) != v {
				t.Error(len(s))
			}
		}
	})
}

func TestCrazyLong(t *testing.T) {
	s := strings.Builder{}
	for i := 0; i < 10000; i++ {
		s.WriteRune('a')
		s.WriteRune('b')
	}
	f := Feed{}
	f.Merge([]item{{Title: "test", Body: body{Text: s.String()}}})
	// s.WriteTo(os.Stdout)
	if strings.Contains(f.String(), "aa") {
		t.Error("dupe a!")
	}
	if strings.Contains(f.String(), "bb") {
		t.Error("dupe b!")
	}
}
