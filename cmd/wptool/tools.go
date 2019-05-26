package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"

	"repo.local/wputil"
)

func errs(e error, m string) {
	if e != nil {
		log.SetPrefix("[   error] ")
		log.Fatalf("%s: %s", m, e)
	}
}

func openFileR(s, m string) io.ReadSeeker {
	if s == "-" {
		return os.Stdin
	}
	b, err := ioutil.ReadFile(s)
	errs(err, m)

	r := bytes.NewReader(b)
	return r
}

func mergeAndSave(f wputil.Feed, pretty bool, path string) (int, int, error) {
	b, err := ioutil.ReadFile(path)
	if err != nil {
		// a nil []byte is ok to use, don't return
		log.Printf("%s does not exist, skipping", path)
	}

	_, err = f.Write(b)
	if err != nil {
		return 0, 0, fmt.Errorf("load existing: %s", err)
	}

	b, err = formatFeed(f, pretty)
	if err != nil {
		return 0, 0, fmt.Errorf("json format: %s", err)
	}

	err = ioutil.WriteFile(path, b, 0644)
	if err != nil {
		return 0, 0, fmt.Errorf("json save: %s", err)
	}
	return len(b), f.Len(), nil
}

func formatFeed(f wputil.Feed, pp bool) ([]byte, error) {
	if f.Len() == 0 {
		return nil, fmt.Errorf("no items to merge")
	}

	b := bytes.Buffer{}
	enc := json.NewEncoder(&b)
	if pp {
		enc.SetIndent("", "\t")
		enc.SetEscapeHTML(false)
	}
	err := enc.Encode(f.List())
	if err != nil {
		return nil, err
	}
	return b.Bytes(), nil
}

func size(b int) string {
	const unit = 1024
	if b < unit {
		return fmt.Sprintf("%dB", b)
	}

	div, exp := int64(unit), 0
	for n := b / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f%c", float64(b)/float64(div), "KMG"[exp])
}

func trim(l int, s string) string {
	if len(s) <= l {
		return s
	}
	cut := len(s) - l
	return "..." + s[cut+3:]
}
