package main

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"repo.local/wputil/wpfeed"
)

var version string

var flagVers = flag.Bool("v", false, "Print version number")
var flagConfig = flag.String("c", "config.json", "Config file to check")

func init() {
	flag.Parse()
}

func main() {
	if *flagVers {
		fmt.Printf("%s version: %s\n", os.Args[0], version)
		return
	}

	f, err := os.Open(*flagConfig)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	fmt.Print(checkConfig(f))
}

func checkConfig(c io.Reader) string {
	report := bytes.Buffer{}
	fmt.Fprintln(&report, "Config report...")

	conf, err := wpfeed.ReadConfig(c)
	if err != nil {
		fmt.Fprintln(&report, "invalid JSON")
		fmt.Fprintln(&report, err)
		fmt.Fprintln(&report, "...done")
		return report.String()
	} else {
		fmt.Fprintln(&report, "...valid JSON...")
	}

	for _, v := range conf {
		fmt.Fprintln(&report,
			"----------------------------------------------------------------------------------------------------")
		fmt.Fprintf(&report, "...checking %q\n", v.Name)
		fmt.Fprintf(&report, "%s’s number: %q\n", v.Name, v.Number)
		fmt.Fprintf(&report, "%s’s deadline: %q\n", v.Name, v.Deadline.Format(time.RFC1123))
		fmt.Fprintf(&report, "%s’s range: %+d days\n", v.Name, v.Days)
		if v.Days > 0 {
			fmt.Fprintf(&report, "%s’s date range: %q to %q\n", v.Name, v.Deadline.Format(time.RFC1123), v.Deadline.AddDate(0, 0, v.Days).Format(time.RFC1123))
		} else {
			fmt.Fprintf(&report, "%s’s date range: %q to %q\n", v.Name, v.Deadline.AddDate(0, 0, v.Days).Format(time.RFC1123), v.Deadline.Format(time.RFC1123))
		}
		fmt.Fprintf(&report, "%s’s working directory is: %q\n", v.Name, v.WorkDir)
		path := filepath.Join(v.WorkDir, v.JSONDir)
		fmt.Fprintf(&report, "%s’s json saved to: %q\n", v.Name, path)
		path = filepath.Join(v.WorkDir, v.RSSDir)
		fmt.Fprintf(&report, "%s’s xml  saved to: %q\n", v.Name, path)
		fmt.Fprintf(&report, "%s’s language is %q\n", v.Name, v.Language)
		fmt.Fprintf(&report, "%s’s site URL is %q\n", v.Name, v.SiteURL)
		fmt.Fprintf(&report, "%s’s is collecting from tags:\n", v.Name)
		for _, v := range v.Tags {
			fmt.Fprintf(&report, "  - %q, limit: %d, priority: %d\n", v.Text, v.Limit, v.Priority)
		}
		fmt.Fprintf(&report, "%s is excluding:\n", v.Name)
		fmt.Fprintf(&report, "%q\n", strings.Join(v.Exclude, ","))

		fmt.Fprintf(&report, "%s’s links:\n", v.Name)
		for _, f := range v.Feeds {
			fmt.Fprintf(&report, "%q, %q\n", f.Name, f.URL)
		}
		fmt.Fprintln(&report,
			"----------------------------------------------------------------------------------------------------")
	}

	fmt.Fprintln(&report, "...done")
	return report.String()
}
