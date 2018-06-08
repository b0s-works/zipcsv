package main

import (
	"flag"
	"github.com/sirupsen/logrus"
	"zipcsv/config"
	"zipcsv/input"
	"zipcsv/summarizer"
)

func main() {
	cfg := config.New()

	stdinPtr := flag.String("-stdin", "", "Using when you want channel data from Stdin")
	flag.Parse()

	var summaryString summarizer.Result

	if stdinPtr != nil {
		summaryString = input.UseStdInAsDataInput(cfg)
	} else {
		summaryString = input.UseFilesAsDataInput(cfg)
	}
	if len(summaryString) > 0 {
		logrus.Debugf(summaryString.GetSummary())
	}
}
