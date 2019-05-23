package wpimage

import (
	"bytes"
	"fmt"
	"image"
	"image/jpeg"
	"io/ioutil"
	"net/http"
	"net/url"

	_ "image/gif"
	_ "image/png"

	_ "golang.org/x/image/bmp"
	_ "golang.org/x/image/tiff"
	_ "golang.org/x/image/webp"
)

type ImageData struct {
	rawpath string
	path    string
	host    string
	valid   bool
	resp    int
	err     error
}

func (i *ImageData) Parse(u string) error {
	data, err := url.Parse(u)
	if err != nil {
		return fmt.Errorf("url parse: %s", err)
	}

	if data.Host == "" {
		data.Host = i.host
	}
	if data.Scheme == "http" || data.Scheme == "" {
		data.Scheme = "https"
	}

	i.rawpath = u
	data.RawQuery = ""
	i.path = data.String()
	return nil
}

func testImagePath(u string) (int, error) {
	resp, err := http.Head(u)
	if err != nil {
		return 0, err
	}
	return resp.StatusCode, err
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

func makeJPEG(d []byte, q int) (b []byte, e error) {
	defer func() {
		if r := recover(); r != nil {
			b, e = nil, fmt.Errorf("decode: %v", r)
		}
	}()

	img, _, err := image.Decode(bytes.NewReader(d))
	if err != nil {
		return nil, err
	}

	jpg := bytes.Buffer{}
	//resize
	err = jpeg.Encode(&jpg, img, &jpeg.Options{Quality: q})
	return jpg.Bytes(), nil
}
