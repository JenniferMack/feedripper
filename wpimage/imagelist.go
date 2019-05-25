package wpimage

import (
	"encoding/json"
	"fmt"
	"io"
)

type ImageList []ImageData

func (i ImageList) SetDefaults(q int, w uint, tls bool) {
	for k := range i {
		i[k].ImgQual = q
		i[k].ImgWidth = w
		i[k].UseTLS = tls
	}
}

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

func (i ImageList) Marshal(out io.Writer) error {
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
