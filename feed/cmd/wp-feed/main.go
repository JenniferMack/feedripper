package main

import (
	"flag"
	"log"
	"os"

	"repo.local/wp-pub/feed"
)

var flagConfig = flag.String("c", "config.json", "config file")

func init() {
	flag.Parse()
}

func main() {
	f, err := os.Open(*flagConfig)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	err = feed.Run(f)
	if err != nil {
		log.Fatal(err)
	}
}
