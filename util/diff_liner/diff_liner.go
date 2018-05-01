package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"strconv"
	"regexp"
	"errors"
	"encoding/json"
	"github.com/ewgRa/gocsfixer/util/diff_liner/response"
)

// DiffLiner parse github diff content, taken by github api call (e.g. curl -H "Accept: application/vnd.github.v3.diff.json" https://api.github.com/repos/ru-de/faq/pulls/377)
// As output it respond with json data, that says which one position in diff belongs to changed line in file.
// This position needed for call comment API endpoint - https://developer.github.com/v3/pulls/comments/#create-a-comment
func main() {
	scanner := bufio.NewScanner(os.Stdin)

	var file string
	var jsonData []byte
	var err error
	var fileLineNumber, diffLineNumber int

	var lineRegexp = regexp.MustCompile(`^@@ -\d+,\d+ \+(?P<newLines>\d+),\d+ @@`)

	for scanner.Scan() {
		line := scanner.Text()

		diffLineNumber++

		if strings.HasPrefix(line, "+++") {
			diffLineNumber = -1

			file = line[4:]

			if strings.HasPrefix(file, "\"") {
				file, err = strconv.Unquote(file)

				if err != nil {
					panic(err)
				}

			}

			file = strings.Replace(file, "b/", "", 1)
		} else if strings.HasPrefix(line, "@@") {
			match := lineRegexp.FindStringSubmatch(line)

			if len(match) != 2 {
				panic(errors.New("Can't parse line numbers"))
			}

			fileLineNumber, err = strconv.Atoi(match[1])

			if err != nil {
				panic(err)
			}
		} else {
			if strings.HasPrefix(line, " ") {
				fileLineNumber++
			} else if strings.HasPrefix(line, "+") {
				data := &response.DiffLinerResponse{
					File: file,
					Line: fileLineNumber,
					DiffLine: diffLineNumber,
				}

				jsonData, err = json.Marshal(data)

				if err != nil {
					panic(err)
				}

				fmt.Println(string(jsonData))

				fileLineNumber++
			}
		}
	}

	if err := scanner.Err(); err != nil {
		panic(err)

	}
}
