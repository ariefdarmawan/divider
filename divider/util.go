package divider

import (
	"bufio"
	"fmt"
	"io/fs"
	"os"
	"sync"
)

var (
	cByLine chan string
	cByFile chan string
)

func getLineCount(fp string) (int, error) {
	f, err := os.Open(fp)
	if err != nil {
		return 0, err
	}
	defer f.Close()

	n := 0
	scanner := bufio.NewScanner(f)
	scanner.Split(bufio.ScanLines)
	for scanner.Scan() {
		n++
	}

	return n, nil
}

func CloseAppError(errTxt string) {
	Logger.Error(errTxt)
	os.Exit(-1)
}

func Divide() {
	cByLine = make(chan string, 1000)
	cByFile = make(chan string, 1000)

	f, err := os.Open(Param.FilePath)
	if err != nil {
		CloseAppError(err.Error())
	}
	defer f.Close()

	wg := new(sync.WaitGroup)
	wg.Add(1)
	if Param.SplitByLine > 0 {
		go processByLine(wg)
	} else {
		go processByFile(wg)
	}

	scanner := bufio.NewScanner(f)
	scanner.Split(bufio.ScanLines)
	nRead := 0
	nSend := 0
	startOfContentIndex := Param.KeepHeader + Param.Skip

loopLine:
	for scanner.Scan() {
		if Param.KeepHeader > 0 && nRead < Param.KeepHeader {
			txt := scanner.Text()
			headers = append(headers, txt)
		}
		if nRead >= startOfContentIndex {
			txt := scanner.Text()
			nSend++
			if Param.SplitByFileCount > 0 {
				cByFile <- txt
			} else {
				cByLine <- txt
			}
		}
		nRead++

		if Param.Limit > 0 && nSend == Param.Limit {
			break loopLine
		}
	}
	close(cByFile)
	close(cByLine)

	wg.Wait()
}

func processByLine(wg *sync.WaitGroup) {
	var (
		f   *os.File
		err error
	)
	defer wg.Done()
	n := 0
	fileIndex := 1
	fileName := fmt.Sprintf(Param.Output, fileIndex)
	read := 0

	for txt := range cByLine {
		n++
		read++
		f, err = appendFile(f, fileName, txt, 0600, false)
		if err != nil {
			CloseAppError(err.Error())
		}

		if n == Param.SplitByLine {
			f.Close()
			f = nil
			Logger.Infof("commit file-%d", fileIndex)
			fileIndex++
			fileName = fmt.Sprintf(Param.Output, fileIndex)
			n = 0
		}
	}

	if Param.SplitByLine > 0 {
		if n != 0 {
			Logger.Infof("commit file-%d", fileIndex)
			f.Close()
		}
		Logger.Info(Printer.Sprintf("done dividing %d lines", read))
	}
}

func processByFile(wg *sync.WaitGroup) {
	defer wg.Done()
	var err error
	n := 0
	fileIndex := 1
	read := 0
	fs := make([]*os.File, Param.SplitByFileCount)

	for txt := range cByFile {
		n++
		read++

		name := ""
		f := fs[fileIndex-1]
		if f == nil {
			name = fmt.Sprintf(Param.Output, fileIndex)
		}
		f, err = appendFile(f, name, txt, 0600, false)
		if err != nil {
			CloseAppError(err.Error())
		}
		fs[fileIndex-1] = f

		fileIndex++
		if fileIndex > Param.SplitByFileCount {
			fileIndex = 1
		}
	}

	for _, f := range fs {
		f.Close()
	}

	if Param.SplitByFileCount > 0 {
		Logger.Info(Printer.Sprintf("done dividing %d lines", read))
	}
}

func appendFile(f *os.File, name, txt string, mode fs.FileMode, closeIt bool) (*os.File, error) {
	var err error
	if f == nil {
		f, err = os.OpenFile(name, os.O_APPEND|os.O_WRONLY|os.O_CREATE, mode)
		if err != nil {
			return nil, err
		}
		if len(headers) > 0 {
			for _, txt := range headers {
				f.WriteString(txt + "\n")
			}
		}
	}
	if closeIt {
		defer f.Close()
	}

	if _, err = f.WriteString(txt + "\n"); err != nil {
		return nil, err
	}

	return f, nil
}
