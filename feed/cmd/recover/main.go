package main

import (
	"flag"
	"fmt"
	"path/filepath"
)

var flagGlob = flag.String("m", "", "file pattern to match")

func init() {
	flag.Parse()
}
func main() {
	if *flagGlob == "" {
		return
	}

	x, _ := filepath.Glob(*flagGlob)
	fmt.Println(x)
}
