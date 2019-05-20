package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
)

var version string

var flagVers = flag.Bool("v", false, "Print version number")
var flagGlob = flag.String("m", "", "file pattern to match")

func init() {
	flag.Parse()
}
func main() {
	if *flagVers {
		fmt.Printf("%s version: %s\n", os.Args[0], version)
		return
	}

	if *flagGlob == "" {
		return
	}

	x, _ := filepath.Glob(*flagGlob)
	fmt.Println(x)
}
