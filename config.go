package gocsfixer

import (
	"fmt"
	"github.com/ewgRa/gocsfixer/fixers"
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

func NewCsFixerConfig(recommend, lint, fix bool, csFixer fixers.CsFixer) *CsFixerConfig {
	return &CsFixerConfig{recommend: recommend, lint: lint, fix: fix, CsFixer: csFixer}
}

type CsFixerConfig struct {
	recommend bool
	lint      bool
	fix       bool
	CsFixer   fixers.CsFixer
}

func (c *CsFixerConfig) Recommend() bool {
	return c.recommend
}

func (c *CsFixerConfig) Lint() bool {
	return c.lint
}

func (c *CsFixerConfig) Fix() bool {
	return c.fix
}

func (c *CsFixerConfig) String() string {
	return fmt.Sprintf("Config, fix: %v, csFixer: %v", c.fix, c.CsFixer)
}

// FIXME XXX: config validation?
func ReadConfig(file string) ([]*CsFixerConfig, error) {
	content, err := ioutil.ReadFile(file)

	if nil != err {
		return []*CsFixerConfig{}, fmt.Errorf("Can't read config file %s: %s", file, err)
	}

	config := map[string]map[string]map[string]interface{}{}

	err = yaml.Unmarshal(content, config)

	if nil != err {
		return []*CsFixerConfig{}, fmt.Errorf("Can't parse config file %s: %s", file, err)
	}

	configs := []*CsFixerConfig{}

	for name, settings := range config["fixers"] {
		createFunc, ok := fixers.FixersMap[name]

		if !ok {
			return []*CsFixerConfig{}, fmt.Errorf("Unknown fixer %s", name)
		}

		recommend, err := extractBool(settings["recommend"])

		if nil != err {
			return []*CsFixerConfig{}, fmt.Errorf("Wrong fixer %s recommend setting: %s", name, settings["recommend"])
		}

		lint, err := extractBool(settings["lint"])

		if nil != err {
			return []*CsFixerConfig{}, fmt.Errorf("Wrong fixer %s lint setting: %s", name, settings["lint"])
		}

		fix, err := extractBool(settings["fix"])

		if nil != err {
			return []*CsFixerConfig{}, fmt.Errorf("Wrong fixer %s fix setting: %s", name, settings["fix"])
		}

		var options fixers.FixerOptions

		if settingsOptions, ok := settings["options"]; ok {
			options, ok = settingsOptions.(map[interface{}]interface{})

			if !ok {
				return []*CsFixerConfig{}, fmt.Errorf("Wrong fixer %s options settings: %s", name, settings["options"])
			}
		}

		fixer, err := createFunc(options)

		if err != nil {
			return []*CsFixerConfig{}, fmt.Errorf("Can't create fixer %s, error is: %s", name, err)
		}

		configs = append(configs, NewCsFixerConfig(recommend, lint, fix, fixer))
	}

	return configs, nil
}

func extractBool(v interface{}) (bool, error) {
	value, ok := v.(bool)

	if !ok {
		return false, fmt.Errorf("%s not a bool value", v)
	}

	return value, nil
}
