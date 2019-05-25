package wpimage

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"path/filepath"
	"strings"
)

type ImageList []ImageData

func (i ImageList) SavedNum() int {
	n := 0
	for _, v := range i {
		if v.Saved {
			n += 1
		}
	}
	return n
}

func (i ImageList) ValidNum() int {
	n := 0
	for _, v := range i {
		if v.Valid {
			n += 1
		}
	}
	return n
}

func (i *ImageList) Marshal(out io.Writer) error {
	enc := json.NewEncoder(out)
	enc.SetIndent("", "  ")
	err := enc.Encode(i)
	if err != nil {
		return fmt.Errorf("json encode: %s", err)
	}
	return nil
}

func (i *ImageList) Unmarshal(in io.Reader) error {
	return json.NewDecoder(in).Decode(i)
}

func (i ImageList) Merge(in ImageList) ImageList {
	mer := append(i, in...)
	tmp := make(map[string]ImageData)

	for _, v := range mer {
		if _, ok := tmp[v.Rawpath]; ok {
			if v.Valid || v.Path != "" {
				tmp[v.Rawpath] = v
				continue
			}
		}
		tmp[v.Rawpath] = v
	}

	out := []ImageData{}
	for _, v := range tmp {
		out = append(out, v)
	}
	return ImageList(out)
}

func (i *ImageList) FetchImages(d string) (int, error) {
	// 	num := 0
	// 	list := []ImageData{}
	// 	for _, v := range *i {
	// 		n, err := v.fetchImage(d)
	// 		if err != nil {
	// 			v.Err = err.Error()
	// 		}
	// 		list = append(list, v)
	// 		num += n
	// 	}
	// 	*i = ImageList(list)
	return 0, nil
}

type ImageData struct {
	Rawpath   string
	Path      string
	LocalPath string
	Host      string
	Valid     bool
	Saved     bool
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

	i.Saved = true
	p := filepath.Base(i.Path)
	e := filepath.Ext(p)
	p = strings.TrimSuffix(p, e) + ".jpg"
	i.LocalPath = filepath.Join(d, p)
	return b, nil
}
