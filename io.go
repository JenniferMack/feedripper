package feedpub

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

func ReadConfig(file string) (*Config, error) {
	b, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}

	c := Config{}
	err = json.Unmarshal(b, &c)
	if err != nil {
		return nil, err
	}

	if c.Tags.priOutOfRange() {
		return nil, fmt.Errorf("tag priority out of range")
	}
	return &c, nil
}

func FetchItem(url, typ string) ([]byte, error) {
	resp, err := http.Head(url)
	if err != nil {
		resp.Body.Close()
		return nil, fmt.Errorf("http head: %s", err)
	}

	if resp.StatusCode >= 400 {
		resp.Body.Close()
		return nil, fmt.Errorf("status: %d", resp.StatusCode)
	}

	ht := resp.Header.Get("Content-Type")
	if !strings.Contains(ht, typ) {
		resp.Body.Close()
		return nil, fmt.Errorf("content-type: %s", ht)
	}

	resp, err = http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("http get: %s", err)
	}
	defer resp.Body.Close()

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("body read: %s", err)
	}
	return b, nil
}
