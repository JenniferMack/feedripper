package main

import (
	"feedpub"
	"flag"
	"fmt"
	"log"
	"os"
)

var (
	lg          *log.Logger
	flagVersion = flag.Bool("v", false, "print version")

	unicodeCmd  = flag.NewFlagSet("subcommand unicode", flag.ExitOnError)
	flagUnifont = unicodeCmd.String("tex", "unifont", "LaTex `command` for Unicode font")
	flagUniFile = unicodeCmd.String("f", "", "LaTeX `file` to use")

	feedCmd        = flag.NewFlagSet("subcommand feed", flag.ExitOnError)
	flagFeedConfig = feedCmd.String("c", "config.json", "config `file` to use")
	flagFeedFetch  = feedCmd.Bool("fetch", false, "fetch feed data")
	flagFeedMerge  = feedCmd.Bool("merge", false, "build json feeds from raw XML")
	flagFeedJSON   = feedCmd.Bool("json", false, "save current feed items to JSON")
	flagFeedPretty = feedCmd.Bool("pp", false, "pretty print output")
	flagFeedTitles = feedCmd.Bool("titles", false, "print article titles")

	imageCmd        = flag.NewFlagSet("subcommand image", flag.ExitOnError)
	flagImageConfig = imageCmd.String("c", "config.json", "config `file` to use")
	flagImagePretty = imageCmd.Bool("pp", false, "pretty print output")
	flagImageLoud   = imageCmd.Bool("loud", false, "verbose output")
	flagImageExt    = imageCmd.Bool("extract", false, "extract images from feed")
	flagImageFetch  = imageCmd.Bool("fetch", false, "download images")
	flagImageHTML   = imageCmd.Bool("render", false, "render HTML with local image links")
)

func init() {
	flag.Parse()

	if *flagVersion {
		fmt.Printf("%s version: %s", os.Args[0], "[vers]")
		os.Exit(0)
	}
}

func main() {
	if len(os.Args) < 2 {
		flag.Usage()
		subCmdHelp()
		os.Exit(1)
	}

	name := fmt.Sprintf("[%-8s] ", os.Args[0])
	lg = log.New(os.Stderr, name, 0)

	switch os.Args[1] {
	case "feed":
		feedCmd.Parse(os.Args[2:])
		errs(doFeedCmd())

	case "image":
		imageCmd.Parse(os.Args[2:])
		errs(doImageCmd())

	case "unicode":
		unicodeCmd.Parse(os.Args[2:])
		errs(doUnicodeCmd())

	default:
		subCmdHelp()
	}
}

func subCmdHelp() {
	fmt.Fprintln(os.Stderr, `Available sub-commands:
  feed    - fetch, build and save RSS feeds
  image   - extact and download feed images, render HTML
  unicode - mark unicode glyphs in LaTeX files`)
}

func errs(e error) {
	if e != nil {
		lg.Fatal(e)
	}
}

func doUnicodeCmd() error {
	return nil
}

func doImageCmd() error {
	conf, err := feedpub.ReadConfig(*flagImageConfig)
	if err != nil {
		return err
	}

	if *flagImageExt {
		if err := feedpub.ExtractImages(*conf, *flagImagePretty, lg); err != nil {
			return err
		}
	}

	if *flagImageFetch {
		if err := feedpub.FetchImages(*conf, *flagImageLoud, lg); err != nil {
			return err
		}
	}
	return nil
}

func doFeedCmd() error {
	conf, err := feedpub.ReadConfig(*flagFeedConfig)
	if err != nil {
		return err
	}

	if *flagFeedTitles {
		if err := feedpub.Titles(*conf, os.Stdout); err != nil {
			return err
		}
		return nil
	}

	if *flagFeedFetch {
		if err := feedpub.FetchFeeds(*conf, lg); err != nil {
			return err
		}
	}

	if *flagFeedMerge {
		if err := feedpub.BuildFeeds(*conf, *flagFeedPretty, lg); err != nil {
			return err
		}
	}

	if *flagFeedJSON {
		if err := feedpub.WriteItemList(*conf, *flagFeedPretty, lg); err != nil {
			return err
		}
	}
	return nil
}
