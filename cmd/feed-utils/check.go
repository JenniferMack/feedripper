package main

import (
	"bytes"
	"fmt"
	"io"
	"strings"
	"time"

	"repo.local/wputil"
)

func checkConfig(c io.Reader) string {
	report := bytes.Buffer{}
	fmt.Fprintln(&report, "Config report...")

	conf, err := wputil.NewConfigList(c)
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
		fmt.Fprintf(&report, "%s’s JSON saved to: %q\n", v.Paths("name"), v.Paths("json-dir"))
		fmt.Fprintf(&report, "%s’s XML  saved to: %q\n", v.Paths("name"), v.Paths("rss-dir"))
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
