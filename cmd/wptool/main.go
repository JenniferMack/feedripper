package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"os"
)

const nameFmt = `%s-%s.%s`

var (
	version   string
	flagVers  = flag.Bool("v", false, "print version number")
	flagRegex = flag.Bool("D", false, "print default regex patterns (json)")

	feedCmd        = flag.NewFlagSet("feed", flag.ExitOnError)
	flagFeedConfig = feedCmd.String("c", "config.json", "config file location")
	flagFeedFetch  = feedCmd.Bool("fetch", false, "retrieve the feeds")
	flagFeedMerge  = feedCmd.Bool("merge", false, "merge the feeds")
	flagFeedFormat = feedCmd.Bool("pp", false, "pretty print feeds")
	flagFeedHTML   = feedCmd.Bool("html", false, "generate html output")

	imageCmd         = flag.NewFlagSet("image", flag.ExitOnError)
	flagImageConfig  = imageCmd.String("c", "config.json", "config file location")
	flagImageParse   = imageCmd.Bool("parse", false, "parse HTML for images")
	flagImageFilter  = imageCmd.Bool("filter", false, "filter image URLs")
	flagImageVerify  = imageCmd.Bool("verify", false, "verify images are downloadable")
	flagImageFetch   = imageCmd.Bool("fetch", false, "fetch images")
	flagImageVerbose = imageCmd.Bool("v", false, "prints status of each download")
)

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

	case "image":
		imageCmd.Parse(os.Args[2:])
		confFile := openFileR(*flagImageConfig, "image config")

		actions := []string{}
		if *flagImageParse {
			actions = append(actions, "parse")
		}
		if *flagImageFilter {
			actions = append(actions, "filter")
		}
		if *flagImageVerify {
			actions = append(actions, "verify")
		}
		if *flagImageFetch {
			actions = append(actions, "fetch")
		}
		errs(images(confFile, actions), "image")
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
