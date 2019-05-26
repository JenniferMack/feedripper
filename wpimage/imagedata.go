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
	Rawpath   string `json:"rawpath"`
	Path      string `json:"path"`
	LocalPath string `json:"local_path"`
	Host      string `json:"host"`
	Valid     bool   `json:"valid"`
	Saved     bool   `json:"saved"`
	ImgQual   int    `json:"img_qual"`
	ImgWidth  uint   `json:"img_width"`
	UseTLS    bool   `json:"use_tls"`
	Resp      int    `json:"resp"`
	Err       string `json:"err"`
}

func (i *ImageData) ParseImageURL(h string) int {
	if i.Saved {
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
	if i.UseTLS && data.Scheme == "http" {
		data.Scheme = "https"
	}

	data.RawQuery = ""
	i.Path = data.String()
	i.LocalPath = makeLocalPath("images", "404.jpg")
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

func fetchImageData(u string) ([]byte, int, error) {
	resp, err := http.Get(u)
	if err != nil {
		return nil, 0, err
	}
	defer resp.Body.Close()

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, 0, err
	}
	return data, resp.StatusCode, nil
}

func (i *ImageData) FetchImage(d string) ([]byte, error) {
	if !i.Valid {
		return nil, nil
	}
	if i.Saved {
		return nil, nil
	}
	// do downlaod
	b, c, err := fetchImageData(i.Path)
	if err != nil {
		return nil, err
	}

	if c != 200 {
		i.Resp = c
		return nil, fmt.Errorf("%d: %s", c, filepath.Base(i.Path))
	}

	i.LocalPath = makeLocalPath(d, i.Path)
	return b, nil
}

func makeLocalPath(dir, path string) string {
	p := filepath.Base(path)
	e := filepath.Ext(p)
	p = strings.TrimSuffix(p, e) + ".jpg"
	return filepath.Join(dir, p)
}
