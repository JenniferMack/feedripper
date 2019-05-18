package main

import (
	"encoding/json"
	"fmt"
	"io"

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

func outputHTMLByTags(f, c, re io.Reader, w io.Writer) error {
	feed, err := wputil.ReadWPJSON(f)
	if err != nil {
		return fmt.Errorf("loading feed: %s", err)
	}
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
		html, err := wphtml.TaggedOutput(feed, v.Tags, "<br/>", regex)
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
