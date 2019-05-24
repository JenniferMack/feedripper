package wpimage

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"path/filepath"
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
	num := 0
	list := []ImageData{}
	for _, v := range *i {
		n, err := v.fetchImage(d)
		if err != nil {
			v.Err = err.Error()
		}
		list = append(list, v)
		num += n
	}
	*i = ImageList(list)
	return num, nil
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

func (i *ImageData) fetchImage(d string) (int, error) {
	if !i.Valid {
		i.LocalPath = filepath.Join(d, "404.jpg")
		return 0, nil
	}
	if i.Saved {
		return 0, nil
	}
	// do downlaod
	i.Saved = true
	i.LocalPath = filepath.Join(d, filepath.Base(i.Path))
	return 1, nil
}

func (i *ImageData) ParseImageURL(u string) error {
	if i.Resp != 0 || i.Valid {
		return nil
	}

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

func (i *ImageData) CheckImageStatus() (int, error) {
	if i.Resp != 0 || i.Valid {
		return 0, nil
	}

	resp, err := http.Head(i.Path)
	if err != nil {
		i.Err = err.Error()
		return 0, fmt.Errorf("http head: %s", err)
	}

	sc := resp.StatusCode
	if sc < 400 {
		i.Valid = true
	}
	i.Resp = sc
	return 1, nil
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
