package wpimage

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
)

type ImageData struct {
	Rawpath string
	Path    string
	Host    string
	Valid   bool
	Resp    int
	Err     error
}

func (i *ImageData) Parse(u string) error {
	data, err := url.Parse(u)
	if err != nil {
		return fmt.Errorf("url parse: %s", err)
	}

	if data.Host == "" {
		data.Host = i.Host
	}
	if data.Scheme == "http" || data.Scheme == "" {
		data.Scheme = "https"
	}

	i.Rawpath = u
	data.RawQuery = ""
	i.Path = data.String()
	return nil
}

func (i *ImageData) CheckImageStatus() error {
	resp, err := http.Head(i.Path)
	if err != nil {
		i.Err = err
		return err
	}

	sc := resp.StatusCode
	if sc < 400 {
		i.Valid = true
	}
	i.Resp = sc
	return nil
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
