package main

import (
	"feedpub"
	"flag"
	"fmt"
	"os"
)

var (
	utilCmd       = flag.NewFlagSet("subcommand util", flag.ExitOnError)
	flagUtilName  = utilCmd.Bool("name", false, "print project name")
	flagUtilSeq   = utilCmd.Bool("seq", false, "print sequence number")
	flagUtilRange = utilCmd.Bool("range", false, "print date range")
	flagUtilUni   = utilCmd.Bool("unicode", false, "mark unicode glyphs in LaTeX (stdin/stdout)")
)

func doUtilCmd() error {
	conf, err := feedpub.ReadConfig(*flagFeedConfig)
	if err != nil {
		return err
	}

	if *flagUtilUni {
		unicodeRegex(os.Stdin, os.Stdout)
		return nil
	}

	if *flagUtilName {
		fmt.Print(conf.Names("name"))
		return nil
	}

	if *flagUtilSeq {
		fmt.Printf("%s %s", conf.SeqName, conf.Number)
		return nil
	}

	if *flagUtilRange {
		srt := conf.Deadline
		end := srt.AddDate(0, 0, conf.Days)
		if conf.Days < 0 {
			srt, end = end, srt
		}

		fm := "02"
		if srt.Month() < end.Month() {
			fm += " Jan"
		}
		if srt.Year() < end.Year() {
			fm = "02 Jan 2006"
		}

		fmt.Printf("%s–%s", srt.Format(fm), end.Format("02 Jan 2006"))
	}
	return nil
}
