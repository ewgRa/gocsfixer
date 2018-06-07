package gocsfixer

import (
	"encoding/json"
	"io/ioutil"
)

func ReadResults(fileName string) ([]*Result, error) {
	var results []*Result

	content, err := ioutil.ReadFile(fileName)

	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(content, &results)

	if err != nil {
		return nil, err
	}

	return results, nil
}