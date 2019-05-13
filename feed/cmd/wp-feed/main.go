package main

import (
	"flag"
	"log"

	"repo.local/wp-pub/feed"
)

var flagConfig = flag.String("c", "config.json", "config file")

func init() {
	flag.Parse()
}

func main() {
	// f := []feed.Config{
	// 	{
	// 		Name: "foo",
	// 		Feeds: []feed.Feed{
	// 			{Name: "sun", URL: "foo"},
	// 		},
	// 	},
	// }
	// enc := json.NewEncoder(os.Stdout)
	// enc.SetIndent("", "  ")
	// enc.Encode(f)
	err := feed.Run(*flagConfig)
	if err != nil {
		log.Fatal(err)
	}
}
