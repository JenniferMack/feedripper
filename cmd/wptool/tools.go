package main

import (
	"bytes"
	"fmt"
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

func mergeAndSave(f wputil.Feed, ind bool, p string) (int, int, error) {
	b, err := ioutil.ReadFile(p)
	if err != nil {
		// a nil []byte is ok to use, don't return
		log.Printf("%s does not exist, skipping", p)
	}

	_, err = f.Write(b)
	if err != nil {
		return 0, 0, fmt.Errorf("load existing: %s", err)
	}

	saved := bytes.Buffer{}
	num := 0
	if ind {
		num, err = saved.WriteString(f.String())
	} else {
		n, e := saved.ReadFrom(&f)
		num, err = int(n), e
	}
	if err != nil {
		return 0, 0, fmt.Errorf("json write: %s", err)
	}

	err = ioutil.WriteFile(p, saved.Bytes(), 0644)
	if err != nil {
		return 0, 0, fmt.Errorf("json save: %s", err)
	}
	return num, f.Len(), nil
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
