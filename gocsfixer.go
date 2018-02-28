package main

import (
	"github.com/ewgRa/gocsfixer/fixers"
	"fmt"
	"io/ioutil"
	"os"
	"bufio"
	"reflect"
	"flag"
)

func main() {
	returnValue := 0

	recommend := flag.Bool("recommend", false, "Show recommends")
	lint := flag.Bool("lint", false, "Perform lint checks")
	fix := flag.Bool("fix", false, "Perform fixes")

	flag.Parse()

	configs, err := readConfig()

	if nil != err {
		handleError(err)
	}

	for _, file := range getFiles() {
		fmt.Println("File", file)
		c, err := ioutil.ReadFile(file)
		content := string(c)

		if nil != err  {
			handleError(fmt.Errorf("Error reading file %s", file))
		}

		for _, config := range configs {
			if *fix {
				if config.Fix() {
					fixer, ok := config.csFixer.(fixers.Fixer)

					if !ok {
						handleError(fmt.Errorf("%s is not a fixer, check your config", reflect.TypeOf(config.csFixer)))
					}

					fmt.Println("fixer", fixer)
				}
			} else if *recommend || *lint {
				linter, ok := config.csFixer.(fixers.Linter)

				if !ok {
					handleError(fmt.Errorf("%s is not a linter, check your config", reflect.TypeOf(config.csFixer)))
				}

				lintMode := *lint && config.Lint()
				if (lintMode || (*recommend && config.Recommend())) {
					problems, err := linter.Lint(content)

					if nil != err {
						handleError(fmt.Errorf("Error during lint file %s: %s", file, err))
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

	scanner := bufio.NewScanner(os.Stdin)

	for scanner.Scan() {
		files = append(files, scanner.Text())
	}

	if nil != scanner.Err() {
		handleError(scanner.Err())
	}

	return files
}

func handleError(e error) {
	panic(e)
}