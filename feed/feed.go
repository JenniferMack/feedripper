package feed

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

func fetch(u string) ([]byte, error) {
	resp, err := http.Get(u)
	if err != nil {
		return nil, fmt.Errorf("getting feed: %s", err)
	}

	b := bytes.Buffer{}

	defer resp.Body.Close()
	n, err := b.ReadFrom(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("reading feed: %s", err)
	}

	log.Printf("[%d] %s", n, u)
	return b.Bytes(), nil
}

func writeRaw(b []byte, dir, name, ext string) error {
	stamp := "_" + string(time.Now().Unix())
	path := filepath.Join(dir, name+stamp+ext)

	f, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("creating file: %s", err)
	}
	defer f.Close()

	_, err = f.Write(b)
	if err != nil {
		return fmt.Errorf("writing file: %s", err)
	}
	return nil
}
