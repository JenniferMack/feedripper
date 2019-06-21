package main

import (
	"flag"
	"fmt"
	"log"
	"os"
)

var (
	lg          *log.Logger
	version     string
	flagVersion = flag.Bool("v", false, "print version")
	flagCheck   = flag.Bool("check", false, "config status report")
)

func init() {
	flag.Parse()

	if *flagVersion {
		fmt.Printf("%s version: %s\n", os.Args[0], version)
		os.Exit(0)
	}

	if *flagCheck {
		f := "config.json"
		if flag.Arg(0) != "" {
			f = flag.Arg(0)
		}
		r := checkConfig(f)
		fmt.Println(r)
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

	case "util":
		utilCmd.Parse(os.Args[2:])
		errs(doUtilCmd())

	default:
		subCmdHelp()
	}
}

func subCmdHelp() {
	fmt.Fprintln(os.Stderr, `Available sub-commands:
  feed   - fetch, build and save RSS feeds
  image  - extact and download feed images, render HTML
  util   - utility functions`)
}

func errs(e error) {
	if e != nil {
		lg.Fatal(e)
	}
}
