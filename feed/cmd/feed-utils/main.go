package main

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"time"

	"repo.local/wputil/feed"
)

var flagConfig = flag.String("c", "config.json", "config file")
var flagCheckConfig = flag.Bool("check", false, "print config file report")

func init() {
	flag.Parse()
}

func main() {
	f, err := os.Open(*flagConfig)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	if *flagCheckConfig {
		fmt.Print(checkConfig(f))
		return
	}
}

func checkConfig(c io.Reader) string {
	report := bytes.Buffer{}
	fmt.Fprintln(&report, "Config report...")

	conf, err := feed.ReadConfig(c)
	if err != nil {
		fmt.Fprintln(&report, "invalid JSON")
		fmt.Fprintln(&report, "...done")
		return report.String()
	} else {
		fmt.Fprintln(&report, "...valid JSON...")
	}

	for _, v := range conf {
		fmt.Fprintln(&report, "--------------------------------------------------")
		fmt.Fprintf(&report, "...checking %q\n", v.Name)
		fmt.Fprintf(&report, "%s's number: %q\n", v.Name, v.Number)
		fmt.Fprintf(&report, "%s's deadline: %q\n", v.Name, v.Deadline.Format(time.RFC1123))
		fmt.Fprintf(&report, "%s's range: %d days\n", v.Name, v.Days)
		if v.Days > 0 {
			fmt.Fprintf(&report, "%s's date range: %q to %q\n", v.Name, v.Deadline.Format(time.RFC1123), v.Deadline.AddDate(0, 0, v.Days).Format(time.RFC1123))
		} else {
			fmt.Fprintf(&report, "%s's date range: %q to %q\n", v.Name, v.Deadline.AddDate(0, 0, v.Days).Format(time.RFC1123), v.Deadline.Format(time.RFC1123))
		}
		fmt.Fprintf(&report, "%s's json saved to: %q\n", v.Name, v.JSONDir)
		fmt.Fprintf(&report, "%s's xml  saved to: %q\n", v.Name, v.RSSDir)
		fmt.Fprintf(&report, "%s's language is %q\n", v.Name, v.Language)
		fmt.Fprintf(&report, "%s's site URL is %q\n", v.Name, v.SiteURL)
		fmt.Fprintf(&report, "%s's is collecting from %d tags:\n", v.Name, v.MainTagNum)
		fmt.Fprintf(&report, "%q\n", strings.Join(v.Tags, ","))
		fmt.Fprintf(&report, "%s is excluding:\n", v.Name)
		fmt.Fprintf(&report, "%q\n", strings.Join(v.Exclude, ","))

		fmt.Fprintf(&report, "%s's links:\n", v.Name)
		for _, f := range v.Feeds {
			fmt.Fprintf(&report, "%q, %q, %q\n", f.Name, f.URL, f.Type)
		}
		fmt.Fprintln(&report, "--------------------------------------------------")
	}

	fmt.Fprintln(&report, "...done")
	return report.String()
}
