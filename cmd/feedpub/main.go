package main

import (
	"bytes"
	"feedpub"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"path"
	"regexp"
	"strconv"
	"strings"

	"golang.org/x/net/html"
)

var (
	lg          *log.Logger
	version     string
	flagVersion = flag.Bool("v", false, "print version")

	feedCmd        = flag.NewFlagSet("subcommand feed", flag.ExitOnError)
	flagFeedConfig = feedCmd.String("c", "config.json", "config `file` to use")
	flagFeedFetch  = feedCmd.Bool("fetch", false, "fetch feed data")
	flagFeedMerge  = feedCmd.Bool("merge", false, "build json feeds from raw XML")
	flagFeedJSON   = feedCmd.Bool("json", false, "save current feed items to JSON")
	flagFeedPretty = feedCmd.Bool("pp", false, "pretty print output")
	flagFeedTitles = feedCmd.Bool("titles", false, "print article titles")
	flagFeedTags   = feedCmd.Bool("tags", false, "print feed tags")

	imageCmd        = flag.NewFlagSet("subcommand image", flag.ExitOnError)
	flagImageConfig = imageCmd.String("c", "config.json", "config `file` to use")
	flagImagePretty = imageCmd.Bool("pp", false, "pretty print output")
	flagImageLoud   = imageCmd.Bool("loud", false, "verbose output")
	flagImageExt    = imageCmd.Bool("extract", false, "extract images from feed")
	flagImageFetch  = imageCmd.Bool("fetch", false, "download images")
	flagImageHTML   = imageCmd.Bool("render", false, "render HTML with local image links")

	utilCmd       = flag.NewFlagSet("subcommand util", flag.ExitOnError)
	flagUtilName  = utilCmd.Bool("name", false, "print name")
	flagUtilSeq   = utilCmd.Bool("seq", false, "print sequence number")
	flagUtilRange = utilCmd.Bool("range", false, "print date range")
	flagUtilUni   = utilCmd.Bool("unicode", false, "mark unicode glyphs in LaTeX")
)

func init() {
	flag.Parse()

	if *flagVersion {
		fmt.Printf("%s version: %s\n", os.Args[0], version)
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

	// case "unicode":
	// 	unicodeCmd.Parse(os.Args[2:])
	// 	errs(doUnicodeCmd())

	default:
		subCmdHelp()
	}
}

func subCmdHelp() {
	fmt.Fprintln(os.Stderr, `Available sub-commands:
  feed    - fetch, build and save RSS feeds
  image   - extact and download feed images, render HTML
  unicode - mark unicode glyphs in LaTeX files`)
}

func errs(e error) {
	if e != nil {
		lg.Fatal(e)
	}
}

func doUnicodeCmd() error {
	return nil
}

func doImageCmd() error {
	conf, err := feedpub.ReadConfig(*flagImageConfig)
	if err != nil {
		return err
	}

	if *flagImageExt {
		if err := feedpub.ExtractImages(*conf, *flagImagePretty, lg,
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
			feedpub.ConvertElemIf("iframe", "img", "src", "youtube.com"),
		); err != nil {
			return err
		}
	}

	if *flagImageFetch {
		if err := feedpub.FetchImages(*conf, *flagImageLoud, lg); err != nil {
			return err
		}
	}

	if *flagImageHTML {
		if err := feedpub.ExportHTML(*conf, lg); err != nil {
			return err
		}
	}
	return nil
}

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

func unicodeRegex(in io.Reader, out io.Writer) {

	var b bytes.Buffer
	b.ReadFrom(in)

	re := regexp.MustCompile(`(\p{Cf}+|\p{Co}+)`)
	x := re.ReplaceAllString(b.String(), "")

	re = regexp.MustCompile(`(\p{Devanagari}+)`)
	x = re.ReplaceAllString(x, "{\\sanskrit $1}")

	re = regexp.MustCompile(`(\p{Runic}+)`)
	x = re.ReplaceAllString(x, "{\\runic $1}")

	re = regexp.MustCompile(`(\p{So}+|\p{No}+)`)
	x = re.ReplaceAllString(x, "{\\unisymbol $1}")

	re = regexp.MustCompile(`(\p{Greek}+|\p{Arabic}+|\p{Hebrew}+|\p{Armenian}+|\p{Georgian}+)`)
	x = re.ReplaceAllString(x, "{\\eastern $1}")

	fmt.Fprint(out, x)
}
