package wputil

import (
	"encoding/json"
	"fmt"
	"io"
	"sort"
)

type Feed struct {
	items []item
	json  []byte
	index int
}

// List returns a slice of Items in the Feed.
func (f Feed) List() []item {
	return f.items
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
func (f *Feed) Merge(n []item) {
	i := append(f.items, n...)
	list := make(map[string]item)

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

	sort.Slice(f.items, func(i, j int) bool {
		return f.items[i].PubDate.After(f.items[j].PubDate.Time)
	})
}

// Write appends the contents of `p` (JSON encoded slice of `Item`) to the Feed.
// Duplicates are removed and the internal list is sorted newest first.
func (f *Feed) Write(p []byte) (n int, err error) {
	t := []item{}

	err = json.Unmarshal(p, &t)
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

// Reset resets the Feed to be empty,
// but it retains the underlying storage for use by future writes.
func (f *Feed) Reset() {
	f.reset()
	f.items = f.items[:0]
}
