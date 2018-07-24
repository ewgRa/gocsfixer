package gocsfixer

import (
	"encoding/json"
	"io/ioutil"
)

// ReadResults for import it in external tools for easy parse results and use them in CI for example
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
