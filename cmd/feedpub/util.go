package main

import (
	"flag"
	"fmt"
	"os"
	"feedripper"
)

var (
	utilCmd       = flag.NewFlagSet("subcommand util", flag.ExitOnError)
	flagUtilName  = utilCmd.Bool("name", false, "print project name")
	flagUtilSeq   = utilCmd.Bool("seq", false, "print sequence number")
	flagUtilRange = utilCmd.Bool("range", false, "print date range")
	flagUtilUni   = utilCmd.Bool("unicode", false, "mark unicode glyphs in LaTeX (stdin/stdout)")
)

func doUtilCmd() error {
	conf, err := feedripper.ReadConfig(*flagFeedConfig)
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
		fmt.Print(conf.DateRange())
	}
	return nil
}
