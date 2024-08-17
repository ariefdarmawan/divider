package divider

import "fmt"

func WriteHeader() {
	fmt.Printf("FileDivider v1.0 by ariefdarmawan\n")
	fmt.Printf("File: %s\n", Param.FilePath)
	if Param.SplitByLine > 0 {
		Printer.Printf("Split by line: %d\n", Param.SplitByLine)
	} else {
		Printer.Printf("Split by file count: %d\n", Param.SplitByFileCount)
	}
	fmt.Printf("Send output to: %s\n\n", Param.Output)
}

var (
	headers = []string{}
)
