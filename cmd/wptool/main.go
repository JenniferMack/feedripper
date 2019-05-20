package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"os"
)

const nameFmt = `%s-%s.%s`

var version string
var flagVers = flag.Bool("v", false, "print version number")
var flagRegex = flag.Bool("D", false, "print default regex patterns (json)")

var feedCmd = flag.NewFlagSet("feed", flag.ExitOnError)
var flagFeedConfig = feedCmd.String("c", "config.json", "config file location")
var flagFeedFetch = feedCmd.Bool("fetch", false, "retrieve the feeds")
var flagFeedMerge = feedCmd.Bool("merge", false, "merge the feeds")
var flagFeedFormat = feedCmd.Bool("pp", false, "pretty print feeds")
var flagFeedHTML = feedCmd.Bool("html", false, "generate html output")

func init() {
	flag.Parse()
}

func main() {
	if len(os.Args) < 2 {
		flag.Usage()
		subCmdHelp()
		return
	}

	// Top level flags
	if *flagVers {
		fmt.Printf("%s version: %s\n", os.Args[0], version)
		return
	}

	if *flagRegex {
		enc := json.NewEncoder(os.Stdout)
		enc.SetIndent("", "  ")
		enc.SetEscapeHTML(false)
		enc.Encode(regexDefault())
		return
	}

	switch os.Args[1] {
	case "feed":
		feedCmd.Parse(os.Args[2:])
		confFile := openFileR(*flagFeedConfig, "feed config")

		if *flagFeedFetch {
			errs(getFeeds(confFile, *flagFeedFormat), "fetching")
		}
		if *flagFeedMerge {
			// reset file pointer
			_, err := confFile.Seek(0, io.SeekStart)
			errs(err, "config seek")
			errs(mergeFeeds(confFile, *flagFeedFormat), "merging")
		}
		if *flagFeedHTML {
			// reset file pointer
			_, err := confFile.Seek(0, io.SeekStart)
			errs(err, "config seek")
			errs(outputHTMLByTags(confFile, nil, os.Stdout), "html")
		}
		return

	default:
		flag.Usage()
		subCmdHelp()
	}
}

func subCmdHelp() {
	fmt.Println(`Usage of wptool subcommands:
    feed    Download and process RSS feeds

    wptool <cmd> -h for subcommand help
    (Stdin/Stdout is the default for subcommands)`)
}
