package wpimage

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"path/filepath"
	"strings"
)

type ImageData struct {
	Rawpath   string
	Path      string
	LocalPath string
	Host      string
	Valid     bool
	Saved     bool
	ImgQual   int
	ImgWidth  uint
	Resp      int
	Err       string
}

func (i *ImageData) ParseImageURL(h string, tls bool) int {
	if i.Resp != 0 || i.Valid {
		return 1
	}

	data, err := url.Parse(i.Rawpath)
	if err != nil {
		i.Err = err.Error()
		return 0
	}

	if data.Host == "" {
		data.Host = h
	}
	if data.Scheme == "" {
		data.Scheme = "http"
	}
	if tls && data.Scheme == "http" {
		data.Scheme = "https"
	}

	data.RawQuery = ""
	i.Path = data.String()
	return 1
}

func (i *ImageData) CheckImageStatus() (int, error) {
	if i.Resp != 0 || i.Valid {
		i.Err = ""
		return 0, nil
	}

	resp, err := http.Head(i.Path)
	if err != nil {
		i.Err = err.Error()
		return 1, fmt.Errorf("http head: %s", err)
	}
	defer resp.Body.Close()

	sc := resp.StatusCode
	if sc < 400 {
		i.Valid = true
	}
	i.Resp = sc
	i.Err = ""
	return 0, nil
}

func fetchImageData(u string) ([]byte, error) {
	resp, err := http.Get(u)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (i *ImageData) FetchImage(d string) ([]byte, error) {
	if !i.Valid {
		i.LocalPath = filepath.Join(d, "404.jpg")
		return nil, nil
	}
	if i.Saved {
		return nil, nil
	}
	// do downlaod
	b, err := fetchImageData(i.Path)
	if err != nil {
		return nil, err
	}

	p := filepath.Base(i.Path)
	e := filepath.Ext(p)
	p = strings.TrimSuffix(p, e) + ".jpg"
	i.LocalPath = filepath.Join(d, p)
	return b, nil
}
