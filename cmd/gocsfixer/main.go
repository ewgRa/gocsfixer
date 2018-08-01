package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"reflect"

	"github.com/ewgRa/gocsfixer"
	"github.com/ewgRa/gocsfixer/fixers"
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

		if nil != err {
			handleError(fmt.Errorf("Error reading file %s", file))
		}

		content := string(c)
		fixContent := content

		for _, config := range configs {
			if *fix {
				if !config.Fix() {
					continue
				}

				fixer, ok := config.CsFixer.(fixers.Fixer)

				if !ok {
					handleError(fmt.Errorf("%s is not a fixer, check your config", reflect.TypeOf(config.CsFixer)))
				}

				fixContent, err = fixer.Fix(fixContent)

				if nil != err {
					handleError(fmt.Errorf("Error during fix file %s: %s", file, err))
				}
			} else if *recommend || *lint {
				lintMode := *lint && config.Lint()
				recommendMode := *recommend && config.Recommend()

				if !lintMode && !recommendMode {
					continue
				}

				linter, ok := config.CsFixer.(fixers.Linter)

				if !ok {
					handleError(fmt.Errorf("%s is not a linter, check your config", reflect.TypeOf(config.CsFixer)))
				}

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

		if fixContent != content {
			err = ioutil.WriteFile(file, []byte(fixContent), 0644)
		}
	}

	res := "[]"

	if len(results) > 0 {
		data, err := json.Marshal(results)

		if err != nil {
			handleError(err)
		}

		res = string(data)
	}

	fmt.Println(res)
	os.Exit(returnValue)
}

// Read file names for processing
func getFiles() []string {
	files := []string{}

	scanner := bufio.NewScanner(os.Stdin)

	for scanner.Scan() {
		file := scanner.Text()

		if file == "" {
			continue
		}

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
