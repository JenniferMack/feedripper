// Package wputil provides tools for working with WordPress RSS feeds.
package wputil

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"sort"
	"time"
)

type items []Item

// Sort interface

func (i items) Len() int           { return len(i) }
func (i items) Less(j, k int) bool { return i[j].PubDate.Before(i[k].PubDate.Time) }
func (i items) Swap(j, k int)      { i[j], i[k] = i[k], i[j] }

type Feed struct {
	Name  string `json:"name"`
	URL   string `json:"url"`
	items items
	json  []byte
	index int
}

// List returns a slice of Items in the feed.
func (f Feed) List() []Item { return f.items }

// Reverse returns a reversed slice of Items in the feed
func (f Feed) Reverse() []Item { sort.Sort(f.items); return f.items }

// Len returns the number of Feed items.
func (f Feed) Len() int { return f.items.Len() }

// // Append adds items without checking for duplicates or sorting.
// func (f *Feed) Append(i Feed) {
// 	f.items = append(f.items, i.items...)
// }
//
// // AppendItem adds a single item to a feed.
// func (f *Feed) AppendItem(i Item) {
// 	f.items = append(f.items, i)
// }

// Include returns a Feed containing only the posts with given tags
func (f Feed) Include(l []string) Feed {
	var out, inc Feed

	for _, v := range f.items {
		if v.hasTagList(l) {
			inc.items = append(inc.items, v)
		}
	}
	out.Merge(inc.items)
	return out
}

// Exclude returns a Feed with excluded posts removed.
func (f Feed) Exclude(l []string) Feed {
	var out, exc Feed

	for _, v := range f.items {
		if !v.hasTagList(l) {
			exc.items = append(exc.items, v)
		}
	}
	out.Merge(exc.items)
	return out
}

// Deadline removes items that are not within `r` days of date `d`.
// `r` can be either positive or negative.
// If `r` is zero, an error is returned.
func (f *Feed) Deadline(d time.Time, r int) error {
	if r == 0 {
		return fmt.Errorf("unable to use range of %d days", r)
	}

	var start, end time.Time
	e := d.AddDate(0, 0, r)

	if d.Before(e) {
		start, end = d, e
	} else {
		start, end = e, d
	}

	ok := []Item{}
	for _, v := range f.items {
		if v.PubDate.After(start) && v.PubDate.Before(end) {
			ok = append(ok, v)
		}
	}
	f.items = ok
	return nil
}

// String returns the contents of the Feed as tab indented JSON.
func (f Feed) String() string {
	b, err := json.MarshalIndent(f.items, "", "\t")
	if err != nil {
		return fmt.Sprintf("json string: %v", err)
	}
	return string(b)
}

// Merge adds the slice of Item `p` to the feed.
// Duplicates are removed and the internal list is sorted newest first.
func (f *Feed) Merge(n []Item) {
	i := append(f.items, n...)
	list := make(map[string]Item)

	// deduplicate
	for _, v := range i {
		if p, ok := list[v.GUID]; ok {
			if v.PubDate.After(p.PubDate.Time) {
				list[v.GUID] = v
			}
		} else {
			list[v.GUID] = v
		}
	}

	f.items = nil
	for _, v := range list {
		f.items = append(f.items, v)
	}
	// default sort is newest first
	sort.Sort(sort.Reverse(f.items))
}

// Write appends the contents of `p` (JSON encoded slice of `Item`) to the Feed.
// Duplicates are removed and the internal list is sorted newest first.
func (f *Feed) Write(p []byte) (n int, err error) {
	if len(p) == 0 {
		return 0, nil
	}
	t := []Item{}

	err = json.NewDecoder(bytes.NewReader(p)).Decode(&t)
	if err != nil {
		return 0, fmt.Errorf("json write: %v", err)
	}

	f.Merge(t)
	return len(p), nil
}

// Read reads the next len(p) bytes from the buffer or until the Feed is drained.
// The return value n is the number of bytes read. If the buffer has no data to return,
// err is io.EOF (unless len(p) is zero); otherwise it is nil.
func (f *Feed) Read(p []byte) (n int, err error) {
	if len(p) == 0 {
		return 0, nil
	}

	if f.json == nil {
		b, err := json.Marshal(f.items)
		if err != nil {
			return 0, fmt.Errorf("json read: %v", err)
		}
		f.json = b
	}

	if len(f.json)-f.index > len(p) {
		n = copy(p, f.json[f.index:f.index+len(p)])
		f.index += len(p)
		return n, nil
	}

	n = copy(p, f.json[f.index:])
	f.reset()
	return n, io.EOF
}

func (f *Feed) reset() {
	f.index = 0
	f.json = nil
}

// Reset resets the feed to be empty,
// but it retains the underlying storage for use by future writes.
func (f *Feed) Reset() {
	f.reset()
	f.items = f.items[:0]
}
