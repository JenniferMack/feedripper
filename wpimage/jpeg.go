package wpimage

import (
	"bytes"
	"fmt"
	"image"
	"image/jpeg"

	_ "image/gif"
	_ "image/png"

	_ "golang.org/x/image/bmp"
	_ "golang.org/x/image/tiff"
	_ "golang.org/x/image/webp"
)

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
