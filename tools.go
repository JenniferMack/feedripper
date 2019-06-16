package feedpub

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

func Titles(conf Config, path string) error {
	b, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}

	feed := items{}
	err = json.Unmarshal(b, &feed)
	if err != nil {
		return err
	}

	for k, v := range feed {
		fmt.Printf("%02d. [%s] %.60s\n", k+1, v.PubDate.Format("0102|15:16"), v.Title)
	}
	return nil
}
