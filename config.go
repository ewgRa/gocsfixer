package gocsfixer

import (
	"github.com/ewgRa/gocsfixer/fixers"
	"fmt"
	"io/ioutil"
	"gopkg.in/yaml.v2"
)

func NewCsFixerConfig(recommend, lint, fix bool, csFixer fixers.CsFixer) *CsFixerConfig {
	return &CsFixerConfig{recommend: recommend, lint: lint, fix: fix, CsFixer: csFixer}
}

type CsFixerConfig struct {
	recommend bool
	lint bool
	fix bool
	CsFixer fixers.CsFixer
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

	if nil != err  {
		return []*CsFixerConfig{}, fmt.Errorf("Can't read config file %s: %s", file, err)
	}

	config := map[string]map[string]map[string]string{}

	err = yaml.Unmarshal(content, config)

	if nil != err  {
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

		configs = append(configs, NewCsFixerConfig(recommend, lint, fix, createFunc()))
	}

	return  configs, nil
}

func extractBool(v string) (bool, error) {
	if v == "true" {
		return true, nil
	} else if v == "false" {
		return false, nil
	}

	return false, fmt.Errorf("%s not a bool value", v)
}