package main

import (
	"bytes"
	"fmt"
	"io"
	"regexp"
)

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
