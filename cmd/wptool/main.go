package main

import (
	"flag"
	"fmt"
	"io"
	"os"
)

var version = "foo"
var flagVers = flag.Bool("v", false, "print version number")

var feedCmd = flag.NewFlagSet("feed", flag.ExitOnError)
var flagFeedConfig = feedCmd.String("c", "config.json", "config file location")
var flagFeedFetch = feedCmd.Bool("fetch", false, "retrieve the feeds")
var flagFeedMerge = feedCmd.Bool("merge", false, "merge the feeds")
var flagFeedFormat = feedCmd.Bool("pp", false, "pretty print feeds")

func init() {
	flag.Parse()
}

func main() {
	if len(os.Args) < 2 {
		flag.Usage()
		subCmdHelp()
		return
	}

	if *flagVers {
		fmt.Printf("%s version: %s\n", os.Args[0], version)
		return
	}

	switch os.Args[1] {
	case "feed":
		feedCmd.Parse(os.Args[2:])

		confFile := openFileR(*flagFeedConfig, "feed config")
		defer confFile.Close()

		if *flagFeedFetch {
			errs(getFeeds(confFile, *flagFeedFormat), "fetching")
		}
		if *flagFeedMerge {
			// reset file pointer
			_, err := confFile.Seek(io.SeekStart, 0)
			errs(err, "config seek")
			errs(mergeFeeds(confFile, *flagFeedFormat), "merging")
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
