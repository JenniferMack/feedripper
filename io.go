package feedpub

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

func ReadConfig(file string) (*Config, error) {
	b, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}

	c := Config{}
	err = json.Unmarshal(b, &c)
	if err != nil {
		return nil, err
	}

	if c.Tags.priOutOfRange() {
		return nil, fmt.Errorf("tag priority out of range")
	}
	return &c, nil
}
