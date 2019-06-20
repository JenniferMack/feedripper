package main

import (
	"flag"
	"path"
	"strconv"
	"strings"
	"wputil"

	"golang.org/x/net/html"
)

var (
	imageCmd        = flag.NewFlagSet("subcommand image", flag.ExitOnError)
	flagImageConfig = imageCmd.String("c", "config.json", "config `file` to use")
	flagImagePretty = imageCmd.Bool("pp", false, "pretty print output")
	flagImageLoud   = imageCmd.Bool("loud", false, "verbose output")
	flagImageExt    = imageCmd.Bool("extract", false, "extract images from feed")
	flagImageFetch  = imageCmd.Bool("fetch", false, "download images")
	flagImageHTML   = imageCmd.Bool("render", false, "render HTML with local image links")
)

func doImageCmd() error {
	conf, err := wputil.ReadConfig(*flagImageConfig)
	if err != nil {
		return err
	}

	if *flagImageExt {
		if err := wputil.ExtractImages(*conf, *flagImagePretty, lg,
			func(n *html.Node) {
				if n.Type == html.ElementNode && n.Data == "img" {
					for _, v := range n.Attr {
						if v.Key == "src" && strings.Contains(v.Val, "s.w.org") {
							emo := path.Base(v.Val)
							ext := path.Ext(emo)
							emo = strings.ReplaceAll(strings.TrimSuffix(emo, ext), "-", "")
							ucp, _ := strconv.ParseInt(emo, 16, 64)

							n.Attr = nil
							n.Type = html.TextNode
							n.Data = string(ucp)
							break
						}
					}
				}
			},
			wputil.ConvertElemIf("iframe", "img", "src", "youtube.com"),
		); err != nil {
			return err
		}
	}

	if *flagImageFetch {
		if err := wputil.FetchImages(*conf, *flagImageLoud, lg); err != nil {
			return err
		}
	}

	if *flagImageHTML {
		if err := wputil.ExportHTML(*conf, lg); err != nil {
			return err
		}
	}
	return nil
}
