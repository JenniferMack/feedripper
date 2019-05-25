package wpimage

import (
	"bytes"
	"fmt"
	"image"
	"image/jpeg"

	_ "image/gif"
	_ "image/png"

	"github.com/nfnt/resize"
	_ "golang.org/x/image/bmp"
	_ "golang.org/x/image/tiff"
	_ "golang.org/x/image/webp"
)

func MakeJPEG(d []byte, q int, w uint) (b []byte, e error) {
	defer func() {
		if r := recover(); r != nil {
			b, e = nil, fmt.Errorf("jpg panic: %v", r)
		}
	}()

	img, _, err := image.Decode(bytes.NewReader(d))
	if err != nil {
		return nil, fmt.Errorf("img decode: %s", err)
	}

	sized := resize.Resize(w, 0, img, resize.Lanczos3)

	jpg := bytes.Buffer{}
	err = jpeg.Encode(&jpg, sized, &jpeg.Options{Quality: q})
	if err != nil {
		return nil, fmt.Errorf("jpg encode: %s", err)
	}
	return jpg.Bytes(), nil
}
