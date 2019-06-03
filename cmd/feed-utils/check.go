package main

import (
	"fmt"
	"io"
	"log"
	"strings"
	"text/template"
	"time"

	"repo.local/wputil"
)

func timeFmt(t time.Time, add int) string {
	return t.AddDate(0, 0, add).Format(time.RFC3339)
}

func checkConfig(c io.Reader) string {
	report := strings.Builder{}
	fmt.Fprintln(&report, "Config report...")

	conf, err := wputil.NewConfigList(c)
	if err != nil {
		fmt.Fprintln(&report, "invalid JSON")
		fmt.Fprintln(&report, err)
		return report.String()
	} else {
		fmt.Fprintln(&report, "...valid JSON...")
	}

	tmpl := template.Must(template.New("config").Funcs(
		template.FuncMap{
			"timeFmt": timeFmt},
	).Parse(body))

	for _, v := range conf {
		err := tmpl.Execute(&report, v)
		if err != nil {
			log.Fatal(err)
		}
	}
	fmt.Fprintln(&report, "...done")
	return report.String()
}

const body = `
* Feed Name: {{.Name}}
-------------------------------------------------------------------------------
Number:     {{.Number}}
Language:   {{.Language}}
File base:  {{.Paths "name"}}
Site URL:   {{.SiteURL}}
Deadline:   {{timeFmt .Deadline 0 }}
Date range: {{.Days | printf "%+d"}} days
{{if lt .Days 0 -}}
Deadline range:   {{timeFmt .Deadline .Days}} to {{timeFmt .Deadline 0}}
{{else}}
Deadline range:   {{timeFmt .Deadline 0}} to {{timeFmt .Deadline .Days}}
{{end -}}
Header format:    {{.Range .SeqName}}
XML is saved to:  {{.RSSDir}}
JSON is saved to: {{.JSONDir}}
The collecting tags are:
Num    L P Name
{{range $key, $value := .Tags -}}
{{$key | printf " %2d"}}.  {{$value.Limit | printf "%2d"}} {{$value.Priority}} {{$value.Name}}
{{end -}}
Excluded tags: {{.Exclude | printf "%v"}}
Feed list:
{{range $value := .Feeds}}
    {{- printf "%q	- %q" $value.Name $value.URL}}
{{end -}}
-------------------------------------------------------------------------------
`
