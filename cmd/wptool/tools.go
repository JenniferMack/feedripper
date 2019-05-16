package main

import (
	"fmt"
	"log"
	"os"
)

func errs(e error, m string) {
	if e != nil {
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
