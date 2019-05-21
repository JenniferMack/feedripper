package main

import (
	"flag"
	"fmt"
	"log"
	"os"
)

var version string

var flagVers = flag.Bool("v", false, "Print version number")
var flagConfig = flag.String("c", "config.json", "Config file to check")

func init() {
	flag.Parse()
}

func main() {
	if *flagVers {
		fmt.Printf("%s version: %s\n", os.Args[0], version)
		return
	}

	f, err := os.Open(*flagConfig)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	fmt.Print(checkConfig(f))
}
