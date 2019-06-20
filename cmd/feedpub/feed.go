package main

import (
	"feedpub"
	"flag"
	"os"
)

var (
	feedCmd        = flag.NewFlagSet("subcommand feed", flag.ExitOnError)
	flagFeedConfig = feedCmd.String("c", "config.json", "config `file` to use")
	flagFeedFetch  = feedCmd.Bool("fetch", false, "fetch feed data")
	flagFeedMerge  = feedCmd.Bool("merge", false, "build json feeds from raw XML")
	flagFeedJSON   = feedCmd.Bool("json", false, "save current feed items to JSON")
	flagFeedPretty = feedCmd.Bool("pp", false, "pretty print output")
	flagFeedTitles = feedCmd.Bool("titles", false, "print article titles")
	flagFeedTags   = feedCmd.Bool("tags", false, "print feed tags")
)

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

	if *flagFeedTags {
		if err := feedpub.Tags(*conf, os.Stdout); err != nil {
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
