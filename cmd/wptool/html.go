package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"repo.local/wputil"
	"repo.local/wputil/wpfeed"
	"repo.local/wputil/wphtml"
)

func regexDefault() []wphtml.RegexList {
	re := []wphtml.RegexList{
		{Pattern: "<img .+/core/emoji/.+ />",
			Replace: ""},
		{Pattern: `\[caption .+?\]`,
			Replace: "<figure>\n"},
		{Pattern: `> ?(.*)\[/caption\]`,
			Replace: ">\n<figcaption>$1</figcaption>\n</figure>\n"},
		{Pattern: `\[audio (.+)\]`,
			Replace: `<a href="$1">audio link</a>`},
		{Pattern: `\[video ?\](.+)\[/video\]`,
			Replace: `<a href="$1"><em>Video link</em></a>`},
		{Pattern: `<iframe .*src="(.+?)".*?></iframe>`,
			Replace: "\n" + `<a href="$1"><em>Video link</em></a>` + "\n"},
		{Pattern: `(</?h)\d.*?(>)`,
			Replace: `${1}3$2`},
		{Pattern: ` `,
			Replace: ` `}, // non-breaking space literal
		{Pattern: `&nbsp;`,
			Replace: ` `}, // non-breaking space entity
		{Pattern: `‘`,
			Replace: `'`},
		{Pattern: `’`,
			Replace: `'`},
		{Pattern: `“`,
			Replace: `"`},
		{Pattern: `”`,
			Replace: `"`},
		{Pattern: " ?… ?",
			Replace: "…"},
		{Pattern: " ?(–|—|&#8211;|&#8212) ?",
			Replace: "—"},
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
	log.SetFlags(0)
	log.SetPrefix("[    html] ")

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
		log.Printf("> Writing HTML for %s, #%s...", v.Name, v.Number)

		if v.IsWorkDir(os.Getwd()) {
			v.WorkDir = "."
		}
		dirs := filepath.Join(v.WorkDir, v.Name)
		path := fmt.Sprintf(nameFmt, dirs, v.Number, "json")
		feed, err := loadFeed(path)
		if err != nil {
			return fmt.Errorf("reading %s: %s", path, err)
		}

		html, err := wphtml.TaggedOutput(feed, v.Tags, v.Separator, regex)
		if err != nil {
			return fmt.Errorf("html: %s", err)
		}

		path = fmt.Sprintf(nameFmt, dirs, v.Number, "html")
		err = ioutil.WriteFile(path, html, 0644)
		if err != nil {
			return fmt.Errorf("writing %s: %s", path, err)
		}

		log.Printf("> [%s/%d] %s", size(len(html)), 0, path)
	}
	return nil
}

func outputHTML(feed io.Reader, order bool) error {
	return nil
}
