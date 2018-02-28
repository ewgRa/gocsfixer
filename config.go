package main

import (
	"github.com/ewgRa/gocsfixer/fixers"
	"fmt"
	"io/ioutil"
	"gopkg.in/yaml.v2"
)

func NewCsFixerConfig(recommend, lint, fix bool, csFixer fixers.CsFixer) *CsFixerConfig {
	return &CsFixerConfig{recommend: recommend, lint: lint, fix: fix, csFixer: csFixer}
}

type CsFixerConfig struct {
	recommend bool
	lint bool
	fix bool
	csFixer fixers.CsFixer
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
	return fmt.Sprintf("Config, fix: %v, csFixer: %v", c.fix, c.csFixer)
}

// FIXME XXX: config validation?
func readConfig() []*CsFixerConfig {
	// FIXME XXX: from command line
	content, err := ioutil.ReadFile(".gocsfixer.yml")
	if nil != err  {
		panic("Error reading config file") // FIXME XXX: better message an react
	}

	config := map[string]map[string]map[string]string{}

	err = yaml.Unmarshal(content, config)

	if nil != err  {
		panic("Error parsing config file") // FIXME XXX: better message an react
	}


	configs := []*CsFixerConfig{}

	for name, settings := range config["fixers"] {
		if name == "no_new_line_before_error_check" {
			configs = append(configs, NewCsFixerConfig(extractBool(settings["recommend"]), extractBool(settings["lint"]), extractBool(settings["fix"]), &fixers.NoNewLineBeforeErrorCheck{}))
		} else {
			panic("Unknown fixer") // FIXME XXX: better message an react
		}
	}

	return  configs
}

func extractBool(v string) bool {
	if v == "true" {
		return true
	} else if v == "false" {
		return false
	}

	panic("Not a bool value") // FIXME XXX: better message an react
}