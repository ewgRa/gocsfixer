package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/ewgRa/gocsfixer"
	"github.com/ewgRa/gocsfixer/fixers"
	"io/ioutil"
	"os"
	"reflect"
)

func main() {
	returnValue := 0

	recommend := flag.Bool("recommend", false, "Show recommends")
	lint := flag.Bool("lint", false, "Perform lint checks")
	fix := flag.Bool("fix", false, "Perform fixes")
	configFile := flag.String("config", ".gocsfixer.yml", "Config file")

	flag.Parse()

	configs, err := gocsfixer.ReadConfig(*configFile)

	if nil != err {
		handleError(err)
	}

	var results []*gocsfixer.Result

	for _, file := range getFiles() {
		c, err := ioutil.ReadFile(file)
		content := string(c)
		fixContent := content

		if nil != err {
			handleError(fmt.Errorf("Error reading file %s", file))
		}

		for _, config := range configs {
			if *fix {
				if config.Fix() {
					fixer, ok := config.CsFixer.(fixers.Fixer)

					if !ok {
						handleError(fmt.Errorf("%s is not a fixer, check your config", reflect.TypeOf(config.CsFixer)))
					}

					fixContent, err = fixer.Fix(fixContent)

					if nil != err {
						handleError(fmt.Errorf("Error during fix file %s: %s", file, err))
					}
				}
			} else if *recommend || *lint {
				linter, ok := config.CsFixer.(fixers.Linter)

				if !ok {
					handleError(fmt.Errorf("%s is not a linter, check your config", reflect.TypeOf(config.CsFixer)))
				}

				lintMode := *lint && config.Lint()
				if lintMode || (*recommend && config.Recommend()) {
					problems, err := linter.Lint(content)

					if nil != err {
						handleError(fmt.Errorf("Error during lint file %s: %s", file, err))
					}

					if lintMode && len(problems) != 0 {
						returnValue = 1
					}

					for _, problem := range problems {
						problemType := "recommendation"

						if lintMode {
							problemType = "error"
						}

						results = append(results, &gocsfixer.Result{
							Type: problemType,
							File: file,
							Line: problem.Position.Line,
							Text: problem.Text,
						})
					}
				}
			}
		}

		if fixContent != content {
			err = ioutil.WriteFile(file, []byte(fixContent), 0644)
		}
	}

	if len(results) > 0 {
		data, err := json.Marshal(results)

		if err != nil {
			handleError(err)
		}

		fmt.Println(string(data))
	} else {
		fmt.Println("[]")
	}

	os.Exit(returnValue)
}

// Read files for processing
func getFiles() []string {
	files := []string{}

	scanner := bufio.NewScanner(os.Stdin)

	for scanner.Scan() {
		file := scanner.Text()

		if file != "" {
			files = append(files, scanner.Text())
		}
	}

	if nil != scanner.Err() {
		handleError(scanner.Err())
	}

	return files
}

func handleError(e error) {
	panic(e)
}
