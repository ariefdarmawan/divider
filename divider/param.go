package divider

import (
	"github.com/sebarcode/logger"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

type param struct {
	FilePath         string
	SplitByLine      int
	SplitByFileCount int
	Output           string
	Skip             int
	Limit            int
	KeepHeader       int
}

var (
	Param   = new(param)
	Logger  = logger.NewLogEngine(true, false, "", "", "")
	Printer = message.NewPrinter(language.English)
)
