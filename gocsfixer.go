package main

import (
	"github.com/ewgRa/gocsfixer/fixers"
	"fmt"
	"io/ioutil"
	"os"
)

// gocsfixer --recommend --lint --fix

// Use case 1: pure linter
// Use case 2: pure fixer
// Use case 3: fixer with lint output

func main() {
	returnValue := 0

	recommend := true
	lint := true
	fix := false

	configs := readConfig()

	for _, file := range getFiles() {
		fmt.Println("File", file)
		c, err := ioutil.ReadFile(file)
		content := string(c)

		if nil != err  {
			panic("Error reading file") // FIXME XXX: better message an react
		}

		for _, config := range configs {
			if fix {
				if config.Fix() {
					fixer, ok := config.csFixer.(fixers.Fixer)

					if !ok {
						panic("It is not a fixer") // FIXME XXX: better message an react
					}

					fmt.Println("fixer", fixer)
				}
			} else if recommend || lint {
				linter, ok := config.csFixer.(fixers.Linter)

				if !ok {
					panic("It is not a linter, can't lint or recommend") // FIXME XXX: better message an react
				}

				lintMode := lint && config.Lint()
				if (lintMode || recommend && config.Recommend()) {
					problems, err := linter.Lint(content)

					if nil != err {
						panic("error during lint") // FIXME XXX: better message an react
					}

					if lintMode && len(problems) != 0 {
						returnValue = 1
					}

					for _, problem := range problems {
						if lintMode {
							fmt.Println("    error", problem)
						} else {
							fmt.Println("    recommendation", problem)
						}
					}
				}
			}
		}
	}

	os.Exit(returnValue)
}

// Read files for processing
func getFiles() []string {
	files := []string{}

	files = append(files, "/home/ewgra/go/src/github.com/b/test.go")

	return files
}