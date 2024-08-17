package main

import (
	"divider/divider"
	"flag"
	"os"
)

var (
	FilePath         = flag.String("f", "", "file to be divided")
	SplitByLine      = flag.Int("sl", 0, "split by line")
	SplitByFileCount = flag.Int("sc", 0, "split by file count")
	Skip             = flag.Int("skip", 0, "number of line to be skipped")
	Limit            = flag.Int("limit", 0, "number of line to be read, 0 = all")
	KeepHeader       = flag.Int("header", 0, "number of line to keeped as header")
	Output           = flag.String("o", "divider_%d", "output file pattern, ie: file_%d.txt")
)

func main() {
	flag.Parse()
	divider.Param.FilePath = *FilePath
	divider.Param.SplitByLine = *SplitByLine
	divider.Param.SplitByFileCount = *SplitByFileCount
	divider.Param.Output = *Output
	divider.Param.KeepHeader = *KeepHeader
	divider.Param.Skip = *Skip
	divider.Param.Limit = *Limit

	isErr := false
	if *FilePath == "" {
		divider.Logger.Error("filepath is required")
		isErr = true
	}

	if *SplitByFileCount == 0 && *SplitByLine == 0 {
		divider.Logger.Error("either split by line (sl) or split by file count(sc) should have value")
		isErr = true
	}

	if *SplitByFileCount > 0 && *SplitByLine > 0 {
		divider.Logger.Error("can only have one split")
		isErr = true
	}

	if isErr {
		os.Exit(1)
	}

	divider.WriteHeader()
	divider.Divide()
}
