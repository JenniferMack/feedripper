package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"

	"repo.local/wputil"
	"repo.local/wputil/wpfeed"
	"repo.local/wputil/wphtml"
)

func regexDefault() []wphtml.RegexList {
	re := []wphtml.RegexList{
		{
			Pattern: "foo",
			Replace: "bar",
		},
	}
	return re
}

func loadFeed(f string) (wputil.Feed, error) {
	var feed wputil.Feed

	j, err := ioutil.ReadFile(f)
	if err != nil {
		return feed, fmt.Errorf("reading %s: %s", f, err)
	}

	feed, err = wputil.ReadWPJSON(bytes.NewReader(j))
	if err != nil {
		return feed, fmt.Errorf("loading feed: %s", err)
	}
	return feed, nil
}

func outputHTMLByTags(c, re io.Reader, w io.Writer) error {
	conf, err := wpfeed.ReadConfig(c)
	if err != nil {
		return fmt.Errorf("loading config: %s", err)
	}

	regex := []wphtml.RegexList{}
	if re == nil {
		regex = regexDefault()
	} else {
		err := json.NewDecoder(re).Decode(&regex)
		if err != nil {
			return fmt.Errorf("loading regex: %s", err)
		}
	}

	for k := range regex {
		regex[k].Compile()
	}

	for _, v := range conf {
		feed, err := loadFeed(v.Name + ".json")
		if err != nil {
			return fmt.Errorf("reading %s: %s", v.Name+".json", err)
		}

		html, err := wphtml.TaggedOutput(feed, v.Tags, "<hr>", regex)
		if err != nil {
			return fmt.Errorf("html: %s", err)
		}

		_, err = w.Write(html)
		if err != nil {
			return fmt.Errorf("write html: %s", err)
		}
	}
	return nil
}

func outputHTML(feed io.Reader, order bool) error {
	return nil
}
