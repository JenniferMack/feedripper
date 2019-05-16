package main

import (
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

func openFileR(s, m string) *os.File {
	if s == "-" {
		return os.Stdin
	}
	f, err := os.Open(s)
	if err != nil {
		log.Fatalf("%s: %s", m, err)
	}
	return f
}

func mergeAndSave(f wputil.Feed, p string) (int, int, error) {
	b, err := ioutil.ReadFile(p)
	if err != nil {
		log.Printf("%s does not exist, skipping", p)
	}

	_, err = f.Write(b)
	if err != nil {
		return 0, 0, fmt.Errorf("load existing: %s", err)
	}

	saved, err := os.Create(p)
	if err != nil {
		return 0, 0, fmt.Errorf("file reset: %s", err)
	}
	defer saved.Close()

	n, err := io.Copy(saved, &f)
	if err != nil {
		return 0, 0, fmt.Errorf("json write: %s", err)
	}
	err = saved.Close()
	if err != nil {
		return 0, 0, fmt.Errorf("json save: %s", err)
	}
	return int(n), f.Len(), nil
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
