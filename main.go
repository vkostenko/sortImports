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

var exitCode = 0

func main() {
	defer func() {
		os.Exit(exitCode)
	}()

	write := flag.Bool("w", false, "write result to (source) file instead of stdout")
	local := flag.String("local", "", "put imports beginning with this string after 3rd-party packages; comma-separated list")
	srcdir := flag.String("srcdir", "", "choose imports as if source code is from `dir`. When operating on a single file, dir may instead be the complete file name.")

	flag.Parse()

	if !*write {
		report(errors.New("only write mode is available"))
		return
	}

	if !isGoFile(*srcdir) {
		report(fmt.Errorf("srcdir % is not a go file", srcdir))
		return
	}

	if err := processFile(*srcdir); err != nil {
		report(err)
		return
	}

	cmd := exec.Command("goimports", "-w", "-local", *local, *srcdir)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		fmt.Println(err.Error())
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

func report(err error) {
	scanner.PrintError(os.Stderr, err)
	exitCode = 2
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
