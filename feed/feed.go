package feed

import (
	"bytes"
	"log"
	"os"
	"path/filepath"
	"sort"

	wppub "repo.local/wp-pub"
)

// May return an empty slice, and that's ok.
func getOldJSON(p string) []wppub.WPItem {
	if _, err := os.Stat(p); os.IsNotExist(err) {
		return nil
	}

	f, err := os.Open(p)
	if err != nil {
		return nil
	}
	defer f.Close()

	i, err := wppub.ReadWPJSON(f)
	if err != nil {
		return nil
	}
	return i
}

func mergeJSON(bn []byte, dir, name string) error {
	path := filepath.Join(dir, name+".json")
	o := getOldJSON(path)

	n, err := wppub.ReadWPXML(bytes.NewReader(bn))
	if err != nil {
		return err
	}

	m := mergeItems(o, n)

	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()

	x, err := wppub.WriteWPJSON(m, f)
	log.Printf("[%d/%d] %s", x, len(m), path)
	if err != nil {
		return err
	}
	return nil
}

func mergeItems(o, n []wppub.WPItem) []wppub.WPItem {
	o = append(o, n...)
	list := make(map[string]wppub.WPItem)

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
