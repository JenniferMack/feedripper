package feed

import (
	"bytes"
	"log"
	"os"
	"path/filepath"
	"sort"
	"time"

	wp "repo.local/wputil"
)

// May return an empty slice, and that's ok.
func getOldJSON(p string) []wp.Item {
	if _, err := os.Stat(p); os.IsNotExist(err) {
		return nil
	}

	f, err := os.Open(p)
	if err != nil {
		return nil
	}
	defer f.Close()

	i, err := wp.ReadWPJSON(f)
	if err != nil {
		return nil
	}
	return i
}

func mergeJSON(bn []byte, dir, name string) error {
	path := filepath.Join(dir, name+".json")
	o := getOldJSON(path)

	n, err := wp.ReadWPXML(bytes.NewReader(bn))
	if err != nil {
		return err
	}

	m := mergeItems(o, n)

	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()

	x, err := wp.WriteWPJSON(m, f)
	log.Printf("[%d/%d] %s", x, len(m), path)
	if err != nil {
		return err
	}
	return nil
}

func mergeItems(o, n []wp.Item) []wp.Item {
	o = append(o, n...)
	list := make(map[string]wp.Item)

	for _, v := range o {
		if p, ok := list[v.GUID]; ok {
			if v.PubDate.After(p.PubDate.Time) {
				list[v.GUID] = v
			}
		} else {
			list[v.GUID] = v
		}
	}

	n = nil
	for _, v := range list {
		n = append(n, v)
	}

	sort.Slice(n, func(i, j int) bool { return n[i].PubDate.After(n[j].PubDate.Time) })
	return n
}

func dropExpired(l []wp.Item, end time.Time, days int) ([]wp.Item, error) {
	start := end.AddDate(0, 0, days)

	list := []wp.Item{}
	for _, v := range l {
		if v.PubDate.Before(end) && v.PubDate.After(start) {
			list = append(list, v)
		}
	}
	return list, nil
}
