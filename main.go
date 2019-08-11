package main

import (
	"errors"
	"flag"
	"fmt"
	"go/scanner"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
	"syscall"
)

func main() {
	var exitCode = 0
	var err error
	defer func() {
		if err != nil {
			scanner.PrintError(os.Stderr, err)
			if exitCode == 0 {
				exitCode = 2
			}
		}
		os.Exit(exitCode)
	}()

	write := flag.Bool("w", false, "write result to (source) file instead of stdout")
	local := flag.String("local", "", "put imports beginning with this string after 3rd-party packages; comma-separated list")
	srcdir := flag.String("srcdir", "", "choose imports as if source code is from `dir`. When operating on a single file, dir may instead be the complete file name.")

	flag.Parse()

	if !*write {
		err = errors.New("only write mode is available")
		return
	}

	if !isGoFile(*srcdir) {
		err = fmt.Errorf("srcdir % is not a go file", srcdir)
		return
	}

	err = processFile(*srcdir)
	if err != nil {
		return
	}

	cmd := exec.Command("goimports", "-w", "-local", *local, *srcdir)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err = cmd.Run()
	if err != nil {
		if exiterr, ok := err.(*exec.ExitError); ok {
			if status, ok := exiterr.Sys().(syscall.WaitStatus); ok {
				exitCode = status.ExitStatus()
			}
		}
	}
}

// isGoFile reports whether name is a go file.
func isGoFile(name string) bool {
	fi, err := os.Stat(name)
	return err == nil && fi.Mode().IsRegular() && !strings.HasPrefix(name, ".") && strings.HasSuffix(name, ".go")
}

func processFile(filename string) error {
	f, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer f.Close()

	src, err := ioutil.ReadAll(f)
	if err != nil {
		return err
	}

	content := string(src)

	importStart := strings.Index(content, "import (")
	if importStart >= 0 {
		importEnd := strings.Index(content, ")")

		importsStr := strings.Trim(content[importStart+len("import ("):importEnd], "\n\r")
		importsArr := strings.Split(importsStr, "\n")
		result := ""
		rowsCount := 0
		for _, row := range importsArr {
			if len(strings.TrimSpace(row)) > 0 {
				// if there is a comment inside imports block - no action
				if isCommentRow(row) {
					return nil
				}

				result += row + "\n"
				rowsCount++
			}
		}

		if rowsCount == len(importsArr) {
			return nil
		}

		newContent := content[:importStart] + "import (\n" + result + content[importEnd:]
		err = ioutil.WriteFile(filename, []byte(newContent), 0)
		if err != nil {
			return err
		}

		return err
	}

	return nil
}

func isCommentRow(row string) bool {
	row = strings.TrimSpace(row)

	return row[:2] == "//" || row[:2] == "/*"
}
